package service

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/idirall22/grpc/pb"
)

// ImageStore struct
type ImageStore interface {
	// Save image
	Save(laptopID, imageType string, imageData bytes.Buffer) (string, error)
}

// DiskImageStore struct
type DiskImageStore struct {
	mutex       sync.Mutex
	imageFolder string
	images      map[string]*pb.ImageInfo
}

// NewImageStore create a new image store
func NewImageStore(imageFolder string) *DiskImageStore {
	return &DiskImageStore{
		imageFolder: imageFolder,
		images:      make(map[string]*pb.ImageInfo),
	}
}

// Save save image
func (d *DiskImageStore) Save(laptopID, imageType string, imageData bytes.Buffer) (string, error) {
	imageUUID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("Could not generate image uuid %v", err)
	}
	imageID := imageUUID.String()
	imagePath := fmt.Sprintf("%s/%s.%s", d.imageFolder, imageID, imageType)

	file, err := os.Create(imagePath)

	if err != nil {
		return "", fmt.Errorf("Could not create image file %v", err)
	}
	defer file.Close()

	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("Could not Write image file %v", err)
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.images[imageID] = &pb.ImageInfo{
		ImageType: imageType,
		LaptopId:  laptopID,
		Path:      imagePath,
	}

	return imageID, nil
}
