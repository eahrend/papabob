package nflapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)


// High level wrapper function that gives us a client so we can talk to the endpoints
type NFLClient 	struct {
	endpoint  	string
	client 		*http.Client
	nflTeams map[string]string
}

// JSON Response from http://www.nfl.com/liveupdate/scores/scores.json
// Top level response and what not
type NFLScoresResponse struct {
	NFLScores map[string]NFLGameStats
}

type NFLGameStats struct {
	Home 		NFLTeamScore 	`json:"home"`
	Away 		NFLTeamScore 	`json:"away"`
	Down 		int 			`json:"down"`
	Togo 		int 			`json:"togo"`
	Clock 		string 			`json:"clock"`
	RedZone 	bool 			`json:"redzone"`
	Stadium 	string 			`json:"stadium"`
	YardLine 	string 			`json:"yl"`
	Quarter 	string 			`json:"qtr"`
}

type NFLTeamScore struct {
	Score 				NFLScore 	`json:"score"`
	TeamAbbreviation 	string 		`json:"abbr"`
	TeamName 			string 		`json:"team_name"`
	// Fuck if I know what to means here
	To 					int 		`json:"to"`
}

type NFLScore struct {
	FirstQuarter 	int `json:"1"`
	SecondQuarter 	int `json:"2"`
	ThirdQuarter 	int `json:"3"`
	FourthQuarter 	int `json:"4"`
	OverTime 		int `json:"5"`
	Total 			int `json:"T"`
}

// Probably should add support for a backend datastore so you can trigger lambda functions if the dataset changes, for subscription based services and what not
func NewNFLClient()(*NFLClient, error){
	client := &http.Client{}
	return NewNFLClientCustom(client)
}

func NewNFLClientCustom(httpClient *http.Client) (*NFLClient,error) {
	// If your client is nil, there is no reason to be calling this
	if httpClient == nil {
		return nil, errors.New("client is nil, please use NewNFLClient")
	}
	nflClient := &NFLClient{
		client:    httpClient,
		endpoint: "http://www.nfl.com/liveupdate/scores/scores.json",
	}
	nflTeams := map[string]string{
		"ARI": "Arizona Cardinals",
		"ATL": "Atlanta Falcons",
		"BAL": "Baltimore Ravens",
		"BUF": "Buffalo Bills",
		"CAR": "Carolina Panthers",
		"CHI": "Chicago Bears",
		"CIN": "Cincinnati Bengals",
		"CLE": "Cleveland Browns",
		"DAL": "Dallas Cowboys",
		"DEN": "Denver Broncos",
		"DET": "Detroit Lions",
		"GB": "Green Bay Packers",
		"HOU": "Houston Texans",
		"IND": "Indianapolis Colts",
		"JAX": "Jacksonville Jaguars",
		"KC": "Kansas City Chiefs",
		"MIA": "Miami Dolphins",
		"MIN": "Minnesota Vikings",
		"NE": "New England Patriots",
		"NO": "New Orleans Saints",
		"NYG": "New York Giants",
		"NYJ": "New York Jets",
		"OAK": "Oakland Raiders",
		"PHI": "Philadelphia Eagles",
		"PIT": "Pittsburgh Steelers",
		"SD": "San Diego Chargers",
		"SEA": "Seattle Seahawks",
		"SF": "San Francisco 49ers",
		"STL": "Saint Louis Rams",
		"TB": "Tampa Bay Buccaneers",
		"TEN": "Tennessee Titans",
		"WAS": "Washington Redskins",
	}
	nflClient.nflTeams = nflTeams
	return nflClient,nil
}


func (nfl *NFLClient) Get() (*NFLScoresResponse, error) {
	resp, err := nfl.client.Get(nfl.endpoint)
	if err != nil {
		return nil, err
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	gameStats := map[string]NFLGameStats{}
	nflScoreStats := &NFLScoresResponse{}
	err = json.Unmarshal(respBytes,&gameStats)
	//err = json.Unmarshal(respBytes, nflScoreStats)
	if err != nil {
		return nil, err
	}
	nflScoreStats.NFLScores = gameStats
	return nflScoreStats, nil
}

func(nfl *NFLClient) GetTeamFromAbbr(teamAbbr string) (string,error) {
	if teamName, ok := nfl.nflTeams[strings.ToUpper(teamAbbr)]; ok {
		return teamName, nil
	}
	return "", fmt.Errorf("no team could be found with the abbreviations %s",teamAbbr)
}



func(nfl *NFLClient) GetTodayScores() ([]NFLGameStats, error) {
	games, err := nfl.Get()
	if err != nil {
		return nil, err
	}
	t := time.Now()
	log.Println(games)
	todayDate := t.Format("20060102")
	log.Println(todayDate)
	var nflStats []NFLGameStats
	log.Println(len(games.NFLScores))
	for scoreDay, scoreStats := range games.NFLScores {
		if scoreDay[0:8] == todayDate {
			nflStats = append(nflStats,scoreStats)
		}
	}
	if len(nflStats) == 0 {
		return nil, fmt.Errorf("no games found today")
	}
	return nflStats,nil
}


func (nfl *NFLClient) BasicReport(nflGames []NFLGameStats)(string, error){

	if len(nflGames) == 0 {
		return "",fmt.Errorf("no games to report on")
	}
	var gameArray []string
	for _, game := range nflGames {
		homeTeam, _  := nfl.GetTeamFromAbbr(game.Home.TeamAbbreviation)
		awayTeam, _ := nfl.GetTeamFromAbbr(game.Away.TeamAbbreviation)
		gameString := fmt.Sprintf("%s playing at %s. %s has %d points and %s has %d points. Quarter is %v and time left is %s",
			homeTeam,
			awayTeam,
			homeTeam,
			game.Home.Score.Total,
			awayTeam,
			game.Away.Score.Total,
			game.Quarter,
			game.Clock)
		gameArray = append(gameArray,gameString)
	}
	outputString := strings.Join(gameArray,"\n")
	return outputString, nil
}