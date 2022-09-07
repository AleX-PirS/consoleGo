package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Набор enum для поверхностей/мест с предметами
const (
	TABLE = iota
	CHAIR
)

const WRONG = "неверная команда"

// Room - структура примитивной комнаты
type Room struct {
	title   string
	info    string
	ways    []*Room
	objects map[Object]int
	isOpen  bool
}

// Object - структура для объекта взаимодействия
type Object struct {
	title, buff           string
	location              int
	toWear, toTake, toUse bool
	contact               string
	result                string
}

// Player - структура, описывающая игрока
type Player struct {
	inventory  []Object
	isWithBag  bool
	actualRoom *Room
	startRoom  *Room
}

// Игровые объекты
var (
	tea = Object{
		title:    "чай",
		location: TABLE,
		toWear:   false,
		toTake:   false,
		toUse:    true,
		contact:  "я",
		result:   "чай выпит",
	}
	key = Object{
		title:    "ключи",
		location: TABLE,
		toWear:   false,
		toTake:   true,
		toUse:    true,
		contact:  "дверь",
		result:   "дверь открыта",
	}
	conspectus = Object{
		title:    "конспекты",
		location: TABLE,
		toWear:   false,
		toTake:   true,
		toUse:    true,
		contact:  "я",
		result:   "всё выучено",
	}
	bag = Object{
		title:    "рюкзак",
		buff:     "курьер",
		location: CHAIR,
		toWear:   true,
		toTake:   false,
		toUse:    false,
		contact:  "",
	}
)

// Игровые комнаты и игрок
var (
	home     = Room{}
	street   = Room{}
	kitchen  = Room{}
	myRoom   = Room{}
	corridor = Room{}
	player   = Player{}
)

func main() {
	initGame()
	for {
		in := bufio.NewReader(os.Stdin)
		comm, err := in.ReadString('\n')
		comm = comm[:len(comm)-2]
		if err != nil {
			fmt.Println("Ошибка ввода: ", err)
		}
		ans := handleCommand(comm)
		fmt.Println(ans)
	}
}

func decodeLocation(location int) string {
	switch location {
	case CHAIR:
		return "на стуле: "
	case TABLE:
		return "на столе: "
	}
	return "неизвестное место"
}

func initGame() {
	home = Room{
		title:   "домой",
		info:    "ты дома, поздравляю!",
		objects: make(map[Object]int),
		ways:    []*Room{},
		isOpen:  true,
	}
	// Чтобы открыть улицу, нужно применить ключ на дверь в коридоре
	street = Room{
		title:   "улица",
		info:    "на улице весна.",
		ways:    []*Room{&home},
		isOpen:  false,
		objects: nil,
	}
	kitchen = Room{
		title:   "кухня",
		info:    "ты находишься на кухне,",
		ways:    []*Room{},
		objects: map[Object]int{tea: 1},
		isOpen:  true,
	}
	myRoom = Room{
		title:   "комната",
		info:    "ты в своей комнате.",
		ways:    []*Room{},
		objects: map[Object]int{key: 1, conspectus: 1, bag: 1},
		isOpen:  true,
	}
	corridor = Room{
		title:   "коридор",
		info:    "ничего интересного.",
		ways:    []*Room{&kitchen, &myRoom, &street},
		objects: map[Object]int{},
		isOpen:  true,
	}
	home.ways = append(home.ways, &corridor)
	kitchen.ways = append(kitchen.ways, &corridor)
	myRoom.ways = append(myRoom.ways, &corridor)

	player = Player{
		inventory:  make([]Object, 0),
		isWithBag:  false,
		actualRoom: &kitchen,
		startRoom:  &kitchen,
	}
}

func handleCommand(command string) string {
	result := ""
	arrOfArgs := strings.Split(command, " ")
	switch arrOfArgs[0] {
	case "осмотреться":
		result = commCheck(result)
	case "идти":
		result = commGo(result, arrOfArgs)
	case "надеть":
		result = commWear(result, arrOfArgs)
	case "взять":
		result = commTake(result, arrOfArgs)
	case "применить":
		result = commUse(result, arrOfArgs)
	default:
		result = "неизвестная команда"
		fmt.Println(command)
	}
	return result
}

func actualWays(result string) string {
	result += "можно пройти - "

	for _, room := range player.actualRoom.ways {
		result = result + room.title + ", "
	}
	result = result[:len(result)-2]
	return result
}

// commCheck - обработка команды "осмотреться"
func commCheck(result string) string {
	if player.actualRoom.title == player.startRoom.title {
		result += player.actualRoom.info + " "
	}
	sumOfObj := 0
	roomInventory := make(map[int][]string)
	for obj, countOfObj := range player.actualRoom.objects {
		sumOfObj += countOfObj
		if countOfObj > 0 {
			roomInventory[obj.location] = append(roomInventory[obj.location], obj.title)
		}
	}
	locations := make([]int, 0, len(roomInventory))
	for location := range roomInventory {
		locations = append(locations, location)
	}
	// Sort location by name
	sort.Ints(locations)
	for _, location := range locations {
		result += decodeLocation(location)
		sort.Strings(roomInventory[location])
		for _, elem := range roomInventory[location] {
			result = result + elem + ", "
		}
	}
	if sumOfObj == 0 {
		result = "пустая комната. "
	}
	if player.actualRoom.title == player.startRoom.title && player.isWithBag {
		result += "надо идти в универ. "
	} else if player.actualRoom.title == player.startRoom.title && !player.isWithBag {
		result += "надо собрать рюкзак и идти в универ. "
	}
	result = result[:len(result)-2] + ". "
	result = actualWays(result)
	return result
}

// commGo - обработка команды "идти"
func commGo(result string, arrOfArgs []string) string {
	if len(arrOfArgs) == 1 {
		result = WRONG
	}
	direction := arrOfArgs[1]
	detectedDirection := 0
	doorFlag := 0
	for _, room := range player.actualRoom.ways {
		if direction == room.title {
			detectedDirection = 1
			if room.isOpen {
				player.actualRoom = room
				if room.title == player.startRoom.title {
					result = result + player.startRoom.title + ", " + "ничего интересного. "
				} else {
					result = result + room.info + " "
				}
			} else {
				doorFlag = 1
				result = "дверь закрыта"
			}
			break
		}
	}
	if doorFlag == 0 {
		result = actualWays(result)
	}
	if detectedDirection == 0 {
		result = "нет пути в " + arrOfArgs[1]
	}
	return result
}

// commUse - обработка команды "применить"
func commUse(result string, arrOfArgs []string) string {
	if len(arrOfArgs) == 1 || len(arrOfArgs) == 2 {
		result = WRONG
	}
	object := arrOfArgs[1]
	toWhich := arrOfArgs[2]
	detectedObject := 0
	for _, obj := range player.inventory {
		if object == obj.title {
			detectedObject = 1
			if obj.toUse {
				if toWhich == obj.contact {
					result = obj.result
					if object == "ключи" {
						for _, room := range player.actualRoom.ways {
							room.isOpen = true
						}
					}
				} else {
					result = "не к чему применить"
				}
			} else {
				result = "нельзя использовать: " + obj.title
			}
			break
		}
	}
	if detectedObject == 0 {
		result = "нет предмета в инвентаре - " + object
	}
	return result
}

// commTake - обработка команды "взять"
func commTake(result string, arrOfArgs []string) string {
	if len(arrOfArgs) == 1 {
		result = WRONG
	}
	object := arrOfArgs[1]
	detectedObject := 0
	for obj, count := range player.actualRoom.objects {
		if object == obj.title && count > 0 {
			detectedObject = 1
			if obj.toTake {
				if !player.isWithBag {
					result = "некуда класть"
				} else {
					player.actualRoom.objects[obj]--
					result = "предмет добавлен в инвентарь: " + obj.title
					player.inventory = append(player.inventory, obj)
				}
			} else {
				result = "нельзя взять: " + obj.title
			}
			break
		}
	}
	if detectedObject == 0 {
		result = "нет такого"
	}
	return result
}

// commWear - обработка команды "надеть"
func commWear(result string, arrOfArgs []string) string {
	if len(arrOfArgs) == 1 {
		result = WRONG
	}
	object := arrOfArgs[1]
	detectedObject := 0
	for obj, count := range player.actualRoom.objects {
		if object == obj.title && count > 0 {
			detectedObject = 1
			if obj.toWear {
				if obj.title == "рюкзак" {
					player.isWithBag = true
				}
				result = "вы надели: " + obj.title
				player.actualRoom.objects[obj]--
				player.inventory = append(player.inventory, obj)
			} else {
				result = "нельзя надеть: " + obj.title
			}
			break
		}
	}
	if detectedObject == 0 {
		result = "нет такого предмета: " + arrOfArgs[1]
	}
	return result
}
