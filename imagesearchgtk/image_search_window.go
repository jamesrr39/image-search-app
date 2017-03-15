package imagesearchgtk

import (
	"image"
	"image-search-app/imagesearch"
	"image-search-app/imagesearch/searchers"
	"image-search-app/imagesearchdal"
	_ "image/gif"  // decode
	_ "image/jpeg" // decode
	_ "image/png"  // decode
	"log"
	"os"

	"github.com/mattn/go-gtk/gdk"

	gtkimageextra "github.com/jamesrr39/goutil/image_gtk_image_bridge"

	"github.com/disintegration/imaging"

	"image-processing/imageprocessingutil"

	"github.com/mattn/go-gtk/gtk"
)

type Window struct {
	dal               *imagesearchdal.ImageSearchDAL
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
	scorer searchers.ImageScorer
}

var scoringAlgorithms = []scoringAlgorithmDisplay{
	scoringAlgorithmDisplay{"Basic", &searchers.BasicScorer{}},
	scoringAlgorithmDisplay{"Chi1", &searchers.ChiDistanceSearchImpl1{}},
}

func NewWindow(dal *imagesearchdal.ImageSearchDAL, options *WindowOptions) *Window {
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

func (window *Window) getScoringAlgorithm() searchers.ImageScorer {
	index := window.algorithmComboBox.GetActive()
	return scoringAlgorithms[index].scorer
}

func getImage(path string, xSize, ySize int) (image.Image, error) {
	file, err := os.Open(path)
	if nil != err {
		return nil, err
	}
	defer file.Close()

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

	descriptor, err := imagesearch.NewImageDescriptorFromFile(path, qtyBins)
	if nil != err {
		panic(err) // todo handle in gui
	}
	window.window.ShowAll()

	log.Printf("image descriptor: %v\n", descriptor)
	go func() {
		matches := window.dal.Search(descriptor, window.getScoringAlgorithm())
		gdk.ThreadsEnter()
		window.matchesContainer.SetMatchesPictures(matches)
		gdk.ThreadsLeave()
	}()
}
