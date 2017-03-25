package imagesearchfscache

import (
	"image-search-app/imagesearch"
	"sync"
)

type imageSearchCacheMap struct {
	mu *sync.Mutex
	m  map[string]*imagesearch.PersistedImageDescriptor
}

func newImageSearchCacheMap() *imageSearchCacheMap {
	return &imageSearchCacheMap{new(sync.Mutex), make(map[string]*imagesearch.PersistedImageDescriptor)}
}

func (cacheMap *imageSearchCacheMap) Add(id string, descriptor *imagesearch.PersistedImageDescriptor) {
	cacheMap.mu.Lock()
	cacheMap.m[id] = descriptor
	cacheMap.mu.Unlock()
}

func (cacheMap *imageSearchCacheMap) Get(id string) *imagesearch.PersistedImageDescriptor {
	return cacheMap.m[id]
}

func (cacheMap *imageSearchCacheMap) GetAll() []*imagesearch.PersistedImageDescriptor {
	var descriptors []*imagesearch.PersistedImageDescriptor
	for _, descriptor := range cacheMap.m {
		descriptors = append(descriptors, descriptor)
	}
	return descriptors
}
