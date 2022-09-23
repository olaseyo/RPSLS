package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

type Choice struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `json:"choice"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type Choices struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

type RandomChoice struct {
	ID   uint32 `json:"id"`
	Name string `json:"name"`
}

func (c *Choice) Prepare() {
	c.ID = 0
	c.Name = c.Name
}

func (c *Choice) Create(db *gorm.DB) (bool, error) {
	err := db.Debug().Create(&c).Error
	if err != nil {
		return true, err
	}
	return true, nil
}

func (c *Choice) FindChoiceById(db *gorm.DB, id uint32) (RandomChoice, error) {
	choice := RandomChoice{}
	err := db.Model(&Choice{}).Where("choices.id=?", id).Scan(&choice).Error

	if err != nil {
		return RandomChoice{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return RandomChoice{}, errors.New("record not found")
	}
	return choice, err
}

func (c *Choice) FindAll(db *gorm.DB) (*[]Choices, error) {
	var err error
	choices := []Choices{}
	err = db.Debug().Model(&Choice{}).Limit(10).Scan(&choices).Error
	if err != nil {
		return &[]Choices{}, err
	}
	return &choices, err
}
