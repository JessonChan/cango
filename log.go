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

var logLevel = 0
var logger = log.New(os.Stdout, "CANGO ", log.LstdFlags)

func line(level int, v ...interface{}) {
	if level >= logLevel {
		logger.Println(v...)
	}
}

func debug(v ...interface{}) {
	line(DEBUG, v...)
}

func info(v ...interface{}) {
	line(INFO, v...)
}
func warn(v ...interface{}) {
	line(WARN, v...)
}
func err(v ...interface{}) {
	line(ERROR, v...)
}
func fatal(v ...interface{}) {
	panic(fmt.Sprint(v...))
}
