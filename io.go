package svmlight

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (file *SVMFile) WriteBinary(label int, fileName string) {
	newF := &SVMFile{}
	for _, instance := range file.Instances {
		newInstance := instance
		if newInstance.Label == label {
			newInstance.Label = 1
		} else {
			newInstance.Label = 1
		}
		newF.Instances = append(newF.Instances, newInstance)
	}
	newF.WriteSVMFile(fileName)
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
	file := SVMFile{}
	file.FileName = fileName
	file.Instances = []SVMInstance{}

	for buffer, _, err := reader.ReadLine(); err == nil; buffer, _, err = reader.ReadLine() {
		newInstance := SVMInstance{}
		lbl_feat_split := strings.Split(string(buffer), " ")
		if len(lbl_feat_split) <= 1 {
			panic("This file does not follow the conventional SVMlight format")
		}
		newInstance.Label, _ = strconv.Atoi(lbl_feat_split[0])
		newInstance.Features = make(map[int]SVMFeature)
		featureSplit := lbl_feat_split[1:]
		for _, featurePair := range featureSplit {
			kvPair := strings.Split(featurePair, ":")
			if len(kvPair) != 2 {
				panic("This file does not follow the conventional SVMlight format")
			}
			key, _ := strconv.Atoi(kvPair[0])
			value, _ := strconv.ParseFloat(kvPair[1], 64)
			newInstance.Features[key] = SVMFeature(value)
		}
		file.Instances = append(file.Instances, newInstance)
	}
	return &file
}

func (file *SVMFile) WriteSVMFile(fileName string) {
	output := bytes.NewBufferString("")
	fs := ""
	for _, instance := range file.Instances {
		output.WriteString(fmt.Sprintf("%d", instance.Label))
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

func ParseModelFile(fileName string) *ModelFile {
	modelFile := ModelFile{}
	modelFile.Kernel = KernelInfo{}
	modelFile.Alphas = []Alpha{}
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
	for buffer, _, err := reader.ReadLine(); err == nil; buffer, _, err = reader.ReadLine() {
		str := strings.TrimRight(strings.Split(string(buffer), "#")[0], " ")
		switch lineNumber {
		case 0:
			break
		case 1:
			modelFile.Kernel.T, _ = strconv.Atoi(str)
		case 2:
			modelFile.Kernel.D, _ = strconv.Atoi(str)
		case 3:
			modelFile.Kernel.G, _ = strconv.Atoi(str)
		case 4:
			modelFile.Kernel.S, _ = strconv.Atoi(str)
		case 5:
			modelFile.Kernel.R, _ = strconv.Atoi(str)
		case 6:
			modelFile.Kernel.U, _ = strconv.Atoi(str)
		case 7:
			modelFile.MaxFeatureIndex, _ = strconv.Atoi(str)
		case 8:
			modelFile.TrainingSize, _ = strconv.Atoi(str)
		case 9:
			break
		case 10:
			modelFile.Bias, _ = strconv.ParseFloat(str, 64)
		default:
			alpha := Alpha{}
			alpha.Vector = make(map[int]SVMFeature)
			lbl_feat_split := strings.Split(str, " ")
			if len(lbl_feat_split) <= 1 {
				panic("This file does not follow the conventional SVMlight format")
			}
			alpha.Weight, _ = strconv.ParseFloat(lbl_feat_split[0], 64)
			featureSplit := lbl_feat_split[1:]
			for _, featurePair := range featureSplit {
				kvPair := strings.Split(featurePair, ":")
				if len(kvPair) != 2 {
					panic("This file does not follow the conventional SVMlight format")
				}
				key, _ := strconv.Atoi(kvPair[0])
				value, _ := strconv.ParseFloat(kvPair[1], 64)
				alpha.Vector[key] = SVMFeature(value)
			}
			modelFile.Alphas = append(modelFile.Alphas, alpha)
		}
		lineNumber++
	}
	return &modelFile
}

func ParseClassificationFile(fileName string) *ClassificationFile {
	file := ClassificationFile{}
	file.Results = []ClassificationResult{}

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

	for buffer, _, err := reader.ReadLine(); err == nil; buffer, _, err = reader.ReadLine() {
		val, _ := strconv.ParseFloat(string(buffer), 64)
		file.Results = append(file.Results, ClassificationResult(val))
	}

	return &file
}

func (f *ModelFile) WriteModelFile(fileName string) {
	output := bytes.NewBufferString(f.ToString())

	sysfile, _ := os.Create(fileName)
	defer func() {
		if err := sysfile.Close(); err != nil {
			panic(err)
		}
	}()
	sysfile.Write(output.Bytes())
}

func (f *ClassificationFile) WriteClassificationFile(fileName string) {
	output := bytes.NewBufferString(f.ToString())
	sysfile, _ := os.Create(fileName)
	defer func() {
		if err := sysfile.Close(); err != nil {
			panic(err)
		}
	}()
	sysfile.Write(output.Bytes())
}
