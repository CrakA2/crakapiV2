package radiant

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

type MMRResponse struct {
	Data struct {
		Elo int `json:"elo"`
	} `json:"data"`
}

type LeaderboardResponse struct {
	Players []struct {
		RankedRating int `json:"rankedRating"`
	} `json:"players"`
}

func FetchAndParseMMR(region string, puuid string) (int, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	apiKey := os.Getenv("HENRIK_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("API key not found in environment variables")
	}

	url := fmt.Sprintf("https://api.henrikdev.xyz/valorant/v1/by-puuid/mmr/%s/%s", region, puuid)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("received non-200 response status: %s", resp.Status)
	}

	var mmrResp MMRResponse
	if err := json.NewDecoder(resp.Body).Decode(&mmrResp); err != nil {
		return 0, err
	}

	return mmrResp.Data.Elo - 2100, nil
}

func FetchRadiantMMR() (int, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	apiKey := os.Getenv("HENRIK_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("API key not found in environment variables")
	}

	url := "https://api.henrikdev.xyz/valorant/v2/leaderboard/ap"
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("received non-200 response status: %s", resp.Status)
	}

	var leaderboardResp LeaderboardResponse
	if err := json.NewDecoder(resp.Body).Decode(&leaderboardResp); err != nil {
		return 0, err
	}

	if len(leaderboardResp.Players) > 499 {
		return leaderboardResp.Players[499].RankedRating, nil
	}

	return 0, fmt.Errorf("not enough players in leaderboard")
}

func CalculateRRRequired(mmrCurrent int) (interface{}, error) {
	radiantMMR, err := FetchRadiantMMR()
	if err != nil {
		return nil, err
	}

	if mmrCurrent < 0 {
		return "Player is not Immortal", nil
	}

	rrRequired := radiantMMR - mmrCurrent
	if rrRequired < 0 {
		return "Player is Radiant", nil
	}

	return fmt.Sprintf("%dRR required for Radiant", rrRequired), nil
}

func MMRHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	region := vars["region"]
	puuid := vars["puuid"]

	// Fetch current MMR
	mmrCurrent, err := FetchAndParseMMR(region, puuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate RR required
	rrRequired, err := CalculateRRRequired(mmrCurrent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the query parameter 'fs' is set to 'json'
	if r.URL.Query().Get("fs") == "json" {
		response := map[string]interface{}{
			"rr_required": rrRequired,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "%s", rrRequired)
	}
}
