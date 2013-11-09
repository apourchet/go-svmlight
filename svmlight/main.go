package main

import (
	"bufio"
	"exec"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type SVMFile struct {
	Instances []SVMInstance
	FileName  string
}
type SVMInstance struct {
	Features map[int]SVMFeature
	Label    string
}

type SVMFeature float64

type KernelInfo struct {
	t, d, g, s, r, u int
}

type Aplha struct {
	Weight float64
	Vector map[int]SVMFeature
}

type ModelFile struct {
	Kernel       KernelInfo
	TrainingSize int
	Bias         float64
	Alphas       []Aplha
}
type ClassificationResult float64

type ClassificationFile struct {
	Results []ClassificationResult
}

func (f *SVMFile) CountLabels(label string) int {
	acc := 0
	for _, instance := range f.Instances {
		if instance.Label == label {
			acc++
		}
	}
	return acc
}

func (f *SVMFile) MakeBinary(label string) {
	for _, instance := range f.Instances {
		if instance.Label == label {
			instance.Label = "1"
		} else {
			instance.Label = "-1"
		}
	}
}

func (f *SVMFile) WriteBinary(label, fileName string) {
	newF := &SVMFile{}
	for _, instance := range file.Instances {
		newInstance := instance
		if newInstance.Label == label {
			newInstance.Label = "1"
		} else {
			newInstance.Label = "-1"
		}
		newF.Instances = append(newF.Instances, newInstance)
	}
	newF.OutputFile(fileName)
}

func ParseSVMFile(fileName string) *SVMFile {
	fi, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	reader := bufio.NewReader(fi)
	file := SVMFile
	file.FileName = fileName
	file.Instances = []SVMInstance{}

	for buffer, err := reader.ReadBytes("\n"); len(buffer) != 0; {
		newInstance := SVMInstance{}
		lbl_feat_split := strings.Split(string(buffer), " ")
		if len(lbl_feat_split) <= 1 {
			panic("This file does not follow the conventional SVMlight format")
		}
		newInstance.Label = lbl_feat_split[0]
		featureSplit := lbl_feat_split[1:]
		for _, featurePair := range featureSplit {
			kvPair := strings.Split(featurePair, ":")
			if len(kvPair) != 2 {
				panic("This file does not follow the conventional SVMlight format")
			}
			key, _ := strconv.Atoi(kvPair[0])
			value, _ := strconv.ParseFloat(kvPair[1], 64)
			newInstance.Features[key] = value
		}
		file.Instances = append(file.Instances, newInstance)
	}
	return &file
}

func (f *SVMFile) WriteSVMFile(fileName string) {
	output := bytes.NewBufferString("")
	fs := ""
	for _, instance := range file.Instances {
		output.WriteString(instance.Label)
		for k, v := range instance.Features {
			fs = fmt.Sprintf(" %d:%f", k, v)
			output.WriteString(fs)
		}
		output.WriteString("\n")
	}
	sysfile, _ := os.Create(fileName)
	defer func() {
		if err := sysfile.Close(); err != nil {
			panic(err)
		}
	}()
	sysfile.Write(output.Bytes())
}

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

func Classify(testFileName, modelFileName, resultFileName string) string {
	out2, _ := exec.Command("svm_classify", "-v", "3", testFileName, modelFileName, resultFileName).Output()
	return fmt.Sprintf("%s\n", out2)
}

func ParseModelFile(fileName string) *ModelFile {
	modelFile := ModelFile{}
	modelFile.Kernel = KernelInfo{}
	modelFile.Alphas = []Aplha{}
	fi, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	reader := bufio.NewReader(fi)
	lineNumber := 0
	for buffer, err := reader.ReadBytes("\n"); len(buffer) != 0; lineNumber++ {
		str := strings.TrimRight(strings.Split(string(buffer), "#"), " ")[0]
		switch lineNumber {
		case 0:
			break
		case 1:
			modelFile.Kernel.t, _ = strconv.Atoi(str)
		case 2:
			modelFile.Kernel.d, _ = strconv.Atoi(str)
		case 3:
			modelFile.Kernel.g, _ = strconv.Atoi(str)
		case 4:
			modelFile.Kernel.s, _ = strconv.Atoi(str)
		case 5:
			modelFile.Kernel.r, _ = strconv.Atoi(str)
		case 6:
			modelFile.Kernel.u, _ = strconv.Atoi(str)
		case 7:
			break
		case 8:
			modelFile.TrainingSize, _ = strconv.Atoi(str)
		case 9:
			break
		case 10:
			modelFile.Bias, _ = strconv.ParseFloat(str, 64)
		default:
			alpha := Alpha{}
			featureLabelSplit := 
			modelFile.Alphas = append(modelFile.Alphas, alpha)
		}

	}
	return &modelFile
}

func ParseResultFile() {

}

func WriteResultFile() {

}
