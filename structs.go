package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	bareHands       = "Bare Hands"
	knife           = "Knife"
	sword           = "Sword"
	ninjaku         = "ninjaku"
	wand            = "wand"
	gophermourne    = "gophermourne"
	healthPotion    = "health potion"
	strengthPotion  = "strength potion"
	agilityPotion   = "agility potion"
	intellectPotion = "intellect potion"
)

var BareHands = Weapon{[]int{1, 1}, bareHands, 0, 0, 0, 0}
var Knife = Weapon{[]int{2, 3}, knife, 10, 0, 0, 0}
var Sword = Weapon{[]int{3, 5}, sword, 35, 2, 0, 0}
var Ninjaku = Weapon{[]int{1, 7}, ninjaku, 25, 0, 0, 2}
var Wand = Weapon{[]int{3, 3}, wand, 30, 0, 2, 0}
var Gophermourne = Weapon{[]int{6, 7}, gophermourne, 65, 3, 2, 0}

var Health = Consumables{healthPotion, 0, 5, 0, 0, 0, 0}
var Strength = Consumables{strengthPotion, 3, 0, 2, 0, 0, 0}
var Intellect = Consumables{agilityPotion, 3, 0, 0, 2, 0, 0}
var Agility = Consumables{intellectPotion, 3, 0, 0, 0, 2, 0}

type Weapon struct {
	damage      []int
	weaponType  string
	gold        int
	strengthReq int
	intellReq   int
	agilityReq  int
}

type Consumables struct {
	consumableType string
	duration       int
	hpEffect       int
	strengthEffect int
	intelEffect    int
	agilityEffect  int
	startTime      int
}

type Gophers struct {
	name         string
	healthpoints int
	gold         int
	currTurn     int
	strength     int
	intel        int
	agility      int
	weapon       Weapon
	inventory    []Consumables
}

type actions interface {
	attack()
	buy()
	work()
	use()
	train()
	exit()
}

func validate(itemType string, g *Gophers) bool {
	if itemType == knife {
		return g.gold >= 10
	} else if itemType == sword {
		return (g.gold >= 35 && g.strength >= 2)
	} else if itemType == ninjaku {
		return (g.gold >= 25 && g.agility >= 2)
	} else if itemType == wand {
		return (g.gold >= 30 && g.intel >= 2)
	} else if itemType == gophermourne {
		return (g.gold >= 65 && g.strength >= 2 && g.intel >= 2)
	} else if itemType == healthPotion {
		return g.gold >= 5
	} else if itemType == strengthPotion {
		return g.gold >= 10
	} else if itemType == agilityPotion {
		return g.gold >= 10
	} else if itemType == intellectPotion {
		return g.gold >= 10
	}
	fmt.Println("You did not buy anything! Please try again")
	return false
}

func (g *Gophers) buy(item interface{}, itemType string) {
	if !validate(itemType, g) {
		fmt.Println("Insufficient funds! Please try again")
		return
	}
	switch item.(type) {
	case Weapon:
		fmt.Println("You bought a weapon")
		g.weapon = item.(Weapon)
	case Consumables:
		fmt.Println("You bought a consumable")
		g.inventory = append(g.inventory, item.(Consumables))
	default:
		fmt.Println("You did not buy anything! Please try again")
	}
}

func (g1 *Gophers) attack(g2 *Gophers) {
	fmt.Println("You attacked", g2.name)
	rand.NewSource(time.Now().UnixNano())
	damage := rand.Intn(g1.weapon.damage[1]-g1.weapon.damage[0]+1) + g1.weapon.damage[0]
	g2.healthpoints -= damage
	fmt.Println(g2.name, "now has", g2.healthpoints, "health")
	if g2.healthpoints <= 0 {
		fmt.Println(g2.name, "has died and "+g1.name+" has won the game!")
		os.Exit(0)
	}
}

func (g1 *Gophers) work() {
	fmt.Println("You worked")
	rand.NewSource(time.Now().UnixNano())
	coins := rand.Intn(11) + 5
	g1.gold += coins
}

func (g *Gophers) train(stat string) {
	if g.gold < 5 {
		fmt.Println("You do not have enough gold to train")
		return
	}
	switch stat {
	case "strength":
		g.strength += 2
	case "intel":
		g.intel += 2
	case "agility":
		g.agility += 2
	default:
		fmt.Println("You did not train anything!")
		return
	}
	g.gold -= 5
}

func (g *Gophers) use(itemType string) {
	for i, item := range g.inventory {
		if item.consumableType == itemType {
			if itemType == healthPotion {
				if g.healthpoints == 30 {
					fmt.Println("You are already at max health")
					return
				} else {
					g.healthpoints += 5
				}
			} else if itemType == strengthPotion {
				g.strength += item.strengthEffect
			} else if itemType == agilityPotion {
				g.agility += item.agilityEffect
			} else if itemType == intellectPotion {
				g.intel += item.intelEffect
			}
			fmt.Println("You used", item.consumableType)
			g.inventory = append(g.inventory[:i], g.inventory[i+1:]...)
			return
		}
	}
	fmt.Println("You do not have any", itemType, "in your inventory")
	return
}

func (g *Gophers) removeExpiredConsumables() {
	for i, item := range g.inventory {
		if item.duration > 0 && item.startTime+item.duration < g.currTurn {
			if item.consumableType == strengthPotion {
				g.strength -= item.strengthEffect
			} else if item.consumableType == agilityPotion {
				g.agility -= item.agilityEffect
			} else if item.consumableType == intellectPotion {
				g.intel -= item.intelEffect
			}
			g.inventory = append(g.inventory[:i], g.inventory[i+1:]...)
		}
	}
}

func displayList() {
	fmt.Println("1. Attack")
	fmt.Println("2. Buy")
	fmt.Println("3. Work")
	fmt.Println("4. Use")
	fmt.Println("5. Train")
	fmt.Println("6. Exit")
}

func buyOptions() {
	fmt.Println("1. Weapons")
	fmt.Println("2. Consumables")
}

func chooseWeapons() {
	fmt.Println("1. Knife")
	fmt.Println("2. Sword")
	fmt.Println("3. Ninjaku")
	fmt.Println("4. Wand")
	fmt.Println("5. Gophermourne")
}

func chooseConsumables() {
	fmt.Println("1. Health Potion")
	fmt.Println("2. Strength Potion")
	fmt.Println("3. Agility Potion")
	fmt.Println("4. Intellect Potion")
}

func chooseTrainOptions() {
	fmt.Println("1. Strength")
	fmt.Println("2. Intellect")
	fmt.Println("3. Agility")
}

func (g1 *Gophers) game(g2 *Gophers) {

	g1.removeExpiredConsumables()
	fmt.Println(g1.name + "What would you like to do?")
	displayList()
	var choice int
	fmt.Scanln(&choice)
	switch choice {
	case 1:
		g1.attack(g2)
	case 2:
		fmt.Println("What would you like to buy?")
		buyOptions()
		var buyOption int
		fmt.Scanln(&buyOption)
		switch buyOption {
		case 1:
			fmt.Println("What weapon would you like to buy?")
			chooseWeapons()
			var weaponOption int
			fmt.Scanln(&weaponOption)
			switch weaponOption {
			case 1:
				g1.buy(Knife, knife)
			case 2:
				g1.buy(Sword, sword)
			case 3:
				g1.buy(Ninjaku, ninjaku)
			case 4:
				g1.buy(Wand, wand)
			case 5:
				g1.buy(Gophermourne, gophermourne)
			default:
				fmt.Println("You did not buy any weapon")
			}
		case 2:
			fmt.Println("What consumable would you like to buy?")
			chooseConsumables()
			var consumableOption int
			fmt.Scanln(&consumableOption)
			switch consumableOption {
			case 1:
				g1.buy(Health, healthPotion)
			case 2:
				g1.buy(Strength, strengthPotion)
			case 3:
				g1.buy(Agility, agilityPotion)
			case 4:
				g1.buy(Intellect, intellectPotion)
			default:
				fmt.Println("You did not buy any consumable")
			}
		default:
			fmt.Println("Invalid Choice")
		}
	case 3:
		g1.work()
	case 4:
		fmt.Println("What consumable would you like to use?")
		chooseConsumables()
		var consumableOption int
		fmt.Scanln(&consumableOption)
		switch consumableOption {
		case 1:
			g1.use(healthPotion)
		case 2:
			g1.use(strengthPotion)
		case 3:
			g1.use(agilityPotion)
		case 4:
			g1.use(intellectPotion)
		default:
			fmt.Println("Invalid Choice")
		}
	case 5:
		fmt.Println("What stat would you like to train?")
		chooseTrainOptions()
		var trainOption int
		fmt.Scanln(&trainOption)
		switch trainOption {
		case 1:
			g1.train("strength")
		case 2:
			g1.train("intel")
		case 3:
			g1.train("agility")
		default:
			fmt.Println("Invalid Choice")
		}
	case 6:
		fmt.Println(g1.name + "have exited the game")
		fmt.Println(g2.name + "has won the game!")
		os.Exit(0)
	default:
		fmt.Println("Invalid Choice")
	}

}

func main() {
	fmt.Println("Welcome to Archaemania")

	g1 := Gophers{"Gopher1", 30, 20, 0, 0, 0, 0, BareHands, []Consumables{}}
	g2 := Gophers{"Gopher2", 30, 20, 0, 0, 0, 0, BareHands, []Consumables{}}
	fmt.Printf("Gopher1: %v\n", g1)
	for true {
		g1.game(&g2)
		g2.game(&g1)
		g1.currTurn++
		g2.currTurn++

	}
}
