package utils

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func ProcessImageUpload(c *gin.Context, formKey string, oldPath string, folderName string) (string, error) {
	file, err := c.FormFile(formKey)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return oldPath, nil
		}
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

func IfEmpty[T comparable](newVal, oldVal T, zeroValue T) T {
	if newVal == zeroValue {
		return oldVal
	}
	return newVal
}