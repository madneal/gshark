package util

import (
	"fmt"
	"github.com/madneal/gshark/logger"
	"github.com/madneal/gshark/vars"
	"net/http"
	url2 "net/url"
	"strings"
	"time"
)

func GetPreAndNext(p int) (currentPage int, pre int, next int) {
	if p < 1 {
		currentPage = 1
		p = 1
	} else {
		currentPage = p
	}

	if p <= 1 {
		pre = 1
	} else {
		pre = p - 1
	}
	next = p + 1
	return currentPage, pre, next
}

func GetPageList(p, step, pages int) []int {
	pageList := make([]int, 0)
	startIndex := p - step
	endIndex := p + step

	if startIndex < 1 && endIndex <= pages {
		startIndex = 1
		endIndex = startIndex + 2*step
	} else if startIndex >= 1 && endIndex > pages {
		endIndex = pages
		startIndex = pages - 2*step
	} else if startIndex < 1 && endIndex > pages {
		startIndex = 1
		endIndex = pages
	}

	if startIndex < 1 {
		startIndex = 1
	}

	if endIndex > pages {
		endIndex = pages
	}

	for i := startIndex; i <= endIndex; i++ {
		pageList = append(pageList, i)
	}

	return pageList
}

func GetLastPage(pageList *[]int) int {
	lastPage := 0
	if len(*pageList) >= 1 {
		lastPage = (*pageList)[len(*pageList)-1]
	}
	return lastPage
}

func GetPageAndPagesByCount(page, count int) (int, int) {
	var pages int
	if count%vars.PAGE_SIZE == 0 {
		pages = count / vars.PAGE_SIZE
	} else {
		pages = count/vars.PAGE_SIZE + 1
	}

	if page >= pages {
		page = pages
	}

	if page < 1 {
		page = 1
	}
	return page, pages
}

func FormatTimeStamp(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// Get the repoName By the url of the repository
func GetRepoNameByUrl(url string) string {
	length := len(strings.Split(url, "/"))
	if length > 1 {
		return strings.Split(url, "/")[length-1]
	} else {
		return ""
	}
}

// Send message to server 酱
func SendMessage(key, title, msg string) {
	url := fmt.Sprintf("https://sc.ftqq.com/%s.send?text=%s&desp=%s", key, title, url2.QueryEscape(msg))
	_, err := http.Get(url)
	if err != nil {
		logger.Log.Error(err)
	}
}
