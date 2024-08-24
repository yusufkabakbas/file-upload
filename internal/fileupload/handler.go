package fileupload

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const defaultUploadPath = "./uploads"

func ensureDirExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to create directory: %w", err)
		}
	}
	return nil
}

func processFile(file io.Reader, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(out)

	_, err = io.Copy(out, file)
	return err
}

func UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("File error: %s", err.Error()))
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	path := c.DefaultPostForm("path", defaultUploadPath)
	err = ensureDirExists(path)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error ensuring directory exists: %s", err.Error()))
		return
	}

	filename := filepath.Join(path, header.Filename)

	go func() {
		err := processFile(file, filename)
		if err != nil {
			fmt.Printf("Error saving file %s: %v\n", filename, err)
		}
	}()

	c.String(http.StatusOK, fmt.Sprintf("File %s is being uploaded.", header.Filename))
}

func UpdateFile(c *gin.Context) {
	filename := c.Param("filename")
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("File error: %s", err.Error()))
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	path := filepath.Join(defaultUploadPath, filename)
	err = ensureDirExists(filepath.Dir(path))
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error ensuring directory exists: %s", err.Error()))
		return
	}

	go func() {
		err := processFile(file, path)
		if err != nil {
			fmt.Printf("Error updating file %s: %v\n", filename, err)
		}
	}()

	c.String(http.StatusOK, fmt.Sprintf("File %s is being updated.", filename))
}

func DeleteFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(defaultUploadPath, filename)

	go func() {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Printf("Error deleting file %s: %v\n", filename, err)
		}
	}()

	c.String(http.StatusOK, fmt.Sprintf("File %s is being deleted.", filename))
}

func DownloadFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join(defaultUploadPath, filename)

	c.File(filePath)
}
