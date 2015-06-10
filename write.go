package mps

import (
	"fmt"
	"os"
)


func WriteMPS(lp *LPData, fname string) error {
	fp, err := os.Create(fname)
	if err != nil { return err }
	defer fp.Close()

	fp.WriteString(fmt.Sprintf("NAME  %s\n", lp.Name))
	fp.WriteString("ROWS\n")
	fp.WriteString(" N obj\n")
	for ii := 0; ii < lp.A.Rows(); ii++ {
		fp.WriteString(fmt.Sprintf(" L c%d\n", ii))
	}

	fp.WriteString("COLUMNS\n")
	for jj := 0; jj < lp.A.Cols(); jj++ {
		cj := lp.C.Get(jj, 0)
		if cj != 0.0 {
			fp.WriteString(fmt.Sprintf("   x%d obj %.5f\n", jj, cj))
		}
		for ii := 0; ii < lp.A.Rows(); ii++ {
			
			if lp.A.Get(ii, jj) == 0.0 { continue }
			fp.WriteString(fmt.Sprintf("   x%d c%d %.5f\n", jj, ii, lp.A.Get(ii, jj)))
		}
	}

	fp.WriteString("RHS\n")
	for ii := 0; ii < lp.A.Rows(); ii++ {
		if lp.B.Get(ii, 0) == 0.0 { continue }
		fp.WriteString(fmt.Sprintf("   rhs  c%d  %.5f\n", ii, lp.B.Get(ii, 0)))
	}
	fp.WriteString("ENDATA\n")
	return nil
}
