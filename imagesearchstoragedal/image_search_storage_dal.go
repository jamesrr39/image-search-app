package imagesearchstoragedal

import (
	"image-search-app/imagesearch"
	"io"
)

type ImageSearchStorageDAL interface {
	Ensure(file io.Reader, qtyBins imagesearch.QtyBins, location imagesearch.Location) (*imagesearch.PersistedImageDescriptor, error)
	Search(file io.Reader, qtyBIns imagesearch.QtyBins, scoringAlgorithm imagesearch.ImageScorer, location imagesearch.Location) ([]*imagesearch.PersistedDescriptorWithMatchScore, error)
}
