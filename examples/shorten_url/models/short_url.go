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
	su.UniqueId = fmt.Sprint(len(urlMap))
	su.Name = name
	urlMap[su.UniqueId] = su
	return &su
}
