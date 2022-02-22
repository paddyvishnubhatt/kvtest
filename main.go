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

var mykv *kv.KeyVal

var (
	serverAddr = flag.String("address", "localhost:10001", "TCP host+port for this node")
	nodeId     = flag.String("id", "node1", "Node ID")
	bs         = flag.Bool("bs", false, "Bootstrap")
)

func main() {
	flag.Parse()
	mykv = &kv.KeyVal{}

	if *bs {
		go runServer()
	}

	mykv.InitRaft(*serverAddr, *nodeId, *bs)
}

func runServer() {
	fmt.Println("Starting on " + string(*serverAddr))

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
	router.HandleFunc("/AddVoter/{voter}/{id}", handleAddVoter)

	log.Println("KV Server is running!")
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("KV Server Main Listen and serve: %v", err)
		}
	}()
	server.WaitShutdown()

	server.CloseMain()

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
	mykv.AddVoter(voter, id)
	response := map[string]string{
		"message": "Welcome to KV - AddVoter " + voter + " " + id,
	}
	json.NewEncoder(rw).Encode(response)
}
