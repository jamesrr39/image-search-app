package imagesearchfscache

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"image-search-app/imagesearch"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ImageSearchFSCache struct {
	cachesLocation string
	memoryCache    *imageSearchCacheMap
}

func NewImageSearchFSCache(cachesLocation string) *ImageSearchFSCache {
	return &ImageSearchFSCache{cachesLocation, newImageSearchCacheMap()}
}

func (cache *ImageSearchFSCache) Ensure(file io.Reader, qtyBins imagesearch.QtyBins) (*imagesearch.ImageDescriptor, error) {
	fileBytes, err := ioutil.ReadAll(file)
	if nil != err {
		return nil, err
	}

	hasher := sha1.New()
	_, err = hasher.Write(fileBytes)
	if nil != err {
		return nil, err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))

	// try and find it in the in-memory cache
	descriptor := cache.memoryCache.Get(hash)
	if nil != descriptor {
		return descriptor, nil
	}

	// try and find it in the on-disk cache
	descriptor, err = cache.readFromDisk(hash)
	if nil != err {
		return nil, err
	}

	if nil != descriptor {
		return descriptor, nil
	}

	// otherwise make a new descriptor and persist it
	descriptor, err = imagesearch.FileDescriptorFromFile(bytes.NewBuffer(fileBytes), qtyBins)
	if nil != err {
		return nil, err
	}

	err = cache.writeToDisk(descriptor)
	if nil != err {
		return nil, err
	}

	return descriptor, nil

}

func (cache *ImageSearchFSCache) Search(file io.Reader, qtyBins imagesearch.QtyBins, scoringAlgorithm imagesearch.ImageScorer) ([]*imagesearch.DescriptorWithMatchScore, error) {
	seedDescriptor, err := cache.Ensure(file, qtyBins)
	if nil != err {
		return nil, err
	}

	var descriptorsWithScores []*imagesearch.DescriptorWithMatchScore
	for _, descriptor := range cache.memoryCache.GetAll() {
		score := scoringAlgorithm.Score(seedDescriptor, descriptor)

		descriptorsWithScores = append(descriptorsWithScores, imagesearch.NewDescriptorWithMatchScore(score, descriptor))
	}
	return descriptorsWithScores, nil
}

func (cache *ImageSearchFSCache) writeToDisk(descriptor *imagesearch.ImageDescriptor) error {
	file, err := os.Create(filepath.Join(cache.cachesLocation, descriptor.Sha1))
	if nil != err {
		return err
	}
	defer file.Close()

	err = gob.NewEncoder(file).Encode(descriptor)
	if nil != err {
		return err
	}

	return nil
}

func (cache *ImageSearchFSCache) readFromDisk(id string) (*imagesearch.ImageDescriptor, error) {
	file, err := os.Open(filepath.Join(cache.cachesLocation, id))
	if nil != err {
		if os.IsNotExist(err) {
			return nil, nil // not in on-disk cache yet
		}
		return nil, err
	}
	defer file.Close()

	var descriptor imagesearch.ImageDescriptor
	err = gob.NewDecoder(file).Decode(&descriptor)
	if nil != err {
		return nil, err
	}

	return &descriptor, nil
}
