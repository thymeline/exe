package main

import (
	"fmt"
)

type PVP struct{}

func NewHotSeat() *PVP {
	return &PVP{}
}

func (h *PVP) Start() {
	fmt.Println("\n====== (PvP) ======")

	var player1Name, player2Name string
	fmt.Print("\nИмя первого игрока: ")
	fmt.Scan(&player1Name)
	fmt.Print("Имя второго игрока: ")
	fmt.Scan(&player2Name)

	player1 := &Fighter{
		Name:     player1Name,
		HP:       100,
		MaxHP:    100,
		Strength: 10,
		BaseStr:  10,
	}

	player2 := &Fighter{
		Name:     player2Name,
		HP:       100,
		MaxHP:    100,
		Strength: 10,
		BaseStr:  10,
	}

	RunBattle(player1, player2, false)
}

type NetworkPVP struct{}

func NewNetworkPVP() *NetworkPVP {
	return &NetworkPVP{}
}

func (n *NetworkPVP) Start() {
	ClearConsole()
	fmt.Println("\n=== СЕТЕВАЯ ИГРА PvP ===")
	fmt.Println("1. Создать игру (сервер)")
	fmt.Println("2. Подключиться к игре (клиент)")

	var choice int
	fmt.Print("Выберите действие: ")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		StartServer()
	case 2:
		StartClient()
	default:
		fmt.Println("Неверный выбор")
		WaitForEnter()
	}
}
