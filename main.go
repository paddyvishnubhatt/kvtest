package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"kvtest/kv"
	"kvtest/server"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var mykv *kv.KV

var (
	serverAddr = flag.String("address", "localhost:10001", "TCP host+port for this node")
	nodeId     = flag.String("id", "node1", "Node ID")
)

func main() {
	flag.Parse()
	mykv = &kv.KV{}
	fmt.Println("Starting on " + string(*serverAddr))

	mykv.InitRaft(*serverAddr, *nodeId)

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
	router.HandleFunc("/AddVoter/{voter}/{id}/{bs}", handleAddVoter)

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

func handleAddVoter(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	voter := vars["voter"]
	id := vars["id"]
	bs := vars["bs"]
	mykv.AddVoter(voter, id, bs)
	response := map[string]string{
		"message": "Welcome to KV - AddVoter " + voter + " " + id,
	}
	json.NewEncoder(rw).Encode(response)
}
