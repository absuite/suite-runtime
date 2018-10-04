package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/ggoop/gocast"
	"github.com/spf13/viper"
)

type Config struct {
	App       App
	Endpoints map[string]Endpoint
}
type App struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
}
type Endpoint struct {
	*sql.DB
	Name     string
	Driver   string
	Host     string
	Port     uint
	Database string
	Username string
	Password string
	Instance string
	Windows  bool
}

type ParamValue struct {
	Name  string
	Value interface{}
}
type ErrorResult struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func (m *Endpoint) Open() (err error) {

	//odbc:server=localhost;user id=sa;password={foobar }
	//sqlserver://sa:mypass@localhost:1234?database=master&connection+timeout=30
	//server=localhost\\SQLExpress;user id=sa;database=master;app name=MyAppNam
	//server=localhost;user id=sa;database=master;app name=MyAppName
	//Provider=SQLOLEDB;Data Source=host;Initial Catalog=Database;user id=Username;password=Password

	query := url.Values{}
	query.Add("encrypt", "disable")
	if m.Database != "" {
		query.Add("database", m.Database)
	}
	query.Add("connection timeout", "0")
	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(m.Username, m.Password),
		Host:     m.Host,
		RawQuery: query.Encode(),
	}
	if m.Port > 0 {
		u.Host = fmt.Sprintf("%s:%d", m.Host, m.Port)
	}
	if m.Instance != "" {
		u.Path = m.Instance
	}
	m.DB, err = sql.Open("sqlserver", u.String())
	if err != nil {
		return err
	}
	return nil
}
func toError(w http.ResponseWriter, txt string, err error) {
	fmt.Printf(txt, err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	result := &ErrorResult{Msg: txt, Code: 0}
	resultByte, err := json.Marshal(result)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(resultByte)
	}

}
func (db *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	oldPath := r.URL.Path
	newPath := ""
	if len(oldPath) > len(db.Name) {
		newPath = oldPath[len(db.Name):len(oldPath)]
	}
	if newPath == "" {
		toError(w, `bad request! path is empty`, errors.New("bad request!"))
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		toError(w, `read body error!`, err)
		return
	}
	input := make(map[string]interface{})
	if len(body) > 0 {
		if err = json.Unmarshal(body, &input); err != nil {
			toError(w, "unmarshal body error!", err)
			return
		}
	} else if len(r.Form) > 0 {
		for k, v := range r.Form {
			input[k] = v[0]
		}
	}
	jsonString, err := handRequest(db, newPath, input)
	if err != nil {
		toError(w, err.Error(), err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}
func handRequest(db *Endpoint, spName string, args map[string]interface{}) ([]byte, error) {
	// 连接数据库
	if err := db.Open(); err != nil {
		return nil, err
	}
	defer db.Close()
	ctx := context.Background()
	err := db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	//查询存储过程参数
	cmd := fmt.Sprintf(`select o.object_id,substring(p.name,2,100) as name,lt.name as type,lt.max_length as length 
	from  sys.objects as o 
	left join sys.parameters as p on p.object_id=o.object_id
	left join sys.types as lt on p.system_type_id=lt.user_type_id
	where o.type='P' and o.name='%s'
	Order By p.parameter_id`, spName)
	spParams, err := execQuery(db, cmd)
	if err != nil {
		return nil, err
	}
	if spParams == nil || len(spParams) == 0 {
		return nil, errors.New(fmt.Sprintf("can not found proc :%s", spName))
	}
	paramsIn := []ParamValue{}
	var paramName string
	for _, pv := range spParams {
		if pv["name"] == nil {
			break
		}
		paramName = pv["name"].(string)
		for ik, iv := range args {
			if paramName != ik {
				continue
			}
			if iv != nil {
				paramsIn = append(paramsIn, ParamValue{Name: paramName, Value: iv})
			}
		}
	}
	// 执行SQL语句
	maps, err := execProc(db, spName, paramsIn)
	if err != nil {
		return nil, err
	}
	jsonString, err := json.Marshal(maps)
	return jsonString, err
}
func execProc(db *Endpoint, spName string, params []ParamValue) ([]map[string]interface{}, error) {
	cmd := fmt.Sprintf("exec %s", spName)
	paramValues := make([]interface{}, 0)
	for i, k := range params {
		if i == 0 {
			cmd = fmt.Sprintf("%v @%v=@p%v", cmd, k.Name, i+1)
		} else {
			cmd = fmt.Sprintf("%v,@%v=@p%v", cmd, k.Name, i+1)
		}
		paramValues = append(paramValues, k.Value)
	}
	stmt, err := db.Prepare(cmd)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	result := paramValues[:]
	rows, err := stmt.Query(result...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToMap(rows)
}
func execQuery(db *Endpoint, query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToMap(rows)
}
func rowsToMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	var maps = make([]map[string]interface{}, 0)
	colNames, _ := rows.Columns()
	var cols = make([]interface{}, len(colNames))
	for i := 0; i < len(colNames); i++ {
		cols[i] = new(interface{})
	}
	for rows.Next() {
		err := rows.Scan(cols...)
		if err != nil {
			return nil, err
		}
		var rowMap = make(map[string]interface{})
		for i := 0; i < len(colNames); i++ {
			rowMap[colNames[i]] = convertRow(*(cols[i].(*interface{})))
		}
		maps = append(maps, rowMap)
	}
	return maps, nil
}
func convertRow(row interface{}) interface{} {
	switch row.(type) {
	case int:
		return gocast.ToInt(row)
	case string:
		return gocast.ToString(row)
	case []byte:
		return gocast.ToString(row)
	case bool:
		return gocast.ToBool(row)
	}
	return row
}

const cmdRoot = "config"

func main() {
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix(cmdRoot)
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetConfigName(cmdRoot)
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Fatal error when reading %s config file:%s", cmdRoot, err)
		os.Exit(1)
	}
	var config Config
	viper.Unmarshal(&config)

	if config.App.Port == "" {
		config.App.Port = "8080"
	}
	if len(config.Endpoints) > 0 {
		for k, v := range config.Endpoints {
			endpoint := &Endpoint{Name: "/" + k + "/", Database: v.Database, Host: v.Host, Windows: v.Windows, Username: v.Username, Password: v.Password, Port: v.Port}
			http.Handle(endpoint.Name, endpoint)
			log.Printf("proxy http://127.0.0.1:%s%s => %s", config.App.Port, endpoint.Name, endpoint.Host)
		}
	} else {
		var v Endpoint
		viper.UnmarshalKey("endpoints", &v)
		if v.Host != "" && v.Database != "" {
			endpoint := &Endpoint{Name: "/", Database: v.Database, Host: v.Host, Windows: v.Windows, Username: v.Username, Password: v.Password, Port: v.Port}
			http.Handle(endpoint.Name, endpoint)
			log.Printf("proxy http://127.0.0.1:%s%s => %s", config.App.Port, endpoint.Name, endpoint.Host)
		}
	}
	log.Printf("server listening at: http://127.0.0.1:%s", config.App.Port)
	http.ListenAndServe(":"+config.App.Port, nil)
}
