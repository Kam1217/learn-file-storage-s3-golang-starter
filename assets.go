package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
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

func getAssetPath(videoID uuid.UUID, mediaType string) string {
	fileExtention := mediaTypeToExtention(mediaType)
	return fmt.Sprintf("%s.%s", videoID.String(), fileExtention)
}

func (cfg apiConfig) getAssetDiskPath(assetPath string) string {
	return filepath.Join(cfg.assetsRoot, assetPath)
}

func (cfg apiConfig) getAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", cfg.port, assetPath)
}
