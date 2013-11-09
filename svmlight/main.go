package svmlight

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
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
	T, D, G, S, R, U int
}

type Alpha struct {
	Weight float64
	Vector map[int]SVMFeature
}

type ModelFile struct {
	Kernel       KernelInfo
	TrainingSize int
	Bias         float64
	Alphas       []Alpha
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

func (file *SVMFile) WriteBinary(label, fileName string) {
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
		newInstance.Label = lbl_feat_split[0]
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

func Learn(trainFile, modelFile string, c float64, j float64, d int) *ModelFile {
	t := 0
	if d != 0 {
		t = 1
	}
	if c == 0. {
		c = 1.
	}
	if j == 0. {
		j = 1.
	}
	// out, _ := exec.Command("svm_learn", "-c", fmt.Sprintf("%f", c), "-v", "1", trainFile, modelFile).Output()

	exec.Command("svm_learn",
		"-c", fmt.Sprintf("%f", c),
		"-j", fmt.Sprintf("%f", j),
		"-t", fmt.Sprintf("%d", t),
		"-d", fmt.Sprintf("%d", d),
		trainFile, modelFile).Output()

	return ParseModelFile(modelFile)
}

func Classify(testFileName, modelFileName, resultFileName string) *ClassificationFile {
	exec.Command("svm_classify", "-v", "3", testFileName, modelFileName, resultFileName).Output()
	return ParseClassificationFile(resultFileName)
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
			break
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

func (k *KernelInfo) ToString() string {
	nu := "empty"
	if k.U != -1 {
		nu = fmt.Sprintf("%d ", k.U)
	}
	return fmt.Sprintf("%d #\n%d #\n%d #\n%d #\n%d #\n%s#\n", k.T, k.D, k.G, k.S, k.R, nu)
}

func (alpha *Alpha) ToString() string {
	output := bytes.NewBufferString(fmt.Sprintf("%f ", alpha.Weight))
	for featureId, featureValue := range alpha.Vector {
		output.WriteString(fmt.Sprintf("%d:%f ", featureId, featureValue))
	}
	return output.String()
}

func (file *ModelFile) ToString() string {
	output := bytes.NewBufferString("SVM-light Version V6.02\n" + file.Kernel.ToString())
	maxIndex := 0
	for _, alpha := range file.Alphas {
		for key, _ := range alpha.Vector {
			if key > maxIndex {
				maxIndex = key
			}
		}
	}
	output.WriteString(fmt.Sprintf("%d #\n", maxIndex))
	output.WriteString(fmt.Sprintf("%d #\n", file.TrainingSize))
	output.WriteString(fmt.Sprintf("%d #\n", len(file.Alphas)+1))
	output.WriteString(fmt.Sprintf("%f #\n", file.Bias))
	for _, alpha := range file.Alphas {
		output.WriteString(alpha.ToString() + "#\n")
	}
	return output.String()
}

func (res *ClassificationResult) ToString() string {
	return fmt.Sprintf("%f\n", res)
}

func (file *ClassificationFile) ToString() string {
	output := bytes.NewBufferString("")
	for _, val := range file.Results {
		output.WriteString(val.ToString())
	}
	return output.String()
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
