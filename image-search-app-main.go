package main

import (
	"image-search-app/imagesearchdal"
	"image-search-app/imagesearchgtk"

	"github.com/mattn/go-gtk/gtk"

	"image-search-app/imagesearch"

	"github.com/jamesrr39/goutil/user"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	rootPath        = kingpin.Arg("rootpath", "set a path to scan under.").Required().String()
	seedPicturePath = kingpin.Arg("seed picture path", "filepath to the seed picture to search by. If blank, no picture is selected").String()
	cachesLocation  string
)

var defaultBins = imagesearch.NewQtyBins(8, 12, 3)

func main() {
	kingpin.Parse()

	var err error
	cachesLocation, err = user.ExpandUser("~/.cache/github.com/jamesrr39/image-search-app")
	if nil != err {
		panic(err)
	}

	err = run()
	if nil != err {
		panic(err)
	}
}

func run() error {

	dal, err := imagesearchdal.NewImageSearchDAL(cachesLocation)
	if nil != err {
		return err
	}

	if "" != *rootPath {
		expandedRootPath, err := user.ExpandUser(*rootPath)
		if nil != err {
			return err
		}
		err = dal.ScanDir(expandedRootPath, defaultBins)
		if nil != err {
			return err
		}
	}

	options := &imagesearchgtk.WindowOptions{SeedPicturePath: *seedPicturePath, QtyBins: defaultBins}

	gtk.Init(nil)

	glib.ThreadInit()
	gdk.ThreadsEnter()
	imagesearchgtk.NewWindow(dal, options)

	gtk.Main()
	return nil
}
