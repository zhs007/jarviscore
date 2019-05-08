package jarviscore

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/zhs007/jarviscore/basedef"
)

// CMDStdOutErr - command stdout & stderr
type CMDStdOutErr struct {
	LenStdOut   int64
	OutFileName string
	ErrFileName string
	fout        *os.File
	ferr        *os.File
}

// NewCMDStdOutErr - new CMDStdOutErr
func NewCMDStdOutErr(fnout string, fnerr string) (*CMDStdOutErr, error) {
	if fnerr == "" {
		return nil, ErrInvalidCMDStdOutErrErrFile
	}

	fperr, err := os.Create(fnerr)
	if err != nil {
		return nil, err
	}

	if fnout == "" {
		return &CMDStdOutErr{
			ErrFileName: fnerr,
			ferr:        fperr,
		}, nil
	}

	fp, err := os.Create(fnout)
	if err != nil {
		return nil, err
	}

	return &CMDStdOutErr{
		OutFileName: fnout,
		fout:        fp,
		ErrFileName: fnerr,
		ferr:        fperr,
	}, nil
}

func (out *CMDStdOutErr) countStdOut(r io.Reader) error {
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			if out.fout != nil {
				out.fout.Write(buf[:n])
			}

			out.LenStdOut += int64(n)
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				out.fout.Close()

				err = nil
			}

			return err
		}
	}
}

func (out *CMDStdOutErr) countStdErr(r io.Reader) error {
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			if out.ferr != nil {
				out.ferr.Write(buf[:n])
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				out.ferr.Close()

				err = nil
			}

			return err
		}
	}
}

// GetStdErr - get stderr
func (out *CMDStdOutErr) GetStdErr() ([]byte, error) {
	return ioutil.ReadFile(out.ErrFileName)
}

// GetStdOut - get stdout
func (out *CMDStdOutErr) GetStdOut() ([]byte, error) {
	fi, _ := os.Stat(out.OutFileName)
	if fi.Size() < basedef.BigLogFileLength {
		return ioutil.ReadFile(out.OutFileName)
	}

	return []byte(fmt.Sprintf("The log file is too large, please check it through %v.", out.OutFileName)), nil
}
