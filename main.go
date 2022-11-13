package web

import (
	"fmt"
	"strings"
)

func test() {
	testStr := "/usr/home"
	for _, s := range strings.Split(testStr, "/") {
		fmt.Println(s)
	}
}