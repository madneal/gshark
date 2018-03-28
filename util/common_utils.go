package util


func GetPreAndNext(p int) (p_ int, pre int, next int) {
	if p < 1 {
		p_ = 1
	} else {
		p_ = p
	}

	if p <= 1 {
		pre = 1
	} else {
		pre = p - 1
	}
	next = p + 1
	return p_, pre, next
}

func GetPageList(p, step, pages int) ([]int) {
	pageList := make([]int, 0)
	startIndex := p - step
	endIndex := p + step

	if startIndex < 1 && endIndex <= pages {
		startIndex = 1
		endIndex = startIndex + 2 * step
	} else if startIndex >= 1 && endIndex > pages {
		endIndex = pages
		startIndex = pages - 2 * step
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