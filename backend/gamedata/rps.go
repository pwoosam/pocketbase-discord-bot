package gamedata

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type RPSGameStatus int

const (
	RPSGameStatusWaitingForPlayers RPSGameStatus = iota
	RPSGameStatusInProgress
	RPSGameStatusFinished
)

type RPSChoice int

const (
	Undecided RPSChoice = iota
	Rock
	Paper
	Scissors
)

func RPSGetGameByInteractionId(dao *daos.Dao, interactionId string) (*models.Record, error) {
	gameId, err := rpsGetGameIdFromInteractionId(dao, interactionId)
	if err != nil {
		return nil, err
	}

	game, err := rpsGetGameByGameId(dao, gameId)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func rpsGetGameByGameId(dao *daos.Dao, gameId string) (*models.Record, error) {
	game, err := dao.FindRecordById(rpsCollection.Name, gameId)
	if err != nil {
		return nil, err
	}

	return game, nil
}

func rpsGetGameIdFromInteractionId(dao *daos.Dao, interactionId string) (string, error) {
	interaction, err := dao.FindRecordById(rpsInteractionCollection.Name, interactionId)
	if err != nil {
		return "", err
	}

	return interaction.GetString("game_id"), nil
}

func RPSUpsertGameInteraction(dao *daos.Dao, interactionId string, gameId string) error {
	if _, err := dao.FindRecordById(rpsInteractionCollection.Name, interactionId); err != nil {
		log.Println("rps interaction not found, creating new one")
		interactionRecord := models.NewRecord(rpsInteractionCollection)
		interactionRecord.Load(map[string]any{
			"id":      interactionId,
			"game_id": gameId,
		})

		if err := dao.SaveRecord(interactionRecord); err != nil {
			return err
		}
	} else {
		log.Println("rps interaction already exists, nothing to do here")
	}

	return nil
}

func RPSCreateGame(dao *daos.Dao, userID string) (*models.Record, error) {
	gameRecord := models.NewRecord(rpsCollection)
	gameRecord.Load(map[string]any{
		"id":         uuid.NewString(),
		"player1_id": userID,
		"status":     RPSGameStatusWaitingForPlayers,
	})

	if err := dao.SaveRecord(gameRecord); err != nil {
		return nil, err
	}

	return gameRecord, nil
}

func RPSJoinGame(dao *daos.Dao, gameId string, userId string) error {
	game, err := rpsGetGameByGameId(dao, gameId)
	if err != nil {
		return err
	}

	if game.GetString("player1_id") == userId {
		return errors.New("you are player 1, you cannot join your own game")
	}

	if game.GetString("player2_id") != "" {
		return errors.New("game already has two players")
	}

	game.Set("player2_id", userId)
	game.Set("status", int(RPSGameStatusInProgress))

	if err := dao.SaveRecord(game); err != nil {
		return err
	}

	return nil
}

func RPSMakeChoice(dao *daos.Dao, gameId string, userId string, choice RPSChoice) error {
	game, err := rpsGetGameByGameId(dao, gameId)
	if err != nil {
		return err
	}

	if RPSGameStatus(game.GetInt("status")) != RPSGameStatusInProgress {
		return errors.New("game is not in progress")
	}

	if game.GetString("player1_id") == userId {
		game.Set("player1_choice", int(choice))
	} else if game.GetString("player2_id") == userId {
		game.Set("player2_choice", int(choice))
	} else {
		return errors.New("user is not a player in this game")
	}

	player1Choice := RPSChoice(game.GetInt("player1_choice"))
	player2Choice := RPSChoice(game.GetInt("player2_choice"))
	if RPSChoice(player1Choice) != Undecided && RPSChoice(player2Choice) != Undecided {
		winner, loser := rpsGetWinnerAndLoser(player1Choice, player2Choice, game.GetString("player1_id"), game.GetString("player2_id"))
		game.Set("player_id_winner", winner)
		game.Set("player_id_loser", loser)
		game.Set("status", int(RPSGameStatusFinished))
	}

	if err := dao.SaveRecord(game); err != nil {
		return err
	}

	return nil
}

func rpsGetWinnerAndLoser(player1Choice, player2Choice RPSChoice, player1Id, player2Id string) (string, string) {
	if player1Choice == player2Choice {
		return "", ""
	}

	if player1Choice == Rock && player2Choice == Scissors {
		return player1Id, player2Id
	}

	if player1Choice == Paper && player2Choice == Rock {
		return player1Id, player2Id
	}

	if player1Choice == Scissors && player2Choice == Paper {
		return player1Id, player2Id
	}

	return player2Id, player1Id
}
