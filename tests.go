package main

import (
	"./svmlight"
	"fmt"
)

func main() {
	fmt.Println("Compiled fine")
	svmlight.ParseSVMFile("data/boxes.train")
	modelFile := svmlight.Learn("data/boxes.train", "data/boxes.model", 0., 0.5, 0)
	fmt.Println(modelFile.ToString())
}
