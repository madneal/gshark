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