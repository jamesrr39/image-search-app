package imagesearch

type QtyBins struct {
	HBins uint
	SBins uint
	VBins uint
}

func NewQtyBins(hBins, sBins, vBins uint) QtyBins {
	return QtyBins{hBins, sBins, vBins}
}
