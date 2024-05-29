package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type store struct {
	lock  sync.Mutex
	store map[string]interface{}
}

func newStore() store {
	return store{
		lock:  sync.Mutex{},
		store: map[string]interface{}{},
	}
}

var keyValueStore = newStore()

func add(w http.ResponseWriter, req *http.Request) {
	key := req.FormValue("key")
	defer req.Body.Close()
	resBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	fmt.Println("The body is ", resBody)
	keyValueStore.lock.Lock()
	defer keyValueStore.lock.Unlock()
	if _, ok := keyValueStore.store[key]; ok {
		http.Error(w, "key already exist", 400)
	}
	keyValueStore.store[key] = resBody
}

func del(w http.ResponseWriter, req *http.Request) {
	key := req.FormValue("key")
	keyValueStore.lock.Lock()
	defer keyValueStore.lock.Unlock()
	if _, ok := keyValueStore.store[key]; !ok {
		http.Error(w, "key does note exist", 400)
	}
	delete(keyValueStore.store, key)
}

func update(w http.ResponseWriter, req *http.Request) {
	key := req.FormValue("key")
	defer req.Body.Close()
	resBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	keyValueStore.lock.Lock()
	defer keyValueStore.lock.Unlock()
	keyValueStore.store[key] = resBody
}

func main() {

	http.HandleFunc("/add", add)
	http.HandleFunc("/delete", del)
	http.HandleFunc("/update", update)
	http.ListenAndServe(":8080", nil)
}
