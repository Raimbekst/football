package media

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"mime/multipart"
	"os"
	"path/filepath"
)

func GetFileName(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
	ImageIsValid := false

	filetype := filepath.Ext(file.Filename)

	fileExt := []string{".png", ".jpg", ".jpeg", ".webp", ".svg", ".jfif"}

	for _, value := range fileExt {
		if filetype == value {
			ImageIsValid = true
		}
	}

	if !ImageIsValid {
		return "", fmt.Errorf("media.GetFileName: %w", errors.New("the file is not image type"))
	}

	filename := uuid.New().String() + filetype
	err := c.SaveFile(file, "media/"+filename)

	if err != nil {
		return "", fmt.Errorf("media.GetFileName: %w", err)
	}

	return filename, nil
}

func DeleteImage(image string) error {
	mkdir := "./media"

	return os.Remove(mkdir + "/" + image)
}
