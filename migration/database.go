package migration

import (
	"game/migration/seed"
	"game/models"
	"github.com/jinzhu/gorm"
	"log"
)

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.ScoreBoard{}, &models.ChoiceMap{}, &models.Player{}, &models.Choice{}, &models.Room{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.Choice{}, &models.Room{}, &models.ScoreBoard{}, &models.ChoiceMap{}, &models.Player{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	db.Debug().Model(&models.Choice{})
	db.Debug().Model(&models.Room{})
	db.Debug().Model(&models.ChoiceMap{}).AddForeignKey("choice_id", "choices(id)", "cascade", "cascade")
	db.Debug().Model(&models.ScoreBoard{}).AddForeignKey("room_id", "rooms(id)", "cascade", "cascade")
	db.Debug().Model(&models.ScoreBoard{}).AddForeignKey("player_id", "players(id)", "cascade", "cascade")
	db.Debug().Model(&models.Player{}).AddForeignKey("room_id", "rooms(id)", "cascade", "cascade")
	db.Debug().Model(&models.Player{}).AddForeignKey("choice_id", "choices(id)", "cascade", "cascade")

	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}
}

func Seed(db *gorm.DB) {

	var choices = []models.Choice{
		{
			Name: "rock",
		}, {
			Name: "paper",
		}, {
			Name: "scissors",
		}, {
			Name: "lizard",
		}, {
			Name: "spock",
		},
	}

	for i, _ := range choices {
		err := db.Debug().Model(&models.Choice{}).Create(&choices[i]).Error
		seed.SeedChoiceMap(db, i)
		if err != nil {
			log.Fatalf("cannot seed choice table: %v", err)
		}

	}
}
