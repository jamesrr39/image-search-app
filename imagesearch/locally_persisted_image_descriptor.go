package imagesearch

type LocallyPersistedImageDescriptor struct {
	lastKnownLocation string
}

func NewLocallyPersistedImageDescriptor(lastKnownLocation string) *LocallyPersistedImageDescriptor {
	return &LocallyPersistedImageDescriptor{lastKnownLocation}
}

func (d *LocallyPersistedImageDescriptor) LastKnownLocation() string {
	return d.lastKnownLocation
}

func (d *LocallyPersistedImageDescriptor) Protocol() string {
	return "file://"
}
