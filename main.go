package main

import "fmt"

func main() {
	_, format, openMS, convertMS := GetGrid("./samples/1.jpg")
	fmt.Println(format, openMS, convertMS)
}
