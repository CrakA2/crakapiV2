package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	leaderboard "crakapiV2/leaderboard"
	account "crakapiV2/v1"
	winloss "crakapiV2/winloss"

	"github.com/gorilla/mux"
)

func AllDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	tag := vars["tag"]

	playerData, err := account.GetAllPlayerData(name, tag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(playerData)
}

func main() {
	// Handle routes here
	fmt.Println("Starting server on : http://localhost:3000")
	router := mux.NewRouter()
	router.HandleFunc("/v1/account/{name}/{tag}", account.AccountHandler)
	router.HandleFunc("/v1/hs/{region}/{puuid}", account.HsHandler)
	router.HandleFunc("/v1/wl/{region}/{puuid}", winloss.WLHandler)
	router.HandleFunc("/v1/kd/{region}/{puuid}", winloss.KDAHandler)
	router.HandleFunc("/v1/rr/{region}/{puuid}", account.RrHandler)
	router.HandleFunc("/v1/lb/{region}/{puuid}", leaderboard.Handler)
	router.HandleFunc("/v1/all/{name}/{tag}", AllDataHandler)
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))

	log.Fatal(http.ListenAndServe(":3000", router))
}
