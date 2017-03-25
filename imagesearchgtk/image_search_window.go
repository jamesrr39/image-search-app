package imagesearchgtk

import (
	"image"
	"github.com/jamesrr39/image-search-app/imagesearch"
	"github.com/jamesrr39/image-search-app/imagesearch/chisquaredscorers"
	"github.com/jamesrr39/image-search-app/imagesearchstoragedal"
	_ "image/gif"  // decode
	_ "image/jpeg" // decode
	_ "image/png"  // decode
	"log"
	"os"

	"github.com/mattn/go-gtk/gdk"

	gtkimageextra "github.com/jamesrr39/goutil/image_gtk_image_bridge"

	"github.com/disintegration/imaging"

	"github.com/jamesrr39/goutil/image-processing/imageprocessingutil"

	"github.com/mattn/go-gtk/gtk"
)

type Window struct {
	dal               imagesearchstoragedal.ImageSearchStorageDAL
	window            *gtk.Window
	seedPicture       *gtk.Image
	matchesContainer  *MatchesContainer
	algorithmComboBox *gtk.ComboBox
}

type WindowOptions struct {
	SeedPicturePath string
	QtyBins         imagesearch.QtyBins
}

type scoringAlgorithmDisplay struct {
	text   string
	scorer imagesearch.ImageScorer
}

var scoringAlgorithms = []scoringAlgorithmDisplay{
	scoringAlgorithmDisplay{"Chi1", &chisquaredscorers.ChiDistanceSearchImpl1{}},
}

func NewWindow(dal imagesearchstoragedal.ImageSearchStorageDAL, options *WindowOptions) *Window {
	window := &Window{dal, gtk.NewWindow(gtk.WINDOW_TOPLEVEL), gtk.NewImage(), NewMatchesContainer(), gtk.NewComboBoxNewText()}
	for i := 0; i < len(scoringAlgorithms); i++ {
		window.algorithmComboBox.AppendText(scoringAlgorithms[i].text)
	}
	window.algorithmComboBox.SetActive(0)

	fileChooserBtn := gtk.NewFileChooserButton("choose file", gtk.FILE_CHOOSER_ACTION_OPEN)
	fileChooserBtn.Connect("file-set", func() {
		window.setMainPicture(fileChooserBtn.GetFilename(), options.QtyBins)
	})

	windowVBox := gtk.NewVBox(false, 0)

	window.window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	windowVBox.PackStart(fileChooserBtn, false, false, 0)
	windowVBox.PackStart(window.seedPicture, false, true, 0)
	windowVBox.PackStart(window.algorithmComboBox, false, false, 0)
	windowVBox.PackStart(window.matchesContainer.Container, true, true, 0)

	window.window.Add(windowVBox)

	window.window.SetSizeRequest(1024, 960)
	window.window.ShowAll()

	return window

}

func (window *Window) getScoringAlgorithm() imagesearch.ImageScorer {
	index := window.algorithmComboBox.GetActive()
	return scoringAlgorithms[index].scorer
}

func getImage(path string, xSize, ySize int) (image.Image, error) {
	file, err := os.Open(path)
	if nil != err {
		return nil, err
	}
	defer file.Close()
	log.Printf("getting image %s\n", path)
	picture, err := imageprocessingutil.RotateAndTransformPicture(file)
	if nil != err {
		return nil, err
	}

	if nil == picture {
		panic("nil picture: " + path)
	}

	return imaging.Resize(picture, xSize, ySize, imaging.Lanczos), nil

}

func (window *Window) setMainPicture(path string, qtyBins imagesearch.QtyBins) {

	picture, err := getImage(path, 400, 400)
	if nil != err {
		panic(err)
	}
	pixBuf := gtkimageextra.NewGdkPixBufFromImage(picture)
	window.seedPicture.SetFromPixbuf(pixBuf)

	window.window.ShowAll()

	go func(seedPicturePath string, qtyBins imagesearch.QtyBins) {
		file, err := os.Open(seedPicturePath)
		if nil != err {
			panic(err)
		}
		defer file.Close()
		log.Println("searching")
		matches, err := window.dal.Search(file, qtyBins, window.getScoringAlgorithm(), imagesearch.NewLocalLocation(path))
		if nil != err {
			panic(err)
		}
		log.Printf("matches: %v\n", matches)
		gdk.ThreadsEnter()
		window.matchesContainer.SetMatchesPictures(matches)
		gdk.ThreadsLeave()
	}(path, qtyBins)
}
