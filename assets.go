package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (cfg apiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.assetsRoot, 0755)
	}
	return nil
}

func mediaTypeToExtention(mediaType string) string {
	fileExtention := strings.Split(mediaType, "/")
	if len(fileExtention) != 2 {
		return "extention should only have 2 parts"
	}
	return fileExtention[1]
}

func getAssetPath(mediaType string) string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic("failed to generate random bytes")
	}
	id := base64.RawURLEncoding.EncodeToString(key)
	fileExtention := mediaTypeToExtention(mediaType)

	return fmt.Sprintf("%s.%s", id, fileExtention)
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}

func getVideoAspectRatio(filePath string) (string, error) {
	type FFProbeOutput struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}

	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	buf := bytes.Buffer{}
	cmd.Stdout = &buf

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run the command: %w", err)
	}

	var probeOutput FFProbeOutput
	if err = json.Unmarshal(buf.Bytes(), &probeOutput); err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	if len(probeOutput.Streams) == 0 {
		return "", errors.New("no videos found")
	}
	width := probeOutput.Streams[0].Width
	height := probeOutput.Streams[0].Height

	result := float64(width) / float64(height)

	if result >= 1.7 && result <= 1.8 {
		return "16:9", nil
	} else if result >= 0.5 && result <= 0.6 {
		return "9:16", nil
	}
	return "other", nil
}
