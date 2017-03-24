package imagesearchfscache

import (
	"image-search-app/imagesearch"
	"sync"
)

type imageSearchCacheMap struct {
	mu *sync.Mutex
	m  map[string]*imagesearch.ImageDescriptor
}

func newImageSearchCacheMap() *imageSearchCacheMap {
	return &imageSearchCacheMap{new(sync.Mutex), make(map[string]*imagesearch.ImageDescriptor)}
}

func (cacheMap *imageSearchCacheMap) Add(id string, descriptor *imagesearch.ImageDescriptor) {
	cacheMap.mu.Lock()
	cacheMap.m[id] = descriptor
	cacheMap.mu.Unlock()
}

func (cacheMap *imageSearchCacheMap) Get(id string) *imagesearch.ImageDescriptor {
	return cacheMap.m[id]
}

func (cacheMap *imageSearchCacheMap) GetAll() []*imagesearch.ImageDescriptor {
	var descriptors []*imagesearch.ImageDescriptor
	for _, descriptor := range cacheMap.m {
		descriptors = append(descriptors, descriptor)
	}
	return descriptors
}
