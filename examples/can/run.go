package main

import (
	"github.com/JessonChan/cango"
)

type SnowController struct {
	cango.URI `value:"/weather/{day}/how-heavy/{heavy}"`
	Day       string
	Heavy     string
}

func (s *SnowController) Snowing(param struct {
	cango.URI `value:"/snowing/{color}.json"`
	Color     int
}) interface{} {
	return map[string]interface{}{
		"day":   s.Day,
		"heavy": s.Heavy,
		"color": param.Color,
		"show":  true,
	}
}

func main() {
	can := cango.NewCan()
	can.Route(&SnowController{})
	can.Run(cango.Addr{Port: 8081})
}
