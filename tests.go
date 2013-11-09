package main

import (
	"./svmlight"
	"fmt"
)

func main() {
	fmt.Println("Compiled fine")
	trainSet := svmlight.ParseSVMFile("data/boxes.train")
	modelFile := svmlight.Learn("data/boxes.train", "data/boxes.model", 0., 0.5, 0)
	classificationFile := svmlight.Classify("data/boxes.train", "data/boxes.model", "data/boxes.classification")
	fmt.Println(modelFile.ToString())
	fmt.Println(classificationFile.Accuracy(trainSet))
}
