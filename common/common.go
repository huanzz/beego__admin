package common

import (
	"crypto/md5"
	"encoding/hex"
	"math"
	"strconv"
	"strings"
)

func StrToMd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	res := hex.EncodeToString(h.Sum(nil))
	return res
}

func PwdHash(s string) string {
	return StrToMd5(s)
}
func Paginator(page, pagesize int, nums int64) map[string]interface{} {
	var prepage int
	var nextpage int

	totalPages := int(math.Ceil(float64(nums) / float64(pagesize)))
	if page > totalPages {
		page = totalPages
	}
	if page <= 0 {
		page = 1
	}
	var pages []int
	switch {
	case page >= totalPages - 8 && totalPages > 8:
		start := totalPages -8 + 1
		prepage = page - 1
		nextpage = int(math.Min(float64(totalPages), float64(page+1)))
		pages = make([]int, 8)
		for i,_ := range pages{
			pages[i] = start + i
		}
		prepage = prepage - 1
		nextpage = prepage + 1
	default:
		pages = make([]int, int(math.Min(8, float64(totalPages))))
		for i,_ := range pages{
			pages[i] = i + 1
		}
		prepage = int(math.Max(float64(1), float64(page - 1)))
		nextpage = page + 1
	}
	paginatorMap := make(map[string]interface{})
	paginatorMap["pages"] = pages
	paginatorMap["totalpages"] = totalPages
	paginatorMap["prepage"] = prepage
	paginatorMap["nextpage"] = nextpage
	paginatorMap["currpage"] = page
	return paginatorMap
}

func NumInIds(num int, ids []int) bool {
	for _,v := range ids {
		if v == num {
			return true
		}
	}
	return false
}

func StrToIntArr(str string) (ids []int) {
	list := strings.Split(str,",")
	for _,v := range list{
		menuId,_ :=strconv.Atoi(v)
		ids = append(ids, menuId)
	}
	return ids
}

func GetFullUrl(url string) (res string) {
	var list []string
	urlStr := strings.Split(url,"?")
	if urlStr != nil {
		str := urlStr[0]
		urlList := strings.Split(str,"/")
		for _,v := range urlList{
			if strings.Trim(v," ") != "" {
				list = append(list, v)
			}
		}
	}
	res = strings.Join(list,"/")
	return "/" + res
}
