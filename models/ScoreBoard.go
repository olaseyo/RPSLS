package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type ScoreBoard struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Player    Player    `json:"player"`
	PlayerId  uint32    `json:"player_id"`
	RoomId    uint32    `json:"room_id"`
	Room      Room      `json:"room"`
	Status    string    `json:"status"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type SingleResultQuery struct {
	PlayerId uint32 `json:"player_id"`
	Status   string `json:"status"`
}

type SingleResult struct {
	Results  string `json:"results"`
	Player   uint32 `json:"player"`
	Computer uint32 `json:"computer"`
}

type PlayerResultObject struct {
	Name     string `json:"name"`
	RoomId   uint32 `json:"room_id"`
	ChoiceId uint32 `json:"choice_id"`
}
type MultipleResult struct {
	Winner        string               `json:"winner"`
	WinningChoice int                  `json:"winning_choice"`
	Players       []PlayerResultObject `json:"players"`
}

func (s *ScoreBoard) Prepare() {
	s.ID = 0
	s.PlayerId = 0
	s.RoomId = 0
	s.Status = ""
}

func (u *ScoreBoard) Save(db *gorm.DB) (*ScoreBoard, error) {
	err := db.Debug().Create(&u).Error
	if err != nil {
		return &ScoreBoard{}, err
	}
	return u, nil
}

func (c *Room) GetWinnerByRoomId(db *gorm.DB, room_id uint32) (*Room, error) {
	err := db.Debug().Model(ScoreBoard{}).Where("room_id = ?", room_id).Take(&c).Error
	if err != nil {
		return &Room{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Room{}, errors.New("Room Not Found")
	}
	return c, err
}

func (c *ScoreBoard) IsScoreRecorded(db *gorm.DB, room_id uint32) bool {
	board := []ScoreBoard{}
	_ = db.Debug().Model(ScoreBoard{}).Where("room_id = ?", room_id).Limit(2).Take(&board).Error

	return len(board) == 2
}

func (c *ScoreBoard) GetSingleResult(db *gorm.DB, room_id uint32, computer_choice uint32) (*SingleResult, error) {
	queryResult := SingleResultQuery{}
	err := db.Debug().Model(&ScoreBoard{}).
		Joins("join players on players.id = score_boards.player_id").
		Where("score_boards.room_id = ?", room_id).Where("name != ?", Botname).
		Scan(&queryResult).Error
	if err != nil {
		return &SingleResult{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &SingleResult{}, errors.New("Room Not Found")
	}
	singleResult := SingleResult{}
	player := Player{}
	getPlayerData, _ := player.GetPlayerChoiceByRoomId(db, room_id)
	singleResult.Results = queryResult.Status
	singleResult.Player = getPlayerData.ChoiceId
	singleResult.Computer = computer_choice
	return &singleResult, err
}

func (c *ScoreBoard) GetMultipleResult(db *gorm.DB, room_id uint32) ([]MultipleResult, error) {
	queryResults := []ScoreBoard{}
	var boardDataMapping []MultipleResult
	err := db.Debug().Preload("Player").Preload("Room").
		Where("score_boards.room_id = ?", room_id).
		Find(&queryResults).Error
	if err != nil {
		return []MultipleResult{}, err
	}
	//fmt.Printf("result %+v\n", queryResults)
	if gorm.IsRecordNotFoundError(err) {
		return []MultipleResult{}, errors.New("record not found")
	}
	player := make([]PlayerResultObject, 0)
	var i = 0
	var winner = ""
	var winningChoice = 0
	for _, score := range queryResults {
		fmt.Printf("result %+v\n", queryResults)
		player = append(player, PlayerResultObject{Name: score.Player.Name, RoomId: score.Player.RoomId, ChoiceId: score.Player.ChoiceId})

		if score.Status == "won" {
			winner = score.Player.Name
			winningChoice = int(score.Player.ChoiceId)
		}
		if score.Status == "tie" {
			winner = "tie"
			winningChoice = 0
		}
		i++
	}
	boardDataMapping = append(boardDataMapping, MultipleResult{Players: player, Winner: winner, WinningChoice: winningChoice})
	player = []PlayerResultObject{}
	return boardDataMapping, err
}

func (c *ScoreBoard) GetScoreBoardData(db *gorm.DB) ([]MultipleResult, error) {
	queryResults := []ScoreBoard{}
	distinctRoom := []ScoreBoard{}
	err := db.Debug().Model(&queryResults).Select("room_id").Order("room_id desc").Group("room_id").Limit(10).Find(&distinctRoom).Error
	fmt.Printf("distinctRoom  %+v\n", (distinctRoom))
	player := make([]PlayerResultObject, 0)
	var boardDataMapping []MultipleResult
	var i = 0
	var winner = ""
	var winningChoice = 0

	for _, room := range distinctRoom {
		db.Debug().Preload("Player").Preload("Room").
			Order("score_boards.id desc").
			Where("room_id=?", room.RoomId).
			Find(&queryResults)
		j := 0
		for _, score := range queryResults {
			roomId := room.RoomId
			playerName := score.Player.Name
			player = append(player, PlayerResultObject{Name: playerName, RoomId: roomId, ChoiceId: score.Player.ChoiceId})
			if score.Status == "won" {
				winner = score.Player.Name
				winningChoice = int(score.Player.ChoiceId)
			}
			if score.Status == "tie" {
				winner = "tie"
				winningChoice = 0
			}
			j++

		}
		boardDataMapping = append(boardDataMapping, MultipleResult{Players: player, Winner: winner, WinningChoice: winningChoice})
		i++
		player = []PlayerResultObject{}
	}
	return boardDataMapping, err
}

func (s *ScoreBoard) Delete(db *gorm.DB) {
	db.Debug().Exec("DELETE FROM score_boards")
}
