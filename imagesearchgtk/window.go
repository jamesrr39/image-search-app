package imagesearchgtk

import (
	"fmt"
	"image"
	"image-search-app/imagesearch"
	"image-search-app/imagesearchdal"
	_ "image/gif"  // decode
	_ "image/jpeg" // decode
	_ "image/png"  // decode
	"log"
	"os"
	"unsafe"

	"github.com/mattn/go-gtk/gdk"

	gtkimageextra "github.com/jamesrr39/goutil/image_gtk_image_bridge"

	"github.com/disintegration/imaging"

	"github.com/mattn/go-gtk/gtk"
)

type Window struct {
	dal         *imagesearchdal.ImageSearchDAL
	window      *gtk.Window
	seedPicture *gtk.Image
	matchesHbox *gtk.HBox
}

type WindowOptions struct {
	SeedPicturePath string
}

func NewWindow(dal *imagesearchdal.ImageSearchDAL, options *WindowOptions) *Window {
	window := &Window{dal, gtk.NewWindow(gtk.WINDOW_TOPLEVEL), gtk.NewImage(), gtk.NewHBox(true, 1)}

	matchesViewport := gtk.NewViewport(nil, nil)

	matchesViewport.Add(window.matchesHbox)

	fileChooserBtn := gtk.NewFileChooserButton("choose file", gtk.FILE_CHOOSER_ACTION_OPEN)
	fileChooserBtn.Connect("file-set", func() {
		window.setMainPicture(fileChooserBtn.GetFilename())
	})

	windowVBox := gtk.NewVBox(false, 0)

	window.window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	windowVBox.PackStart(fileChooserBtn, false, false, 0)
	windowVBox.PackStart(window.seedPicture, false, true, 0)
	windowVBox.PackStart(matchesViewport, true, true, 0)

	window.window.Add(windowVBox)

	window.window.SetSizeRequest(1024, 960)
	window.window.ShowAll()

	return window

}

func getImage(path string, xSize, ySize int) (image.Image, error) {
	file, err := os.Open(path)
	if nil != err {
		return nil, err
	}

	picture, _, err := image.Decode(file)
	if nil != err {
		return nil, err
	}

	return imaging.Resize(picture, xSize, ySize, imaging.Lanczos), err

}

func (window *Window) setMainPicture(path string) {

	picture, err := getImage(path, 400, 400)
	if nil != err {
		panic(err)
	}
	pixBuf := gtkimageextra.NewGdkPixBufFromImage(picture)
	window.seedPicture.SetFromPixbuf(pixBuf)

	descriptor, err := imagesearch.NewImageDescriptorFromFile(path)
	if nil != err {
		panic(err) // todo handle in gui
	}
	window.window.ShowAll()

	log.Printf("image descriptor: %v\n", descriptor)
	go func() {
		matches := window.dal.Search(descriptor)
		gdk.ThreadsEnter()
		window.setMatchesPictures(matches)
		gdk.ThreadsLeave()
	}()
}

func (window *Window) setMatchesPictures(matches []*imagesearchdal.DescriptorWithMatchScore) {

	children := window.matchesHbox.GetChildren()
	children.ForEach(func(data unsafe.Pointer, d interface{}) {
		log.Printf("child pointer: %v\n", data)
		children.Remove(data)
	})

	matchXSize := 150
	matchYSize := 150

	first10matches := matches[:10]
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
		vbox.PackStart(gtk.NewLabel(fmt.Sprintf("%f", match.Score)), false, false, 0)
		window.matchesHbox.PackStart(vbox, false, false, 0)
	}
	window.window.ShowAll()
}
