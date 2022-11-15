package web

import (
	"fmt"
	"regexp"
)

func test() {
	testStr := "1"
	matchRule := "([0-9]+)"
	matched, err := regexp.Match(matchRule, []byte(testStr))
	if err != nil {
		fmt.Println(err.Error())
	}
	if matched {
		fmt.Println("matched")
	}
}
