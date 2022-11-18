package web

type user struct {
	Name    string
	Country string
}

func test() {
	//val := user{}
	//err := bindJSON(&val)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("bind result is: ", val)

	// wrong
	//var es []string
	//fmt.Println(es[0])
}

//func bindJSON(val any) error {
//	resource := user{
//		Name:    "a",
//		Country: "b",
//	}
//
//	input, err := json.Marshal(resource)
//	if err != nil {
//		return err
//	}
//	fmt.Println("input is: " + string(input))
//
//	err = json.Unmarshal(input, val)
//	return err
//}
