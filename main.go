package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Player struct {
	lastX   int
	lastY   int
	isStuck bool
}

func main() {
	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	var p = &Player{}
	http.HandleFunc("/", p.handler)

	log.Printf("starting server on port :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatalf("http listen error: %v", err)
}

func (u Player) handler(w http.ResponseWriter, req *http.Request) {
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

	resp := play(v, &u.lastX, &u.lastY, &u.isStuck)
	fmt.Fprint(w, resp)
}

func play(input ArenaUpdate, lastX *int, lastY *int, isStuck *bool) (response string) {
	log.Printf("IN: %#v", input)

	var dir = input.Arena.State["https://cloud-run-hackathon-go-7dzaoqbgzq-uc.a.run.app"].Direction
	var posX = input.Arena.State["https://cloud-run-hackathon-go-7dzaoqbgzq-uc.a.run.app"].X
	var posY = input.Arena.State["https://cloud-run-hackathon-go-7dzaoqbgzq-uc.a.run.app"].Y
	var dimX = input.Arena.Dimensions[0] - 1
	var dimY = input.Arena.Dimensions[1] - 1
	var wasHit = input.Arena.State["https://cloud-run-hackathon-go-7dzaoqbgzq-uc.a.run.app"].WasHit

	var ret string

	if dir == "E" && posX == dimX || dir == "S" && posY == dimY {
		ret = "L"
		return ret
	}

	if dir == "W" && posX == 0 || dir == "N" && posY == 0 {
		ret = "L"
		return ret
	}

	if (*lastX == posX) && (*lastY == posY) {
		*isStuck = true
		log.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@Stuck Shot")
		ret = "T"
	}

	*lastX = posX
	*lastY = posY

	if !wasHit {
		ret = "F"
	}

	// commands := []string{"F", "T"}
	// rand := rand2.Intn(2)
	// log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	// return commands[rand]
	log.Println("pppppppppppppppppppppppppppppppp Final Shot - " + ret)
	return ret
}
