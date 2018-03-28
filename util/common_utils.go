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

func GetPageList(p, pages int) ([]int) {
	pageList := make([]int, 0)
	pList := 0

	if pages - p > 10 {
		if p == 1 {
			pList = p + 10
		} else {
			pList = p + 6
		}
	} else {
		pList = pages
	}

	if pages <= 10 {
		for i := 1; i <= pList; i++ {
			pageList = append(pageList, i)
		}
	} else {
		if p <= 10 {
			for i := pList - 10; i <= pList; i++ {
				pageList = append(pageList, i)
			}
		} else {
			t := p + 5
			if t > pages {
				t = pages
			}
			for i := p - 5; i <= t; i++ {
				pageList = append(pageList, i)
			}
		}
	}
	return pageList
}