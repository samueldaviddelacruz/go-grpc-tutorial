package service

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

type ImageStore interface {
	Save(LaptopId string, imageType string, imageData bytes.Buffer) (string, error)
}

type DiskImageStore struct {
	mutex       sync.Mutex
	imageFolder string
	images      map[string]*ImageInfo
}

type ImageInfo struct {
	LaptopId string
	Type     string
	Path     string
}

func NewDiskImageStore(imageFolder string) *DiskImageStore {
	return &DiskImageStore{imageFolder: imageFolder, images: make(map[string]*ImageInfo)}
}

func (store *DiskImageStore) Save(LaptopId string, imageType string, imageData bytes.Buffer) (string, error) {
	imageId, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot generate image id: %w", err)
	}

	imagePath := fmt.Sprintf("%s/%s%s", store.imageFolder, imageId, imageType)

	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("cannot create image file: %w", err)
	}
	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("cannot write image to fine: %w", err)
	}
	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.images[imageId.String()] = &ImageInfo{
		LaptopId: LaptopId,
		Type:     imageType,
		Path:     imagePath,
	}

	return imageId.String(), nil
}
