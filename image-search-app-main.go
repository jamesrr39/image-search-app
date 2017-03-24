package main

import (
	"image-search-app/imagesearchgtk"
	"image-search-app/imagesearchstoragedal/imagesearchfscache"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-gtk/gtk"

	"image-search-app/imagesearch"

	"github.com/jamesrr39/goutil/user"
	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	rootPath        = kingpin.Arg("rootpath", "set a path to scan under").Required().String()
	seedPicturePath = kingpin.Arg("seed picture path", "filepath to the seed picture to search by. If blank, no picture is selected").String()
	cachesLocation  string
)

var qtyBins = imagesearch.NewQtyBins(8, 12, 3)

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

	cache := imagesearchfscache.NewImageSearchFSCache(cachesLocation)

	expandedRootPath, err := user.ExpandUser(*rootPath)
	if nil != err {
		return err
	}

	err = filepath.Walk(expandedRootPath, func(path string, fileInfo os.FileInfo, err error) error {
		if nil != err {
			return err
		}

		if fileInfo.IsDir() {
			return nil
		}

		switch strings.ToLower(filepath.Ext(fileInfo.Name())) {
		case ".png", ".jpg", ".jpeg", ".gif":
			file, err := os.Open(path)
			if nil != err {
				return err
			}
			defer file.Close()
			_, err = cache.Ensure(file, qtyBins)
			if nil != err {
				return err
			}
		default:
			log.Printf("skipping %s\n", path)
		}
		return nil
	})
	if nil != err {
		return err
	}

	options := &imagesearchgtk.WindowOptions{SeedPicturePath: *seedPicturePath, QtyBins: qtyBins}

	gtk.Init(nil)

	glib.ThreadInit()
	gdk.ThreadsEnter()
	imagesearchgtk.NewWindow(cache, options)

	gtk.Main()
	return nil
}
