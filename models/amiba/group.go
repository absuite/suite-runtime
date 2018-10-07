package amibaModels

type Group struct {
	Id       string `xorm:"varchar(200) 'id'"`
	Code     string `xorm:"varchar(200) 'code'"`
	Name     string `xorm:"varchar(200) 'name'"`
	TypeEnum string `xorm:"varchar(200) 'type_enum'"`
	Datas    []GroupData
}
type GroupData struct {
	GroupId  string `xorm:"varchar(200) 'group_id'"`
	Id       string `xorm:"varchar(200) 'id'"`
	Code     string `xorm:"varchar(200) 'code'"`
	Name     string `xorm:"varchar(200) 'name'"`
	TypeEnum string `xorm:"varchar(200) 'type_enum'"`
}

func (s *Group) AddData(item GroupData) {
	s.Datas = append(s.Datas, item)
}
func (s *Group) GetDataCodes() []string {
	if s.Datas == nil || len(s.Datas) == 0 {
		return nil
	}
	codes := make([]string, 0)
	for _, d := range s.Datas {
		if d.Code != "" {
			codes = append(codes, d.Code)
		}
	}
	return codes
}