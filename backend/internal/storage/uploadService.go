package storage

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Paveluts42/bookreader/backend/internal/shared"
	"github.com/google/uuid"
)

type UploadService struct{}

func NewUploadService() *UploadService {
	return &UploadService{}
}

func ExtractTextFromPDF(pdfPath string) (string, error) {
	txtFile := pdfPath + ".txt"
	cmd := exec.Command("pdftotext", pdfPath, txtFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("pdftotext error: %v, output: %s", err, string(out))
		return "", err
	}
	data, err := ioutil.ReadFile(txtFile)
	if err != nil {
		log.Printf("Read txt file error: %v", err)
		return "", err
	}
	os.Remove(txtFile)
	return string(data), nil
}
func (s *UploadService) SavePDF(bookID, title, author, userID string, chunk []byte, pageCount int, filePath, coverPath string) (*Book, error) {
    log.Printf("Saving PDF for book: %s", bookID)

    uploadDir := "/uploads/" + bookID
    audioDir := uploadDir + "/audio"
    pagesDir := audioDir + "/pages"

    wd, err := os.Getwd()
    if err != nil {
        log.Printf("Error getting working dir: %v", err)
    } else {
        log.Printf("Current working dir: %s", wd)
    }
    log.Printf("Trying to create: %s", pagesDir)

    if err := os.MkdirAll(pagesDir, 0777); err != nil {
        log.Printf("Error creating upload dirs: %v", err)
        return nil, err
    } else {
        log.Printf("Created upload dirs: %s", pagesDir)
    }

    pdfPath := uploadDir + "/" + bookID + ".pdf"
    coverPath = uploadDir + "/" + bookID + ".png"
    textPath := uploadDir + "/" + bookID + ".txt"
    audioPath := audioDir + "/full.mp3"

    log.Printf("Writing PDF: %s", pdfPath)
    if err := os.WriteFile(pdfPath, chunk, 0666); err != nil {
        log.Printf("Error writing PDF: %v", err)
        return nil, err
    }

    log.Printf("Generating cover: %s", coverPath)
    if err := shared.GenerateCover(pdfPath, coverPath); err != nil {
        log.Printf("Failed to generate cover: %v", err)
    }

    log.Printf("Extracting text from PDF: %s", pdfPath)
    text, err := ExtractTextFromPDF(pdfPath)
    if err != nil {
        log.Printf("Failed to extract text: %v", err)
        return nil, err
    }
    log.Printf("Writing text file: %s", textPath)
    if err := os.WriteFile(textPath, []byte(text), 0666); err != nil {
        log.Printf("Error writing text file: %v", err)
    }

    go func() {
        pages := strings.Split(text, "\f")
        var wg sync.WaitGroup
        audioFiles := []string{}
        for i, pageText := range pages {
            pageText = strings.TrimSpace(pageText)
            if len(pageText) < 5 {
                continue
            }
            wg.Add(1)
            pageAudioPath := fmt.Sprintf("%s/page_%d.mp3", pagesDir, i+1)
            audioFiles = append(audioFiles, pageAudioPath)
            go func(i int, pageText string, pageAudioPath string) {
                defer wg.Done()
                log.Printf("Generating audio for page %d: %s", i+1, pageAudioPath)
                // Генерируем WAV, затем конвертируем в MP3
                tmpWav := fmt.Sprintf("%s/page_%d_tmp.wav", pagesDir, i+1)
                if err := shared.GenerateAudioFromText(pageText, tmpWav, 1.0, 1.0, 1.0); err != nil {
                    log.Printf("Audio gen error for page %d: %v", i+1, err)
                    return
                }
                cmd := exec.Command("ffmpeg", "-y", "-i", tmpWav, "-acodec", "libmp3lame", "-ab", "192k", pageAudioPath)
                out, err := cmd.CombinedOutput()
                if err != nil {
                    log.Printf("ffmpeg mp3 error for page %d: %v, output: %s", i+1, err, string(out))
                }
                os.Remove(tmpWav)
            }(i, pageText, pageAudioPath)
        }
        wg.Wait()
        log.Printf("All page audio files generated for book: %s", bookID)

        listPath := audioDir + "/pages_list.txt"
        f, err := os.Create(listPath)
        if err != nil {
            log.Printf("Error creating concat list: %v", err)
            return
        }
        for _, fpath := range audioFiles {
            absPath, _ := filepath.Abs(fpath)
            f.WriteString(fmt.Sprintf("file '%s'\n", absPath))
        }
        f.Close()

        fullAudioPath := audioDir + "/full.mp3"
        cmd := exec.Command(
            "ffmpeg",
            "-y",
            "-f", "concat",
            "-safe", "0",
            "-i", listPath,
            "-acodec", "libmp3lame",
            "-ab", "192k",
            fullAudioPath,
        )
        out, err := cmd.CombinedOutput()
        if err != nil {
            log.Printf("ffmpeg error: %v, output: %s", err, string(out))
        } else {
            log.Printf("Full audio generated: %s", fullAudioPath)
            DB.Model(&Book{}).
                Where("id = ?", bookID).
                Update("audio_path", fullAudioPath)
        }
    }()

    book := Book{
        ID:        uuid.MustParse(bookID),
        Title:     title,
        Author:    author,
        Page:      int32(0),
        PageAll:   int32(pageCount),
        FilePath:  pdfPath,
        CoverPath: coverPath,
        AudioPath: audioPath,
        UserID:    uuid.MustParse(userID),
        CreatedAt: time.Now(),
    }
    if err := DB.Create(&book).Error; err != nil {
        log.Printf("Error saving book to DB: %v", err)
        return nil, err
    }
    log.Printf("Book saved to DB: %s", bookID)
    return &book, nil
}