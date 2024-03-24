package games

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

type RPSGame struct {
	Id              string        `db:"id" json:"id"`
	Status          RPSGameStatus `db:"status" json:"status"`
	Player1Id       string        `db:"player1_id" json:"player1_id"`
	Player2Id       string        `db:"player2_id" json:"player2_id"`
	Choice1         RPSChoice     `db:"choice1" json:"choice1"`
	Choice2         RPSChoice     `db:"choice2" json:"choice2"`
	WinningPlayerId string        `db:"winning_player_id" json:"winning_player_id"`
	LosingPlayerId  string        `db:"losing_player_id" json:"losing_player_id"`
}
