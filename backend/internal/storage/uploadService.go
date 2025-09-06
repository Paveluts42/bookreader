package storage

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

func (s *UploadService) SavePDF(bookID, title, author, userID string, username string, chunk []byte, pageCount int, filePath, coverPath string) (*Book, error) {
	log.Printf("Saving PDF for book: %s", bookID)
	userDir := "/uploads/" + username
	uploadDir := userDir + "/" + bookID
	audioDir := uploadDir + "/audio"
	pagesDir := audioDir + "/pages"

	if err := os.MkdirAll(pagesDir, 0777); err != nil {
		log.Printf("Error creating upload dirs: %v", err)
		return nil, err
	}

	pdfPath := uploadDir + "/" + bookID + ".pdf"
	coverPath = uploadDir + "/" + bookID + ".png"
	textPath := uploadDir + "/" + bookID + ".txt"
	audioPath := audioDir + "/full.mp3"

	if err := os.WriteFile(pdfPath, chunk, 0666); err != nil {
		log.Printf("Error writing PDF: %v", err)
		return nil, err
	}

	if err := shared.GenerateCover(pdfPath, coverPath); err != nil {
		log.Printf("Failed to generate cover: %v", err)
	}

	text, err := ExtractTextFromPDF(pdfPath)
	if err != nil {
		log.Printf("Failed to extract text: %v", err)
		return nil, err
	}
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

		// --- Склейка по батчам ---
		batchSize := 100
		var batchFiles []string
		var batchPaths []string
		for i, f := range audioFiles {
			batchFiles = append(batchFiles, f)
			if len(batchFiles) == batchSize || i == len(audioFiles)-1 {
				batchOut := fmt.Sprintf("%s/batch_%d.mp3", audioDir, len(batchPaths)+1)
				shared.ConcatMp3Batch(batchFiles, batchOut)
				batchPaths = append(batchPaths, batchOut)
				batchFiles = []string{}
			}
		}
		shared.ConcatMp3Batch(batchPaths, audioPath)

		log.Printf("Full audio generated: %s", audioPath)
        
		for _, f := range batchPaths {
			os.Remove(f)
		}

		DB.Model(&Book{}).
			Where("id = ?", bookID).
			Update("audio_path", audioPath)
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

func ConcatMp3Batch(batchFiles []string, batchOut string) {
	panic("unimplemented")
}
