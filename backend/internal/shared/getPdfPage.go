package shared

import (
	"bytes"

    "rsc.io/pdf"
)

func GetPDFPageCountFromBytes(pdfData []byte) (int, error) {
    r, err := pdf.NewReader(bytes.NewReader(pdfData), int64(len(pdfData)))
    if err != nil {
        return 0, err
    }
    return r.NumPage(), nil
}