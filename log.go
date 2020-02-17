// Copyright 2020 Cango Author.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cango

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	ZERO = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

var logPrefix = []string{"ZERO", "DEBUG", " INFO", " WARN", "ERROR"}

var logLevel = 0
var logger = log.New(os.Stdout, "CANGO ", log.LstdFlags|log.Lshortfile)

func InitLogger(rw io.Writer, prefix string) {
	logger = log.New(rw, prefix, log.LstdFlags|log.Lshortfile)
}

func canLine(level int, v ...interface{}) {
	if level >= logLevel {
		_ = logger.Output(3, logPrefix[level]+" "+fmt.Sprintln(v...))
	}
}

func canDebug(v ...interface{}) {
	canLine(DEBUG, v...)
}

func canInfo(v ...interface{}) {
	canLine(INFO, v...)
}
func canWarn(v ...interface{}) {
	canLine(WARN, v...)
}
func canError(v ...interface{}) {
	canLine(ERROR, v...)
}
func canFatal(v ...interface{}) {
	panic(fmt.Sprint(v...))
}
