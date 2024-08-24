package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"file-upload-app/internal/fileupload"
)

func main() {
	r := gin.Default()

	r.POST("/upload", fileupload.UploadFile)
	r.PUT("/update/:filename", fileupload.UpdateFile)
	r.DELETE("/delete/:filename", fileupload.DeleteFile)
	r.GET("/download/:filename", fileupload.DownloadFile)

	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Server run failed: %v", err)
	}
}
