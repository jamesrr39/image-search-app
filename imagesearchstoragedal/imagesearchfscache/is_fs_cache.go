package imagesearchfscache

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"encoding/hex"
	"image-search-app/imagesearch"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/bradfitz/slice"
)

func init() {
	gob.RegisterName("image-search-app/imagesearch.LocalLocation", &imagesearch.LocalLocation{}) // locallocation
}

type ImageSearchFSCache struct {
	cachesLocation string
	memoryCache    *imageSearchCacheMap
}

func NewImageSearchFSCache(cachesLocation string) *ImageSearchFSCache {
	return &ImageSearchFSCache{cachesLocation, newImageSearchCacheMap()}
}

func (cache *ImageSearchFSCache) Ensure(file io.Reader, qtyBins imagesearch.QtyBins, location imagesearch.Location) (*imagesearch.PersistedImageDescriptor, error) {
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

	if nil == err && nil != descriptor {
		cache.memoryCache.Add(descriptor.Sha1, descriptor)
		return descriptor, nil
	}

	if nil != descriptor {
		return descriptor, nil
	}

	// otherwise make a new descriptor and persist it
	descriptor, err = imagesearch.FileDescriptorFromFile(bytes.NewBuffer(fileBytes), qtyBins, location)
	if nil != err {
		return nil, err
	}

	err = cache.writeToDisk(descriptor)
	if nil != err {
		return nil, err
	}

	return descriptor, nil

}

func (cache *ImageSearchFSCache) Search(file io.Reader, qtyBins imagesearch.QtyBins, scoringAlgorithm imagesearch.ImageScorer, location imagesearch.Location) ([]*imagesearch.PersistedDescriptorWithMatchScore, error) {
	seedDescriptor, err := cache.Ensure(file, qtyBins, location)
	if nil != err {
		return nil, err
	}

	var descriptorsWithScores []*imagesearch.PersistedDescriptorWithMatchScore
	for _, descriptor := range cache.memoryCache.GetAll() {
		score := scoringAlgorithm.Score(seedDescriptor.ImageDescriptor, descriptor.ImageDescriptor)

		descriptorsWithScores = append(descriptorsWithScores, imagesearch.NewPersistedDescriptorWithMatchScore(score, descriptor))
	}

	slice.Sort(descriptorsWithScores, func(i, j int) bool {
		return descriptorsWithScores[i].MatchScore < descriptorsWithScores[j].MatchScore
	})
	return descriptorsWithScores, nil
}

func (cache *ImageSearchFSCache) writeToDisk(descriptor *imagesearch.PersistedImageDescriptor) error {
	descriptorPath := filepath.Join(cache.cachesLocation, descriptor.Sha1)

	log.Printf("writing descriptor to disk at %s\n", descriptorPath)
	file, err := os.Create(descriptorPath)
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

func (cache *ImageSearchFSCache) readFromDisk(id string) (*imagesearch.PersistedImageDescriptor, error) {
	descriptorPath := filepath.Join(cache.cachesLocation, id)

	file, err := os.Open(descriptorPath)
	if nil != err {
		if os.IsNotExist(err) {
			return nil, nil // not in on-disk cache yet
		}
		return nil, err
	}
	defer file.Close()

	log.Printf("reading descriptor from %s\n", descriptorPath)

	var descriptor imagesearch.PersistedImageDescriptor
	err = gob.NewDecoder(file).Decode(&descriptor)
	if nil != err {
		return nil, err
	}

	return &descriptor, nil
}
