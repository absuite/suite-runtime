package routes

import (
	"go.uber.org/dig"
)

func Register(container *dig.Container) {
	registerAmiba(container)
}
