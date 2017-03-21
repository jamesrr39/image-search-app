package imagesearchcache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"image-search-app/imagesearch"
	"io"
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
	descriptorMap map[string]*imagesearch.PersistedImageDescriptor // map[sha1]*PersistedImageDescriptor
}

func NewImageSearchCache(cachesLocation string) (*ImageSearchCache, error) {
	cache := &ImageSearchCache{
		cachesLocation: cachesLocation,
		mu:             sync.Mutex{},
		descriptorMap:  make(map[string]*imagesearch.PersistedImageDescriptor),
	}

	err := cache.ensureCacheFolders()
	if nil != err {
		return nil, err
	}

	return cache, nil
}

func (cache *ImageSearchCache) getCacheFilePath(sha1 string) string {
	return filepath.Join(cache.cachesLocation, descriptorCachesFolderName, sha1)
}

// filePath to persist it
func (cache *ImageSearchCache) EnsureInCache(file io.Reader, qtyBins imagesearch.QtyBins, filePath string) error {
	fileBytes, err := ioutil.ReadAll(file)
	if nil != err {
		return err
	}

	// generate sha1
	sha1, err := imagesearch.HashOfFile(bytes.NewBuffer(fileBytes))
	if nil != err {
		return err
	}

	// look in in-memory cache
	descriptorFromCache := cache.Get(sha1)
	if nil != descriptorFromCache {
		log.Printf("descriptor for %s is already in cache (sha1: %s)\n", descriptorFromCache.LastDiskLocation, descriptorFromCache.Sha1)
		return nil // already in cache
	}

	descriptorCachePath := cache.getCacheFilePath(sha1)

	// look in on-disk cache (and load into in-memory cache if found)
	descriptorFile, _ := os.Open(descriptorCachePath)
	if nil != descriptorFile {
		// descriptorFile was found
		defer descriptorFile.Close()
		var descriptor *imagesearch.PersistedImageDescriptor
		err := gob.NewDecoder(descriptorFile).Decode(&descriptor)
		if nil != err {
			return err
		}
		cache.addToCache(descriptor)
		return nil
	}

	// if not found, generate descriptor and save
	descriptor, err := imagesearch.FileDescriptorFromFile(bytes.NewBuffer(fileBytes), qtyBins)
	if nil != err {
		return err
	}

	persistedDescriptor := imagesearch.NewPersistedImageDescriptor(descriptor, filePath)

	err = cache.persistCacheFile(persistedDescriptor)
	if nil != err {
		return err
	}

	cache.addToCache(persistedDescriptor)

	return nil

}

func (cache *ImageSearchCache) addToCache(descriptor *imagesearch.PersistedImageDescriptor) {

	cache.mu.Lock()
	cache.descriptorMap[descriptor.Sha1] = descriptor
	cache.mu.Unlock()
}

func (cache *ImageSearchCache) persistCacheFile(descriptor *imagesearch.PersistedImageDescriptor) error {
	descriptorCachePath := cache.getCacheFilePath(descriptor.ImageDescriptor.Sha1)

	log.Printf("writing cache to %s\n", descriptorCachePath)
	file, err := os.Create(descriptorCachePath)
	if nil != err {
		return err
	}

	err = gob.NewEncoder(file).Encode(descriptor)
	if nil != err {
		return err
	}

	return nil
}

func (cache *ImageSearchCache) Get(sha1 string) *imagesearch.PersistedImageDescriptor {
	return cache.descriptorMap[sha1]
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

func (cache *ImageSearchCache) DebugDump() {
	s := fmt.Sprintf("debug cache dump:\nitems in cache: %d\n", len(cache.descriptorMap))
	for _, item := range cache.GetAll() {
		s += fmt.Sprintf("at %s\n", item.LastDiskLocation)
	}
	log.Print(s)
}
