package mps


import (
	"os"
	//"compress/gzip"
	"io"
)



func ParseMPS(fname string) (*LPData, error) {

	rd, err := openMPSFile(fname)
	if err != nil { return nil, err }
	return mpsParser(rd)
}

func openMPSFile(fname string) (io.Reader, error) {
	fp, err := os.Open(fname)
	if err != nil { return nil, err }
	//fp1, err1 := gzip.NewReader(fp)
	//if err1 != nil { return fp, nil }
	return fp, nil
}
