package imagesearchgtk

import (
	"fmt"
	"image-search-app/imagesearchdal"
	"log"

	gtkimageextra "github.com/jamesrr39/goutil/image_gtk_image_bridge"
	"github.com/mattn/go-gtk/gtk"
)

type MatchesContainer struct {
	Container             *gtk.ScrolledWindow
	imageWidgetsContainer *gtk.HBox
}

func NewMatchesContainer() *MatchesContainer {

	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin.SetShadowType(gtk.SHADOW_IN)

	return &MatchesContainer{swin, nil}
}

func (matchesContainer *MatchesContainer) SetMatchesPictures(matches []*imagesearchdal.DescriptorWithMatchScore) {
	if nil != matchesContainer.imageWidgetsContainer {
		log.Printf("about to detach\n")
		matchesContainer.imageWidgetsContainer.Destroy()
	}

	matchXSize := 150
	matchYSize := 150

	hbox := gtk.NewHBox(true, 0)

	var first10matches []*imagesearchdal.DescriptorWithMatchScore
	if len(matches) < 10 {
		first10matches = matches
	} else {
		first10matches = matches[:10]
	}
	for _, match := range first10matches {

		picture, err := getImage(match.Descriptor.LastDiskLocation, matchXSize, matchYSize)
		if nil != err {
			panic(err)
		}
		matchImageWidget := gtkimageextra.NewGtkImageFromImage(picture)
		matchImageWidget.SetSizeRequest(matchXSize, matchYSize)
		vbox := gtk.NewVBox(false, 2)
		vbox.SetSizeRequest(matchXSize, matchYSize+40)
		vbox.PackStart(matchImageWidget, false, false, 0)
		vbox.PackStart(gtk.NewLabel(fmt.Sprintf("score: %f", match.MatchScore)), false, false, 0)
		hbox.PackStart(vbox, false, false, 0)
	}
	matchesContainer.imageWidgetsContainer = hbox
	matchesContainer.Container.AddWithViewPort(matchesContainer.imageWidgetsContainer)

	matchesContainer.Container.ShowAll()

}
