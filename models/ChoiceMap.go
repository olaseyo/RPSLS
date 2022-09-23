package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

type ChoiceMap struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	ChoiceId  uint32    `json:"choice_id"`
	Choice    Choice    `json:"choice"`
	Slave     uint32    `json:"slave"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *ChoiceMap) FindMatch(db *gorm.DB, player1_choice_id uint32, player2_choice_id uint32) (bool, error) {
	choices := []ChoiceMap{}
	err := db.Debug().Model(&ChoiceMap{}).Where("choice_maps.choice_id=?", player1_choice_id).Limit(2).Scan(&choices).Error

	for _, element := range choices {
		fmt.Printf("player2Choice %+v\n", element.Slave)
		if element.Slave == player2_choice_id {
			//player 1 on won
			return true, err
		}
	}
	//player 2 won
	return false, err
}
