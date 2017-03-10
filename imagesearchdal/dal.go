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
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

type ImageSearchDAL struct {
	cache *imagesearchcache.ImageSearchCache
}

func NewImageSearchDAL(cachesFolderLocation string) (*ImageSearchDAL, error) {
	cache, err := imagesearchcache.NewImageSearchCacheAndScan(cachesFolderLocation)
	if nil != err {
		return nil, err
	}
	return &ImageSearchDAL{cache}, nil
}

func (dal *ImageSearchDAL) ScanDir(dirpath string) error {

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

	var maxOps int32 = int32(runtime.NumCPU())
	var opsRunning int32
	for {

		if nil != err {
			break
		}

		isFinished := len(paths) == 0 && opsRunning == 0
		if isFinished {
			break
		}

		if maxOps <= opsRunning {
			continue
		}

		path := paths[0]
		log.Printf("p: %s\n", path)
		paths = paths[1:]

		atomic.AddInt32(&opsRunning, 1)
		go func(path string) {
			//defer atomic.AddInt32(&opsRunning, -1)
			log.Printf("Now %d ops running. Running %s\n", opsRunning, path)

			start := time.Now()
			innerErr := dal.EnsureInCache(path)
			if nil != innerErr {
				err = innerErr
			}
			end := time.Now()
			log.Printf("took %v to ensure %s in cache\n", end.Sub(start).Nanoseconds()/1000, path)

			atomic.AddInt32(&opsRunning, -1)
		}(path)
	}

	return nil
}

func (dal *ImageSearchDAL) EnsureInCache(path string) error {

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

	err = dal.cache.EnsureInCache(&imagesearch.PersistedImageDescriptor{descriptor, path})

	return nil
}
