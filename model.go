package gophy

import "gonum.org/v1/gonum/mat"

//interfact for model (DNAModel, MULTModel)

// StateModel ...
type StateModel interface {
	GetPCalc(float64) *mat.Dense
	GetPMap(float64) *mat.Dense
	SetMap()
	GetNumStates() int
	GetBF() []float64
	EmptyPDict()
	GetStochMapMatrices(float64, int, int) (*mat.Dense, *mat.Dense)
	SetupQJC1Rate(float64)
}
