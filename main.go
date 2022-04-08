package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var lastX int
var lastY int
var isStuck bool

func main() {
	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	http.HandleFunc("/", handler)

	log.Printf("starting server on port :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatalf("http listen error: %v", err)
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		fmt.Fprint(w, "Let the battle begin!")
		return
	}

	var v ArenaUpdate
	defer req.Body.Close()
	d := json.NewDecoder(req.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&v); err != nil {
		log.Printf("WARN: failed to decode ArenaUpdate in response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := play(v)
	fmt.Fprint(w, resp)
}

func play(input ArenaUpdate) (response string) {
	log.Printf("IN: %#v", input)

	var dir = input.Arena.State["https://cloud-run-hackathon-go-7dzaoqbgzq-uc.a.run.app"].Direction
	var posX = input.Arena.State["https://cloud-run-hackathon-go-7dzaoqbgzq-uc.a.run.app"].X
	var posY = input.Arena.State["https://cloud-run-hackathon-go-7dzaoqbgzq-uc.a.run.app"].Y
	var dimX = input.Arena.Dimensions[0] - 1
	var dimY = input.Arena.Dimensions[1] - 1
	var wasHit = input.Arena.State["https://cloud-run-hackathon-go-7dzaoqbgzq-uc.a.run.app"].WasHit
	if lastX == posX || lastY == posY {
		isStuck = true
	} else {
		isStuck = false
	}

	log.Println("#######################################################")
	log.Println("DATA")
	log.Printf("dir:%v\n", dir)
	log.Printf("posX:%v\n", posX)
	log.Printf("posY:%v\n", posY)
	log.Printf("dimX:%v\n", dimX)
	log.Printf("dimY:%v\n", dimY)
	log.Printf("wasHit:%v\n", wasHit)
	log.Printf("lastX:%v\n", lastX)
	log.Printf("lastY:%v\n", lastY)
	log.Println("#######################################################")

	if dir == "E" && posX == dimX || dir == "S" && posY == dimY {
		return "L"
	}

	if dir == "W" && posX == 0 || dir == "N" && posY == 0 {
		return "L"
	}

	if !wasHit || !isStuck {
		return "F"
	}

	lastX = posX
	lastY = posY

	return "T"
}
