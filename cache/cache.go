package cache

import (
	"net/http"
	"strings"
	"sync"
)

type response struct {
	header http.Header
	code   int
	body   []byte
}

type memCache struct {
	lock sync.RWMutex
	data map[string]response
}

var (
	cache = memCache{data: map[string]response{}}
)

func set(resource string, response *response) {
	cache.lock.Lock()

	if response == nil {
		delete(cache.data, resource)
	} else {
		cache.data[resource] = *response
	}

	cache.lock.Unlock()
}

func get(resource string) *response {
	cache.lock.RLock()
	resp, _ := cache.data[resource]
	cache.lock.RUnlock()

	return &resp
}

func copyHeader(source, destination http.Header) {
	for key, list := range source {
		for _, listItem := range list {
			destination.Add(key, listItem)
		}
	}
}

//MakeResource function, map uri(key on the cache map) to its response value from cache
func MakeResource(request *http.Request) string {
	if request == nil {
		return ""
	}

	return strings.TrimSuffix(request.URL.RequestURI(), "/")
}

//Clean function, remove all entries on the cache map
func Clean() {
	cache.lock.Lock()
	cache.data = map[string]response{}
	cache.lock.Unlock()
}

//Drop function, remove a path(key on the cache map)
func Drop(resource string) {
	set(resource, nil)
}

//Serve function, return a bool value depending if the incoming request was cached and update the writter if so
func Serve(writter http.ResponseWriter, request *http.Request) bool {
	if writter == nil || request == nil {
		return false
	}

	if request.Header.Get("Cache-Control") == "no-cache" {
		return false
	}

	response := get(MakeResource(request))

	if response == nil {
		return false
	}

	copyHeader(response.header, writter.Header())
	writter.WriteHeader(response.code)

	if request.Method != http.MethodHead {
		writter.Write(response.body)
	}

	return true
}
