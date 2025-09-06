package shared

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func ConcatMp3Batch(files []string, outPath string) error {
    listPath := outPath + "_list.txt"
    f, err := os.Create(listPath)
    if err != nil {
        return err
    }
    for _, fpath := range files {
        absPath, _ := filepath.Abs(fpath)
        f.WriteString(fmt.Sprintf("file '%s'\n", absPath))
    }
    f.Close()
    cmd := exec.Command(
        "ffmpeg",
        "-y",
        "-f", "concat",
        "-safe", "0",
        "-i", listPath,
        "-acodec", "libmp3lame",
        "-ab", "192k",
        outPath,
    )
    out, err := cmd.CombinedOutput()
    if err != nil {
        log.Printf("ffmpeg batch error: %v, output: %s", err, string(out))
        return err
    }
    os.Remove(listPath)
    return nil
}