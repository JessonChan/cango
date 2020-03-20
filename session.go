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
	"net/http"
	"time"

	"github.com/JessonChan/canlog"
	"github.com/gorilla/sessions"
)

var gorillaStore sessions.Store

func SetGorillaSessionStore(store sessions.Store) {
	gorillaStore = store
}

func newCookieSession(key, secure string) sessions.Store {
	switch len(secure) {
	case 16, 24, 32:
		return sessions.NewCookieStore([]byte(key), []byte(secure))
	default:
		return sessions.NewCookieStore([]byte(key))
	}
}

const cangoSessionKey = "__cango_session_id"

func SessionGet(r *http.Request, key string, value interface{}) {
	gs, _ := gorillaStore.Get(r, cangoSessionKey)
	if i, ok := gs.Values[key]; ok {
		err := gob.NewDecoder(bytes.NewReader(i.([]byte))).Decode(value)
		if err != nil {
			canlog.CanError(err)
		}
	}
}

func (wr *WebRequest) SessionGet(key string, value interface{}) {
	SessionGet(wr.Request, key, value)
}

func SessionPut(r *http.Request, rw http.ResponseWriter, key string, value interface{}, timeOut ...time.Time) {
	bb := &bytes.Buffer{}
	err := gob.NewEncoder(bb).Encode(value)
	if err != nil {
		canlog.CanError(err)
		return
	}
	gs, _ := gorillaStore.Get(r, cangoSessionKey)
	gs.Values[key] = bb.Bytes()
	err = gorillaStore.Save(r, rw, gs)
	if err != nil {
		canlog.CanError(err)
	}
}

func (wr *WebRequest) SessionPut(key string, value interface{}, timeOut ...time.Time) {
	SessionPut(wr.Request, wr.ResponseWriter, key, value, timeOut...)
}
