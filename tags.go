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

const (
	uriTagName       string = "value"
	pathFormName     string = "name"
	cookieTagName    string = "cookie"
	pathValueTagName string = "path"
	formValueTagName string = "form"
	sessionTagName   string = "session"
	nameTagName      string = "name"
	holderLen               = 12
)

var cookieHolderKey = makeHolderKey(cookieTagName)
var formPathHolderKey = makeHolderKey(formValueTagName + pathValueTagName)
var sessionHolderKey = makeHolderKey(sessionTagName)

func makeHolderKey(key string) string {
	if len(key) > holderLen {
		return key[0:holderLen]
	}
	newKey := make([]byte, holderLen)
	for i := 0; i < holderLen; i++ {
		if i < len(key) {
			newKey[i] = key[i]
		} else {
			newKey[i] = '-'
		}
	}
	return string(newKey)
}
