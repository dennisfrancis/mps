package mps


import (
	"os"
	"compress/gzip"
	"io"
)



func ParseMPS(fname string) (*LPData, error) {

	rd, err := openMPSFile(fname, false)
	if err != nil { return nil, err }
	return mpsParser(rd)
}

func ParseCompressedMPS(fname string) (*LPData, error) {

	rd, err := openMPSFile(fname, true)
	if err != nil { return nil, err }
	return mpsParser(rd)
}


func openMPSFile(fname string, isCompressed bool) (io.Reader, error) {
	fp, err := os.Open(fname)
	if err != nil { return nil, err }
	if isCompressed {
		fp1, err1 := gzip.NewReader(fp)
		if err1 != nil { return nil, err1 }
		return fp1, nil
	}
	return fp, nil
}
