package main

import (
	"exec"
	"fmt"
)

func Learn(trainFile, modelFile string, c float64, j float64, d int) {
	t := 0
	if d != 0 {
		t := 1
	}
	if c == 0. {
		c = 1.
	}
	if j == 0. {
		j = 1.
	}
	out, _ := exec.Command("svm_learn", "-c", fmt.Sprintf("%f", c),
		"-j", fmt.Sprintf("%f", j),
		"-t", fmt.Sprintf("%f", t),
		"-d", fmt.Sprintf("%f", d))
	return fmt.Sprintf("%s\n", out)
}

func Classify() {

}

func GetAccuracy() {

}
