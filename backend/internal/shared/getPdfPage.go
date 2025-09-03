package shared

import (
    "bytes"
    "github.com/pdfcpu/pdfcpu/pkg/api"
)

func GetPDFPageCountFromBytes(pdfData []byte) (int, error) {
    ctx, err := api.ReadContext(bytes.NewReader(pdfData), nil)
    if err != nil {
        return 0, err
    }
    return ctx.PageCount, nil
}