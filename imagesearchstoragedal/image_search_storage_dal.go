package imagesearchstoragedal

import (
	"image-search-app/imagesearch"
	"io"
)

type ImageSearchStorageDAL interface {
	Ensure(file io.Reader, qtyBins imagesearch.QtyBins) (*imagesearch.ImageDescriptor, error) // (*imagesearch.PersistedImageDescriptor, error)
	Search(file io.Reader, qtyBIns imagesearch.QtyBins, scoringAlgorithm imagesearch.ImageScorer) ([]*imagesearch.DescriptorWithMatchScore, error)
}
