package svmlight

import (
	"code.google.com/p/gomat/vec"
	"fmt"
	"math"
	"os/exec"
)

func (f *SVMFile) CountLabels(label int) int {
	acc := 0
	for _, instance := range f.Instances {
		if instance.Label == label {
			acc++
		}
	}
	return acc
}

func (instance *SVMInstance) Norm() SVMFeature {
	cumSum := 0.
	for _, val := range instance.Features {
		cumSum += float64(val)
	}
	return SVMFeature(math.Sqrt(float64(cumSum)))
}

func (f *SVMFile) MakeBinary(label int) {
	for _, instance := range f.Instances {
		if instance.Label == label {
			instance.Label = 1
		} else {
			instance.Label = -1
		}
	}
}

func (f *SVMFile) NormalizeInstances() {
	for _, instance := range f.Instances {
		norm := instance.Norm()
		for k, v := range instance.Features {
			instance.Features[k] = v / norm
		}
	}
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

func (file *ClassificationFile) Accuracy(set *SVMFile) float64 {
	totalInstances := float64(len(set.Instances))
	correctPreds := 0.
	if len(file.Results) <= len(set.Instances) {
		for i, res := range file.Results {
			if math.Signbit(float64(res)) == math.Signbit(float64(set.Instances[i].Label)) {
				correctPreds++
			}
		}
	} else {
		for i, instance := range set.Instances {
			if math.Signbit(float64(file.Results[i])) == math.Signbit(float64(instance.Label)) {
				correctPreds++
			}
		}
	}

	return correctPreds / totalInstances
}

func (f *ModelFile) ComputeWeightVector() *WeightVector {
	w := WeightVector{}
	w.Vec = vec.New(f.MaxFeatureIndex + 1)
	for _, alpha := range f.Alphas {
		for featureId, featureVal := range alpha.Vector {
			w.Vec[featureId] += alpha.Weight * float64(featureVal)
		}
	}
	return &w
}
