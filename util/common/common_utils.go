package common

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"time"
	"github.com/neal1991/gshark/vars"
)

// Utility function for producing a hex encoded sha1 hash for a string.
func HashFor(name string) string {
	h := sha1.New()
	h.Write([]byte(name))
	return hex.EncodeToString(h.Sum(nil))
}

func GetPreAndNext(p int) (currentPage int, pre int, next int) {
	if p < 1 {
		currentPage = 1
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

func GetPageAndPagesByTotalPages(page, totalPages int) (int, int) {
	var pages int
	if int(totalPages)%vars.PAGE_SIZE == 0 {
		pages = int(totalPages) / vars.PAGE_SIZE
	} else {
		pages = int(totalPages)/vars.PAGE_SIZE + 1
	}

	if page >= pages {
		page = pages
	}

	if page < 1 {
		page = 1
	}
	return page,pages
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
