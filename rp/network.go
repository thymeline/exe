package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	port = ":8080"
)

// ChatHistory хранит последние N сообщений чата
type ChatHistory struct {
	messages []string
	capacity int
}

func NewChatHistory(capacity int) *ChatHistory {
	return &ChatHistory{
		messages: make([]string, 0, capacity),
		capacity: capacity,
	}
}

func (ch *ChatHistory) Add(msg string) {
	if len(ch.messages) >= ch.capacity {
		// сдвигаем влево
		ch.messages = append(ch.messages[1:], msg)
	} else {
		ch.messages = append(ch.messages, msg)
	}
}

func (ch *ChatHistory) Display() {
	if len(ch.messages) == 0 {
		fmt.Println("   (чат пуст)")
		return
	}
	for _, msg := range ch.messages {
		fmt.Printf("  %s\n", msg)
	}
}

// displayBattleScreen очищает экран и показывает состояние боя + историю чата
func displayBattleScreen(localName string, localHP, localMaxHP int,
	remoteName string, remoteHP, remoteMaxHP int,
	round int, chat *ChatHistory) {

	ClearConsole()
	fmt.Printf("=== Бой: %s vs %s ===\n", localName, remoteName)
	fmt.Printf("Раунд: %d\n", round)
	fmt.Printf("%s: %d/%d HP | %s: %d/%d HP\n",
		localName, localHP, localMaxHP,
		remoteName, remoteHP, remoteMaxHP)
	fmt.Println("\n--- Чат ---")
	chat.Display()
	fmt.Println("------------")
}

// ------------------------------------------------------------------
// SERVER
// ------------------------------------------------------------------

func StartServer() {
	ClearConsole()
	fmt.Println("=== СЕРВЕР ===")

	// Имя локального игрока
	var localName string
	fmt.Print("Ваше имя: ")
	fmt.Scan(&localName)

	// Создание бойца
	localPlayer := &Fighter{
		Name:     localName,
		HP:       100,
		MaxHP:    100,
		Strength: 10,
		BaseStr:  10,
	}
	localPlayer.Inventory = CreateTestItems()

	// Запуск TCP-сервера
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		WaitForEnter()
		return
	}
	defer listener.Close()

	fmt.Println("Сервер запущен. Ожидание подключения противника...")
	fmt.Println("Адрес для подключения: 127.0.0.1" + port)
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Ошибка принятия соединения:", err)
		WaitForEnter()
		return
	}
	defer conn.Close()
	fmt.Println("Противник подключился!")

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Обмен именами
	// 1. Получаем имя клиента
	nameLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения имени клиента")
		WaitForEnter()
		return
	}
	remoteName := strings.TrimSpace(strings.TrimPrefix(nameLine, "NAME:"))

	// 2. Отправляем своё имя
	writer.WriteString("NAME:" + localName + "\n")
	writer.Flush()

	// Создаём удалённого бойца (противника)
	remotePlayer := &Fighter{
		Name:     remoteName,
		HP:       100,
		MaxHP:    100,
		Strength: 10,
		BaseStr:  10,
	}

	// История чата
	chat := NewChatHistory(5)

	// Каналы для синхронизации
	actionChan := make(chan [2]int, 1) // действия удалённого игрока
	chatChan := make(chan string, 10)  // сообщения чата от клиента
	continueChan := make(chan bool, 1) // сигнал о нажатии Enter после раунда
	done := make(chan bool)

	// Горутина чтения от клиента
	go func() {
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("\nСоединение с противником разорвано.")
				done <- true
				return
			}
			msg = strings.TrimSpace(msg)

			switch {
			case strings.HasPrefix(msg, "CHAT:"):
				text := strings.TrimPrefix(msg, "CHAT:")
				chatChan <- text
			case strings.HasPrefix(msg, "ACTION:"):
				var a, b int
				fmt.Sscanf(strings.TrimPrefix(msg, "ACTION:"), "%d %d", &a, &b)
				actionChan <- [2]int{a, b}
			case msg == "CONTINUE":
				continueChan <- true
			}
		}
	}()

	// Горутина обработки чата
	go func() {
		for text := range chatChan {
			chat.Add("[" + remoteName + "] " + text)
			// Не выводим сразу, так как экран будет обновлён
		}
	}()

	round := 1
	for !localPlayer.IsDead() && !remotePlayer.IsDead() {
		// Отображаем экран перед ходом
		displayBattleScreen(localPlayer.Name, localPlayer.HP, localPlayer.MaxHP,
			remotePlayer.Name, remotePlayer.HP, remotePlayer.MaxHP,
			round, chat)

		// Ход локального игрока (с поддержкой чата)
		attackPart, blockPart := getNetworkPlayerChoices(localPlayer.Name, writer, chat)

		// Отправляем запрос хода клиенту
		writer.WriteString("TURN\n")
		writer.Flush()

		// Показываем сообщение о ожидании
		fmt.Println("\n⏳ Ждём хода противника...")

		// Ждём действие от клиента
		var remoteAction [2]int
		select {
		case remoteAction = <-actionChan:
		case <-done:
			fmt.Println("Игра прервана.")
			WaitForEnter()
			return
		case <-time.After(120 * time.Second):
			fmt.Println("Противник не отвечает. Игра завершена.")
			WaitForEnter()
			return
		}

		// Преобразуем числа в BodyPart
		remoteAttack := intToBodyPart(remoteAction[0])
		remoteBlock := intToBodyPart(remoteAction[1])

		// Применяем действия
		localPlayer.SetBlocking(blockPart)
		remotePlayer.SetBlocking(remoteBlock)

		// Выполняем атаки (сначала локальный, потом удалённый)
		PerformAttack(localPlayer, remotePlayer, attackPart)
		if remotePlayer.IsDead() {
			writer.WriteString("GAME_OVER:you_lose\n")
			writer.Flush()
			break
		}

		PerformAttack(remotePlayer, localPlayer, remoteAttack)
		if localPlayer.IsDead() {
			writer.WriteString("GAME_OVER:you_win\n")
			writer.Flush()
			break
		}

		// Отправляем состояние клиенту
		writer.WriteString(fmt.Sprintf("STATE:%d:%d\n", localPlayer.HP, remotePlayer.HP))
		writer.Flush()

		// Отображаем результаты раунда
		displayBattleScreen(localPlayer.Name, localPlayer.HP, localPlayer.MaxHP,
			remotePlayer.Name, remotePlayer.HP, remotePlayer.MaxHP,
			round, chat)
		fmt.Println(PromptRoundOver)
		fmt.Print(PromptEnterToContinue)
		// Ожидаем ввод локального игрока (Enter)
		bufio.NewReader(os.Stdin).ReadBytes('\n')

		// Сообщаем клиенту, что мы готовы к следующему раунду
		writer.WriteString("ROUND_OVER\n")
		writer.Flush()

		// Показываем сообщение о ожидании готовности противника
		fmt.Println("\n⏳ Ждём, пока противник нажмёт Enter...")

		// Ждём подтверждения от клиента
		select {
		case <-continueChan:
			// клиент нажал Enter
		case <-done:
			fmt.Println("Игра прервана.")
			WaitForEnter()
			return
		case <-time.After(120 * time.Second):
			fmt.Println("Противник не отвечает. Игра завершена.")
			WaitForEnter()
			return
		}

		round++
	}

	// Определяем победителя
	ClearConsole()
	fmt.Println("\n=== Бой окончен ===")
	if localPlayer.IsDead() {
		fmt.Printf("Победитель: %s\n", remotePlayer.Name)
	} else {
		fmt.Printf("Победитель: %s\n", localPlayer.Name)
	}
	WaitForEnter()
}

// ------------------------------------------------------------------
// CLIENT
// ------------------------------------------------------------------

func StartClient() {
	ClearConsole()
	fmt.Println("=== КЛИЕНТ ===")

	// Имя локального игрока
	var localName string
	fmt.Print("Ваше имя: ")
	fmt.Scan(&localName)

	// Автоматически пробуем подключиться к localhost
	defaultAddr := "127.0.0.1" + port
	fmt.Printf("Пытаюсь подключиться к %s...\n", defaultAddr)

	conn, err := net.Dial("tcp", defaultAddr)
	if err != nil {
		fmt.Println("Не удалось подключиться к localhost.")
		fmt.Print("Введите адрес сервера (например, 192.168.1.100:8080): ")
		var addr string
		fmt.Scan(&addr)
		if addr == "" {
			fmt.Println("Адрес не введён. Выход.")
			WaitForEnter()
			return
		}

		conn, err = net.Dial("tcp", addr)
		if err != nil {
			fmt.Println("Ошибка подключения:", err)
			WaitForEnter()
			return
		}
	} else {
		fmt.Println("Подключено к localhost!")
	}

	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Отправляем своё имя
	writer.WriteString("NAME:" + localName + "\n")
	writer.Flush()

	// Получаем имя сервера
	nameLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка получения имени сервера")
		WaitForEnter()
		return
	}
	remoteName := strings.TrimSpace(strings.TrimPrefix(nameLine, "NAME:"))

	// Создаём локального бойца
	localPlayer := &Fighter{
		Name:     localName,
		HP:       100,
		MaxHP:    100,
		Strength: 10,
		BaseStr:  10,
	}
	localPlayer.Inventory = CreateTestItems()

	// История чата
	chat := NewChatHistory(5)

	// Каналы
	turnChan := make(chan bool, 1)
	stateChan := make(chan [2]int, 1) // [localHP, remoteHP]
	gameOverChan := make(chan string, 1)
	roundOverChan := make(chan bool, 1)
	chatChan := make(chan string, 5)
	done := make(chan bool)

	// Горутина чтения от сервера
	go func() {
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("\nСоединение с сервером разорвано.")
				done <- true
				return
			}
			msg = strings.TrimSpace(msg)

			switch {
			case strings.HasPrefix(msg, "CHAT:"):
				text := strings.TrimPrefix(msg, "CHAT:")
				chatChan <- text
			case msg == "TURN":
				turnChan <- true
			case strings.HasPrefix(msg, "STATE:"):
				var lhp, rhp int
				fmt.Sscanf(strings.TrimPrefix(msg, "STATE:"), "%d:%d", &lhp, &rhp)
				stateChan <- [2]int{lhp, rhp}
			case strings.HasPrefix(msg, "GAME_OVER:"):
				res := strings.TrimPrefix(msg, "GAME_OVER:")
				gameOverChan <- res
			case msg == "ROUND_OVER":
				roundOverChan <- true
			}
		}
	}()

	// Горутина обработки чата
	go func() {
		for text := range chatChan {
			chat.Add("[" + remoteName + "] " + text)
		}
	}()

	// Основной цикл
	round := 1
	for {
		select {
		case <-turnChan:
			// Наш ход - показываем состояние и запрашиваем действие
			displayBattleScreen(localPlayer.Name, localPlayer.HP, localPlayer.MaxHP,
				remoteName, 0, 100,
				round, chat)

			attackPart, blockPart := getNetworkPlayerChoices(localPlayer.Name, writer, chat)

			// Отправляем действие
			writer.WriteString(fmt.Sprintf("ACTION:%d %d\n", bodyPartToInt(attackPart), bodyPartToInt(blockPart)))
			writer.Flush()

		case state := <-stateChan:
			// Обновляем HP
			localPlayer.HP = state[0]
			remoteHP := state[1]
			displayBattleScreen(localPlayer.Name, localPlayer.HP, localPlayer.MaxHP,
				remoteName, remoteHP, 100,
				round, chat)

			// Показываем сообщение о ожидании следующего хода
			fmt.Println("\n⏳ Ждём хода противника...")

		case <-roundOverChan:
			fmt.Println(PromptRoundOver)
			fmt.Print(PromptEnterToContinue)
			bufio.NewReader(os.Stdin).ReadBytes('\n')
			writer.WriteString("CONTINUE\n")
			writer.Flush()
			round++

		case result := <-gameOverChan:
			ClearConsole()
			if result == "you_win" {
				fmt.Println("\n*** ВЫ ПОБЕДИЛИ! ***")
			} else {
				fmt.Println("\n*** ВЫ ПРОИГРАЛИ! ***")
			}
			WaitForEnter()
			return

		case <-done:
			fmt.Println("\nИгра завершена.")
			WaitForEnter()
			return
		}
	}
}

// ------------------------------------------------------------------
// Вспомогательные функции
// ------------------------------------------------------------------

// getNetworkPlayerChoices запрашивает у игрока атаку и блок, позволяя отправить сообщение чата.
// При получении сообщения чата оно отправляется на сервер и добавляется в локальную историю.
func getNetworkPlayerChoices(playerName string, writer *bufio.Writer, chat *ChatHistory) (BodyPart, BodyPart) {
	reader := bufio.NewReader(os.Stdin)

	// Показываем приглашение один раз перед циклом
	fmt.Printf(PromptChooseAction, playerName)
	fmt.Println(PromptAttackBlock)
	fmt.Println(PromptChatMessage)

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Проверяем, является ли ввод двумя числами
		var a, b int
		n, err := fmt.Sscanf(input, "%d %d", &a, &b)
		if err == nil && n == 2 && a >= 1 && a <= 3 && b >= 1 && b <= 3 {
			// Это действие
			return intToBodyPart(a), intToBodyPart(b)
		}

		// Иначе считаем это сообщением чата
		if input != "" {
			// Отправляем на сервер
			writer.WriteString("CHAT:" + input + "\n")
			writer.Flush()
			// Добавляем в локальную историю
			chat.Add("[" + playerName + "] " + input)
		}
	}
}

// Преобразования
func intToBodyPart(i int) BodyPart {
	switch i {
	case 1:
		return Head
	case 2:
		return Torso
	case 3:
		return Legs
	default:
		return Torso
	}
}

func bodyPartToInt(part BodyPart) int {
	switch part {
	case Head:
		return 1
	case Torso:
		return 2
	case Legs:
		return 3
	default:
		return 2
	}
}
