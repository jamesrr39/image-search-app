package imagesearch

type PersistedImageDescriptor struct {
	*ImageDescriptor
	LastDiskLocation string
}

func NewPersistedImageDescriptor(descriptor *ImageDescriptor, lastDiskLocation string) *PersistedImageDescriptor {
	return &PersistedImageDescriptor{descriptor, lastDiskLocation}
}
