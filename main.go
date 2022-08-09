package main

import (
	"fmt"
  "os"
	//"net/http"

	"github.com/KnutZuidema/golio"
	"github.com/KnutZuidema/golio/riot/lol"
  "github.com/KnutZuidema/golio/api"
	"github.com/sirupsen/logrus"
)
const (
   API_KEY = "RGAPI-7ee91f55-688e-4162-b9bb-66d4a4094b8c"
   RIOT_MAX_MATCHES_RETURNED = 100
)

// Map of Player Name to map of MatchID to MatchData
var matchDataByAlly = make(map[string]map[string]*MatchData)

// Map of Player Name to WinData
var winDataByAlly = make(map[string]*WinData)

type MatchData struct {

}

type WinData struct {
  Wins int
  Losses int
  Total int
}

func (w *WinData) Ratio() float64 {
  return float64(w.Wins)/float64(w.Total)
}

func main() {
	client := golio.NewClient(API_KEY,
                golio.WithRegion(api.RegionNorthAmerica),
                golio.WithLogger(logrus.New().WithField("foo", "bar")))
	summoner, _ := client.Riot.Summoner.GetByName("Chickaloo")
	fmt.Printf("%s is a level %d summoner\n", summoner.Name, summoner.SummonerLevel)

  RANKED_SOLO_5x5 := int(420)
  GAMES_PLAYED := 561
  for i := 0; i < 6; i++ {
    startIndex := i * RIOT_MAX_MATCHES_RETURNED
    endIndex := ((1+i)*RIOT_MAX_MATCHES_RETURNED)-1
    fmt.Printf("Retrieving matches %d to %d\n", startIndex, endIndex)


    gamesToRequest := 100
    if GAMES_PLAYED < 100 {
      gamesToRequest = GAMES_PLAYED
    }

    fmt.Printf("Requesting %d games, starting from %d", gamesToRequest, startIndex)

    matchIds, matchListErr := client.Riot.Match.List(summoner.PUUID, startIndex, gamesToRequest, &lol.MatchListOptions{Queue:&RANKED_SOLO_5x5})
    if matchListErr != nil {
      fmt.Println(matchListErr.Error())
      os.Exit(1)
    }

    for matchIndex := 0; matchIndex < RIOT_MAX_MATCHES_RETURNED; matchIndex ++ {
      matchId := matchIds[matchIndex]

      matchData, _ := client.Riot.Match.Get(matchId)

      participants := matchData.Info.Participants

      var teamId int

      // Find the team that the main is on, and then cache team data
      for participantIndex := 0; participantIndex < 10; participantIndex++ {
        participant := participants[participantIndex]

        if participant.PUUID == summoner.PUUID {
          teamId = participant.TeamID
        }
      }

      // Now, we can cache the win data for each ally.
      for participantIndex := 0; participantIndex < 10; participantIndex++ {
        participant := participants[participantIndex]

        if participant.TeamID == teamId && participant.PUUID != summoner.PUUID {
          if winDataByAlly[participant.SummonerName] == nil {
            winDataByAlly[participant.SummonerName] = new(WinData)
          }

          if participant.Win {
            winDataByAlly[participant.SummonerName].Wins += 1
          } else {
            winDataByAlly[participant.SummonerName].Losses += 1
          }
          winDataByAlly[participant.SummonerName].Total += 1
        }
      }

    }

    GAMES_PLAYED -= gamesToRequest
    // END MATCH PROCESSING
  }

  // Until the data is empty
  for len(winDataByAlly) > 0 {
    bestWR := float64(0)
    bestWRTotal := int(0)
    bestData := new(WinData)
    bestSummoner := ""

    // For each summoner
    for summonerName, data := range(winDataByAlly) {
      if data.Ratio() > bestWR {
        bestWR = data.Ratio()
        bestWRTotal = data.Total
        bestData = data
        bestSummoner = summonerName
      } else if data.Ratio() == bestWR {
        if data.Total > bestWRTotal {
          bestWR = data.Ratio()
          bestWRTotal = data.Total
          bestData = data
          bestSummoner = summonerName
        }
      }
    }
    if bestData.Total > 3 {
      fmt.Printf("Summoner: %s - W/L: %d - %d Ratio: %f\n", bestSummoner, bestData.Wins, bestData.Losses, bestData.Ratio())
    }
    delete(winDataByAlly, bestSummoner)
  }

}
