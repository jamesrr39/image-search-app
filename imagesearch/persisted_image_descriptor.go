package imagesearch

type PersistedImageDescriptor interface {
	LastKnownLocation() string
	Protocol() string
}
