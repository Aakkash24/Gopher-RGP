package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

// Constants for item types
const (
	BareHands         = "Bare Hands"
	Knife             = "Knife"
	Sword             = "Sword"
	Ninjaku           = "Ninjaku"
	Wand              = "Wand"
	Gophermourne      = "Gophermourne"
	HealthPotion      = "Health Potion"
	StrengthPotion    = "Strength Potion"
	AgilityPotion     = "Agility Potion"
	IntellectPotion   = "Intellect Potion"
	MaxHealth         = 30
	WorkMinGold       = 5
	WorkMaxGold       = 15
	MaxStrength       = 10
	MaxIntellect      = 10
	MaxAgility        = 10
	TrainCost         = 5
	InitialGoldCoins  = 20
	InitialAttributes = 0
)

// Constants for consumable effects
const (
	HealthPotionHP        = 5
	StrengthPotionEffect  = 2
	AgilityPotionEffect   = 2
	IntellectPotionEffect = 2
)

// Weapon represents a weapon in the game
type Weapon struct {
	damage      []int
	weaponType  string
	gold        int
	strengthReq int
	intellReq   int
	agilityReq  int
}

// Consumables represent consumable items in the game
type Consumables struct {
	consumableType string
	duration       int
	hpEffect       int
	strengthEffect int
	intelEffect    int
	agilityEffect  int
	startTime      int
}

// Gopher represents a player (or gopher) in the game
type Gopher struct {
	name            string
	healthpoints    int
	gold            int
	currTurn        int
	strength        int
	intel           int
	agility         int
	weapon          Weapon
	inventory       []Consumables
	currentPortions []Consumables
}

// Constants for item requirements
var itemRequirements = map[string]struct {
	gold int
	reqs map[string]int
}{
	Knife: {
		gold: 10,
	},
	Sword: {
		gold: 35,
		reqs: map[string]int{"strength": 2},
	},
	Ninjaku: {
		gold: 25,
		reqs: map[string]int{"agility": 2},
	},
	Wand: {
		gold: 30,
		reqs: map[string]int{"intel": 2},
	},
	Gophermourne: {
		gold: 65,
		reqs: map[string]int{"strength": 2, "intel": 2},
	},
	HealthPotion: {
		gold: 5,
	},
	StrengthPotion: {
		gold: 10,
	},
	AgilityPotion: {
		gold: 10,
	},
	IntellectPotion: {
		gold: 10,
	},
}

// Helper function to check requirements for a specific item type
func hasRequirements(g *Gopher, itemType string) bool {
	reqs, exists := itemRequirements[itemType]
	if !exists {
		fmt.Println("Invalid item type!")
		return false
	}

	for stat, value := range reqs.reqs {
		switch stat {
		case "strength":
			if g.strength < value {
				fmt.Printf("Insufficient %s to buy %s\n", stat, itemType)
				return false
			}
		case "intel":
			if g.intel < value {
				fmt.Printf("Insufficient %s to buy %s\n", stat, itemType)
				return false
			}
		case "agility":
			if g.agility < value {
				fmt.Printf("Insufficient %s to buy %s\n", stat, itemType)
				return false
			}
		default:
			fmt.Println("Unknown requirement type!")
			return false
		}
	}

	if g.gold < reqs.gold {
		fmt.Printf("Insufficient gold to buy %s\n", itemType)
		return false
	}

	return true
}

// Function to handle buying items
func (g *Gopher) buy(itemType string) {
	if !hasRequirements(g, itemType) {
		fmt.Println("Purchase failed!")
		return
	}

	switch itemType {
	case Knife:
		g.weapon = Weapon{[]int{2, 3}, Knife, 10, 0, 0, 0}
	case Sword:
		g.weapon = Weapon{[]int{3, 5}, Sword, 35, 2, 0, 0}
	case Ninjaku:
		g.weapon = Weapon{[]int{1, 7}, Ninjaku, 25, 0, 0, 2}
	case Wand:
		g.weapon = Weapon{[]int{3, 3}, Wand, 30, 0, 2, 0}
	case Gophermourne:
		g.weapon = Weapon{[]int{6, 7}, Gophermourne, 65, 3, 2, 0}
	case HealthPotion:
		g.inventory = append(g.inventory, Consumables{itemType, 0, HealthPotionHP, 0, 0, 0, -1})
	case StrengthPotion:
		g.inventory = append(g.inventory, Consumables{itemType, 3, 0, 2, 0, 0, -1})
	case AgilityPotion:
		g.inventory = append(g.inventory, Consumables{itemType, 3, 0, 0, 0, 2, -1})
	case IntellectPotion:
		g.inventory = append(g.inventory, Consumables{itemType, 3, 0, 0, 2, 0, -1})
	}

	g.gold -= itemRequirements[itemType].gold
	fmt.Printf("You bought a %s\n", itemType)
}

// Function to handle attacking
func (g1 *Gopher) attack(g2 *Gopher) {
	fmt.Println("You attacked", g2.name)
	rand.NewSource(time.Now().UnixNano())
	if len(g1.weapon.damage) == 1 {
		g2.healthpoints -= g1.weapon.damage[0]
	} else {
		damage := rand.Intn(g1.weapon.damage[1]-g1.weapon.damage[0]+1) + g1.weapon.damage[0]
		g2.healthpoints -= damage
	}
	if g2.healthpoints < 0 {
		g2.healthpoints = 0
	}

	fmt.Println(g2.name, "now has", g2.healthpoints, "health")

	if g2.healthpoints <= 0 {
		fmt.Println(g2.name, "has died and", g1.name, "has won the game!")
		os.Exit(0)
	}
}

// Function to handle working
func (g *Gopher) work() {
	rand.NewSource(time.Now().UnixNano())
	coins := rand.Intn(WorkMaxGold-WorkMinGold+1) + WorkMinGold
	g.gold += coins
	fmt.Printf("You worked and earned %d gold\n", coins)
}

// Function to handle using consumables
func (g *Gopher) use(itemType string) {
	for i, item := range g.inventory {
		if item.consumableType == itemType {
			switch itemType {
			case HealthPotion:
				if g.healthpoints == MaxHealth {
					fmt.Println("You are already at max health")
				} else {
					g.healthpoints += HealthPotionHP
					if g.healthpoints > MaxHealth {
						g.healthpoints = MaxHealth
					}
				}
			case StrengthPotion:
				g.strength += StrengthPotionEffect
			case AgilityPotion:
				g.agility += AgilityPotionEffect
			case IntellectPotion:
				g.intel += IntellectPotionEffect
			}

			fmt.Println("You used", item.consumableType)
			g.inventory[i].startTime = g.currTurn
			temp := g.inventory[i]
			temp.startTime = g.currTurn
			g.currentPortions = append(g.currentPortions, temp)
			g.inventory = append(g.inventory[:i], g.inventory[i+1:]...)
			return
		}
	}
	fmt.Println("You do not have any", itemType, "in your inventory")
}

func (g *Gopher) removePortionFromInventory(item string) {
	var temp []Consumables
	for _, ele := range g.currentPortions {
		if ele.consumableType == item {
			continue
		}
		temp = append(temp, ele)
	}
	g.currentPortions = temp
}

// Function to remove expired consumables
func (g *Gopher) removeExpiredConsumables() {
	currentTime := g.currTurn
	var updatedInventory []Consumables
	for _, item := range g.currentPortions {
		if item.duration == 0 || item.startTime == -1 || currentTime-item.startTime != item.duration-1 {
			updatedInventory = append(updatedInventory, item)
		} else {
			switch item.consumableType {
			case StrengthPotion:
				g.strength -= item.strengthEffect
				g.removePortionFromInventory(StrengthPotion)
			case AgilityPotion:
				g.agility -= item.agilityEffect
				g.removePortionFromInventory(AgilityPotion)
			case IntellectPotion:
				g.intel -= item.intelEffect
				g.removePortionFromInventory(IntellectPotion)
			}
		}
	}

	g.currentPortions = updatedInventory
}

// Function to handle training a specific stat
func (g *Gopher) train(stat string) {
	if g.gold < TrainCost {
		fmt.Println("You do not have enough gold to train")
		return
	}

	switch stat {
	case "strength":
		if g.strength+StrengthPotionEffect < MaxStrength {
			g.strength += StrengthPotionEffect
			g.gold -= TrainCost
			fmt.Println("You trained strength.")
		} else {
			fmt.Println("Strength is already at maximum.")
		}
	case "intel":
		if g.intel+IntellectPotionEffect < MaxIntellect {
			g.intel += IntellectPotionEffect
			g.gold -= TrainCost
			fmt.Println("You trained intellect.")
		} else {
			fmt.Println("Intellect is already at maximum.")
		}
	case "agility":
		if g.agility+AgilityPotionEffect < MaxAgility {
			g.agility += AgilityPotionEffect
			g.gold -= TrainCost
			fmt.Println("You trained agility.")
		} else {
			fmt.Println("Agility is already at maximum.")
		}
	default:
		fmt.Println("You did not train anything!")
		return
	}
}

func chooseTrainOptions() {
	fmt.Println("1. Strength")
	fmt.Println("2. Intellect")
	fmt.Println("3. Agility")
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

func checkInventory(item string, array []Consumables) bool {
	for _, ele := range array {
		fmt.Println(ele.consumableType)
		if ele.consumableType == item {
			return true
		}
	}
	return false
}

// Function to handle the game loop and player actions
func (g *Gopher) game(opponent *Gopher) {
	g.removeExpiredConsumables()

	fmt.Printf("%s, what would you like to do?\n", g.name)
	displayList()

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		g.attack(opponent)
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
				g.buy(Knife)
			case 2:
				g.buy(Sword)
			case 3:
				g.buy(Ninjaku)
			case 4:
				g.buy(Wand)
			case 5:
				g.buy(Gophermourne)
			default:
				fmt.Println("You did not buy any weapon")
				g.game(opponent)
			}
		case 2:
			fmt.Println("What consumable would you like to buy?")
			chooseConsumables()
			var consumableOption int
			fmt.Scanln(&consumableOption)
			switch consumableOption {
			case 1:
				g.buy(HealthPotion)
			case 2:
				g.buy(StrengthPotion)
			case 3:
				g.buy(AgilityPotion)
			case 4:
				g.buy(IntellectPotion)
			default:
				fmt.Println("You did not buy any consumable")
			}
		default:
			fmt.Println("Invalid Choice")
			g.game(opponent)
		}
	case 3:
		g.work()
	case 4:
		fmt.Println("What consumable would you like to use?")
		chooseConsumables()
		var consumableOption int
		fmt.Scanln(&consumableOption)
		switch consumableOption {
		case 1:
			g.use(HealthPotion)
		case 2:
			g.use(StrengthPotion)
		case 3:
			g.use(AgilityPotion)
		case 4:
			g.use(IntellectPotion)
		default:
			fmt.Println("Invalid Choice")
			g.game(opponent)
		}
	case 5:
		fmt.Println("What stat would you like to train?")
		chooseTrainOptions()
		var trainOption int
		fmt.Scanln(&trainOption)
		switch trainOption {
		case 1:
			g.train("strength")
		case 2:
			g.train("intel")
		case 3:
			g.train("agility")
		default:
			fmt.Println("Invalid Choice")
			g.game(opponent)
		}
	case 6:
		fmt.Printf("%s has exited the game\n", g.name)
		fmt.Printf("%s has won the game!\n", opponent.name)
		os.Exit(0)
	default:
		fmt.Println("Invalid Choice")
		g.game(opponent)
	}
}

func printDetails(g Gopher) {
	fmt.Println("\nName:", g.name)
	fmt.Println("Health:", g.healthpoints)
	fmt.Println("Gold:", g.gold)
	fmt.Println("Agility:", g.agility)
	fmt.Println("Strength:", g.strength)
	fmt.Println("Intelligence:", g.intel)
	fmt.Println("Weapon:", g.weapon.weaponType)
	fmt.Println("Consumables:", g.inventory)
}

func clearConsole() {
	time.Sleep(2 * time.Second)
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	fmt.Printf("Welcome to Archaemania\n\n\n")

	g1 := Gopher{"Gopher1", MaxHealth, InitialGoldCoins, InitialAttributes, InitialAttributes, InitialAttributes, InitialAttributes, Weapon{[]int{1}, BareHands, 0, 0, 0, 0}, []Consumables{}, []Consumables{}}
	g2 := Gopher{"Gopher2", MaxHealth, InitialGoldCoins, InitialAttributes, InitialAttributes, InitialAttributes, InitialAttributes, Weapon{[]int{1}, BareHands, 0, 0, 0, 0}, []Consumables{}, []Consumables{}}

	for {
		clearConsole()
		fmt.Println("\nTurn:", g1.currTurn)
		printDetails(g1)
		g1.game(&g2)
		printDetails(g2)
		g2.game(&g1)
		g1.currTurn++
		g2.currTurn++
	}
}
