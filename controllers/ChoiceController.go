package controllers

import (
	//"encoding/json"
	"fmt"
	"game/helpers"
	"game/models"

	//"io/ioutil"
	"game/responses"
	"net/http"
)

func (server *Server) FindAll(w http.ResponseWriter, r *http.Request) {

	choice := models.Choice{}

	choices, err := choice.FindAll(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusNoContent, err)
		return
	}
	responses.JSON(w, http.StatusOK, choices)
}

func (server *Server) FindChoice(w http.ResponseWriter, r *http.Request) {
	cid, err := helpers.MapNumber()
	fmt.Printf("Random Choice  %+v\n", cid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	choice := models.Choice{}
	choiceFound, err := choice.FindChoiceById(server.DB, uint32(cid))
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	responses.JSON(w, http.StatusOK, choiceFound)
}
