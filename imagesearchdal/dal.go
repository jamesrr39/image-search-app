package imagesearchdal

import (
	"bytes"
	"errors"
	"image-search-app/imagesearch"
	"image-search-app/imagesearchcache"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ImageSearchDAL struct {
	cache *imagesearchcache.ImageSearchCache
}

func NewImageSearchDAL(cache *imagesearchcache.ImageSearchCache) *ImageSearchDAL {
	return &ImageSearchDAL{cache}
}

func (dal *ImageSearchDAL) ScanDir(dirpath string) error {
	err := filepath.Walk(dirpath, func(path string, fileinfo os.FileInfo, err error) error {
		if nil != err {
			return err
		}

		if fileinfo.IsDir() {
			return nil
		}

		err = dal.EnsureInCache(path)

		return err
	})
	if nil != err {
		return err
	}
	return nil
}

func (dal *ImageSearchDAL) EnsureInCache(path string) error {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpeg", ".jpg":

		fileBytes, err := ioutil.ReadFile(path)
		if nil != err {
			return err
		}

		// get sha1 hash
		fileHash, err := imagesearch.HashOfFile(bytes.NewBuffer(fileBytes))
		if nil != err {
			return err
		}

		// if not in cache, build file descriptor
		descriptorFromCache := dal.cache.Get(fileHash)
		if nil != descriptorFromCache {
			log.Printf("%s (%s) already in cache\n", fileHash, path)
			return nil // already in cache
		}

		descriptor, err := imagesearch.FileDescriptorFromFileBytes(fileBytes)
		if nil != err {
			log.Printf("error reading %s. Error: %s\n", path, err)
			return errors.New("error reading " + path + ". Error: " + err.Error())
		}

		err = dal.cache.EnsureInCache(&imagesearchcache.PersistedImageDescriptor{descriptor, path})
		if nil != err {
			log.Printf("error writing to cache: %s\n", err)
			return err
		}

	default:
		log.Printf("Skipped %v\n", path)
	}
	return nil
}
