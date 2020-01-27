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
	"net/http"
	"sync"
	"time"

	"github.com/JessonChan/cango"
)

type SnowController struct {
	cango.URI `value:"/weather/{day}/how-heavy/{heavy}"`
	Day       string
	Heavy     string
}

func (s *SnowController) Snowing(param struct {
	cango.URI `value:"/snowing/{color}.json"`
	Color     int
}) interface{} {
	return map[string]interface{}{
		"day":   s.Day,
		"heavy": s.Heavy,
		"color": param.Color,
		"show":  true,
	}
}

func (s *SnowController) Wind(param struct {
	cango.URI `value:"/winding.html"`
}) interface{} {
	return cango.ModelView{Tpl: "/view/wind.tpl", Model: map[string]string{
		"Day":   s.Day,
		"Heavy": s.Heavy,
	}}
}

func (s *SnowController) Raining(param struct {
	cango.URI `value:"/raining/{dropSize}.json"`
	cango.PostMethod
	DropSize int
}) interface{} {
	return map[string]interface{}{
		"day":   s.Day,
		"heavy": s.Heavy,
		"color": param.DropSize,
		"show":  false,
	}
}

type LogFilter struct {
	cango.Filter
}

func (m *LogFilter) PreHandle(r *http.Request) interface{} {
	fmt.Println(time.Now().Format("2006-03-02 15:04:05"), r.Method, r.URL)
	return nil
}

var count struct {
	sync.RWMutex
	cnt int64
}

type StatFilter struct {
	cango.Filter
	*SnowController
}

func (m *StatFilter) PreHandle(r *http.Request) interface{} {
	count.Lock()
	count.cnt++
	fmt.Println("web page visit count total:", count.cnt)
	count.Unlock()
	return nil
}

func main() {
	can := cango.NewCan()
	can.Filter(&LogFilter{}, &SnowController{}).
		Filter(&StatFilter{}).
		Route(&SnowController{}).
		RouteWithPrefix("/api/v2", &SnowController{}).
		Run(cango.Addr{Port: 8081})
}
