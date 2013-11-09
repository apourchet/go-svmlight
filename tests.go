package main

import (
	"./svmlight"
	"fmt"
)

func main() {
	fmt.Println("Working??")
	svmlight.ParseSVMFile("data/boxes.train")
}
