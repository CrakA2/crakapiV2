package account

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Account struct {
	PUUID  string `json:"puuid"`
	Region string `json:"region"`
}

func getAccount(name, tag string) (*Account, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	apiKey := os.Getenv("HENRIK_KEY")

	url := fmt.Sprintf("https://api.henrikdev.xyz/valorant/v1/account/%s/%s", name, tag)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data Account `json:"data"`
	}
	err = json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func AccountHandler(w http.ResponseWriter, r *http.Request) { // Exported function
	vars := mux.Vars(r)
	name := vars["name"]
	tag := vars["tag"]
	fs := r.URL.Query().Get("fs")

	account, err := getAccount(name, tag)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if fs == "json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]Account{"account": *account})
	} else {
		fmt.Fprintf(w, "%s", account)
	}
}

type HSApiResponse struct {
	Data []struct {
		Stats struct {
			Shots struct {
				Head int `json:"head"`
				Body int `json:"body"`
				Leg  int `json:"leg"`
			} `json:"shots"`
		} `json:"stats"`
	} `json:"data"`
}

func getHeadshot(region, puuid string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %v", err)
	}
	apiKey := os.Getenv("HENRIK_KEY")

	url := fmt.Sprintf("https://api.henrikdev.xyz/valorant/v1/by-puuid/lifetime/matches/%s/%s", region, puuid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}
	var apiResponse HSApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	if len(apiResponse.Data) > 0 {
		shots := apiResponse.Data[0].Stats.Shots
		totalShots := shots.Head + shots.Body + shots.Leg
		headshots := shots.Head
		var headshotPercentage float64
		if totalShots > 0 {
			headshotPercentage = float64(headshots) / float64(totalShots) * 100
		} else {
			headshotPercentage = 0
		}
		return fmt.Sprintf("%.1f%%", headshotPercentage), nil
	}

	return "", errors.New("no data available")
}

func HsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Region := vars["region"]
	PUUID := vars["puuid"]

	headshotPercentage, err := getHeadshot(Region, PUUID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	headshotPercentage = fmt.Sprintf("%s%%", headshotPercentage)
	fmt.Fprint(w, headshotPercentage)
}

func RrHandler(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		http.Error(w, "error loading .env file", http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	region := vars["region"]
	puuid := vars["puuid"]

	rrResponse, err := getRR(region, puuid, os.Getenv("HENRIK_KEY"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s, %d", rrResponse.CurrentTierPatched, rrResponse.RankingInTier)
}

type MMRResponse struct {
	Data struct {
		CurrentTierPatched string `json:"currenttierpatched"`
		RankingInTier      int    `json:"ranking_in_tier"`
	} `json:"data"`
}

type RRResponse struct {
	CurrentTierPatched string
	RankingInTier      int
}

func getRR(region, puuid, apiKey string) (*RRResponse, error) {
	url := fmt.Sprintf("https://api.henrikdev.xyz/valorant/v1/by-puuid/mmr/%s/%s", region, puuid)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching data, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var mmrResponse MMRResponse
	err = json.Unmarshal(body, &mmrResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	rrResponse := &RRResponse{
		CurrentTierPatched: mmrResponse.Data.CurrentTierPatched,
		RankingInTier:      mmrResponse.Data.RankingInTier,
	}
	return rrResponse, nil
}
