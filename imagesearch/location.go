package imagesearch

// Location is a representation of the last known location of the image (on disk, via http, etc)
// It doesn't really sit so well in the core image search code, but the results are not very useful without a location the top-ranked images can be found at
type Location interface {
	LastKnownLocation() string
	Protocol() string
}
