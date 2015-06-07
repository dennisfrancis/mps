package mps


import (
	"fmt"
	"math"
	"bufio"
	"strings"
	"strconv"
	"io"
	"sort"
	"github.com/skelterjohn/go.matrix"
)

func mpsParser(rd io.Reader) (*LPData, error) {

	scanner := bufio.NewScanner(rd)
	section := UNKNOWN
	raw := LPRaw{
		Name : "",
		Constraints: map[string]*Constraint{},
		Variables : map[string][]string{},
		VariableTransforms : map[string]VarTrans{},
		Objective : map[string]float64{},
	}
	objectiveName := ""
	//varsBoundSeen := map[string]bool{}
	var err error = nil
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), "\n\t\r ")
		line  = strings.Split(line, "*")[0]
		line  = strings.Split(line, "$")[0]
		//fmt.Println(line)
		flds := strings.Fields(line)
		if len(flds) == 0 { continue }
		if flds[0] == "NAME" {
			raw.Name = flds[1]
			continue
		} else if flds[0] == "ROWS" {
			section = ROWS
			continue
		} else if flds[0] == "COLUMNS" {
			section = COLUMNS
			continue
		} else if flds[0] == "RHS" {
			section = RHS
			continue
		} else if flds[0] == "RANGES" {
			section = RANGES
			continue
		} else if flds[0] == "BOUNDS" {
			section = BOUNDS
			continue
		} else if flds[0] == "ENDATA" {
			continue
		}

		switch(section) {
		case ROWS:
			if flds[0] == "N" {
				if objectiveName == "" {
					objectiveName = flds[1]
				}
				continue
			}
			if ineq, ok := InequalityCode2InequalityType[flds[0]]; ok {
				raw.Constraints[flds[1]] = &Constraint{
					Name : flds[1],
					Type : ineq,
					Rhs  : 0.0,
					Coefficients : map[string]float64{},
				}
			}
		case COLUMNS:
			if len(flds) != 3 && len(flds) != 5 {
				return nil,
				fmt.Errorf("In COLUMNS section there is unsupported number of fields")
			}
			// flds[0] is var name, flds[1] is obj/constraint name flds[2] is the coeff
			// if len(flds) is 5, then flds[3] is obj/constraint name flds[4] is the coeff
			raw.Variables[flds[0]] = []string{flds[0]}
			raw.VariableTransforms[flds[0]] = defaultTransform
			if flds[1] == objectiveName {
				raw.Objective[flds[0]], err = strconv.ParseFloat(flds[2], 64)
				if err != nil {
					return nil, err
				}
			} else if con, ok := raw.Constraints[flds[1]]; ok {
				con.Coefficients[flds[0]], err = strconv.ParseFloat(flds[2], 64)
				if err != nil {
					return nil, err
				}
			}
			if len(flds) == 5 {
				if flds[3] == objectiveName {
					raw.Objective[flds[0]], err = strconv.ParseFloat(flds[4], 64)
					if err != nil {
						return nil, err
					}
				} else if con, ok := raw.Constraints[flds[3]]; ok {
					con.Coefficients[flds[0]], err = strconv.ParseFloat(flds[4], 64)
					if err != nil {
						return nil, err
					}
				}	
			}
		case RHS:
			if len(flds) != 3 && len(flds) != 5 {
				return nil,
				fmt.Errorf("In RHS section there is unsupported number of fields")
			}
			if con, ok := raw.Constraints[flds[1]]; ok {
				con.Rhs, err = strconv.ParseFloat(flds[2], 64)
				//fmt.Println("con.Rhs =", con.Rhs)
				if err != nil {
					return nil, err
				}
			}
			if len(flds) == 5 {
				if con, ok := raw.Constraints[flds[3]]; ok {
					con.Rhs, err = strconv.ParseFloat(flds[4], 64)
					if err != nil {
						return nil, err
					}
				}
					
			}
		case RANGES:
			if len(flds) != 3 && len(flds) != 5 {
				return nil,
				fmt.Errorf("In RANGES section there is unsupported number of fields")
			}
			if con, ok := raw.Constraints[flds[1]]; ok {
				rang, err := strconv.ParseFloat(flds[2], 64)
				if err != nil {
					return nil, err
				}
				switch con.Type {
				case GE:
					connew := &Constraint{
						Name         : con.Name + "_addl",
						Coefficients : con.Coefficients,
						Type         : LE,
						Rhs          : con.Rhs + math.Abs(rang),
					}
					raw.Constraints[connew.Name] = connew
				case LE:
					connew := &Constraint{
						Name         : con.Name + "_addl",
						Coefficients : con.Coefficients,
						Type         : GE,
						Rhs          : con.Rhs - math.Abs(rang),
					}
					raw.Constraints[connew.Name] = connew					
				case EQUALITY:
					if rang > 0.0 {
						connew := &Constraint{
							Name         : con.Name + "_addl",
							Coefficients : con.Coefficients,
							Type         : LE,
							Rhs          : con.Rhs + rang,
						}
						raw.Constraints[connew.Name] = connew
						con.Type = GE
					} else {
						connew := &Constraint{
							Name         : con.Name + "_addl",
							Coefficients : con.Coefficients,
							Type         : GE,
							Rhs          : con.Rhs + rang,
						}
						raw.Constraints[connew.Name] = connew
						con.Type = LE
					}
				}
			}
			if len(flds) == 5 {
				if con, ok := raw.Constraints[flds[3]]; ok {
					rang, err := strconv.ParseFloat(flds[4], 64)
					if err != nil {
						return nil, err
					}
					switch con.Type {
					case GE:
						connew := &Constraint{
							Name         : con.Name + "_addl",
							Coefficients : con.Coefficients,
							Type         : LE,
							Rhs          : con.Rhs + math.Abs(rang),
						}
						raw.Constraints[connew.Name] = connew
					case LE:
						connew := &Constraint{
							Name         : con.Name + "_addl",
							Coefficients : con.Coefficients,
							Type         : GE,
							Rhs          : con.Rhs - math.Abs(rang),
						}
						raw.Constraints[connew.Name] = connew					
					case EQUALITY:
						if rang > 0.0 {
							connew := &Constraint{
								Name         : con.Name + "_addl",
								Coefficients : con.Coefficients,
								Type         : LE,
								Rhs          : con.Rhs + rang,
							}
							raw.Constraints[connew.Name] = connew
							con.Type = GE
						} else {
							connew := &Constraint{
								Name         : con.Name + "_addl",
								Coefficients : con.Coefficients,
								Type         : GE,
								Rhs          : con.Rhs + rang,
							}
							raw.Constraints[connew.Name] = connew
							con.Type = LE
						}
					}
				}
			}
		case BOUNDS:
			if len(flds) != 4 && len(flds) != 3 {
				return nil,
				fmt.Errorf("In BOUNDS section there is unsupported number of fields")
			}
			if _, ok := raw.Variables[flds[2]]; !ok {
				continue
			}
			var rhs float64 = 0.0
			if len(flds) == 4 {
				rhs, err = strconv.ParseFloat(flds[3], 64)
				if err != nil {
					return nil, err
				}
			}
			switch flds[0] {
			case "LO":
				if raw.Variables[flds[2]] == nil {
					continue
				}
				// xj >= lb
				// introduce new var xj_prime >= 0 st, xj_prime = xj - lb
				// xj = x_prime + lb
				oldvar := flds[2]
				newvar := oldvar + "_prime"
				raw.Variables[oldvar] = []string{newvar}
				raw.VariableTransforms[oldvar] = func(x_prime, x_prime2 float64) (float64) {
					return x_prime + rhs
				}
				for _, con := range raw.Constraints {
					if coeff, found := con.Coefficients[oldvar]; found {
						delete(con.Coefficients, oldvar)
						con.Coefficients[newvar] = coeff
						con.Rhs -= (coeff*rhs)
					}
				}
				if coeff, found := raw.Objective[oldvar]; found {
					delete(raw.Objective, oldvar)
					raw.Objective[newvar] = coeff
					raw.Objective[OBJECTIVE_CONST_TERM] += (coeff*rhs)
				}
			case "UP":
				// xj <= ub
				// add a new constraint
				newvars, found := raw.Variables[flds[2]]
				if !found { continue }
				newvar := newvars[0]
				constterm := raw.VariableTransforms[flds[2]](0.0, 0.0)
				// xj_prime = xj - constterm
				// xj_prime + constterm <= ub
				// xj_prime <= ub - constterm. this constraint needs to be added
				connew := &Constraint{
					Name             : newvar + "_upperbound",
					Coefficients     : map[string]float64{newvar : 1.0},
					Type             : LE,
					Rhs              : rhs - constterm,
				}
				raw.Constraints[connew.Name] = connew
			case "FX":
				fmt.Println("Warning : in BOUND, FX is not implemented, skipping")
				continue
			case "FR":
				oldvar := flds[2]
				newvars := []string{oldvar + "_prime1", oldvar + "_prime2"}
				raw.Variables[oldvar] = newvars
				raw.VariableTransforms[oldvar] = func(x_prime1, x_prime2 float64) float64 {
					return x_prime1 - x_prime2
				}
				for _, con := range raw.Constraints {
					if coeff, found := con.Coefficients[oldvar]; found {
						delete(con.Coefficients, oldvar)
						con.Coefficients[newvars[0]] = coeff
						con.Coefficients[newvars[1]] = -coeff
					}
				}
				if coeff, found := raw.Objective[oldvar]; found {
					delete(raw.Objective, oldvar)
					raw.Objective[newvars[0]] = coeff
					raw.Objective[newvars[1]] = -coeff
				}
			}
		}
		
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	//fmt.Println("DEBUG : Constraints :")
	//for conname, con := range raw.Constraints {
	//	fmt.Printf("%s = %+v\n", conname, con)
	//}
	// Now convert to standard form
	// All constrains should convert to <= type
	newcons := map[string]*Constraint{}
	for conname, con := range raw.Constraints {
		if con.Type == GE {
			coeffnew := map[string]float64{}
			for varname, coeff := range con.Coefficients {
				coeffnew[varname] = -coeff
			}
			connew := &Constraint{
				Name          : con.Name,
				Coefficients  : coeffnew,
				Type          : LE,
				Rhs           : -con.Rhs,
			}
			newcons[conname] = connew
		} else if con.Type == EQUALITY {
			coeffnew1 := map[string]float64{}
			coeffnew2 := map[string]float64{}
			for varname, coeff := range con.Coefficients {
				coeffnew1[varname] = coeff
				coeffnew2[varname] = -coeff
			}
			connew1 := &Constraint{
				Name          : con.Name,
				Coefficients  : coeffnew1,
				Type          : LE,
				Rhs           : con.Rhs,
			}
			connew2 := &Constraint{
				Name          : con.Name + "_prime",
				Coefficients  : coeffnew2,
				Type          : LE,
				Rhs           : -con.Rhs,
			}
			newcons[connew1.Name] = connew1
			newcons[connew2.Name] = connew2
		} else {
			newcons[con.Name] = con
		}
	}
	raw.Constraints = newcons
	/*fmt.Println("DEBUG : Constraints after standardising2 :")
	for conname, con := range raw.Constraints {
		fmt.Printf("%s = %+v\n", conname, con)
	}*/

	// map varnames used in the constraints/objective to integers (index)
	var varsused = map[string]int{}
	var varslice = []string{}
	for _, vars := range raw.Variables {
		varslice = append(varslice, vars...)
	}
	sort.Strings(varslice)
	for idx, vr := range varslice {
		varsused[vr] = idx
	}

	// Generate LPData struct
	n := len(varsused)
	m := len(raw.Constraints)

	
	A := matrix.ZerosSparse(m, n)
	B := matrix.Zeros(m, 1)
	C := matrix.Zeros(n, 1)
	constraintNames := make([]string, len(raw.Constraints))
	idx := 0
	for cname, _ := range raw.Constraints {
		constraintNames[idx] = cname
		idx += 1
	}
	sort.Strings(constraintNames)
	for rowIdx, cname := range constraintNames {
		cons := raw.Constraints[cname]
		for varname, coeff := range cons.Coefficients {
			A.Set(rowIdx, varsused[varname], coeff)
		}
		B.Set(rowIdx, 0, cons.Rhs)
	}
	for idx, varname := range varslice {
		coeff := raw.Objective[varname]
		C.Set(idx, 0, coeff)
	}
	orig2used := map[string][]int{}
	for orig, used := range raw.Variables {
		usedint := []int{varsused[used[0]]}
		if len(used) == 2 { usedint = append(usedint, varsused[used[1]]) }
		orig2used[orig] = usedint
	}
	return &LPData{
		Name             : raw.Name,
		A                : A,
		B                : B,
		C                : C,
		CConst           : raw.Objective[OBJECTIVE_CONST_TERM],
		Orig2UsedVars    : orig2used,
		OrigVarsMap      : raw.VariableTransforms,
	}, nil
}
