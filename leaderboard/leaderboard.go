package leaderboard

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Leaderboard struct {
	PUUID  string `json:"puuid"`
	Region string `json:"region"`
}

type Data struct {
	Leaderboard Leaderboard `json:"leaderboard"`
}

func getLeaderboard(region, puuid string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("error loading .env file: %w", err)
	}

	apiKey := os.Getenv("HENRIK_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("HENRIK_KEY is not set in the environment")
	}

	url := fmt.Sprintf("https://api.henrikdev.xyz/valorant/v2/leaderboard/%s?puuid=%s", region, puuid)
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

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var response struct {
		Data []struct {
			LeaderboardRank int `json:"leaderboardRank"`
		} `json:"data"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response: %w", err)
	}

	if len(response.Data) == 0 {
		return "", fmt.Errorf("no data found in response")
	}

	leaderboardRank := strconv.Itoa(response.Data[0].LeaderboardRank)

	return leaderboardRank, nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	region := vars["region"]
	puuid := vars["puuid"]

	response, err := getLeaderboard(region, puuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", response)
}
