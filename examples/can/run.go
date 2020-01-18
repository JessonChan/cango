package main

import (
	"github.com/JessonChan/cango"
)

func main() {
	can := cango.Can{}
	can.Run(cango.Addr{Port: 8081})
}
