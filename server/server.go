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
	mykv        *kv.KeyVal
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

func (s *HTTPServer) RunServer(ikv *kv.KeyVal, serverAddr string, port int) {
	fmt.Printf("KV Server Starting on %v %d\n", serverAddr, port)
	p := fmt.Sprintf(":%d", port)
	router := mux.NewRouter()
	s.Addr = p
	s.mykv = ikv
	s.Handler = router
	s.ShutdownReq = make(chan bool)

	router.HandleFunc("/", s.handleMain)
	router.HandleFunc("/Get/{key}", s.handleGet)
	router.HandleFunc("/Put/{key}/{val}", s.handlePut)
	router.HandleFunc("/AddVoter/{voter}/{id}", s.handleAddVoter)

	log.Printf("KV Server is running %v %d\n", serverAddr, port)
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Printf("Error KV Server Main Listen and serve: %v\n", err)
		}
	}()
	s.WaitShutdown()

	s.CloseMain()

	log.Printf("KV Server DONE!")

}

func (s *HTTPServer) handleMain(rw http.ResponseWriter, r *http.Request) {
	log.Println("main.handleMain")

	response := map[string]string{
		"message": "Welcome to KV - Main ",
	}
	json.NewEncoder(rw).Encode(response)
}

func (s *HTTPServer) handleGet(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	val := s.mykv.Get(key)
	response := map[string]string{
		"message": "Welcome to KV - Get (" + key + "): " + val,
	}
	json.NewEncoder(rw).Encode(response)

	fmt.Println("In handleGet " + key + ": " + val)
}

func (s *HTTPServer) handlePut(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	val := vars["val"]
	s.mykv.Put(key, val)
	log.Println("main.handlePut " + key + " " + val)

	response := map[string]string{
		"message": "Welcome to KV - Put ",
	}
	json.NewEncoder(rw).Encode(response)
}

func (s *HTTPServer) handleAddVoter(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	voter := vars["voter"]
	id := vars["id"]
	s.mykv.AddVoter(voter, id)
	response := map[string]string{
		"message": "Welcome to KV - AddVoter " + voter + " " + id,
	}
	json.NewEncoder(rw).Encode(response)
}
