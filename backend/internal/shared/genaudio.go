package shared

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func GenerateAudioFromText(text string, audioPath string, speed, pitch, volume float64) error {
	ttsURL := os.Getenv("TTS_URL")
	if ttsURL == "" {
		ttsURL = "http://tts:8080/say"
	}
	params := url.Values{}
	params.Set("text", text)
	params.Set("voice", "anna")
	fullURL := ttsURL + "?" + params.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		return fmt.Errorf("tts request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("tts error: %s", string(b))
	}
	outFile, err := os.Create(audioPath)
	if err != nil {
		return fmt.Errorf("failed to create audio file: %v", err)
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, resp.Body)
	return err
}
