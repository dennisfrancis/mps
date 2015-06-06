package mps

import (
	"github.com/skelterjohn/go.matrix"
)
const (
	UNKNOWN int = iota
	NAME
	ROWS
	COLUMNS
	RHS
	RANGES
	BOUNDS
	ENDDATA
)

type InequalityType int
const (
	EQUALITY  InequalityType = iota
	LE
	GE
)

type VarTrans func(float64, float64) float64

type LPRaw struct {
	Name                 string
	Variables            map[string][]string
	VariableTransforms   map[string]VarTrans
	Constraints          map[string]*Constraint // Var bounds are to be included as a constraint
	Objective            map[string]float64
}

type Constraint struct {
	Name          string
	Type          InequalityType
	Coefficients  map[string]float64
	Rhs           float64
}

var InequalityCode2InequalityType = map[string]InequalityType{
	"E" : EQUALITY,
	"G" : GE,
	"L" : LE, 
}

func defaultTransform(x_prime, x_prime2 float64) float64 {
	return x_prime
}

type LPData struct {
	Name           string
	A              matrix.MatrixRO      // Constriant coeff matrix
	B              matrix.MatrixRO      // RHS of constriants in std form
	C              matrix.MatrixRO      // Variable coeffs in objective function
	CConst         float64              // Constant term in objective function
	Orig2UsedVars  map[string][]int     // key = orig var name, value = slice of actually used var index
	OrigVarsMap    map[string]VarTrans  // key = orig var name, value = function(actual vars) returning orig var
}

const OBJECTIVE_CONST_TERM string = "const_term"
