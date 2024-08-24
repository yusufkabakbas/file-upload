package fileupload

import (
	"log"
	"os"
	"time"
)

type FileProcessingService struct {
	processingChannel chan string
}

func NewFileProcessingService(channel chan string) *FileProcessingService {
	return &FileProcessingService{processingChannel: channel}
}

func (s *FileProcessingService) StartProcessing() {
	for filePath := range s.processingChannel {
		go func(path string) {
			log.Printf("Processing file: %s", path)
			time.Sleep(5 * time.Second)
			log.Printf("File processed: %s", path)
			os.Remove(path)
		}(filePath)
	}
}
