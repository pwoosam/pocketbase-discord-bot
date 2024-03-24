package games

// import (
// 	"encoding/json"
// 	"slices"

// 	"github.com/google/uuid"
// )

// type RPSTourney struct {
// 	Id              string
// 	Status          RPSTourneyStatus
// 	PlayersLost     []string
// 	Players         []string
// 	Games           [][]*RPSGame
// 	WinningPlayerId string
// 	Round           int
// }

// type RPSTourneyStatus int

// const (
// 	RPSTourneyStatusWaitingForPlayers RPSTourneyStatus = iota
// 	RPSTourneyStatusInProgress
// 	RPSTourneyStatusFinished
// )

// func NewRPSTourney() *RPSTourney {
// 	return &RPSTourney{
// 		Id:     uuid.NewString(),
// 		Status: RPSTourneyStatusWaitingForPlayers,
// 	}
// }

// func (t *RPSTourney) AddPlayer(playerId string) {
// 	t.Players = append(t.Players, playerId)
// }

// func (t *RPSTourney) NextRound() {
// 	mostRecentRound := t.GetMostRecentRound()

// 	// all games must be finished
// 	for _, g := range mostRecentRound {
// 		if g.Status != RPSGameStatusFinished {
// 			return
// 		}
// 	}

// 	// for each losing player, remove them from the tourney
// 	for _, g := range mostRecentRound {
// 		if g.WinningPlayerId != "" {
// 			t.PlayersLost = append(t.PlayersLost, g.LosingPlayerId)
// 			for i := 0; i < len(t.Players); i++ {
// 				if t.Players[i] == g.LosingPlayerId {
// 					t.Players = slices.Delete(t.Players, i, i+1)
// 					break
// 				}
// 			}
// 		}
// 	}

// 	switch len(t.Players) {
// 	case 0:
// 		t.Status = RPSTourneyStatusFinished
// 	case 1:
// 		t.Status = RPSTourneyStatusFinished
// 		t.WinningPlayerId = t.Players[0]
// 	default:
// 		t.Status = RPSTourneyStatusInProgress
// 		t.Round++
// 		// for every 2 players still in the tourney, create a new game
// 		newGames := make([]*RPSGame, 0)
// 		for i := 1; i < len(t.Players); i += 2 {
// 			newGames = append(newGames, NewRPSGame(t.Players[i-1]).Join(t.Players[i]))
// 		}
// 		t.Games = append(t.Games, newGames)
// 	}
// }

// func (t *RPSTourney) GetMostRecentRound() []*RPSGame {
// 	if len(t.Games) == 0 {
// 		return []*RPSGame{}
// 	}

// 	return t.Games[len(t.Games)-1]
// }

// func (t *RPSTourney) MakeChoice(playerId string, choice RPSChoice) {
// 	if t.Status != RPSTourneyStatusInProgress {
// 		return
// 	}

// 	for _, g := range t.GetMostRecentRound() {
// 		if g.Player1Id == playerId || g.Player2Id == playerId {
// 			g.MakeChoice(playerId, choice)
// 		}
// 	}
// }

// func (t *RPSTourney) ToJSON() ([]byte, error) {
// 	return json.Marshal(t)
// }

// func RPSTourneyFromJSON(data []byte) (*RPSTourney, error) {
// 	t := &RPSTourney{}
// 	err := json.Unmarshal(data, t)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return t, nil
// }
