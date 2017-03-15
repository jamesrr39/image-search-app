package qtybinsview

import (
	"image-search-app/imagesearch"
	"strconv"

	"github.com/mattn/go-gtk/gtk"
)

type QtyBinsView struct {
	Container  gtk.IWidget
	hbinsEntry *gtk.Entry
	sbinsEntry *gtk.Entry
	vbinsEntry *gtk.Entry
}

var defaultQtyBins = &(imagesearch.NewQtyBins(8, 12, 3))

// todo
func NewQtyBinsView(startingValues *imagesearch.QtyBins) *QtyBinsView {
	if nil == startingValues {
		startingValues = defaultQtyBins
	}

}

func (q *QtyBinsView) GetAll() (imagesearch.QtyBins, error) {
	hbins, err := strconv.ParseUint(q.hbinsEntry.GetText(), 10, 64)
	if nil != err {
		return nil, err
	}

	sbins, err := strconv.ParseUint(q.sbinsEntry.GetText(), 10, 64)
	if nil != err {
		return nil, err
	}

	vbins, err := strconv.ParseUint(q.vbinsEntry.GetText(), 10, 64)
	if nil != err {
		return nil, err
	}

	return imagesearch.NewQtyBins(uint(hbins), uint(sbins), uint(vbins))
}
