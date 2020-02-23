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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		showHelp()
		return
	}
	switch strings.ToLower(os.Args[1]) {
	case "demo":
		files := map[string][]byte{}
		files["controller/can.go"] = []byte{112, 97, 99, 107, 97, 103, 101, 32, 99, 111, 110, 116, 114, 111, 108, 108, 101, 114, 10, 10, 105, 109, 112, 111, 114, 116, 32, 40, 10, 9, 34, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 74, 101, 115, 115, 111, 110, 67, 104, 97, 110, 47, 99, 97, 110, 103, 111, 34, 10, 41, 10, 10, 116, 121, 112, 101, 32, 67, 97, 110, 67, 116, 114, 108, 32, 115, 116, 114, 117, 99, 116, 32, 123, 10, 9, 99, 97, 110, 103, 111, 46, 85, 82, 73, 32, 96, 118, 97, 108, 117, 101, 58, 34, 47, 34, 96, 10, 125, 10, 10, 118, 97, 114, 32, 95, 32, 61, 32, 99, 97, 110, 103, 111, 46, 82, 101, 103, 105, 115, 116, 101, 114, 85, 82, 73, 40, 38, 67, 97, 110, 67, 116, 114, 108, 123, 125, 41, 10, 10, 102, 117, 110, 99, 32, 40, 99, 32, 42, 67, 97, 110, 67, 116, 114, 108, 41, 32, 73, 110, 100, 101, 120, 40, 112, 115, 32, 115, 116, 114, 117, 99, 116, 32, 123, 10, 9, 99, 97, 110, 103, 111, 46, 85, 82, 73, 32, 96, 118, 97, 108, 117, 101, 58, 34, 47, 59, 47, 105, 110, 100, 101, 120, 47, 105, 110, 100, 101, 120, 46, 104, 116, 109, 108, 34, 96, 10, 125, 41, 32, 105, 110, 116, 101, 114, 102, 97, 99, 101, 123, 125, 32, 123, 10, 9, 114, 101, 116, 117, 114, 110, 32, 99, 97, 110, 103, 111, 46, 77, 111, 100, 101, 108, 86, 105, 101, 119, 123, 10, 9, 9, 84, 112, 108, 58, 32, 32, 32, 34, 105, 110, 100, 101, 120, 46, 104, 116, 109, 108, 34, 44, 10, 9, 9, 77, 111, 100, 101, 108, 58, 32, 34, 72, 101, 108, 108, 111, 44, 99, 97, 110, 103, 111, 33, 34, 44, 10, 9, 125, 10, 125, 10, 10, 102, 117, 110, 99, 32, 40, 99, 32, 42, 67, 97, 110, 67, 116, 114, 108, 41, 32, 65, 112, 105, 40, 112, 115, 32, 115, 116, 114, 117, 99, 116, 32, 123, 10, 9, 99, 97, 110, 103, 111, 46, 85, 82, 73, 32, 96, 118, 97, 108, 117, 101, 58, 34, 47, 97, 112, 105, 59, 97, 112, 105, 46, 106, 115, 111, 110, 34, 96, 10, 125, 41, 32, 105, 110, 116, 101, 114, 102, 97, 99, 101, 123, 125, 32, 123, 10, 9, 114, 101, 116, 117, 114, 110, 32, 109, 97, 112, 91, 115, 116, 114, 105, 110, 103, 93, 115, 116, 114, 105, 110, 103, 123, 10, 9, 9, 34, 104, 101, 108, 108, 111, 34, 58, 32, 34, 99, 97, 110, 103, 111, 34, 44, 10, 9, 125, 10, 125, 10}
		files["filter/visit.go"] = []byte{112, 97, 99, 107, 97, 103, 101, 32, 102, 105, 108, 116, 101, 114, 10, 10, 105, 109, 112, 111, 114, 116, 32, 40, 10, 9, 34, 110, 101, 116, 47, 104, 116, 116, 112, 34, 10, 10, 9, 47, 47, 32, 34, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 74, 101, 115, 115, 111, 110, 67, 104, 97, 110, 47, 99, 97, 110, 108, 111, 103, 34, 10, 10, 9, 34, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 74, 101, 115, 115, 111, 110, 67, 104, 97, 110, 47, 99, 97, 110, 103, 111, 34, 10, 41, 10, 10, 116, 121, 112, 101, 32, 86, 105, 115, 105, 116, 70, 105, 108, 116, 101, 114, 32, 115, 116, 114, 117, 99, 116, 32, 123, 10, 9, 99, 97, 110, 103, 111, 46, 70, 105, 108, 116, 101, 114, 10, 9, 47, 47, 32, 99, 111, 110, 116, 114, 111, 108, 108, 101, 114, 10, 9, 47, 47, 32, 42, 67, 97, 110, 67, 116, 114, 108, 10, 125, 10, 10, 118, 97, 114, 32, 95, 32, 61, 32, 99, 97, 110, 103, 111, 46, 82, 101, 103, 105, 115, 116, 101, 114, 70, 105, 108, 116, 101, 114, 40, 38, 86, 105, 115, 105, 116, 70, 105, 108, 116, 101, 114, 123, 125, 41, 10, 10, 102, 117, 110, 99, 32, 40, 118, 32, 42, 86, 105, 115, 105, 116, 70, 105, 108, 116, 101, 114, 41, 32, 80, 114, 101, 72, 97, 110, 100, 108, 101, 40, 114, 101, 113, 32, 42, 104, 116, 116, 112, 46, 82, 101, 113, 117, 101, 115, 116, 41, 32, 105, 110, 116, 101, 114, 102, 97, 99, 101, 123, 125, 32, 123, 10, 9, 47, 47, 32, 99, 97, 110, 108, 111, 103, 46, 67, 97, 110, 68, 101, 98, 117, 103, 40, 114, 101, 113, 46, 77, 101, 116, 104, 111, 100, 44, 32, 114, 101, 113, 46, 85, 82, 76, 46, 80, 97, 116, 104, 41, 10, 9, 114, 101, 116, 117, 114, 110, 32, 116, 114, 117, 101, 10, 125, 10}
		files["go.mod"] = []byte{109, 111, 100, 117, 108, 101, 32, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 74, 101, 115, 115, 111, 110, 67, 104, 97, 110, 47, 98, 121, 116, 101, 115, 114, 111, 108, 108, 105, 110, 103, 10, 10, 103, 111, 32, 49, 46, 49, 50, 10, 10, 114, 101, 113, 117, 105, 114, 101, 32, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 74, 101, 115, 115, 111, 110, 67, 104, 97, 110, 47, 99, 97, 110, 103, 111, 32, 118, 48, 46, 48, 46, 48, 10, 10, 114, 101, 112, 108, 97, 99, 101, 32, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 74, 101, 115, 115, 111, 110, 67, 104, 97, 110, 47, 99, 97, 110, 103, 111, 32, 61, 62, 32, 47, 85, 115, 101, 114, 115, 47, 106, 101, 115, 115, 111, 110, 99, 104, 97, 110, 47, 99, 111, 100, 101, 47, 103, 111, 47, 112, 97, 116, 104, 47, 115, 114, 99, 47, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 74, 101, 115, 115, 111, 110, 67, 104, 97, 110, 47, 99, 97, 110, 103, 111, 10}
		files["go.sum"] = []byte{103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 74, 101, 115, 115, 111, 110, 67, 104, 97, 110, 47, 99, 97, 110, 108, 111, 103, 32, 118, 48, 46, 48, 46, 48, 45, 50, 48, 50, 48, 48, 50, 50, 49, 48, 57, 53, 53, 50, 50, 45, 51, 49, 49, 55, 55, 98, 49, 99, 101, 54, 99, 48, 32, 104, 49, 58, 75, 78, 57, 104, 55, 83, 55, 88, 84, 100, 87, 100, 68, 83, 115, 73, 48, 89, 112, 117, 84, 81, 99, 113, 117, 65, 69, 55, 66, 69, 88, 53, 55, 99, 77, 82, 48, 103, 89, 49, 67, 104, 99, 61, 10, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 74, 101, 115, 115, 111, 110, 67, 104, 97, 110, 47, 99, 97, 110, 108, 111, 103, 32, 118, 48, 46, 48, 46, 48, 45, 50, 48, 50, 48, 48, 50, 50, 49, 48, 57, 53, 53, 50, 50, 45, 51, 49, 49, 55, 55, 98, 49, 99, 101, 54, 99, 48, 47, 103, 111, 46, 109, 111, 100, 32, 104, 49, 58, 72, 113, 82, 68, 76, 99, 121, 69, 75, 86, 122, 73, 66, 122, 110, 113, 83, 79, 109, 74, 106, 110, 97, 70, 43, 103, 72, 115, 79, 70, 66, 75, 65, 55, 105, 97, 84, 73, 68, 52, 102, 88, 65, 61, 10, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 74, 101, 115, 115, 111, 110, 67, 104, 97, 110, 47, 106, 115, 117, 110, 32, 118, 48, 46, 48, 46, 48, 45, 50, 48, 50, 48, 48, 50, 48, 54, 48, 52, 49, 56, 48, 48, 45, 48, 50, 56, 101, 102, 49, 56, 53, 49, 54, 98, 100, 32, 104, 49, 58, 99, 83, 97, 81, 77, 99, 118, 54, 98, 89, 100, 99, 86, 104, 70, 66, 88, 76, 105, 111, 104, 86, 98, 117, 109, 77, 84, 101, 110, 88, 90, 119, 88, 76, 79, 78, 98, 47, 97, 120, 119, 105, 111, 61, 10, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 74, 101, 115, 115, 111, 110, 67, 104, 97, 110, 47, 106, 115, 117, 110, 32, 118, 48, 46, 48, 46, 48, 45, 50, 48, 50, 48, 48, 50, 48, 54, 48, 52, 49, 56, 48, 48, 45, 48, 50, 56, 101, 102, 49, 56, 53, 49, 54, 98, 100, 47, 103, 111, 46, 109, 111, 100, 32, 104, 49, 58, 79, 75, 50, 47, 113, 112, 71, 50, 100, 86, 99, 81, 116, 67, 119, 66, 50, 75, 88, 104, 119, 71, 77, 98, 98, 73, 118, 43, 56, 84, 82, 56, 53, 110, 49, 78, 75, 118, 101, 82, 104, 57, 89, 61, 10}
		files["main.go"] = []byte{112, 97, 99, 107, 97, 103, 101, 32, 109, 97, 105, 110, 10, 10, 105, 109, 112, 111, 114, 116, 32, 40, 10, 9, 34, 110, 101, 116, 47, 104, 116, 116, 112, 34, 10, 9, 34, 111, 115, 34, 10, 10, 9, 47, 47, 32, 34, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 74, 101, 115, 115, 111, 110, 67, 104, 97, 110, 47, 99, 97, 110, 108, 111, 103, 34, 10, 10, 9, 34, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 74, 101, 115, 115, 111, 110, 67, 104, 97, 110, 47, 99, 97, 110, 108, 111, 103, 34, 10, 10, 9, 34, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 74, 101, 115, 115, 111, 110, 67, 104, 97, 110, 47, 99, 97, 110, 103, 111, 34, 10, 41, 10, 10, 116, 121, 112, 101, 32, 67, 97, 110, 67, 116, 114, 108, 32, 115, 116, 114, 117, 99, 116, 32, 123, 10, 9, 99, 97, 110, 103, 111, 46, 85, 82, 73, 32, 96, 118, 97, 108, 117, 101, 58, 34, 47, 34, 96, 10, 125, 10, 10, 118, 97, 114, 32, 95, 32, 61, 32, 99, 97, 110, 103, 111, 46, 82, 101, 103, 105, 115, 116, 101, 114, 85, 82, 73, 40, 38, 67, 97, 110, 67, 116, 114, 108, 123, 125, 41, 10, 10, 102, 117, 110, 99, 32, 40, 99, 32, 42, 67, 97, 110, 67, 116, 114, 108, 41, 32, 73, 110, 100, 101, 120, 40, 112, 115, 32, 115, 116, 114, 117, 99, 116, 32, 123, 10, 9, 99, 97, 110, 103, 111, 46, 85, 82, 73, 32, 96, 118, 97, 108, 117, 101, 58, 34, 47, 59, 47, 105, 110, 100, 101, 120, 47, 105, 110, 100, 101, 120, 46, 104, 116, 109, 108, 34, 96, 10, 125, 41, 32, 105, 110, 116, 101, 114, 102, 97, 99, 101, 123, 125, 32, 123, 10, 9, 114, 101, 116, 117, 114, 110, 32, 99, 97, 110, 103, 111, 46, 77, 111, 100, 101, 108, 86, 105, 101, 119, 123, 10, 9, 9, 84, 112, 108, 58, 32, 32, 32, 34, 105, 110, 100, 101, 120, 46, 104, 116, 109, 108, 34, 44, 10, 9, 9, 77, 111, 100, 101, 108, 58, 32, 34, 72, 101, 108, 108, 111, 44, 99, 97, 110, 103, 111, 33, 34, 44, 10, 9, 125, 10, 125, 10, 10, 116, 121, 112, 101, 32, 80, 97, 103, 101, 32, 115, 116, 114, 117, 99, 116, 32, 123, 10, 9, 79, 102, 102, 115, 101, 116, 32, 105, 110, 116, 10, 9, 83, 105, 122, 101, 32, 32, 32, 105, 110, 116, 10, 125, 10, 10, 102, 117, 110, 99, 32, 40, 99, 32, 42, 67, 97, 110, 67, 116, 114, 108, 41, 32, 65, 112, 105, 40, 112, 115, 32, 115, 116, 114, 117, 99, 116, 32, 123, 10, 9, 99, 97, 110, 103, 111, 46, 85, 82, 73, 32, 96, 118, 97, 108, 117, 101, 58, 34, 47, 123, 110, 97, 109, 101, 125, 47, 123, 121, 101, 97, 114, 125, 47, 97, 112, 105, 46, 106, 115, 111, 110, 34, 96, 10, 9, 99, 97, 110, 103, 111, 46, 80, 111, 115, 116, 77, 101, 116, 104, 111, 100, 10, 9, 78, 97, 109, 101, 32, 32, 115, 116, 114, 105, 110, 103, 10, 9, 89, 101, 97, 114, 32, 32, 105, 110, 116, 10, 9, 67, 111, 108, 111, 114, 32, 115, 116, 114, 105, 110, 103, 10, 9, 80, 97, 103, 101, 10, 125, 41, 32, 105, 110, 116, 101, 114, 102, 97, 99, 101, 123, 125, 32, 123, 10, 9, 114, 101, 116, 117, 114, 110, 32, 109, 97, 112, 91, 115, 116, 114, 105, 110, 103, 93, 105, 110, 116, 101, 114, 102, 97, 99, 101, 123, 125, 123, 10, 9, 9, 34, 110, 97, 109, 101, 34, 58, 32, 32, 112, 115, 46, 78, 97, 109, 101, 44, 10, 9, 9, 34, 97, 103, 101, 34, 58, 32, 32, 32, 112, 115, 46, 89, 101, 97, 114, 44, 10, 9, 9, 34, 99, 111, 108, 111, 114, 34, 58, 32, 112, 115, 46, 67, 111, 108, 111, 114, 44, 10, 9, 9, 34, 115, 105, 122, 101, 34, 58, 32, 32, 112, 115, 46, 83, 105, 122, 101, 44, 10, 9, 125, 10, 125, 10, 10, 116, 121, 112, 101, 32, 86, 105, 115, 105, 116, 70, 105, 108, 116, 101, 114, 32, 115, 116, 114, 117, 99, 116, 32, 123, 10, 9, 99, 97, 110, 103, 111, 46, 70, 105, 108, 116, 101, 114, 10, 9, 47, 47, 32, 99, 111, 110, 116, 114, 111, 108, 108, 101, 114, 10, 9, 42, 67, 97, 110, 67, 116, 114, 108, 10, 125, 10, 10, 118, 97, 114, 32, 95, 32, 61, 32, 99, 97, 110, 103, 111, 46, 82, 101, 103, 105, 115, 116, 101, 114, 70, 105, 108, 116, 101, 114, 40, 38, 86, 105, 115, 105, 116, 70, 105, 108, 116, 101, 114, 123, 125, 41, 10, 10, 102, 117, 110, 99, 32, 40, 118, 32, 42, 86, 105, 115, 105, 116, 70, 105, 108, 116, 101, 114, 41, 32, 80, 114, 101, 72, 97, 110, 100, 108, 101, 40, 114, 101, 113, 32, 42, 104, 116, 116, 112, 46, 82, 101, 113, 117, 101, 115, 116, 41, 32, 105, 110, 116, 101, 114, 102, 97, 99, 101, 123, 125, 32, 123, 10, 9, 99, 97, 110, 108, 111, 103, 46, 67, 97, 110, 68, 101, 98, 117, 103, 40, 114, 101, 113, 46, 77, 101, 116, 104, 111, 100, 44, 32, 114, 101, 113, 46, 85, 82, 76, 46, 80, 97, 116, 104, 41, 10, 9, 114, 101, 116, 117, 114, 110, 32, 116, 114, 117, 101, 10, 125, 10, 10, 102, 117, 110, 99, 32, 109, 97, 105, 110, 40, 41, 32, 123, 10, 9, 47, 47, 32, 99, 97, 110, 108, 111, 103, 46, 83, 101, 116, 87, 114, 105, 116, 101, 114, 40, 99, 97, 110, 108, 111, 103, 46, 78, 101, 119, 70, 105, 108, 101, 87, 114, 105, 116, 101, 114, 40, 34, 47, 116, 109, 112, 47, 99, 97, 110, 103, 111, 45, 97, 112, 112, 46, 108, 111, 103, 34, 41, 44, 32, 34, 65, 112, 112, 34, 41, 10, 9, 99, 97, 110, 108, 111, 103, 46, 83, 101, 116, 87, 114, 105, 116, 101, 114, 40, 111, 115, 46, 83, 116, 100, 111, 117, 116, 44, 32, 34, 65, 112, 112, 34, 41, 10, 9, 99, 97, 110, 103, 111, 46, 73, 110, 105, 116, 76, 111, 103, 103, 101, 114, 40, 99, 97, 110, 108, 111, 103, 46, 71, 101, 116, 76, 111, 103, 103, 101, 114, 40, 41, 46, 87, 114, 105, 116, 101, 114, 40, 41, 41, 10, 9, 99, 97, 110, 103, 111, 46, 10, 9, 9, 78, 101, 119, 67, 97, 110, 40, 41, 46, 10, 9, 9, 82, 101, 103, 84, 112, 108, 70, 117, 110, 99, 40, 34, 98, 117, 116, 116, 111, 110, 86, 97, 108, 117, 101, 34, 44, 32, 102, 117, 110, 99, 40, 41, 32, 115, 116, 114, 105, 110, 103, 32, 123, 10, 9, 9, 9, 114, 101, 116, 117, 114, 110, 32, 34, 99, 108, 105, 99, 107, 32, 109, 101, 33, 34, 10, 9, 9, 125, 41, 46, 10, 9, 9, 82, 117, 110, 40, 99, 97, 110, 103, 111, 46, 65, 100, 100, 114, 123, 80, 111, 114, 116, 58, 32, 56, 48, 48, 56, 125, 44, 32, 99, 97, 110, 103, 111, 46, 79, 112, 116, 115, 123, 84, 112, 108, 83, 117, 102, 102, 105, 120, 58, 32, 91, 93, 115, 116, 114, 105, 110, 103, 123, 34, 46, 104, 116, 109, 108, 34, 125, 44, 32, 68, 101, 98, 117, 103, 84, 112, 108, 58, 32, 116, 114, 117, 101, 125, 41, 10, 125, 10}
		files["static/css/index.css"] = []byte{46, 32, 114, 101, 98, 101, 99, 99, 97, 112, 117, 114, 112, 108, 101, 123, 10, 32, 32, 32, 32, 99, 111, 108, 111, 114, 58, 32, 114, 101, 98, 101, 99, 99, 97, 112, 117, 114, 112, 108, 101, 59, 10, 125}
		files["static/favicon.ico"] = []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 16, 0, 0, 0, 16, 8, 6, 0, 0, 0, 31, 243, 255, 97, 0, 0, 2, 60, 73, 68, 65, 84, 56, 141, 133, 147, 91, 72, 147, 97, 0, 134, 63, 195, 234, 34, 130, 168, 46, 130, 48, 130, 192, 240, 38, 161, 155, 174, 138, 32, 232, 50, 10, 247, 167, 81, 33, 22, 166, 86, 12, 65, 49, 148, 148, 130, 186, 200, 80, 8, 67, 178, 130, 200, 74, 230, 191, 50, 115, 182, 245, 87, 104, 82, 90, 233, 194, 169, 149, 121, 152, 139, 153, 90, 185, 105, 30, 182, 166, 255, 225, 233, 98, 120, 154, 154, 31, 60, 23, 223, 233, 249, 224, 227, 125, 69, 130, 44, 10, 76, 178, 8, 73, 178, 96, 134, 147, 85, 235, 57, 42, 175, 98, 254, 90, 36, 38, 89, 132, 76, 149, 34, 95, 72, 178, 152, 114, 116, 151, 96, 182, 239, 228, 81, 123, 46, 254, 224, 0, 0, 193, 233, 113, 110, 126, 76, 94, 81, 34, 36, 89, 112, 171, 37, 21, 221, 208, 0, 104, 240, 148, 115, 189, 241, 8, 63, 198, 58, 233, 245, 183, 252, 87, 32, 201, 2, 33, 201, 130, 103, 157, 133, 0, 168, 250, 52, 158, 209, 86, 92, 67, 10, 33, 53, 64, 69, 123, 222, 202, 130, 116, 91, 12, 170, 62, 197, 120, 200, 199, 169, 234, 205, 228, 215, 237, 165, 168, 201, 68, 150, 178, 107, 222, 193, 40, 178, 149, 120, 114, 94, 238, 38, 211, 17, 71, 113, 83, 226, 156, 160, 162, 61, 15, 0, 165, 167, 116, 201, 23, 204, 246, 88, 250, 70, 90, 1, 240, 5, 250, 241, 5, 188, 232, 134, 198, 249, 231, 59, 194, 130, 214, 161, 23, 0, 75, 126, 88, 242, 211, 13, 12, 7, 188, 24, 134, 129, 165, 163, 128, 36, 107, 52, 253, 99, 95, 0, 48, 219, 99, 195, 2, 239, 159, 207, 0, 20, 190, 59, 188, 72, 240, 176, 237, 2, 0, 175, 220, 183, 145, 100, 65, 90, 205, 86, 12, 195, 96, 44, 52, 140, 36, 71, 133, 5, 110, 191, 19, 128, 114, 87, 214, 130, 203, 73, 214, 53, 204, 236, 93, 105, 56, 136, 36, 11, 110, 124, 56, 14, 64, 163, 183, 146, 99, 143, 215, 98, 235, 42, 70, 84, 125, 189, 10, 192, 72, 112, 144, 203, 111, 14, 144, 110, 139, 225, 218, 219, 67, 244, 250, 155, 153, 25, 101, 45, 103, 200, 116, 196, 49, 56, 222, 13, 64, 105, 243, 105, 234, 251, 238, 161, 244, 148, 34, 78, 60, 89, 135, 107, 72, 33, 114, 124, 31, 109, 163, 230, 91, 209, 236, 220, 48, 12, 166, 180, 32, 0, 206, 1, 27, 63, 39, 122, 73, 169, 222, 24, 206, 129, 36, 11, 114, 95, 239, 161, 204, 153, 198, 29, 103, 6, 5, 117, 251, 102, 163, 124, 169, 126, 63, 15, 92, 57, 220, 119, 101, 241, 107, 194, 13, 128, 219, 239, 36, 163, 118, 219, 92, 144, 150, 35, 91, 137, 167, 203, 247, 30, 221, 208, 208, 116, 21, 0, 71, 79, 9, 137, 214, 232, 133, 73, 92, 14, 75, 199, 69, 126, 79, 122, 48, 219, 99, 105, 244, 90, 208, 116, 149, 179, 181, 219, 23, 38, 49, 178, 137, 243, 73, 173, 217, 130, 219, 239, 68, 211, 85, 254, 170, 147, 220, 253, 116, 110, 113, 153, 18, 44, 139, 235, 28, 73, 74, 245, 38, 146, 172, 171, 151, 172, 243, 63, 171, 252, 163, 137, 13, 228, 3, 54, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130}
		files["static/js/index.js"] = []byte{102, 117, 110, 99, 116, 105, 111, 110, 32, 115, 104, 111, 119, 76, 105, 110, 107, 40, 41, 32, 123, 10, 32, 32, 32, 32, 100, 111, 99, 117, 109, 101, 110, 116, 46, 103, 101, 116, 69, 108, 101, 109, 101, 110, 116, 66, 121, 73, 100, 40, 34, 108, 105, 110, 107, 34, 41, 46, 115, 101, 116, 65, 116, 116, 114, 105, 98, 117, 116, 101, 40, 34, 115, 116, 121, 108, 101, 34, 44, 32, 34, 100, 105, 115, 112, 108, 97, 121, 58, 98, 108, 111, 99, 107, 34, 41, 10, 125, 10, 10, 102, 117, 110, 99, 116, 105, 111, 110, 32, 114, 101, 113, 117, 101, 115, 116, 74, 83, 79, 78, 40, 117, 114, 108, 44, 32, 100, 97, 116, 97, 41, 32, 123, 10, 32, 32, 32, 32, 118, 97, 114, 32, 104, 116, 116, 112, 82, 101, 113, 117, 101, 115, 116, 32, 61, 32, 110, 101, 119, 32, 88, 77, 76, 72, 116, 116, 112, 82, 101, 113, 117, 101, 115, 116, 40, 41, 59, 10, 32, 32, 32, 32, 104, 116, 116, 112, 82, 101, 113, 117, 101, 115, 116, 46, 111, 112, 101, 110, 40, 39, 80, 79, 83, 84, 39, 44, 32, 117, 114, 108, 44, 32, 102, 97, 108, 115, 101, 41, 59, 10, 32, 32, 32, 32, 104, 116, 116, 112, 82, 101, 113, 117, 101, 115, 116, 46, 115, 101, 116, 82, 101, 113, 117, 101, 115, 116, 72, 101, 97, 100, 101, 114, 40, 34, 67, 111, 110, 116, 101, 110, 116, 45, 116, 121, 112, 101, 34, 44, 32, 34, 97, 112, 112, 108, 105, 99, 97, 116, 105, 111, 110, 47, 120, 45, 119, 119, 119, 45, 102, 111, 114, 109, 45, 117, 114, 108, 101, 110, 99, 111, 100, 101, 100, 34, 41, 59, 10, 32, 32, 32, 32, 104, 116, 116, 112, 82, 101, 113, 117, 101, 115, 116, 46, 115, 101, 110, 100, 40, 100, 97, 116, 97, 41, 59, 10, 32, 32, 32, 32, 105, 102, 32, 40, 104, 116, 116, 112, 82, 101, 113, 117, 101, 115, 116, 46, 114, 101, 97, 100, 121, 83, 116, 97, 116, 101, 32, 61, 61, 61, 32, 52, 32, 38, 38, 32, 104, 116, 116, 112, 82, 101, 113, 117, 101, 115, 116, 46, 115, 116, 97, 116, 117, 115, 32, 61, 61, 61, 32, 50, 48, 48, 41, 32, 123, 10, 32, 32, 32, 32, 32, 32, 32, 32, 97, 108, 101, 114, 116, 40, 104, 116, 116, 112, 82, 101, 113, 117, 101, 115, 116, 46, 114, 101, 115, 112, 111, 110, 115, 101, 84, 101, 120, 116, 41, 10, 32, 32, 32, 32, 125, 10, 125}
		files["view/index.html"] = []byte{60, 33, 68, 79, 67, 84, 89, 80, 69, 32, 104, 116, 109, 108, 62, 10, 60, 104, 116, 109, 108, 32, 108, 97, 110, 103, 61, 34, 101, 110, 34, 62, 10, 60, 104, 101, 97, 100, 62, 10, 32, 32, 32, 32, 60, 109, 101, 116, 97, 32, 99, 104, 97, 114, 115, 101, 116, 61, 34, 85, 84, 70, 45, 56, 34, 62, 10, 32, 32, 32, 32, 60, 116, 105, 116, 108, 101, 62, 99, 97, 110, 103, 111, 60, 47, 116, 105, 116, 108, 101, 62, 10, 32, 32, 32, 32, 60, 108, 105, 110, 107, 32, 114, 101, 108, 61, 34, 115, 116, 121, 108, 101, 115, 104, 101, 101, 116, 34, 32, 104, 114, 101, 102, 61, 34, 47, 115, 116, 97, 116, 105, 99, 47, 99, 115, 115, 47, 105, 110, 100, 101, 120, 46, 99, 115, 115, 34, 62, 10, 32, 32, 32, 32, 60, 115, 99, 114, 105, 112, 116, 32, 115, 114, 99, 61, 34, 47, 115, 116, 97, 116, 105, 99, 47, 106, 115, 47, 105, 110, 100, 101, 120, 46, 106, 115, 34, 32, 116, 121, 112, 101, 61, 34, 116, 101, 120, 116, 47, 106, 97, 118, 97, 115, 99, 114, 105, 112, 116, 34, 32, 99, 104, 97, 114, 115, 101, 116, 61, 34, 117, 116, 102, 45, 56, 34, 62, 60, 47, 115, 99, 114, 105, 112, 116, 62, 10, 60, 47, 104, 101, 97, 100, 62, 10, 60, 98, 111, 100, 121, 62, 10, 60, 98, 117, 116, 116, 111, 110, 32, 111, 110, 99, 108, 105, 99, 107, 61, 34, 115, 104, 111, 119, 76, 105, 110, 107, 40, 41, 34, 62, 123, 123, 98, 117, 116, 116, 111, 110, 86, 97, 108, 117, 101, 125, 125, 60, 47, 98, 117, 116, 116, 111, 110, 62, 10, 123, 123, 36, 117, 114, 108, 58, 61, 34, 47, 99, 97, 110, 103, 111, 47, 50, 48, 50, 48, 47, 97, 112, 105, 46, 106, 115, 111, 110, 34, 125, 125, 10, 123, 123, 36, 100, 97, 116, 97, 58, 61, 34, 99, 111, 108, 111, 114, 61, 114, 101, 100, 38, 111, 102, 102, 115, 101, 116, 61, 49, 38, 115, 105, 122, 101, 61, 49, 48, 34, 125, 125, 10, 60, 98, 117, 116, 116, 111, 110, 32, 111, 110, 99, 108, 105, 99, 107, 61, 34, 114, 101, 113, 117, 101, 115, 116, 74, 83, 79, 78, 40, 123, 123, 36, 117, 114, 108, 125, 125, 44, 123, 123, 36, 100, 97, 116, 97, 125, 125, 41, 34, 62, 10, 32, 32, 32, 32, 65, 80, 73, 61, 62, 40, 117, 114, 108, 58, 123, 123, 36, 117, 114, 108, 125, 125, 44, 100, 97, 116, 97, 58, 123, 123, 36, 100, 97, 116, 97, 125, 125, 41, 10, 60, 47, 98, 117, 116, 116, 111, 110, 62, 10, 60, 97, 32, 105, 100, 61, 34, 108, 105, 110, 107, 34, 32, 116, 97, 114, 103, 101, 116, 61, 34, 95, 98, 108, 97, 110, 107, 34, 32, 104, 114, 101, 102, 61, 34, 104, 116, 116, 112, 58, 47, 47, 119, 119, 119, 46, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 74, 101, 115, 115, 111, 110, 67, 104, 97, 110, 47, 99, 97, 110, 103, 111, 34, 32, 99, 108, 97, 115, 115, 61, 34, 114, 101, 98, 101, 99, 99, 97, 112, 117, 114, 112, 108, 101, 34, 10, 32, 32, 32, 115, 116, 121, 108, 101, 61, 34, 100, 105, 115, 112, 108, 97, 121, 58, 32, 110, 111, 110, 101, 34, 62, 123, 123, 46, 125, 125, 60, 47, 97, 62, 10, 60, 47, 98, 111, 100, 121, 62, 10, 60, 47, 104, 116, 109, 108, 62}
		for k, v := range files {
			path := "./" + k
			dir := path[:func() int {
				idx := strings.LastIndex(path, "/")
				if idx == -1 {
					return len(path)
				}
				return idx
			}()]
			_ = os.MkdirAll(dir, os.ModePerm)
			newFile, err := os.Create("./" + k)
			if err != nil {
				log.Println(err)
			}
			_, err = newFile.Write(v)
			if err != nil {
				log.Println(err)
			}
			_ = newFile.Close()
		}
	case "bootstrap":
		files := map[string]string{}
		var fileNames []string
		_ = filepath.Walk("./demo", func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			name := strings.TrimPrefix(path, "demo/")
			fileNames = append(fileNames, name)
			files[name] = bytes(func() []byte {
				bs, _ := ioutil.ReadFile(path)
				return bs
			}())
			return nil
		})
		sort.Strings(fileNames)
		selfName := "./cango-cli.go"
		self := func() string {
			self, err := ioutil.ReadFile(selfName)
			if err != nil {
				log.Println(err)
				return ""
			}
			return string(self)
		}()
		if self == "" {
			return
		}
		prefixMarker := "files := map[string][]byte{}"
		suffixMarker := "for k, v := range files {"
		selfPrefix := strings.Index(self, prefixMarker)
		selfSuffix := strings.Index(self, suffixMarker)
		content := self[:selfPrefix] + func() string {
			builder := "files := map[string][]byte{}\n"
			for _, name := range fileNames {
				builder = builder + repeat(8) + "files[\"" + name + "\"] = " + files[name] + "\n"
			}
			return builder[:len(builder)]
		}() + repeat(8) + self[selfSuffix:]
		err := ioutil.WriteFile(selfName, []byte(content), 0666)
		if err != nil {
			log.Println(err)
		}
	}
}
func showHelp() {
	fmt.Println("cango-cli demo")
}
func bytes(bs []byte) string {
	var builder = "[]byte{"
	for _, b := range bs {
		builder = builder + (fmt.Sprintf("%d, ", b))
	}
	builder = builder[0:len(builder)-2] + "}"
	return builder
}
func repeat(c int) (rs string) {
	for i := 0; i < c; i++ {
		rs = rs + " "
	}
	return rs
}
