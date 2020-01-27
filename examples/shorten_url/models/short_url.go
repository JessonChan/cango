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
package models

import (
	"fmt"
	"sort"
	"strings"
)

var urlMap = map[string]ShortenUrl{
	"bidu": {1, "bidu", "百度", "http://www.baidu.com", 0},
	"sogo": {2, "sogo", "搜狗", "http://www.sogou.com", 0},
	"bing": {3, "bing", "必应", "http://www.bing.com", 0},
	"soso": {4, "soso", "搜搜", "http://www.soso.com", 0},
}

type ShortenUrl struct {
	Id       int64
	UniqueId string
	Name     string
	Url      string
	Count    int64
}

func getShortUrl(uid string) *ShortenUrl {
	v, ok := urlMap[uid]
	if ok {
		return &v
	}
	return nil
}
func GetUrl(uid string) string {
	su := getShortUrl(uid)
	if su != nil {
		su.Count++
		urlMap[uid] = *su
		return su.Url
	}
	return ""
}

type ShortenList []ShortenUrl

func (u ShortenList) Len() int           { return len(u) }
func (u ShortenList) Less(i, j int) bool { return u[i].Id < u[j].Id }
func (u ShortenList) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }

func GetAll() []ShortenUrl {
	ss := make([]ShortenUrl, len(urlMap))
	i := 0
	for _, v := range urlMap {
		ss[i] = v
		i++
	}
	sort.Sort(ShortenList(ss))
	return ss
}

func Insert(info, name string) *ShortenUrl {
	isUrl := strings.HasPrefix(info, "http://") || strings.HasPrefix(info, "https://")
	if !isUrl {
		return nil
	}
	var su ShortenUrl
	su.Id = int64(len(urlMap) + 1)
	su.Url = info
	// todo 短缩算法
	su.UniqueId = fmt.Sprint(len(urlMap))
	su.Name = name
	urlMap[su.UniqueId] = su
	return &su
}
