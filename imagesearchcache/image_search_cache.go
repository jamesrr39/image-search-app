package imagesearchcache

import (
	"encoding/gob"
	"image-search-app/imagesearch"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const descriptorCachesFolderName = "descriptor_caches"

type ImageSearchCache struct {
	cachesLocation string

	mu            sync.Mutex
	descriptorMap map[string]*imagesearch.PersistedImageDescriptor
}

func NewImageSearchCacheAndScan(cachesLocation string) (*ImageSearchCache, error) {
	cache := &ImageSearchCache{
		cachesLocation: cachesLocation,
		mu:             sync.Mutex{},
	}

	err := cache.ensureCacheFolders()
	if nil != err {
		return nil, err
	}

	err = cache.ScanCachesDir()
	if nil != err {
		return nil, err
	}
	return cache, nil
}

func (cache *ImageSearchCache) EnsureInCache(descriptor *imagesearch.PersistedImageDescriptor) error {
	descriptorFromCache := cache.Get(descriptor.Sha1)
	if nil != descriptorFromCache {
		log.Printf("descriptor for %s is already in cache (sha1: %s)\n", descriptor.LastDiskLocation, descriptor.Sha1)
		return nil // already in cache
	}

	descriptorCachePath := filepath.Join(cache.cachesLocation, descriptorCachesFolderName, descriptor.Sha1)

	log.Printf("writing cache to %s\n", descriptorCachePath)
	file, err := os.Create(descriptorCachePath)
	if nil != err {
		return err
	}

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(descriptor)
	if nil != err {
		return err
	}

	cache.mu.Lock()
	cache.descriptorMap[descriptor.Sha1] = descriptor
	cache.mu.Unlock()
	return nil

}

func (cache *ImageSearchCache) Get(sha1 string) *imagesearch.PersistedImageDescriptor {
	return cache.descriptorMap[sha1]
}

func (cache *ImageSearchCache) ScanCachesDir() error {
	dirPath := filepath.Join(cache.cachesLocation, descriptorCachesFolderName)
	fileInfos, err := ioutil.ReadDir(dirPath)
	if nil != err {
		return err
	}

	imageDescriptorMap := make(map[string]*imagesearch.PersistedImageDescriptor)

	for _, fileInfo := range fileInfos {
		file, err := os.Open(filepath.Join(cache.cachesLocation, descriptorCachesFolderName, fileInfo.Name()))
		if nil != err {
			return err
		}

		decoder := gob.NewDecoder(file)
		var descriptor *imagesearch.PersistedImageDescriptor
		err = decoder.Decode(&descriptor)
		if nil != err {
			return err
		}

		imageDescriptorMap[descriptor.Sha1] = descriptor
	}

	log.Printf("scanned %s and built map: %v\n", dirPath, imageDescriptorMap)
	cache.descriptorMap = imageDescriptorMap
	return nil
}

func (cache *ImageSearchCache) GetAll() []*imagesearch.PersistedImageDescriptor {
	var descriptors []*imagesearch.PersistedImageDescriptor

	for _, descriptor := range cache.descriptorMap {
		descriptors = append(descriptors, descriptor)
	}

	return descriptors
}

func (cache *ImageSearchCache) ensureCacheFolders() error {
	err := os.MkdirAll(cache.cachesLocation, 0755)
	if nil != err {
		return err
	}

	return os.MkdirAll(filepath.Join(cache.cachesLocation, descriptorCachesFolderName), 0755)
}
