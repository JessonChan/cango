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
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/JessonChan/canlog"
)

var sessionMap = map[string]sessionValue{}

type sessionValue struct {
	timeOut time.Time
	value   []byte
}

const cangoSessionKey = "__cango_session_id"

func (wr *WebRequest) SessionGet(key string, value interface{}) {
	if sc, err := wr.Cookie(cangoSessionKey); err == nil {
		if v, ok := sessionMap[sc.Value]; ok {
			err := gob.NewDecoder(bytes.NewReader(v.value)).Decode(&value)
			if err != nil {
				canlog.CanError(err)
			}
			return
		}
	}
}

func (wr *WebRequest) SessionPut(key string, value interface{}, timeOut ...time.Time) {
	bb := &bytes.Buffer{}
	err := gob.NewEncoder(bb).Encode(value)
	if err != nil {
		canlog.CanError(err)
		return
	}
	mapKey := "random" + fmt.Sprintf("%d", time.Now().UnixNano())
	sessionMap[mapKey] = sessionValue{value: bb.Bytes()}
	wr.SetCookie(&http.Cookie{
		Name:     cangoSessionKey,
		Value:    mapKey,
		Path:     "/",
		Expires:  time.Now().AddDate(0, 0, 7),
		MaxAge:   int(time.Hour * 24 * 7),
		Secure:   false,
		HttpOnly: false,
		SameSite: 0,
	})
}
