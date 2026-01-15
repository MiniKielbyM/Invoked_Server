package main

// Usage
var lobbyHandler = func() *Handler {
	r := NewHandler()
	r.Register("create", handleLobbyCreate)
	r.Register("join", handleLobbyJoin)
	r.Register("list", handleLobbyList)
	r.Register("start", handleGameStart)
	return r
}()

var gameHandler = func() *Handler {
	r := NewHandler()
	r.RegisterSubHandler("lobby", lobbyHandler)
	return r
}()

var outGameDeckHandler = func() *Handler {
	r := NewHandler()
	r.Register("update", handlePlayerUpdateDeck)
	return r
}()

var playerHandler = func() *Handler {
	r := NewHandler()
	r.Register("connect", handlePlayerConnect)
	r.RegisterSubHandler("deck", outGameDeckHandler)
	return r
}()

var mainHandler = func() *Handler {
	r := NewHandler()
	r.RegisterSubHandler("game", gameHandler)
	r.RegisterSubHandler("player", playerHandler)
	return r
}()
