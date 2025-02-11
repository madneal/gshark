package service

import (
	"fmt"
	"testing"
)

func TestQuestion(t *testing.T) {
	result := Question("You are a helpful assistant",
		"请编写一个Python函数 find_prime_numbers，该函数接受一个整数 n 作为参数，"+
			"并返回一个包含所有小于 n 的质数（素数）的列表。质数是指仅能被1和其自身整除的正整数，如2, 3, 5, 7等。不要输出非代码的内容。")
	fmt.Println(result)
}
