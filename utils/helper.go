package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func ProcessImageUpload(c *gin.Context, oldPath string, folderName string) (string, error) {
	file, err := c.FormFile("image")
	if err != nil {
		return oldPath, err
	}

	newPath := fmt.Sprintf("uploads/%s/%d_%s", folderName, time.Now().Unix(), file.Filename)
	if err := c.SaveUploadedFile(file, newPath); err != nil {
		return "", err
	}

	if oldPath != "" {
		_ = os.Remove(oldPath)
	}

	return newPath, nil
}