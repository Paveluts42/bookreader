package shared

import (
    "github.com/gen2brain/go-fitz"
    "image/png"
    "os"
)

func GenerateCover(pdfPath, coverPath string) error {
    doc, err := fitz.New(pdfPath)
    if err != nil {
        return err
    }
    defer doc.Close()

    img, err := doc.Image(0) // first page
    if err != nil {
        return err
    }

    f, err := os.Create(coverPath)
    if err != nil {
        return err
    }
    defer f.Close()

    return png.Encode(f, img)
}