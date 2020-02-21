package cango

import (
	"io/ioutil"
	"strings"

	"github.com/JessonChan/canlog"
)

// todo 支持对于can全局配置的文件化设定
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
