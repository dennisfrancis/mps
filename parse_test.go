package mps

import (
	"fmt"
	"testing"
	"strings"
)


func TestSimple(t *testing.T) {
	for idx, input := range testInput {
		rd := strings.NewReader(input)
		lp, err := mpsParser(rd)
		if err != nil {
			t.Fatalf("Error returned by mpsParser, err = %s\n", err.Error())
		}
		out := fmt.Sprintf("A = %s\nB = %s\nC = %s\nCConst = %.4f\nLower Bound on x1 = %.4f\n",
			lp.A.String(),
			lp.B.String(),
			lp.C.String(),
			lp.CConst,
			lp.OrigVarsMap["x1"](0.0, 0.0),
		)
		if out != testOutput[idx] {
			
			t.Fatalf("Test failed for input #%d\nActual output:\n%s\nExpected output:\n%s\n",
				idx,
				out,
				testOutput[idx],
			)
		}
	}
}

var testInput = []string{
// INPUT 0
`NAME          example2.mps
ROWS
 N  obj     
 L  c1      
 L  c2      
COLUMNS
    x1        obj                 -1   c1                  -1
    x1        c2                   1
    x2        obj                 -2   c1                   1
    x2        c2                  -3
    x3        obj                 -3   c1                   1
    x3        c2                   1
RHS
    rhs       c1                  20   c2                  30
BOUNDS
 UP BOUND     x1                  40
ENDATA
`,
// INPUT 1
`
NAME          example2.mps
ROWS
 N  obj     
 G  c1      
 L  c2      
COLUMNS
    x1        obj                 -1   c1                  -1
    x1        c2                   1
    x2        obj                 -2   c1                   1
    x2        c2                  -3
    x3        obj                 -3   c1                   1
    x3        c2                   1
RHS
    rhs       c1                  20   c2                  30
BOUNDS
 UP BOUND     x1                  40
ENDATA
`,
// INPUT 2
`
NAME          example2.mps
ROWS
 N  obj     
 L  c1      
 E  c2      
COLUMNS
    x1        obj                 -1   c1                  -1
    x1        c2                   1
    x2        obj                 -2   c1                   1
    x2        c2                  -3
    x3        obj                 -3   c1                   1
    x3        c2                   1
RHS
    rhs       c1                  20   c2                  30
BOUNDS
 UP BOUND     x1                  40
ENDATA
`,
// INPUT 3
`
NAME          example2.mps
ROWS
 N  obj     
 L  c1      
 L  c2      
COLUMNS
    x1        obj                 -1   c1                  -1
    x1        c2                   1
    x2        obj                 -2   c1                   1
    x2        c2                  -3
    x3        obj                 -3   c1                   1
    x3        c2                   1
RHS
    rhs       c1                  20   c2                  30
BOUNDS
 UP BOUND     x1                  40
 LO BOUND     x1                  20
ENDATA
`,
// INPUT 4
`
NAME          example2.mps
ROWS
 N  obj     
 L  c1      
 L  c2      
COLUMNS
    x1        obj                 -1   c1                  -1
    x1        c2                   1
    x2        obj                 -2   c1                   1
    x2        c2                  -3
    x3        obj                 -3   c1                   1
    x3        c2                   1
RHS
    rhs       c1                  20   c2                  30
BOUNDS
 LO BOUND     x1                  20
 UP BOUND     x1                  40
ENDATA
`,
// INPUT 5
`
NAME          example2.mps
ROWS
 N  obj     
 L  c1      
 L  c2      
COLUMNS
    x1        obj                 -1   c1                  -1
    x1        c2                   1
    x2        obj                 -2   c1                   1
    x2        c2                  -3
    x3        obj                 -3   c1                   1
    x3        c2                   1
RHS
    rhs       c1                  20   c2                  30
BOUNDS
 UP BOUND     x1                  40
 FR BOUND     x3
ENDATA
`,
// INPUT 6
`
NAME          example2.mps
ROWS
 N  obj     
 L  c1      
 L  c2      
COLUMNS
    x1        obj                 -1   c1                  -1
    x1        c2                   1
    x2        obj                 -2   c1                   1
    x2        c2                  -3
    x3        obj                 -3   c1                   1
    x3        c2                   1
RHS
    rhs       c1                  20   c2                  30
RANGES
    rhs       c1                  30   c2                 -20
BOUNDS
 UP BOUND     x1                  40
ENDATA
`,
// INPUT 7
`
NAME          example2.mps
ROWS
 N  obj     
 G  c1      
 L  c2      
COLUMNS
    x1        obj                 -1   c1                  -1
    x1        c2                   1
    x2        obj                 -2   c1                   1
    x2        c2                  -3
    x3        obj                 -3   c1                   1
    x3        c2                   1
RHS
    rhs       c1                  20   c2                  30
RANGES
    rhs       c1                  50
BOUNDS
 UP BOUND     x1                  40
ENDATA
`,
// INPUT 8
`
NAME          example2.mps
ROWS
 N  obj     
 L  c1      
 E  c2      
COLUMNS
    x1        obj                 -1   c1                  -1
    x1        c2                   1
    x2        obj                 -2   c1                   1
    x2        c2                  -3
    x3        obj                 -3   c1                   1
    x3        c2                   1
RHS
    rhs       c1                  20   c2                  30
RANGES
    rhs       c2                  50
BOUNDS
 UP BOUND     x1                  40
ENDATA
`,
// INPUT 9
`
NAME          example2.mps
ROWS
 N  obj     
 L  c1      
 E  c2      
COLUMNS
    x1        obj                 -1   c1                  -1
    x1        c2                   1
    x2        obj                 -2   c1                   1
    x2        c2                  -3
    x3        obj                 -3   c1                   1
    x3        c2                   1
RHS
    rhs       c1                  20   c2                  30
RANGES
    rhs       c2                 -50
BOUNDS
 UP BOUND     x1                  40
ENDATA
`,
}

var testOutput = []string{
// OUTPUT 0
`A = {-1,  1,  1,
  1, -3,  1,
  1,  0,  0}
B = {20,
 30,
 40}
C = {-1,
 -2,
 -3}
CConst = 0.0000
Lower Bound on x1 = 0.0000
`,
// OUTPUT 1
`A = { 1, -1, -1,
  1, -3,  1,
  1,  0,  0}
B = {-20,
  30,
  40}
C = {-1,
 -2,
 -3}
CConst = 0.0000
Lower Bound on x1 = 0.0000
`,
// OUTPUT 2
`A = {-1,  1,  1,
  1, -3,  1,
 -1,  3, -1,
  1,  0,  0}
B = { 20,
  30,
 -30,
  40}
C = {-1,
 -2,
 -3}
CConst = 0.0000
Lower Bound on x1 = 0.0000
`,
// OUTPUT 3
`A = {-1,  1,  1,
  1, -3,  1,
  1,  0,  0}
B = {40,
 10,
 20}
C = {-1,
 -2,
 -3}
CConst = -20.0000
Lower Bound on x1 = 20.0000
`,
// OUTPUT 4
`A = {-1,  1,  1,
  1, -3,  1,
  1,  0,  0}
B = {40,
 10,
 20}
C = {-1,
 -2,
 -3}
CConst = -20.0000
Lower Bound on x1 = 20.0000
`,
// OUTPUT 5
`A = {-1,  1,  1, -1,
  1, -3,  1, -1,
  1,  0,  0,  0}
B = {20,
 30,
 40}
C = {-1,
 -2,
 -3,
  3}
CConst = 0.0000
Lower Bound on x1 = 0.0000
`,
// OUTPUT 6
`A = {-1,  1,  1,
  1, -1, -1,
  1, -3,  1,
 -1,  3, -1,
  1,  0,  0}
B = { 20,
  10,
  30,
 -10,
  40}
C = {-1,
 -2,
 -3}
CConst = 0.0000
Lower Bound on x1 = 0.0000
`,
// OUTPUT 7
`A = { 1, -1, -1,
 -1,  1,  1,
  1, -3,  1,
  1,  0,  0}
B = {-20,
  70,
  30,
  40}
C = {-1,
 -2,
 -3}
CConst = 0.0000
Lower Bound on x1 = 0.0000
`,
// OUTPUT 8
`A = {-1,  1,  1,
 -1,  3, -1,
  1, -3,  1,
  1,  0,  0}
B = { 20,
 -30,
  80,
  40}
C = {-1,
 -2,
 -3}
CConst = 0.0000
Lower Bound on x1 = 0.0000
`,
// OUTPUT 9
`A = {-1,  1,  1,
  1, -3,  1,
 -1,  3, -1,
  1,  0,  0}
B = {20,
 30,
 20,
 40}
C = {-1,
 -2,
 -3}
CConst = 0.0000
Lower Bound on x1 = 0.0000
`,
}
