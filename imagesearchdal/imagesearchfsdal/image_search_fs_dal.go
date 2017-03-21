package imagesearchfsdal

import (
	"image-search-app/imagesearch"
	"image-search-app/imagesearchdal/imagesearchfsdal/imagesearchcache"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ImageSearchFsDAL struct {
	cache *imagesearchcache.ImageSearchCache
}

func NewImageSearchDAL(cachesFolderLocation string) (*ImageSearchFsDAL, error) {
	cache, err := imagesearchcache.NewImageSearchCache(cachesFolderLocation)
	if nil != err {
		return nil, err
	}
	cache.DebugDump()
	return &ImageSearchFsDAL{cache}, nil
}

func (dal *ImageSearchFsDAL) ScanDir(dirpath string, qtyBins imagesearch.QtyBins) error {
	log.Println("cache before scan:")
	dal.cache.DebugDump()

	var paths []string

	err := filepath.Walk(dirpath, func(path string, fileinfo os.FileInfo, err error) error {
		if nil != err {
			return err
		}

		if fileinfo.IsDir() {
			return nil
		}
		switch strings.ToLower(filepath.Ext(path)) {
		case ".jpeg", ".jpg", "png":
			paths = append(paths, path)
		default:
			log.Printf("Skipped %v\n", path)
		}
		return nil
	})
	if nil != err {
		return err
	}

	for _, path := range paths {
		//		currentPath = path
		err := dal.EnsureInCache(path, qtyBins)
		if nil != err {
			return err
		}
	}

	log.Println("cache after scan:")
	dal.cache.DebugDump()
	return nil
}

func (dal *ImageSearchFsDAL) EnsureInCache(path string, qtyBins imagesearch.QtyBins) error {
	file, err := os.Open(path)
	if nil != err {
		return err
	}
	defer file.Close()
	err = dal.cache.EnsureInCache(file, qtyBins, path)
	return err
}
