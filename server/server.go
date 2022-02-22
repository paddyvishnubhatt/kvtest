package server

import (
	"context"
	"encoding/json"
	"fmt"
	"kvtest/kv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	http.Server
	ShutdownReq chan bool
}

func (s *HTTPServer) WaitShutdown() {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	//Wait interrupt or shutdown request through /shutdown
	select {
	case sig := <-irqSig:
		log.Printf("Shutdown request (signal: %v)", sig)
	case sig := <-s.ShutdownReq:
		log.Printf("Shutdown request (/shutdown %v)", sig)
	}

	log.Printf("Stopping HTTP server ...")

	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//shutdown the server
	err := s.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutdown request error: %v", err)
	}
}

func (s *HTTPServer) ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Shutdown server"))
	go func() {
		s.ShutdownReq <- true
	}()
}

func (s *HTTPServer) CloseMain() {
	// clean up map
}

var mykv *kv.KeyVal

func RunServer(ikv *kv.KeyVal, serverAddr string, port int) {
	fmt.Println("Starting on " + serverAddr)
	mykv = ikv

	router := mux.NewRouter()

	server := &HTTPServer{
		Server: http.Server{
			Addr:    ":" + string(port),
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
