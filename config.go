package cango

import (
	"io/ioutil"
	"strings"
)

func initConfig(configPath string) {
	bs, err := ioutil.ReadFile(configPath)
	if err != nil {
		canError(err)
		return
	}
	lines := strings.Split(string(bs), "\n")
	for _, line := range lines {
		vs := strings.Split(line, "=")
		if len(vs) != 2 {
			canError("can't parse config line", line)
			continue
		}
	}
}
