package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	"html"
	"strings"
	"time"
)

type Room struct {
	ID         uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Identifier string    `json:"identifier"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (u *Room) Prepare() {
	u.ID = 0
	u.Identifier = html.EscapeString(strings.TrimSpace(u.Identifier))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (c *Room) CreateRoom(db *gorm.DB) (*Room, error) {
	err := db.Debug().Create(&c).Error
	if err != nil {
		return &Room{}, err
	}
	return c, nil
}

func (c *Room) FindAll(db *gorm.DB) (*[]Room, error) {
	rooms := []Room{}
	err := db.Debug().Model(&Room{}).Limit(100).Find(&rooms).Error
	if err != nil {
		return &[]Room{}, err
	}
	return &rooms, err
}

func (c *Room) FindByID(db *gorm.DB, uid uint32) (*Room, error) {
	err := db.Debug().Model(Room{}).Where("id = ?", uid).Take(&c).Error
	if err != nil {
		return &Room{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Room{}, errors.New("Room Not Found")
	}
	return c, err
}

func (c *Room) GetRoomByIdentifier(db *gorm.DB, indentifier string) (*Room, error) {
	err := db.Debug().Model(Room{}).Where("identifier = ?", indentifier).Take(&c).Error
	if err != nil {
		return &Room{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Room{}, errors.New("Room Not Found")
	}
	return c, err
}
