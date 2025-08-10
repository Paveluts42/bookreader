package shared

import "github.com/pdfcpu/pdfcpu/pkg/api"


func GetPDFPageCount(path string) (int, error) {
    ctx, err := api.ReadContextFile(path)
    if err != nil {
        return 0, err
    }
    return ctx.PageCount, nil
}