package controllers

import (
	"game/helpers"
	"game/models"
	"game/responses"
	"net/http"
)

func (server *Server) CreateRoom(w http.ResponseWriter, r *http.Request) {
	room := models.Room{}
	room_id, _ := helpers.New()
	room.Identifier = room_id
	roomCreated, err := room.CreateRoom(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	responses.JSON(w, http.StatusOK, roomCreated)
	return
}
