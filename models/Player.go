package models

import (
	"errors"
	"fmt"
	//"strings"
	"github.com/jinzhu/gorm"
	"time"
)

const (
	Botname = "BillupBot"
)

type Player struct {
	ID       uint32 `gorm:"primary_key;auto_increment" json:"id"`
	Name     string `gorm:"size:255;not null;" json:"name"`
	Room     Room   ` json:"room"`
	RoomId   uint32 `json:"room_id"`
	Choice   Choice ` json:"choice"`
	ChoiceId uint32 `json:"choice_id"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Playload struct {
	Player uint32 `json:"player"`
}

type MultiplayerPlayload struct {
	Player     uint32 `json:"player"`
	PlayerName string `json:"player_name"`
	RoomId     string `json:"room_id"`
}

type ResultPlayload struct {
	RoomId string `json:"room_id"`
}

func (p *ResultPlayload) ResultPlayloadChecker() error {
	if p.RoomId == "" {
		return errors.New("room_id required in payload")
	}
	return nil
}

func (p *Playload) SinglePlayerChecker() error {
	if p.Player == 0 {
		return errors.New("player required in payload")
	}
	return nil
}

func (p *MultiplayerPlayload) MultiPlayerChecker() error {
	if p.Player == 0 {
		return errors.New("player required in payload")
	}
	if p.PlayerName == "" {
		return errors.New("player_name required in payload")
	}
	if p.RoomId == "" {
		return errors.New("room_id required in payload")
	}
	return nil
}

func (p *Player) Prepare() {
	p.ID = 0
	p.Name = p.Name
	p.RoomId = p.RoomId
	p.ChoiceId = p.ChoiceId
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (c *Player) CreatePlayer(db *gorm.DB) (*Player, error) {
	err := db.Debug().Create(&c).Error
	if err != nil {
		return &Player{}, err
	}
	return c, nil
}

func (c *Player) GetPlayerChoiceByRoomId(db *gorm.DB, room_id uint32) (*Player, error) {
	err := db.Debug().Model(&Player{}).Where("room_id=?", room_id).Take(&c).Error
	if err != nil {
		return &Player{}, err
	}
	return c, err
}

func (c *Player) FindAll(db *gorm.DB) (*[]Player, error) {
	player := []Player{}
	err := db.Debug().Model(&Player{}).Preload("Room").Limit(2).Find(&player).Error
	if err != nil {
		return &[]Player{}, err
	}
	return &player, err
}

func (c *Player) IsMovesComplete(db *gorm.DB, room_id uint32) bool {
	player := []Player{}
	_ = db.Debug().Model(&Player{}).Where("room_id = ?", room_id).Limit(2).Find(&player).Error
	return len(player) == 2
}

func (c *Player) FindByID(db *gorm.DB, pid uint32) ([]Player, error) {
	queryResult := []Player{}
	err := db.Debug().Model(&Player{}).Where("room_id = ?", pid).Find(&queryResult).Error
	if err != nil {
		return []Player{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return []Player{}, errors.New("Player Not Found")
	}
	return queryResult, err
}

func (c *Player) GetOutcome(db *gorm.DB, room_id uint32, isBot bool) (map[uint32]string, error) {
	player := []Player{}
	choiceMap := ChoiceMap{}
	outCome := make(map[uint32]string)
	err := db.Debug().Model(&Player{}).Preload("Room").Where("room_id = ?", room_id).Limit(2).Scan(&player).Error
	if err != nil {
		return outCome, err
	}

	result, _ := choiceMap.FindMatch(db, player[0].ChoiceId, player[1].ChoiceId)
	fmt.Printf("outcome %+v\n", result)
	if player[0].ChoiceId == player[1].ChoiceId {
		outCome[player[0].ID] = "tie"
		outCome[player[1].ID] = "tie"
	} else if result {
		outCome[player[0].ID] = "won"
		outCome[player[1].ID] = "lose"
	} else {
		outCome[player[0].ID] = "lose"
		outCome[player[1].ID] = "won"
	}

	return outCome, err
}

func (s *Player) Delete(db *gorm.DB) {
	db.Debug().Exec("DELETE FROM players")
}
