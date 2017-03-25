package imagesearch

type LocalLocation struct {
	LocationOnDisk string
}

func NewLocalLocation(lastKnownLocation string) *LocalLocation {
	return &LocalLocation{lastKnownLocation}
}

func (d *LocalLocation) LastKnownLocation() string {
	return d.LocationOnDisk
}

func (d *LocalLocation) Protocol() string {
	return "file://"
}
