package seed

import (
	"game/models"
	"github.com/jinzhu/gorm"
)

const (
	Rock     = 0
	Paper    = 1
	Scissors = 2
	Lizard   = 3
	Spock    = 4
)

var ChoiceMap = map[int][]int{
	Rock:     {Scissors, Lizard},
	Paper:    {Rock, Spock},
	Scissors: {Paper, Lizard},
	Lizard:   {Paper, Spock},
	Spock:    {Scissors, Rock},
}

func SeedChoiceMap(db *gorm.DB, index int) {
	for _, element := range ChoiceMap[index] {
		choiceMapping := &models.ChoiceMap{}
		choiceMapping.ChoiceId = uint32(index + 1)
		choiceMapping.Slave = uint32(element + 1)
		_ = db.Debug().Model(choiceMapping).Create(choiceMapping).Error
	}

}
