package cango

import (
	"io/ioutil"
	"strings"

	"github.com/JessonChan/canlog"
)

func initConfig(configPath string) {
	bs, err := ioutil.ReadFile(configPath)
	if err != nil {
		canlog.CanError(err)
		return
	}
	lines := strings.Split(string(bs), "\n")
	for _, line := range lines {
		vs := strings.Split(line, "=")
		if len(vs) != 2 {
			canlog.CanError("can't parse config line", line)
			continue
		}
	}
}
