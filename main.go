package web

import "fmt"

func test() {
	testStr := ":id"
	if string(testStr[0]) == ":" {
		fmt.Println("yes")
	}
	fmt.Println(string(testStr[0]))
}
