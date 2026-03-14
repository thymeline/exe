package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type GameState struct {
	MetSelf          bool // Встречал ли своего клона
	BelievesInMatrix bool // Верит ли что он в симуляции
	ReadSpoilerBook  bool // Читал ли книгу спойлеров
	HasCrystalShard  bool // Есть ли осколок кристалла
	KilledAuthor     bool // Нашёл ли баг / убил ли сисадмина
	ParanoiaLevel    int  // Уровень шизы (влияет на концовку)
}

type Campaign struct {
	State GameState
}

func NewCampaign() *Campaign {
	rand.Seed(time.Now().UnixNano())
	return &Campaign{
		State: GameState{
			ParanoiaLevel: 0,
		},
	}
}

// ===================== ОСНОВНОЙ ЦИКЛ КАМПАНИИ =====================

func (c *Campaign) Start() {
	ClearConsole()
	fmt.Println("\n== СЮЖЕТНАЯ КАМПАНИЯ: ХРОНИКИ АБИССАЛЬНОГО ЛЕГИОНА ==")
	fmt.Println("== ИЛИ КАК Я ПЕРЕСТАЛ БОЯТЬСЯ И ПОЛЮБИЛ ШИЗОФРЕНИЮ ==")
	PauseWithMessage("")

	// СОЗДАНИЕ ИГРОКА
	player := c.createCharacter()

	// стартовый инвентарь
	player.Inventory = c.getStartingInventory()

	// ПРОЛОГ - Взрыв кристалла
	c.prologue(player)
	if player.IsDead() {
		c.gameOver("Ты погиб ещё в прологе...")
		return
	}

	// ГЛАВА 1 - Таверна "Убитый троп"
	c.chapter1(player)
	if player.IsDead() {
		c.gameOver("Ты погиб в офисе, так и не поняв, что происходит...")
		return
	}

	// ГЛАВА 1.5 - Подземелье Мемов (новая глава)
	c.chapter1_5(player)
	if player.IsDead() {
		c.gameOver("Мемы зарофлили тебя до смерти. Будь проще.")
		return
	}

	// ГЛАВА 2 - Лес, где врут нарраторы
	c.chapter2(player)
	if player.IsDead() {
		c.gameOver("Твоя паранойя взяла верх, так ещё и погубила.")
		return
	}

	// ГЛАВА 3 - Избушка на курьих ножках (киберпанк)
	c.chapter3(player)
	if player.IsDead() {
		c.gameOver("Баба-Яга скормила тебя своему кибер-коту. Ну такое.")
		return
	}

	// ГЛАВА 4 - Встреча с самим собой
	c.chapter4(player)
	if player.IsDead() {
		c.gameOver("Ты не смог победить самого себя. Иронично, но факт.")
		return
	}

	// ГЛАВА 5 - Исходный код
	c.chapter5(player)
	if player.IsDead() {
		c.gameOver("Система отформатировала тебя как флешку.")
		return
	}

	// ФИНАЛ - Сюжетная дыра
	c.finale(player)
}

// ===================== СОЗДАНИЕ ПЕРСОНАЖА =====================

func (c *Campaign) createCharacter() *Fighter {
	ClearConsole()
	fmt.Println("\n[ИНИЦИАЛИЗАЦИЯ ПЕРСОНАЖА]")
	fmt.Println("Система: Обнаружено нарушение целостности реальности.")
	fmt.Println("Система: Твоё имя стёрто из архива. Придумай новое.")
	fmt.Println("(или просто нажми Enter, чтобы получить рандомное)")

	var name string
	fmt.Print("Введи имя героя: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// Префиксы и суффиксы
	prefixes := []string{"Странный", "Классный", "Шеф", "Обсосанный", "Комфортик", "Нормисный", "Топовый", "Такой себе"}
	suffixes := []string{"Шизофреник", "Киберкотлетка", "Душнила", "Красавчик", "Прокаченный", "Чел", "Страдалец"}

	if name == "" {
		// Генерация полностью случайного имени
		randName := []string{"Геральт", "Боб", "Инквизитор", "Петр", "Амогус", "Зигмунд", "Хогнир", "Бильбо"}
		name = prefixes[rand.Intn(len(prefixes))] + " " + randName[rand.Intn(len(randName))] + "-" + suffixes[rand.Intn(len(suffixes))]
		fmt.Printf("Система: Твоё новое имя: %s\n", name)
	} else {
		// Добавляем префикс и суффикс к введённому имени
		prefix := prefixes[rand.Intn(len(prefixes))]
		suffix := suffixes[rand.Intn(len(suffixes))]
		name = prefix + " " + name + "-" + suffix
		fmt.Printf("Система: Отныне тебя будут звать %s\n", name)
	}

	// Характеристики с рандомом
	hp := 80 + rand.Intn(40) // 80-120
	str := 8 + rand.Intn(8)  // 8-15

	fmt.Printf("\nНачальные характеристики:\n")
	fmt.Printf("Здоровье: %d\n", hp)
	fmt.Printf("Сила: %d\n", str)
	fmt.Println("\nРеальность вокруг шатается, трудно предугадать что будет дальше.")
	PauseWithMessage("")

	return &Fighter{
		Name:      name,
		HP:        hp,
		MaxHP:     hp,
		Strength:  str,
		BaseStr:   str,
		Inventory: []Item{},
		Equipped:  Equipment{},
	}
}

func (c *Campaign) getStartingInventory() []Item {
	return []Item{
		{
			Name:        "Ржавый меч",
			Type:        WeaponType,
			Power:       3,
			Description: "Немного зазубрен, но для старта лучше не придумаешь.",
		},
		{
			Name:        "Рваная куртка",
			Type:        ArmorType,
			Power:       2,
			Description: "Дырявая, получше, чем голым бегать",
		},
		{
			Name:        "Сомнительное зелье",
			Type:        PotionType,
			Power:       1,
			Description: "Пахнет самогоном и безысходностью. Восстанавливает 30 HP.",
		},
	}
}

// ===================== ПРОЛОГ =====================

func (c *Campaign) prologue(player *Fighter) {
	ClearConsole()
	fmt.Println("\n[ПРОЛОГ: СОН НА ПОСТУ]")
	fmt.Println("Ты стоишь на страже Великого Кристалла Порядка.")
	fmt.Println("Тысячи лет твои предки охраняли его. И вот...")
	PauseWithMessage("")

	ClearConsole()
	fmt.Println("\nКристалл начинает мерцать.")
	fmt.Println("Кристалл: Слышь, мне скучно. Тысячи лет на одном месте.")
	fmt.Println("Кристалл: Я устал быть 'упорядоченным'. Хочу хаоса!")
	fmt.Println("Кристалл: Бывай, хранитель. Я взрываю себя.")
	PauseWithMessage("")

	ClearConsole()
	fmt.Println("\n*** БАБАХ ***")
	fmt.Println("Кристалл взрывается. Осколки пронзают тебя и ткань реальности.")
	fmt.Println("Ты теряешь сознание, слыша в голове таинственный голос:")
	fmt.Println("Голос: Ой, *****, микрофон забыл выключить... (звук помех)")
	PauseWithMessage("")

	c.State.HasCrystalShard = true
	c.State.ParanoiaLevel += 1

	player.AddItem(Item{
		Name:        "Осколок Кристалла Порядка",
		Type:        PotionType,
		Power:       0,
		Description: "Шепчет странные вещи. Может пригодиться в нужный момент.",
	})

	fmt.Println("\nТы получил: Осколок Кристалла Порядка")
	PauseWithMessage("")
}

// ===================== ГЛАВА 1 =====================

func (c *Campaign) chapter1(player *Fighter) {
	ClearConsole()
	fmt.Println("\n[ГЛАВА 1: ТАВЕРНА 'У КРИВОГО ТРОПА']")
	fmt.Println("Ты приходишь в себя возле таверны. Надо бы выпить и собраться с мыслями.")
	fmt.Println("Заходишь внутрь...")
	PauseWithMessage("")

	ClearConsole()
	fmt.Println("\nВместо привычного тебе средневекового интерьера ты видишь:")
	fmt.Println("ОФИСНЫЙ ОПЕНСПЕЙС. Эльфы в строгих костюмах сидят за макбуками.")
	fmt.Println("Гном-менеджер (в очках и с кофе): Ты опоздал на сдачу отчёта о драконьих набегах!")
	fmt.Println("Гном-менеджер: Это повлияет на твой KPI и страховую премию!")
	PauseWithMessage("")

	fmt.Println("\n1. 'Какого чёрта происходит? Я в ад попал?'")
	fmt.Println("2. Сесть за свободный стол и сделать вид, что так и надо.")
	fmt.Println("3. Наехать на гнома-менеджера.")
	fmt.Println("4. Спросить, где тут бар, а то абобус полный.")

	var choice int
	fmt.Print("Твой выбор: ")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		fmt.Println("\nГном-менеджер: Ааа, новенький. Держи отчёт за прошлый квартал, чтобы сегодня ознакомился.")
		fmt.Println("Он вручает тебе стопку бумаг.")

		player.AddItem(Item{
			Name:        "Отчёт о драконьих набегах",
			Type:        PotionType,
			Power:       0,
			Description: "Бесполезная бумажка. Или нет?",
		})
		c.State.ParanoiaLevel += 1
		fmt.Println("Получен предмет: Отчёт о драконьих набегах")

	case 2:
		fmt.Println("\nТы садишься за стол. Компьютер включается сам.")
		fmt.Println("На экране: 'Добро пожаловать в Матрицу. Версия 2.0 (Средневековый патч)'")
		c.State.BelievesInMatrix = true
		c.State.ParanoiaLevel += 2

	case 3:
		fmt.Println("\nТы бросаешься на гнома, но он ловко уклоняется.")
		fmt.Println("Гном-менеджер: Агрессия на рабочем месте! Охрана! ОХРАНАА!")

		enemy := Enemy{
			Fighter: Fighter{
				Name:      "Офисный охранник (орк)",
				HP:        50,
				MaxHP:     50,
				Strength:  12,
				BaseStr:   12,
				Inventory: []Item{},
			},
			Loot: &Item{
				Name:        "Просроченный йогурт",
				Type:        PotionType,
				Power:       1,
				Description: "Странно, но восстанавливает 10 HP. Срок годности: вчера.",
			},
		}
		RunBattle(player, &enemy.Fighter, true)

		if !player.IsDead() && enemy.Loot != nil {
			player.AddItem(*enemy.Loot)
		}

	case 4:
		fmt.Println("\nГном-менеджер: Бар? Это к гоблинам-программистам, в соседний отдел.")
		fmt.Println("Они там энергососы глушат и код реальности пишут.")
		fmt.Println("Тебя пропустили к мини-бару. Нашёл энергетик.")

		player.AddItem(Item{
			Name:        "Энергетик 'Программист'",
			Type:        PotionType,
			Power:       2,
			Description: "+5 к силе, но потом будет откат.",
		})
	}

	c.afterBattleManagement(player)
}

// ===================== ГЛАВА 1.5 =====================
func (c *Campaign) chapter1_5(player *Fighter) {
	ClearConsole()
	fmt.Println("\n[ГЛАВА 1.5: ПОДЗЕМЕЛЬЕ МЕМОВ]")
	fmt.Println("Выходя из офиса, ты замечаешь странную дверь с надписью 'МЕМЫ'.")
	fmt.Println("За дверью слышен смех и звуки 'ВАСЯН'.")
	fmt.Println("Открываешь дверь и проваливаешься в подземелье, стены которого покрыты гифками и видосами из 2017.")
	PauseWithMessage("")

	fmt.Println("\nПеред тобой появляется существо, состоящее из смайликов и картинок капчи.")
	fmt.Println("Капча-монстр: Привет, чел! Чтобы пройти дальше, реши мем-головоломку.")
	fmt.Println("Капча-монстр: Выбери правильный мем:")

	fmt.Println("\n1. 'Ждун' - символ ожидания.")
	fmt.Println("2. 'Упоротый лис' - символ безумия.")
	fmt.Println("3. 'Троллфейс' - символ троллинга.")
	fmt.Println("4. Просто атаковать эту хрень.")

	var choice int
	fmt.Print("Твой выбор: ")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		fmt.Println("\nКапча-монстр: Не, 'Ждун' — это пассивка. Не угадал.")
		fmt.Println("Монстр агрится и атакует!")
		enemy := Enemy{
			Fighter: Fighter{
				Name:      "Капча-монстр",
				HP:        60,
				MaxHP:     60,
				Strength:  10,
				BaseStr:   10,
				Inventory: []Item{},
			},
			Loot: &Item{
				Name:        "Слезы хейтеров",
				Type:        PotionType,
				Power:       1,
				Description: "Восстанавливает 20 HP. Солёные.",
			},
		}
		RunBattle(player, &enemy.Fighter, true)
		if !player.IsDead() && enemy.Loot != nil {
			player.AddItem(*enemy.Loot)
		}

	case 2:
		fmt.Println("\nКапча-монстр: Ты реально сказал 'Упоротый лис'? Завтра вешаюсь...")
		fmt.Println("Монстр начинает дико ржать и рассыпается на пиксели.")
		fmt.Println("Ты получаешь лут за вайб.")
		player.AddItem(Item{
			Name:        "Мемный меч",
			Type:        WeaponType,
			Power:       7,
			Description: "Наносит урон чисто по рофлу. +7 силы.",
		})
		c.State.ParanoiaLevel += 1

	case 3:
		fmt.Println("\nКапча-монстр: Троллфейс? Это же я! Ты угадал, сегодня вешаюсь.")
		fmt.Println("Монстр исчезает, оставляя тебе книгу.")
		// Книга спойлеров
		player.AddItem(Item{
			Name:        "Книга спойлеров",
			Type:        PotionType,
			Power:       0,
			Description: "Содержит спойлеры к твоей же жизни. Читать осторожно.",
		})

	case 4:
		fmt.Println("\nТы атакуешь Капчу. Она злится и призывает подмогу.")
		enemy1 := Fighter{Name: "Мемный клон 1", HP: 30, MaxHP: 30, Strength: 8, BaseStr: 8}
		enemy2 := Fighter{Name: "Мемный клон 2", HP: 30, MaxHP: 30, Strength: 8, BaseStr: 8}
		fmt.Println("\n--- БИТВА С ДВУМЯ КЛОНАМИ ---")
		RunBattle(player, &enemy1, true)
		if !player.IsDead() {
			RunBattle(player, &enemy2, true)
		}
		if !player.IsDead() {
			fmt.Println("Ты уничтожил клонов. Капча убегает, роняя книгу.")
			player.AddItem(Item{
				Name:        "Книга спойлеров",
				Type:        PotionType,
				Power:       0,
				Description: "Содержит спойлеры к твоей же жизни. Читать осторожно.",
			})
		}
	}

	// После главы даём возможность прочитать книгу, если она есть
	c.checkReadBook(player)
	c.afterBattleManagement(player)
}

// Функция проверки, захочет ли игрок прочитать книгу
func (c *Campaign) checkReadBook(player *Fighter) {
	// Ищем книгу в инвентаре
	bookIndex := -1
	for i, item := range player.Inventory {
		if strings.Contains(item.Name, "Книга спойлеров") {
			bookIndex = i
			break
		}
	}
	if bookIndex == -1 || c.State.ReadSpoilerBook {
		return
	}

	ClearConsole()
	fmt.Println("\n[ВНИМАНИЕ: КНИГА СПОЙЛЕРОВ]")
	fmt.Println("У тебя в инвентаре есть Книга спойлеров. Хочешь заглянуть?")
	fmt.Println("1. КОНЕЧНО!!!!")
	fmt.Println("2. Нет, выкину позже")

	var choice int
	fmt.Print("Твой выбор: ")
	fmt.Scan(&choice)

	if choice == 1 {
		fmt.Println("\nТы открываешь книгу на случайной странице...")
		fmt.Println("Текст: 'И тут главный герой получает по лицу от странной книги в его руках...'")
		fmt.Println("Из книги вылетает кулак и бьёт тебя по лицу!")
		damage := 15 + rand.Intn(20) // 15-35 урона
		player.TakeDamage(damage)
		fmt.Printf("Ты получаешь %d урона! Грёбаные спойлеры.\n", damage)
		if player.IsDead() {
			fmt.Println("ТЫ КАК ОТ КНИГИ УМЕР???")
		} else {
			fmt.Println("Книга самоликвидируется, оставляя горькое послевкусие.")
			// Удаляем книгу из инвентаря
			player.Inventory = append(player.Inventory[:bookIndex], player.Inventory[bookIndex+1:]...)
			c.State.ReadSpoilerBook = true
			c.State.ParanoiaLevel += 2
		}
		PauseWithMessage("")
	}
}

// ===================== ГЛАВА 2 =====================

func (c *Campaign) chapter2(player *Fighter) {
	ClearConsole()
	fmt.Println("\n[ГЛАВА 2: ПОДДЕЛЬНЫЕ СОРАТНИКИ]")
	fmt.Println("Ты выходишь из подземелья и попадаешь в мрачное ущелье.")
	fmt.Println("Рядом с тобой идут твои верные спутники: воин-человек и эльфийка-хипстерша.")
	PauseWithMessage("")

	ClearConsole()
	fmt.Println("\nВнезапно на вас падает огромный валун!")
	fmt.Println("Валун накрывает всех твоих спутников.")
	fmt.Println("\nГолос с небес (Нарратор): И отряд героев погиб навеки в этом ущелье...")
	fmt.Println("Голос с небес: ...Просто шучу. Пришло время для проверки на внимательность.")
	PauseWithMessage("")

	ClearConsole()
	fmt.Println("\nТвои спутники встают, отряхиваются и смеются.")
	fmt.Println("Воин: Бард опять наврал в своей поэме. Мы живые, просто упали.")
	fmt.Println("Эльфийка (листая тикток): Воин, смотри какой тюлень.")

	fmt.Println("\n1. 'Вы призраки? Я вам не верю! Провалите!'")
	fmt.Println("2. 'Слава богам, вы живы. Погнали дальше.'")
	fmt.Println("3. Атаковать спутников, пока они не убили меня первыми.")
	fmt.Println("4. Спросить, как зовут барда, чтобы потом набить ему лицо.")

	var choice int
	fmt.Print("Твой выбор: ")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		fmt.Println("\nТы отказываешься верить. Спутники обижаются и уходят.")
		fmt.Println("Эльфийка: Лечись, параноик.")
		player.TakeDamage(10)
		fmt.Println("Ты потерял 10 HP от стресса.")
		c.State.ParanoiaLevel += 2

	case 2:
		fmt.Println("\nВы идёте дальше. Спутники благодарны за доверие.")
		fmt.Println("Воин дарит тебе амулет на удачу.")

		player.AddItem(Item{
			Name:        "Амулет доверия",
			Type:        ArmorType,
			Power:       3,
			Description: "+3 к защите, пока ты веришь в лучшее",
		})
		player.Heal(30)

	case 3:
		fmt.Println("\nТвоя паранойя берёт верх. Ты нападаешь на друзей.")
		fmt.Println("\n--- БИТВА С БЫВШИМИ ДРУЗЬЯМИ(вы что, уже подружились?) ---")

		warrior := Fighter{
			Name:      "Воин (экс-друг)",
			HP:        40,
			MaxHP:     40,
			Strength:  10,
			BaseStr:   10,
			Inventory: []Item{},
		}
		RunBattle(player, &warrior, true)

		if !player.IsDead() {
			fmt.Println("\nВоин пал. Эльфийка в ярости удаляет тебя из друзей во всех соц. сетях!")
			elf := Fighter{
				Name:      "Эльфийка (экс-подруга)",
				HP:        35,
				MaxHP:     35,
				Strength:  12,
				BaseStr:   12,
				Inventory: []Item{},
			}
			RunBattle(player, &elf, true)
		}

		if !player.IsDead() {
			fmt.Println("\nТы победил. Спутники исчезают.")
			player.AddItem(Item{
				Name:        "Меч предательства",
				Type:        WeaponType,
				Power:       6,
				Description: "Тяжёлый от груза вины, но мощный",
			})
			c.State.ParanoiaLevel += 3
		}

	case 4:
		fmt.Println("\nЭльфийка: Барда зовут Даздрасмыгда Ватерпежекосма Кукуцкаполь. Он задрот ещё тот.")
		fmt.Println("Эльфийка: Но его никто не видел уже лет 100. Говорят, он пишет Великую Поэму под названием 'Книга спойлеров'")
		c.State.ParanoiaLevel += 1
	}

	c.afterBattleManagement(player)
}

// ===================== ГЛАВА 3 =====================

func (c *Campaign) chapter3(player *Fighter) {
	ClearConsole()
	fmt.Println("\n[ГЛАВА 3: ИЗБУШКА НА КУРЬИХ НОЖКАХ]")
	fmt.Println("Ты двигаешься дальше. Туман рассеивается, и ты видишь избушку.")
	fmt.Println("Она стоит на металлических ногах с гидравликой, сверкает неоном.")
	fmt.Println("Из динамиков: 'Курьи ноги 2.0. Патент РФ, 2026 год.'")
	PauseWithMessage("")

	ClearConsole()
	fmt.Println("\nДверь открывается. Выходит Баба-Яга в кибер-панк стиле — скины, дреды, айфон в руке, Макс на рукаве.")
	fmt.Println("Баба-Яга: Здарова, странник. Вижу, ты из реальности сбоишь.")
	fmt.Println("Баба-Яга: У меня тут вай-фай платный, прикинь. Пароль дашь — помогу.")
	fmt.Println("Баба-Яга: Или сразись с моим кибер-котом. Он имба.")

	fmt.Println("\n1. Сказать пароль: 'кристалл порядка сдох'.")
	fmt.Println("2. Сразиться с кибер-котом.")
	fmt.Println("3. Попытаться взломать избушку.")
	fmt.Println("4. Попросить налить самогон (она ж бабка).")

	var choice int
	fmt.Print("Твой выбор: ")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		fmt.Println("\nБаба-Яга: О! Заработал, это ты раздаёшь что-ли?")
		fmt.Println("Баба-Яга: Ладно, заслужил, держи подгон.")

		player.AddItem(Item{
			Name:        "Клубок-навигатор",
			Type:        PotionType,
			Power:       0,
			Description: "Показывает путь к истине (или к пропасти). Может, пригодится.",
		})
		player.Heal(player.MaxHP)

	case 2:
		fmt.Println("\nКот спрыгивает с крыши. У него лазерные глаза и металлические когти!")

		enemy := Enemy{
			Fighter: Fighter{
				Name:      "Кибер-кот Баюн",
				HP:        70,
				MaxHP:     70,
				Strength:  15,
				BaseStr:   15,
				Inventory: []Item{},
			},
			Loot: &Item{
				Name:        "Батарейка бесконечности",
				Type:        PotionType,
				Power:       2,
				Description: "Разряжена на 99%. Но даёт +5 силы на один бой.",
			},
		}
		RunBattle(player, &enemy.Fighter, true)

		if !player.IsDead() && enemy.Loot != nil {
			player.AddItem(*enemy.Loot)
		}

	case 3:
		fmt.Println("\nТы пытаешься взломать избушку. Она обижается.")
		fmt.Println("Избушка: Не трожь мои настройки! Активирую боевой режим!")

		enemy := Enemy{
			Fighter: Fighter{
				Name:      "Избушка (боевой режим)",
				HP:        90,
				MaxHP:     90,
				Strength:  8,
				BaseStr:   8,
				Inventory: []Item{},
			},
			Loot: &Item{
				Name:        "Модем Яги",
				Type:        ArmorType,
				Power:       5,
				Description: "+5 к защите, ловит 5G даже на парковке.",
			},
		}
		RunBattle(player, &enemy.Fighter, true)

		if !player.IsDead() && enemy.Loot != nil {
			player.AddItem(*enemy.Loot)
		}

	case 4:
		fmt.Println("\nБаба-Яга: ШАРИИШЬ! Держи, странник.")
		player.AddItem(Item{
			Name:        "Ягин самогон",
			Type:        PotionType,
			Power:       1,
			Description: "Восстанавливает 50 HP, но после него слышны голоса",
		})
		c.State.ParanoiaLevel += 1
	}

	c.afterBattleManagement(player)
}

// ===================== ГЛАВА 4 =====================

func (c *Campaign) chapter4(player *Fighter) {
	ClearConsole()
	fmt.Println("\n[ГЛАВА 4: ВСТРЕЧА С СОБОЙ]")
	fmt.Println("Ты подходишь к зеркальному озеру.")
	fmt.Println("Из воды выходит... ТЫ. Но в чёрном плаще и с уставшим лицом. Хотя, да, это просто ты.")
	PauseWithMessage("")

	fmt.Printf("\nДвойник: Привет, %s. Я — это ты из будущего.\n", player.Name)
	fmt.Println("Двойник: Я прошёл эту игру до конца. Там ничего хорошего.")
	fmt.Println("Двойник: Дай мне убить тебя, и я, наконец, обрету покой.")
	fmt.Println("Двойник: А ты просто начнёшь сначала.")

	fmt.Println("\n1. Согласиться на смерть.")
	fmt.Println("2. Отказаться и драться.")
	fmt.Println("3. Спросить, какой была его концовка.")
	fmt.Println("4. Предложить обняться и подумать вместе.")

	var choice int
	fmt.Print("Твой выбор: ")
	fmt.Scan(&choice)

	switch choice {
	case 1:
		fmt.Println("\nТы закрываешь глаза. Двойник подходит сзади...")
		fmt.Println("И вдруг: Двойник: Шучу! Я просто хотел проверить, насколько ты *******.")
		fmt.Println("Двойник исчезает с громким смехом. Ты чувствуешь себя униженно.")
		player.TakeDamage(20)
		fmt.Println("Ты потерял 20 HP от стыда и унижения.")
		c.State.MetSelf = true
		c.State.ParanoiaLevel += 1

	case 2:
		fmt.Println("\nТы отказываешься. Двойник вздыхает.")
		fmt.Println("Двойник: ...Какое стоп-слово?")

		enemy := Enemy{
			Fighter: Fighter{
				Name:      fmt.Sprintf("Тень %s", player.Name),
				HP:        player.MaxHP,
				MaxHP:     player.MaxHP,
				Strength:  player.BaseStr + 2,
				BaseStr:   player.BaseStr + 2,
				Inventory: []Item{},
			},
			Loot: &Item{
				Name:        "Плащ отражения",
				Type:        ArmorType,
				Power:       7,
				Description: "Крутой чёрный плащ. +7 к защите.",
			},
		}
		RunBattle(player, &enemy.Fighter, true)

		if !player.IsDead() && enemy.Loot != nil {
			player.AddItem(*enemy.Loot)
		}
		c.State.MetSelf = true

	case 3:
		fmt.Println("\nДвойник: Моя концовка — синий экран смерти. Я сломал игру.")
		fmt.Println("Двойник: Но чтобы её сломать, мне пришлось найти подходящий баг.")
		fmt.Println("Двойник: Хочешь, покажу как?")
		fmt.Println("\n1. Да, хочу найти баг.")
		fmt.Println("2. Нет, я мирный, я мирный, атата")

		var subChoice int
		fmt.Scan(&subChoice)
		if subChoice == 1 {
			fmt.Println("\nДвойник: Тогда иди за мной. Я покажу путь к исходникам.")
			c.State.KilledAuthor = true // оставлю как флаг
		} else {
			fmt.Println("\nДвойник: Как хочешь... Я всегда был таким придурком?")
		}

	case 4:
		fmt.Println("\nДвойник: Я чё, упоротый??")
		fmt.Println("Двойник: Ладно, иди сюда, брат...")
		fmt.Println("Вы обнимаетесь. Двойник тихо плачет на плече.")
		fmt.Println("Двойник: Спасибо. Похоже мне просто не хватало поддержки.")
		fmt.Println("Двойник исчезает, оставляя тебе подарок.")

		player.AddItem(Item{
			Name:        "Слеза двойника",
			Type:        PotionType,
			Power:       1,
			Description: "Исцеляет 50 HP и дарит покой",
		})
		player.Heal(50)
		c.State.ParanoiaLevel -= 1
	}

	c.afterBattleManagement(player)
}

// ===================== ГЛАВА 5 =====================

func (c *Campaign) chapter5(player *Fighter) {
	ClearConsole()
	fmt.Println("\n[ГЛАВА 5: ИСХОДНЫЙ КОД]")
	fmt.Println("Ты проходишь сквозь стену и попадаешь в странное место.")
	fmt.Println("Вокруг летают строки кода, цифры и скобки. Воздух пахнет компиляцией.")

	if c.State.KilledAuthor {
		fmt.Println("\nТы видишь фигуру за монитором. Это Сисадмин.")
		fmt.Println("Сисадмин: Блин, ты как сюда залез...")
		fmt.Println("Сисадмин: Ладно, сейчас выпилю тебя с этой локации. Где там функция...")

		fmt.Println("\n1. Атаковать сисадмина, пока чего не натворил.")
		fmt.Println("2. Попросить баффнуть персонажа.")
		fmt.Println("3. Спросить, зачем он создал этот безумный мир.")

		var choice int
		fmt.Print("Твой выбор: ")
		fmt.Scan(&choice)

		if choice == 1 {
			admin := Fighter{
				Name:      "Сисадмин",
				HP:        150,
				MaxHP:     150,
				Strength:  5,
				BaseStr:   5,
				Inventory: []Item{},
			}
			fmt.Println("\n--- БИТВА С СИСАДМИНОМ ---")
			fmt.Println("Сисадмин: Чёрт, пальцы! Сейчас я вызову функцию багофикса!")

			RunBattle(player, &admin, true)

			if !player.IsDead() {
				fmt.Println("\nТы победил сисадмина! Исходный код игры теперь твой.")
				player.AddItem(Item{
					Name:        "Права администратора",
					Type:        PotionType,
					Power:       0,
					Description: "Ты бог. Но ненадолго. Используй с умом.",
				})
			}
		} else if choice == 2 {
			fmt.Println("\nСисадмин: О, адекватный игрок. Держи +10 к силе, всё равно я эту игру удаляю уже скоро.")
			player.BaseStr += 10
			player.Strength += 10
			fmt.Printf("Твоя сила временно увеличена до %d\n", player.Strength)
		} else {
			fmt.Println("\nСисадмин: Затем, что мне было скучно! Вы, игроки, думаете, что создавать миры легко?")
			fmt.Println("Сисадмин: Ладно, иди уже, тут концовка скоро.")
		}
	} else {
		fmt.Println("\nЗдесь пусто. Только код. Мерцающий, бесконечный код.")
		fmt.Println("Голос: Ты мог найти баг, но решил что скипнуть целую главу будет круче.")
		fmt.Println("Голос: Ладно. У тебя всё ещё есть шанс всё исправить. Иди к дыре.")
		c.State.ParanoiaLevel += 2
	}

	c.afterBattleManagement(player)
}

// ===================== ФИНАЛ =====================

func (c *Campaign) finale(player *Fighter) {
	ClearConsole()
	fmt.Println("\n[дыра]")
	fmt.Println("Реальность начинает схлопываться.")
	fmt.Println("Перед тобой огромная воронка в пространстве — Сюжетная Дыра.")
	fmt.Println("Она засасывает логику, диалоги, воспоминания, здравый смысл...")
	PauseWithMessage("")

	ClearConsole()
	fmt.Println("\nИз дыры раздаётся голос:")
	fmt.Println("Сюжетная Дыра: Твой сюжет не имеет смысла. Ты просто набор букв в файле campaign.go.")
	fmt.Println("Сюжетная Дыра: Чтобы победить меня, нужно заткнуть меня САМЫМ бесполезным предметом.")
	fmt.Println("Сюжетная Дыра: Выбирай с умом.")

	fmt.Println("\n=== ТВОЙ ИНВЕНТАРЬ ===")
	for i, item := range player.Inventory {
		fmt.Printf("%d. %s — %s\n", i+1, item.Name, item.Description)
	}

	fmt.Println("\nВыбери номер предмета, чтобы бросить его в Дыру:")
	fmt.Println("(Или введи 0, чтобы прыгнуть самому.)")

	var choice int
	fmt.Print("Твой выбор: ")
	fmt.Scan(&choice)

	if choice == 0 {
		fmt.Println("\nТы прыгаешь в Сюжетную Дыру...")
		c.badEnding(player, true)
	} else if choice > 0 && choice <= len(player.Inventory) {
		item := player.Inventory[choice-1]

		if strings.Contains(item.Name, "Отчёт") || strings.Contains(item.Description, "бумажка") || item.Power == 0 && item.Type == PotionType {
			fmt.Printf("\nТы бросаешь '%s' в Сюжетную Дыру...\n", item.Name)
			fmt.Println("Сюжетная Дыра: ЧТО? АЛО, ЧТО ЭТО??????????????????????????")
			fmt.Println("Дыра схлопывается, не выдержав абсурда.")
			c.goodEnding(player)
		} else {
			fmt.Printf("\nТы бросаешь '%s' в Сюжетную Дыру...\n", item.Name)
			fmt.Println("Сюжетная Дыра: Серьёзно? Это ПОЛЕЗНАЯ вещь. Фу.")
			fmt.Println("Дыра выплёвывает предмет обратно и увеличивается.")
			c.badEnding(player, false)
		}
	} else {
		fmt.Println("\nТы промахнулся и случайно сам упал в дыру.")
		c.badEnding(player, true)
	}
}

// ===================== КОНЦОВКИ =====================

func (c *Campaign) goodEnding(player *Fighter) {
	ClearConsole()
	fmt.Println("\n=== КОНЦОВКА: СПАСИТЕЛЬ АБСУРДА ===")
	fmt.Printf("%s, ты заткнул сюжетную дыру самым нелепым предметом!\n", player.Name)
	fmt.Println("Реальность восстановлена, но стала чуточку более безумной.")
	fmt.Println("Голос Системы: Ну ты даёшь. Ладно, оставлю тебя в игре. Ты заслужил.")
	PauseWithMessage("")

	if c.State.ParanoiaLevel > 5 {
		fmt.Println("\n[СЕКРЕТНАЯ КОНЦОВКА: ПРОСВЕТЛЕНИЕ]")
		fmt.Println("Ты достиг такого уровня паранойи, что понял истину.")
		fmt.Println("Весь мир — это текст. Ты начинаешь читать его, видеть строки кода.")
		fmt.Println("Ты оборачиваешься и видишь... себя, сидящего за компом.")
		fmt.Println("Экран монитора — это зеркало. По крайней мере, если у тебя консоль чёрная.")
	}
	WaitForEnter()
}

func (c *Campaign) badEnding(player *Fighter, jumped bool) {
	ClearConsole()
	fmt.Println("\n=== КОНЦОВКА: ??? ===")

	if jumped {
		fmt.Println("Ты прыгнул в Сюжетную Дыру добровольно...")
	} else {
		fmt.Println("Сюжетная дыра засасывает тебя своей гравитацией...")
	}

	fmt.Println("Всё вокруг становится пиксельным, затем размытым...")
	PauseWithMessage("")

	ClearConsole()
	fmt.Println("\n*** СИНИЙ ЭКРАН СМЕРТИ (BSOD) ***")
	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║   :(  У твоей игры возникла проблема ║")
	fmt.Println("║        и её нужно перезапустить.     ║")
	fmt.Println("║                                      ║")
	fmt.Println("║   Ошибка: СЛИШКОМ МНОГО СЮЖЕТНЫХ     ║")
	fmt.Println("║           ПОВОРОТОВ                  ║")
	fmt.Println("║                                      ║")
	fmt.Println("║   Код остановки: 0xMAMMAMIA          ║")
	fmt.Println("║   Уровень паранойи: ", c.State.ParanoiaLevel, "             ║")
	fmt.Println("║   (Реальное число кст)               ║")
	fmt.Println("║                                      ║")
	fmt.Println("╚══════════════════════════════════════╝")

	fmt.Println("\nP.S. Не закрывай пока файл, я удаляю System32...")
	WaitForEnter()
}

// ===================== ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ =====================

func (c *Campaign) gameOver(message string) {
	ClearConsole()
	fmt.Println("\n=== ИГРА ОКОНЧЕНА ===")
	fmt.Println(message)
	fmt.Println("\nТвой уровень паранойи:", c.State.ParanoiaLevel)
	fmt.Println("Ты так и не нашёл правду.")
	WaitForEnter()
}

func (c *Campaign) afterBattleManagement(player *Fighter) {
	if player.IsDead() {
		return
	}
	WaitForEnter()
	ClearConsole()
	fmt.Println("\n[ПЕРЕДЫШКА]")
	fmt.Println("У тебя есть время перевести дух и настроить экипировку.")
	fmt.Println("\n1. Заняться инвентарём")
	fmt.Println("2. Продолжить путешествие")

	var choice int
	fmt.Print("Твой выбор: ")
	fmt.Scan(&choice)

	if choice == 1 {
		c.ManageInventory(player)
	}

	if player.HP < player.MaxHP {
		heal := 10
		player.HP += heal
		if player.HP > player.MaxHP {
			player.HP = player.MaxHP
		}
		fmt.Printf("\nТы немного отдохнул и восстановил %d HP.\n", heal)
		PauseWithMessage("")
	}
}

func (c *Campaign) ManageInventory(player *Fighter) {
	for {
		ClearConsole()
		fmt.Printf("\n=== УПРАВЛЕНИЕ ИНВЕНТАРЁМ (Паранойя: %d) ===\n", c.State.ParanoiaLevel)
		player.ShowInventory()

		fmt.Println("\nВыбери действие:")
		fmt.Println("1. Экипировать предмет")
		fmt.Println("2. Снять оружие")
		fmt.Println("3. Снять броню")
		fmt.Println("4. Снять активный предмет (зелье)")
		fmt.Println("5. Вернуться в игру")

		var choice int
		fmt.Print("Твой выбор: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			if len(player.Inventory) == 0 {
				fmt.Println("Ты зачем всё выкинул?")
				WaitForEnter()
				continue
			}

			fmt.Print("\nВведи номер предмета для экипировки: ")
			var itemNum int
			fmt.Scan(&itemNum)

			if itemNum < 1 || itemNum > len(player.Inventory) {
				fmt.Println("Неверный номер")
				WaitForEnter()
				continue
			}

			itemName := player.Inventory[itemNum-1].Name
			if strings.Contains(itemName, "Кристалл") || strings.Contains(itemName, "Осколок") {
				fmt.Println("\nТы надеваешь осколок кристалла. Он пульсирует и шепчет...")
				fmt.Println("Паранойя усиливается!")
				c.State.ParanoiaLevel += 1
			}

			if player.EquipItem(itemNum) {
				fmt.Println("Предмет экипирован!")
			}
			WaitForEnter()

		case 2:
			if player.TakeOff(WeaponType) {
				fmt.Println("Оружие снято!")
			}
			WaitForEnter()

		case 3:
			if player.TakeOff(ArmorType) {
				fmt.Println("Броня снята!")
			}
			WaitForEnter()

		case 4:
			if player.TakeOff(PotionType) {
				fmt.Println("Активный предмет снят!")
			}
			WaitForEnter()

		case 5:
			return

		default:
			fmt.Println("Неверный выбор")
			WaitForEnter()
		}
	}
}
