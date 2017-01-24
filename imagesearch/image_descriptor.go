package imagesearch

type ImageDescriptor struct {
	Bins uint
}

func NewImageDescriptor(bins uint) *ImageDescriptor {
	return &ImageDescriptor{bins}
}
