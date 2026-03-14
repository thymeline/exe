package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Системные тексты (вынесены сюда)
const (
	PromptChooseAction    = "%s, выберите действие:\n"
	PromptAttackBlock     = "Для атаки и блока введите два числа через пробел (1-голова, 2-туловище, 3-ноги)"
	PromptChatMessage     = "Чтобы отправить сообщение, введите текст (он будет отправлен как чат)"
	PromptEnterToContinue = "\nНажмите Enter для продолжения..."
	PromptRoundOver       = "\n--- Раунд завершён ---"
)

func ClearConsole() {
	fmt.Print("\033[H\033[2J")
}

func WaitForEnter() {
	fmt.Print(PromptEnterToContinue)
	var discard string
	fmt.Scanln(&discard)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func PauseWithMessage(message string) {
	fmt.Println(message)
	WaitForEnter()
}

type BodyPart string

const (
	Head  BodyPart = "голову"
	Torso BodyPart = "туловище"
	Legs  BodyPart = "ноги"
)

type ItemType string

const (
	WeaponType ItemType = "оружие"
	ArmorType  ItemType = "броня"
	PotionType ItemType = "зелье"
)

type Character interface {
	GetName() string
	GetHP() int
	GetStrength() int
	GetBlocking() BodyPart
	SetBlocking(BodyPart)
	TakeDamage(int)
	IsDead() bool
	GetMaxHP() int
	Heal(int)
}

// Предмет
type Item struct {
	Name        string
	Type        ItemType
	Power       int // Для оружия - урон, для брони - защита, для зелья - эффект
	Description string
}

type Equipment struct {
	Weapon *Item
	Armor  *Item
	Potion *Item
}

type Fighter struct {
	Name      string
	HP        int
	MaxHP     int
	Strength  int
	BaseStr   int
	Blocking  BodyPart
	Inventory []Item
	Equipped  Equipment
}

type Enemy struct {
	Fighter
	Loot *Item
}

func (f *Fighter) GetName() string {
	return f.Name
}

func (f *Fighter) GetHP() int {
	return f.HP
}

func (f *Fighter) GetStrength() int {
	totalStr := f.BaseStr
	if f.Equipped.Weapon != nil {
		totalStr += f.Equipped.Weapon.Power
	}
	return totalStr
}

func (f *Fighter) GetBlocking() BodyPart {
	return f.Blocking
}

func (f *Fighter) SetBlocking(part BodyPart) {
	f.Blocking = part
}

func (f *Fighter) TakeDamage(damage int) {
	actualDamage := damage
	if f.Equipped.Armor != nil {
		armorReduction := f.Equipped.Armor.Power
		actualDamage -= armorReduction
		if actualDamage < 1 {
			actualDamage = 1 // Минимальный урон всегда 1
		}
	}

	f.HP -= actualDamage
	if f.HP < 0 {
		f.HP = 0
	}
}

func (f *Fighter) IsDead() bool {
	return f.HP <= 0
}

func (f *Fighter) GetMaxHP() int {
	return f.MaxHP
}

func (f *Fighter) Heal(amount int) {
	f.HP += amount
	if f.HP > f.MaxHP {
		f.HP = f.MaxHP
	}
	fmt.Printf("%s восстанавливает %d HP\n", f.Name, amount)
}

// методы для работы с инвентарем
func (f *Fighter) AddItem(item Item) {
	f.Inventory = append(f.Inventory, item)
	fmt.Printf("%s получает предмет: %s\n", f.Name, item.Name)
}

func (f *Fighter) ShowInventory() {
	ClearConsole()
	fmt.Printf("\n=== Инвентарь %s ===\n", f.Name)
	if len(f.Inventory) == 0 {
		fmt.Println("Инвентарь пуст")
		return
	}

	for i, item := range f.Inventory {
		fmt.Printf("%d. %s (%s) - %s\n", i+1, item.Name, item.Type, item.Description)
	}

	fmt.Println("\nЭкипировано:")
	if f.Equipped.Weapon != nil {
		fmt.Printf("Оружие: %s (+%d к урону)\n", f.Equipped.Weapon.Name, f.Equipped.Weapon.Power)
	}
	if f.Equipped.Armor != nil {
		fmt.Printf("Броня: %s (+%d к защите)\n", f.Equipped.Armor.Name, f.Equipped.Armor.Power)
	}
	if f.Equipped.Potion != nil {
		fmt.Printf("Активный предмет: %s\n", f.Equipped.Potion.Name)
	}
}

func (f *Fighter) EquipItem(index int) bool {
	if index < 1 || index > len(f.Inventory) {
		fmt.Println("Неверный номер предмета")
		return false
	}

	item := f.Inventory[index-1]

	switch item.Type {
	case WeaponType:
		if f.Equipped.Weapon != nil {
			// Возвращаем старое оружие в инвентарь
			f.Inventory = append(f.Inventory, *f.Equipped.Weapon)
		}
		f.Equipped.Weapon = &item
		fmt.Printf("%s экипирует оружие: %s\n", f.Name, item.Name)

	case ArmorType:
		if f.Equipped.Armor != nil {
			f.Inventory = append(f.Inventory, *f.Equipped.Armor)
		}
		f.Equipped.Armor = &item
		fmt.Printf("%s экипирует броню: %s\n", f.Name, item.Name)

	case PotionType:
		if f.Equipped.Potion != nil {
			f.Inventory = append(f.Inventory, *f.Equipped.Potion)
		}
		f.Equipped.Potion = &item
		fmt.Printf("%s берет предмет для использования: %s\n", f.Name, item.Name)
	}

	// Удаляем предмет из инвентаря
	f.Inventory = append(f.Inventory[:index-1], f.Inventory[index:]...)
	return true
}

func (f *Fighter) TakeOff(itemType ItemType) bool {
	var itemToRemove *Item

	switch itemType {
	case WeaponType:
		if f.Equipped.Weapon == nil {
			fmt.Println("Оружие не экипировано")
			return false
		}
		itemToRemove = f.Equipped.Weapon
		f.Equipped.Weapon = nil

	case ArmorType:
		if f.Equipped.Armor == nil {
			fmt.Println("Броня не экипирована")
			return false
		}
		itemToRemove = f.Equipped.Armor
		f.Equipped.Armor = nil

	case PotionType:
		if f.Equipped.Potion == nil {
			fmt.Println("Активный предмет не выбран")
			return false
		}
		itemToRemove = f.Equipped.Potion
		f.Equipped.Potion = nil
	}

	if itemToRemove != nil {
		f.Inventory = append(f.Inventory, *itemToRemove)
		fmt.Printf("%s снимает: %s\n", f.Name, itemToRemove.Name)
		return true
	}

	return false
}

func (f *Fighter) UsePotionInBattle() bool {
	if f.Equipped.Potion == nil {
		fmt.Println("Нет активного предмета для использования")
		return false
	}

	item := f.Equipped.Potion

	switch item.Type {
	case PotionType:
		// Эффекты зелья
		switch item.Power {
		case 1: // Лечебное зелье
			f.Heal(30)
		case 2: // Силовое зелье
			f.BaseStr += 5
			fmt.Printf("%s получает +5 к силе на этот бой\n", f.Name)
		}

		fmt.Printf("%s использует: %s\n", f.Name, item.Name)
		f.Equipped.Potion = nil
		return true
	}

	return false
}

// Создание тестовых предметов
func CreateTestItems() []Item {
	return []Item{
		{
			Name:        "Меч воина",
			Type:        WeaponType,
			Power:       5,
			Description: "Простой стальной меч",
		},
		{
			Name:        "Кожаная броня",
			Type:        ArmorType,
			Power:       3,
			Description: "Легкая кожаная защита",
		},
		{
			Name:        "Лечебное зелье",
			Type:        PotionType,
			Power:       1, // 1 - лечение
			Description: "Восстанавливает 30 HP",
		},
		{
			Name:        "Зелье силы",
			Type:        PotionType,
			Power:       2, // 2 - усиление
			Description: "Увеличивает силу на 5",
		},
	}
}

// Функция для получения случайного лута от врага
func GetRandomLoot() *Item {
	items := CreateTestItems()
	rand.Seed(time.Now().UnixNano())
	loot := items[rand.Intn(len(items))]
	return &loot
}

// Функция для атаки
func PerformAttack(attacker, defender Character, attackPart BodyPart) {
	damage := attacker.GetStrength()

	fmt.Printf("%s атакует %s в %s\n", attacker.GetName(), defender.GetName(), attackPart)

	if defender.GetBlocking() == attackPart {
		fmt.Printf("%s блокирует удар!\n", defender.GetName())
	} else {
		defender.TakeDamage(damage)
		fmt.Printf("%s получает %d урона\n", defender.GetName(), damage)

		if defender.IsDead() {
			fmt.Printf("%s побежден!\n", defender.GetName())
		}
	}
}

// Функция для запуска боя (горячий стул)
func RunBattle(player1, player2 Character, player2IsAI bool) {
	ClearConsole()
	fmt.Printf("\n=== Бой: %s vs %s ===\n", player1.GetName(), player2.GetName())
	fmt.Printf("%s: %d/%d HP | %s: %d/%d HP\n",
		player1.GetName(), player1.GetHP(), player1.(*Fighter).GetMaxHP(),
		player2.GetName(), player2.GetHP(), player2.(*Fighter).GetMaxHP())

	round := 1
	for !player1.IsDead() && !player2.IsDead() {
		ClearConsole()
		fmt.Printf("\n=== Бой: %s vs %s ===\n", player1.GetName(), player2.GetName())
		fmt.Printf("Раунд: %d\n", round)
		fmt.Printf("%s: %d/%d HP | %s: %d/%d HP\n",
			player1.GetName(), player1.GetHP(), player1.(*Fighter).GetMaxHP(),
			player2.GetName(), player2.GetHP(), player2.(*Fighter).GetMaxHP())

		var player1Attack, player1Block BodyPart
		var player2Attack, player2Block BodyPart
		useItem := false

		if player2IsAI {
			// Ход игрока
			fmt.Printf("\n--- Ход: %s ---\n", player1.GetName())
			fmt.Println("1. Атаковать")
			fmt.Println("2. Использовать предмет")

			var choice int
			fmt.Print("Выберите действие: ")
			fmt.Scan(&choice)

			if choice == 2 {
				// Попытка использовать предмет
				if fighter, ok := player1.(*Fighter); ok {
					if fighter.UsePotionInBattle() {
						useItem = true
						WaitForEnter()
					} else {
						fmt.Println("Не удалось использовать предмет, атакуем вместо этого")
						WaitForEnter()
						choice = 1
					}
				}
			}

			if choice == 1 || !useItem {
				player1Attack, player1Block = GetPlayerChoices(player1.GetName())
			}

			// Ход AI
			player2Attack = MakeAIChoice()
			player2Block = MakeAIChoice()

		} else {
			// PvP режим
			fmt.Printf("\n--- Ход: %s ---\n", player1.GetName())
			player1Attack, player1Block = GetPlayerChoices(player1.GetName())

			ClearConsole()
			fmt.Printf("\n=== Бой: %s vs %s ===\n", player1.GetName(), player2.GetName())
			fmt.Printf("Раунд: %d\n", round)
			fmt.Printf("%s: %d/%d HP | %s: %d/%d HP\n",
				player1.GetName(), player1.GetHP(), player1.(*Fighter).GetMaxHP(),
				player2.GetName(), player2.GetHP(), player2.(*Fighter).GetMaxHP())

			fmt.Printf("\n--- Ход: %s ---\n", player2.GetName())
			player2Attack, player2Block = GetPlayerChoices(player2.GetName())
		}

		if !useItem {
			player1.SetBlocking(player1Block)
			player2.SetBlocking(player2Block)

			fmt.Println("\n--- Результаты раунда ---")
			PerformAttack(player1, player2, player1Attack)
			if player2.IsDead() {
				WaitForEnter()
				break
			}

			PerformAttack(player2, player1, player2Attack)
			if player1.IsDead() {
				WaitForEnter()
				break
			}
		}

		fmt.Printf("\nКонец раунда %d\n", round)
		fmt.Printf("%s: %d/%d HP | %s: %d/%d HP\n",
			player1.GetName(), player1.GetHP(), player1.(*Fighter).GetMaxHP(),
			player2.GetName(), player2.GetHP(), player2.(*Fighter).GetMaxHP())

		round++
		WaitForEnter()
	}

	ClearConsole()
	fmt.Println("\n=== Бой окончен ===")
	if player1.IsDead() {
		fmt.Printf("Победитель: %s\n", player2.GetName())
	} else {
		fmt.Printf("Победитель: %s\n", player1.GetName())
	}
	WaitForEnter()
}

func MakeAIChoice() BodyPart {
	parts := []BodyPart{Head, Torso, Legs}
	rand.Seed(time.Now().UnixNano())
	return parts[rand.Intn(len(parts))]
}

func GetPlayerChoices(playerName string) (BodyPart, BodyPart) {
	fmt.Printf(PromptChooseAction, playerName)

	var attackChoice, blockChoice int

	fmt.Println("Атаковать в:")
	fmt.Println("1. Голову")
	fmt.Println("2. Туловище")
	fmt.Println("3. Ноги")
	fmt.Print("Ваш выбор: ")
	fmt.Scan(&attackChoice)

	fmt.Println("\nБлокировать:")
	fmt.Println("1. Голову")
	fmt.Println("2. Туловище")
	fmt.Println("3. Ноги")
	fmt.Print("Ваш выбор: ")
	fmt.Scan(&blockChoice)

	var attackPart, blockPart BodyPart

	switch attackChoice {
	case 1:
		attackPart = Head
	case 2:
		attackPart = Torso
	case 3:
		attackPart = Legs
	default:
		attackPart = Torso
	}

	switch blockChoice {
	case 1:
		blockPart = Head
	case 2:
		blockPart = Torso
	case 3:
		blockPart = Legs
	default:
		blockPart = Torso
	}

	return attackPart, blockPart
}
