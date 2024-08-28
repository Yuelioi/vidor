package main

import (
	"strings"
)

func main() {
	downUrl := "124"

	// 删除?之后的内容
	if i := strings.LastIndex(downUrl, "?"); i >= 0 {
		downUrl = downUrl[:i]
	}

	// 删除最后一个/
	downUrl = strings.TrimSuffix(downUrl, "/")

	splitUrl := strings.Split(downUrl, "/")
	name := splitUrl[len(splitUrl)-1]
	print(name)

}
