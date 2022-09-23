package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"game/helpers"
	"game/models"
	"game/responses"
	"io/ioutil"
	"net/http"
)

func SimulateBillupBotChoice(server *Server) (int, error) {
	cid, err := helpers.MapNumber()
	return cid, err

}

func CreateRoom(server *Server, room models.Room) *models.Room {
	room_id, _ := helpers.New()
	room.Identifier = room_id
	roomCreated, _ := room.CreateRoom(server.DB)
	return roomCreated
}
func GeneratePlayerName(name string) string {
	if name == "" {
		name = "Player 1" //autoname.Generate();
	}
	return name
}
func RecordPlayerMove(server *Server, player models.Player, room_id uint32, choice_id uint32, name string) *models.Player {
	player.RoomId = room_id
	player.Name = GeneratePlayerName(name)
	player.ChoiceId = choice_id
	playerCreated, _ := player.CreatePlayer(server.DB)
	return playerCreated
}

func RecordScores(server *Server, results map[uint32]string, room_id uint32) {

	for key, result := range results {
		scoreBoardModel := models.ScoreBoard{}
		scoreBoardModel.PlayerId = key
		scoreBoardModel.RoomId = room_id
		scoreBoardModel.Status = result
		scoreBoardModel.Save(server.DB)
	}
}
func (server *Server) SinglePlayer(w http.ResponseWriter, r *http.Request) {
	room := models.Room{}
	//choiceMap := models.ChoiceMap{}
	playerModel := models.Player{}
	playerChoice := models.Playload{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = json.Unmarshal(body, &playerChoice)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = playerChoice.SinglePlayerChecker()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	roomCreated := CreateRoom(server, room)
	botChoice, err := SimulateBillupBotChoice(server)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, errors.New("error while making api call. Kindly your network"))
		return
	}
	RecordPlayerMove(server, playerModel, roomCreated.ID, playerChoice.Player, "")
	RecordPlayerMove(server, playerModel, roomCreated.ID, uint32(botChoice), "BillupBot")
	fmt.Printf("roomId %+v\n", roomCreated.ID)
	fmt.Printf("botChoice %+v\n", botChoice)
	fmt.Printf("playerChoice %+v\n", playerChoice.Player)
	results, _ := playerModel.GetOutcome(server.DB, roomCreated.ID, true)
	fmt.Printf("result %+v\n", results)
	RecordScores(server, results, roomCreated.ID)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	scoreBoardModel := models.ScoreBoard{}
	gameResult, err := scoreBoardModel.GetSingleResult(server.DB, roomCreated.ID, uint32(botChoice))

	if err != nil {
		responses.ERROR(w, http.StatusNoContent, err)
		return
	}
	responses.JSON(w, http.StatusOK, gameResult)

}

func (server *Server) MultiPlayer(w http.ResponseWriter, r *http.Request) {
	room := models.Room{}
	playerModel := models.Player{}
	//scoreModel := models.ScoreBoard{}
	playerChoice := models.MultiplayerPlayload{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	err = json.Unmarshal(body, &playerChoice)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = playerChoice.MultiPlayerChecker()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	roomCreated, err := room.GetRoomByIdentifier(server.DB, playerChoice.RoomId)

	if err != nil {
		responses.ERROR(w, http.StatusNoContent, err)
		return
	}

	if playerModel.IsMovesComplete(server.DB, roomCreated.ID) {
		responses.ERROR(w, http.StatusConflict, errors.New("player moves completed"))
		return
	}

	RecordPlayerMove(server, playerModel, roomCreated.ID, playerChoice.Player, playerChoice.PlayerName)
	fmt.Printf("roomId %+v\n", roomCreated.ID)
	fmt.Printf("playerChoice %+v\n", playerChoice.Player)

	if err != nil {
		responses.ERROR(w, http.StatusNoContent, err)
		return
	}
	fmt.Printf("Score recording %+v\n", playerModel.IsMovesComplete(server.DB, roomCreated.ID))
	if playerModel.IsMovesComplete(server.DB, roomCreated.ID) {

		results, _ := playerModel.GetOutcome(server.DB, roomCreated.ID, false)

		fmt.Printf("results %+v\n", results)
		RecordScores(server, results, roomCreated.ID)
	}

	responses.JSON(w, http.StatusOK, playerChoice)
	return

}

func (server *Server) GetResult(w http.ResponseWriter, r *http.Request) {
	room := models.Room{}
	playerChoice := models.ResultPlayload{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusNoContent, errors.New("body error"))
		return
	}
	err = json.Unmarshal(body, &playerChoice)
	if err != nil {
		responses.ERROR(w, http.StatusNoContent, errors.New("could not process"))
		return
	}

	err = playerChoice.ResultPlayloadChecker()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("Validation error"))
		return
	}
	roomCreated, err := room.GetRoomByIdentifier(server.DB, playerChoice.RoomId)
	if err != nil {
		responses.ERROR(w, http.StatusNoContent, err)
		return
	}
	scoreBoardModel := models.ScoreBoard{}
	gameResult, err := scoreBoardModel.GetMultipleResult(server.DB, roomCreated.ID)
	if err != nil {
		responses.ERROR(w, http.StatusNoContent, err)
		return
	}
	responses.JSON(w, http.StatusOK, gameResult)
	return
}

func (server *Server) GetScoreBoardData(w http.ResponseWriter, r *http.Request) {
	scoreBoardModel := models.ScoreBoard{}
	gameResult, err := scoreBoardModel.GetScoreBoardData(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusNoContent, err)
		return
	}
	responses.JSON(w, http.StatusOK, gameResult)
	return
}

func (server *Server) Reset(w http.ResponseWriter, r *http.Request) {
	player := models.Player{}
	scoreBoardModel := models.ScoreBoard{}

	player.Delete(server.DB)
	scoreBoardModel.Delete(server.DB)
	responses.JSON(w, http.StatusOK, "Board reset successfully")
	return
}
