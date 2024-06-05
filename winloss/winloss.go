package winloss

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Response struct {
	Data []MatchData `json:"data"`
}

type MatchData struct {
	Meta  MetaData  `json:"meta"`
	Stats StatsData `json:"stats"`
	Teams TeamsData `json:"teams"`
}

type MetaData struct {
	StartedAt string `json:"started_at"`
	Mode      string `json:"mode"`
}

type StatsData struct {
	Team    string `json:"team"`
	Kills   int    `json:"kills"`
	Deaths  int    `json:"deaths"`
	Assists int    `json:"assists"`
}

type TeamsData struct {
	Red  *int `json:"red"`
	Blue *int `json:"blue"`
}

type Career struct {
	Record []string
}

func CalculateWinsLosses(response Response, modes []string) (Career, int, int, int) {
	wins := 0
	losses := 0
	draws := 0
	career := Career{Record: []string{}}
	modeSet := make(map[string]bool)

	for _, mode := range modes {
		modeSet[strings.ToLower(mode)] = true
	}

	currentTime := time.Now()

	for i := 0; i < len(response.Data); i++ {
		match := response.Data[i]
		matchStartedAt, err := time.Parse(time.RFC3339, match.Meta.StartedAt)
		if err != nil {
			continue
		}

		if i == 0 && currentTime.Sub(matchStartedAt) > 6*time.Hour {
			//fmt.Println("First match is more than 6 hours old. Breaking the loop.")
			break
		}

		if len(modes) > 0 && !modeSet[strings.ToLower(match.Meta.Mode)] {
			continue
		}
		// If the current match is more than 2 hours newer than the next match, break the loop
		if i < len(response.Data)-1 {
			prevMatchStartedAt, err := time.Parse(time.RFC3339, response.Data[i+1].Meta.StartedAt)
			if err != nil {
				continue
			}
			// Calculate the time difference between the current and next match
			timeDifference := matchStartedAt.Sub(prevMatchStartedAt)

			if timeDifference > 2*time.Hour {
				fmt.Println("Current match is more than 2 hours newer than the next match. Breaking the loop.")
				break
			}
		}
		if match.Meta.Mode == "Deathmatch" || match.Meta.Mode == "Custom" || match.Meta.Mode == "Team Deathmatch" {
			continue
		}

		// Check if the first match in the list is more than 6 hours old

		playerTeam := match.Stats.Team
		if match.Teams.Red == nil || match.Teams.Blue == nil {
			continue
		}

		redRoundsWon := *match.Teams.Red
		blueRoundsWon := *match.Teams.Blue

		var winningTeam string
		if redRoundsWon > blueRoundsWon {
			winningTeam = "Red"
		} else if blueRoundsWon > redRoundsWon {
			winningTeam = "Blue"
		} else {
			winningTeam = "Draw"
		}

		if winningTeam == playerTeam {
			wins++
			career.Record = append(career.Record, "W")
		} else if winningTeam == "Draw" {
			draws++
			career.Record = append(career.Record, "D")
		} else {
			losses++
			career.Record = append(career.Record, "L")
		}
	}

	career.Record = reverseSlice(career.Record)

	return career, wins, losses, draws
}

func reverseSlice(slice []string) []string {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func FetchMatchData(region, puuid string) (Response, error) {
	fmt.Println("Fetching match data")
	fmt.Println(region, puuid)
	err := godotenv.Load()
	if err != nil {
		return Response{}, fmt.Errorf("error loading .env file: %v", err)
	}
	apiKey := os.Getenv("HENRIK_KEY")
	if apiKey == "" {
		return Response{}, fmt.Errorf("HENRIK_KEY environment variable not set")
	}

	url := fmt.Sprintf("https://api.henrikdev.xyz/valorant/v1/by-puuid/lifetime/matches/%s/%s", region, puuid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Response{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return Response{}, fmt.Errorf("error decoding response: %w", err)
	}
	fmt.Println("response")
	return response, nil
}

func WLHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	region := vars["region"]
	puuid := vars["puuid"]
	var modes []string
	if modeValues, ok := r.URL.Query()["mode"]; ok {
		modes = modeValues
	}
	format := r.URL.Query().Get("fs")

	response, err := FetchMatchData(region, puuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	career, wins, losses, draws := CalculateWinsLosses(response, modes)

	if format == "json" {
		result := map[string]interface{}{
			"career": career,
			"wins":   wins,
			"losses": losses,
			"draws":  draws,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	} else {
		result := fmt.Sprintf("has won %d matches and lost %d matches.", wins, losses)
		if draws > 0 {
			result += fmt.Sprintf(" %d matches were draw.", draws)
		}
		result += fmt.Sprintf(" Career: %s", strings.Join(career.Record, ""))

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(result))
	}
}

func KDAHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	region := vars["region"]
	puuid := vars["puuid"]

	response, err := FetchMatchData(region, puuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	KDA, _ := CalculateKDA(response)
	result := fmt.Sprintf("%v", KDA)
	fmt.Fprint(w, result)
}

func CalculateKDA(response Response) (string, error) {
	if len(response.Data) == 0 {
		return "", fmt.Errorf("no data available")
	}

	stats := response.Data[0].Stats
	kills := stats.Kills
	deaths := stats.Deaths
	assists := stats.Assists

	return fmt.Sprintf("%d/%d/%d", kills, deaths, assists), nil
}
