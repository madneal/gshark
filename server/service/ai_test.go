package service

import (
	"fmt"
	"testing"
)

func TestQuestion(t *testing.T) {
	result := Question("You are a security operation engineer, you are expected to assistant.please judge if the below content contains sensitive information, and the sensitive information \" +\n\t\t\t\t\"could be exploited. Just answer yes or no",
		"# PopWin_MeiTuan\n仿美团做的一个下拉选择页，类似于电商app中筛选距离的下拉菜单。很常见。\n#效果图\n![](https://github.com/reallin/PopWin_MeiTuan/blob/master/cam.gif)\n#功能点\n* popwindow下拉列表\n* 加载progressBar的生成")
	fmt.Println(result)
}
