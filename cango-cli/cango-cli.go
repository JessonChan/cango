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
package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const pic = `iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAACPElEQVQ4jYWTW0iTYQCGP8PqIoKoLoIwgsDwJqGbroog6DIK96dRIRamVgxBMZSUgrrIUAhDsoLISua/MnO29VdoUlrpwqmVeZiLmVq5aR62pv/h6WJ4mpofPBff6fng431FgiwKTLIISbJghpNV6zkqr2L+WiQmWYRMlSJfSLKYcnSXYLbv5FF7Lv7gAADB6XFufkxeUSIkWXCrJRXd0ABo8JRzvfEIP8Y66fW3/FcgyQIhyYJnnYUAqPo0ntFWXEMKITVARXveyoJ0WwyqPsV4yMep6s3k1+2lqMlElrJr3sEospV4cl7uJtMRR3FT4pygoj0PAKWndMkXzPZY+kZaAfAF+vEFvOiGxvnnO8KC1qEXAEt+WPLTDQwHvBiGgaWjgCRrNP1jXwAw22PDAu+fzwAUvju8SPCw7QIAr9y3kWRBWs1WDMNgLDSMJEeFBW6/E4ByV9aCy0nWNczsXWk4iCQLbnw4DkCjt5Jjj9di6ypGVH29CsBIcJDLbw6Qbovh2ttD9PqbmRllLWfIdMQxON4NQGnzaer77qH0lCJOPFmHa0ghcnwfbaPmW9Hs3DAMprQgAM4BGz8nekmp3hjOgSQLcl/vocyZxh1nBgV1+2ajfKl+Pw9cOdx3ZfFrwg2A2+8ko3bbXJCWI1uJp8v3Ht3Q0HQVAEdPCYnW6IVJXA5Lx0V+T3ow22Np9FrQdJWztdsXJjGyifNJrdmC2+9E01X+qpPc/XRucZkSLIvrHElK9SaSrKuXrPM/q/yjiQ3kAzYAAAAASUVORK5CYII=`

func favicon() {
	_ = ioutil.WriteFile("static/favicon.ico", func() []byte {
		bs, _ := base64.StdEncoding.DecodeString(pic)
		return bs
	}(), os.ModePerm)
}

func main() {
	if len(os.Args) == 1 {
		showHelp()
		return
	}
	switch strings.ToLower(os.Args[1]) {
	case "create":
		dirNames := []string{"controller", "filter", "manager", "model", "static/css", "static/js", "util", "view"}
		for _, dir := range dirNames {
			_ = os.MkdirAll(dir, os.ModePerm)
		}
		favicon()
		_, _ = os.Create("main.go")
		_, _ = os.Create("view/index.html")
	}
}
func showHelp() {
	fmt.Println("cango-cli create")
}
