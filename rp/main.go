package main

import (
	"fmt"
)

type GameMode interface {
	Start()
}

func main() {
	for {
		ClearConsole()
		fmt.Println("==========================================")
		fmt.Println("           БОЕВАЯ СИСТЕМА")
		fmt.Println("==========================================")

		fmt.Println("\n1. Сюжетная кампания")
		fmt.Println("2. PVP")
		fmt.Println("3. Выход")

		var choice int
		fmt.Print("\nВыберите режим: ")
		fmt.Scan(&choice)

		var gameMode GameMode

		switch choice {
		case 1:
			gameMode = NewCampaign()
			gameMode.Start()
		case 2:
			ClearConsole()
			fmt.Println("Выберите формат PvP:")
			fmt.Println("1. Горячий стул (за одним устройством)")
			fmt.Println("2. Сетевая игра")
			var pvpChoice int
			fmt.Scan(&pvpChoice)

			switch pvpChoice {
			case 1:
				gameMode = NewHotSeat()
			case 2:
				gameMode = NewNetworkPVP()
			default:
				fmt.Println("Неверный выбор")
				WaitForEnter()
				continue
			}
			gameMode.Start()
		case 3:
			ClearConsole()
			fmt.Println("\nСпасибо за игру!")
			return
		default:
			fmt.Println("Неверный выбор")
			WaitForEnter()
		}
	}
}
