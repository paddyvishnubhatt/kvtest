package main

import (
	"encoding/json"
	"fmt"
	"kvtest/kv"
	"kvtest/server"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var mykv *kv.KV

func main() {
	mykv = &kv.KV{}
	mykv.InitRaft()

	router := mux.NewRouter()

	server := &server.HTTPServer{
		Server: http.Server{
			Addr:    ":9000",
			Handler: router,
		},
		ShutdownReq: make(chan bool),
	}

	router.HandleFunc("/", handleMain)
	router.HandleFunc("/Get/{key}", handleGet)
	router.HandleFunc("/Put/{key}/{val}", handlePut)

	log.Println("KV Server is running!")
	done := make(chan bool)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("KV Server Main Listen and serve: %v", err)
		}
		done <- true
	}()

	//wait shutdown
	server.WaitShutdown()

	server.CloseMain()

	<-done
	log.Printf("KV Server DONE!")
}

func handleMain(rw http.ResponseWriter, r *http.Request) {
	log.Println("main.handleMain")

	response := map[string]string{
		"message": "Welcome to KV - Main ",
	}
	json.NewEncoder(rw).Encode(response)
}

func handleGet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	val := mykv.Get(key)
	response := map[string]string{
		"message": "Welcome to KV - Get (" + key + "): " + val,
	}
	json.NewEncoder(rw).Encode(response)

	fmt.Println("In handleGet " + key + ": " + val)
}

func handlePut(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	val := vars["val"]
	mykv.Put(key, val)
	log.Println("main.handlePut " + key + " " + val)

	response := map[string]string{
		"message": "Welcome to KV - Put ",
	}
	json.NewEncoder(rw).Encode(response)
}
