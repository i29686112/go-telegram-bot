package main

import (
	"fmt"
	"strconv"
)

func getTempString() (string, int) {
	return "hello world yo", 1
}

func main1() {
	var str string

	str, num := getTempString()
	fmt.Println(str)
	fmt.Println(num)

	test := 123

	fmt.Print(test)
	coolFunc := func() { fmt.Println("cool man") }
	coolFunc()

	fmt.Println(main2())

	fmt.Println(main3())
}

func main2() [10]string {

	var array [10]string

	for i := 0; i < 10; i++ {
		array[i] = "123=>" + strconv.Itoa(i)
	}

	return array

}

func main3() map[string]int {

	good := make(map[string]int)

	good["ian"] = 100
	good["jj"] = 200
	return good
}
