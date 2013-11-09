package svmlight

import (
	"code.google.com/p/gomat/vec"
)

type SVMFile struct {
	Instances []SVMInstance
	FileName  string
}
type SVMInstance struct {
	Features map[int]SVMFeature
	Label    int
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
	Kernel          KernelInfo
	TrainingSize    int
	MaxFeatureIndex int
	Bias            float64
	Alphas          []Alpha
}

type ClassificationResult float64

type ClassificationFile struct {
	Results []ClassificationResult
}

type WeightVector struct {
	B   float64
	Vec vec.DenseVector
}
