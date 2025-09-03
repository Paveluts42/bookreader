package shared

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
)

func ExtractTextPagesFromPDF(pdfPath string, pageCount int) ([]string, error) {
    var pages []string
    for i := 1; i <= pageCount; i++ {
        txtFile := fmt.Sprintf("%s.page%d.txt", pdfPath, i)
        cmd := exec.Command("pdftotext", "-f", strconv.Itoa(i), "-l", strconv.Itoa(i), pdfPath, txtFile)
        out, err := cmd.CombinedOutput()
        if err != nil {
            log.Printf("pdftotext error: %v, output: %s", err, string(out))
            return nil, err
        }
        data, err := ioutil.ReadFile(txtFile)
        if err != nil {
            log.Printf("Read txt file error: %v", err)
            return nil, err
        }
        pages = append(pages, string(data))
    }
    return pages, nil
}