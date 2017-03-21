package imagesearchdal

import (
	"image-search-app/imagesearch"
)

type PersistedImageDescriptor struct {
	*imagesearch.ImageDescriptor
	LastDiskLocation string
}

func NewPersistedImageDescriptor(descriptor *imagesearch.ImageDescriptor, lastDiskLocation string) *PersistedImageDescriptor {
	return &PersistedImageDescriptor{descriptor, lastDiskLocation}
}
