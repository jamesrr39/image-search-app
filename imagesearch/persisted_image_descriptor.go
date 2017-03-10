package imagesearch

type PersistedImageDescriptor struct {
	*ImageDescriptor
	LastDiskLocation string
}
