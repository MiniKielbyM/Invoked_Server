package main
import ( "net" )
var lobbies = make(map[string]*lobby)

var players = make(map[net.Conn]player)