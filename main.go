package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"

	leaderboard "crakapiV2/leaderboard"
	radiant "crakapiV2/radiant"
	account "crakapiV2/v1"
	winloss "crakapiV2/winloss"
)

var (
	c = cache.New(1*time.Minute, 1*time.Minute) // Initialize cache with 1-minute expiration
)

// CacheHandler is a wrapper that adds caching to the provided handler function
func CacheHandler(handler func(w http.ResponseWriter, r *http.Request) (interface{}, string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Generate a unique cache key based on the request URL
		cacheKey := r.URL.String()

		// Check if the data is in the cache
		if cachedData, found := c.Get(cacheKey); found {
			// Set the appropriate content type header
			if contentType, ok := cachedData.(string); ok && strings.HasPrefix(contentType, "text") {
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte(cachedData.(string)))
			} else {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(cachedData)
			}
			return
		}

		// If not found, call the actual handler function
		data, contentType, err := handler(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Store the data in the cache
		c.Set(cacheKey, data, cache.DefaultExpiration)

		// Write the response with the correct content type
		w.Header().Set("Content-Type", contentType)
		if contentType == "application/json" {
			json.NewEncoder(w).Encode(data)
		} else {
			w.Write([]byte(data.(string)))
		}
	}
}

// Wrapper function for existing handlers
func wrapHandler(existingHandler http.HandlerFunc, contentType string) func(w http.ResponseWriter, r *http.Request) (interface{}, string, error) {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, string, error) {
		// Capture the response using a ResponseRecorder
		rr := httptest.NewRecorder()
		existingHandler(rr, r)

		// Read the response body
		data := rr.Body.String()

		// If the content type is JSON, decode the data
		if contentType == "application/json" {
			var jsonData interface{}
			if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
				return nil, "", err
			}
			return jsonData, contentType, nil
		}

		// For text content type, return as string
		return data, contentType, nil
	}
}

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

	// Wrap existing handlers with their appropriate content type
	router.HandleFunc("/v1/account/{name}/{tag}", CacheHandler(wrapHandler(account.AccountHandler, "text/plain"))).Methods("GET")
	router.HandleFunc("/v1/hs/{region}/{puuid}", CacheHandler(wrapHandler(account.HsHandler, "text/plain"))).Methods("GET")
	router.HandleFunc("/v1/wl/{region}/{puuid}", CacheHandler(wrapHandler(winloss.WLHandler, "text/plain"))).Methods("GET")
	router.HandleFunc("/v1/kd/{region}/{puuid}", CacheHandler(wrapHandler(winloss.KDAHandler, "text/plain"))).Methods("GET")
	router.HandleFunc("/v1/rr/{region}/{puuid}", CacheHandler(wrapHandler(account.RrHandler, "text/plain"))).Methods("GET") // Assuming text/plain for demonstration
	router.HandleFunc("/v1/lb/{region}/{puuid}", CacheHandler(wrapHandler(leaderboard.Handler, "application/json"))).Methods("GET")
	router.HandleFunc("/v1/all/{name}/{tag}", AllDataHandler)
	router.HandleFunc("/v1/lbr/{region}/{puuid}", CacheHandler(wrapHandler(radiant.MMRHandler, "text/plain"))).Methods("GET")

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./build/"))))

	log.Fatal(http.ListenAndServe(":3000", router))
}
