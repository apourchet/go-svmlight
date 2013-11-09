package svmlight

import (
	"bytes"
	"fmt"
)

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
	output.WriteString(fmt.Sprintf("%d #\n", file.MaxFeatureIndex))
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
