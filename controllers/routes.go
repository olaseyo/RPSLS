package controllers

import "game/middlewares"

func (s *Server) initializeRoutes() {

	// Choice Route
	s.Router.HandleFunc("/choice", middlewares.SetMiddlewareJSON(s.FindChoice)).Methods("GET")
	// Choices Route
	s.Router.HandleFunc("/choices", middlewares.SetMiddlewareJSON(s.FindAll)).Methods("GET")
	// Play Route
	s.Router.HandleFunc("/play", middlewares.SetMiddlewareJSON(s.SinglePlayer)).Methods("POST")
	
	s.Router.HandleFunc("/room/create", middlewares.SetMiddlewareJSON(s.CreateRoom)).Methods("POST")

	s.Router.HandleFunc("/multiplayer/play", middlewares.SetMiddlewareJSON(s.MultiPlayer)).Methods("POST")

	s.Router.HandleFunc("/multiplayer/result", middlewares.SetMiddlewareJSON(s.GetResult)).Methods("POST")
	s.Router.HandleFunc("/board/result", middlewares.SetMiddlewareJSON(s.GetScoreBoardData)).Methods("GET")
	s.Router.HandleFunc("/board/delete", middlewares.SetMiddlewareJSON(s.Reset)).Methods("DELETE")
	
}
