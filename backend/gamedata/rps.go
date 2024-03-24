package gamedata

import (
	"errors"
	"log"
	"myapp/games"

	"github.com/google/uuid"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
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
		log.Println("rps interaction already existins, nothing to do here")
	}

	return nil
}

func RPSCreateGame(dao *daos.Dao, userID string) (*models.Record, error) {
	gameRecord := models.NewRecord(rpsCollection)
	gameRecord.Load(map[string]any{
		"id":         uuid.NewString(),
		"player1_id": userID,
		"status":     games.RPSGameStatusWaitingForPlayers,
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
	game.Set("status", int(games.RPSGameStatusInProgress))

	if err := dao.SaveRecord(game); err != nil {
		return err
	}

	return nil
}

func RPSMakeChoice(dao *daos.Dao, gameId string, userId string, choice games.RPSChoice) error {
	game, err := rpsGetGameByGameId(dao, gameId)
	if err != nil {
		return err
	}

	if games.RPSGameStatus(game.GetInt("status")) != games.RPSGameStatusInProgress {
		return errors.New("game is not in progress")
	}

	if game.GetString("player1_id") == userId {
		game.Set("player1_choice", int(choice))
	} else if game.GetString("player2_id") == userId {
		game.Set("player2_choice", int(choice))
	} else {
		return errors.New("user is not a player in this game")
	}

	player1Choice := games.RPSChoice(game.GetInt("player1_choice"))
	player2Choice := games.RPSChoice(game.GetInt("player2_choice"))
	if games.RPSChoice(player1Choice) != games.Undecided && games.RPSChoice(player2Choice) != games.Undecided {
		winner, loser := rpsGetWinnerAndLoser(player1Choice, player2Choice, game.GetString("player1_id"), game.GetString("player2_id"))
		game.Set("player_id_winner", winner)
		game.Set("player_id_loser", loser)
		game.Set("status", int(games.RPSGameStatusFinished))
	}

	if err := dao.SaveRecord(game); err != nil {
		return err
	}

	return nil
}

func rpsGetWinnerAndLoser(player1Choice, player2Choice games.RPSChoice, player1Id, player2Id string) (string, string) {
	if player1Choice == player2Choice {
		return "", ""
	}

	if player1Choice == games.Rock && player2Choice == games.Scissors {
		return player1Id, player2Id
	}

	if player1Choice == games.Paper && player2Choice == games.Rock {
		return player1Id, player2Id
	}

	if player1Choice == games.Scissors && player2Choice == games.Paper {
		return player1Id, player2Id
	}

	return player2Id, player1Id
}
