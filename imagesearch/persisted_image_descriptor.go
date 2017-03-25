package imagesearch

type PersistedImageDescriptor struct {
	*ImageDescriptor
	Location
}

func NewPersistedImageDescriptor(descriptor *ImageDescriptor, location Location) *PersistedImageDescriptor {
	return &PersistedImageDescriptor{descriptor, location}
}
