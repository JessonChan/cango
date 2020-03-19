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
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	numRand "math/rand"
	"net/http"
	"time"

	"github.com/JessonChan/canlog"
)

var sessionStore = &memStore{store: map[string]*sessionValue{}}

type sessionValue struct {
	timeOut time.Time
	values  map[string][]byte
}

const cangoSessionKey = "__cango_session_id"

func (wr *WebRequest) SessionGet(key string, value interface{}) {
	if sc, err := wr.Cookie(cangoSessionKey); err == nil {
		if vs, ok := sessionStore.Get(sc.Value); ok {
			if v, ok := vs.values[key]; ok {
				err := gob.NewDecoder(bytes.NewReader(v)).Decode(value)
				if err != nil {
					canlog.CanError(err)
				}
				return
			}
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
	sid := sessionID()
	sessionStore.Put(sid, &sessionValue{values: map[string][]byte{key: bb.Bytes()}})
	wr.SetCookie(&http.Cookie{
		Name:     cangoSessionKey,
		Value:    sid,
		Path:     "/",
		Expires:  time.Now().AddDate(0, 0, 7),
		MaxAge:   int(time.Hour * 24 * 7),
		Secure:   false,
		HttpOnly: false,
		SameSite: 0,
	})
}

var defaultSessionLength = 32

func sessionID() string {
	b := make([]byte, defaultSessionLength)
	n, err := rand.Read(b)
	if n != len(b) || err != nil {
		getRandBytes(&b)
	}
	return hex.EncodeToString(b)
}

func getRandBytes(b *[]byte) {
	rd := numRand.New(numRand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len(*b); i++ {
		(*b)[i] = byte(rd.Int31n(26) + (func() int32 {
			if rd.Int31n(2) == 0 {
				return 'a'
			}
			return 'A'
		}()))
	}
}
