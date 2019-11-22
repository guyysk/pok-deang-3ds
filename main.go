package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type person struct {
	name      string
	cards     []int
	signCards []string
	money     int
	score     int
	indicator int
}

func getCardsScore(cards []int) int {
	var sum int
	for i := 0; i < len(cards); i++ {
		var realScoreCard int
		if cards[i] > 40 {
			realScoreCard = 0
		} else {
			realScoreCard = cards[i] / 4
		}
		sum = sum + realScoreCard
	}
	return sum % 10
}

func remove(f *[]int, i int) []int {
	var s = *f
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	*f = s[:len(s)-1]
	return *f
}

func shuffleCard(cards *[]int) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	elementNo := r1.Intn(len(*cards))
	card := (*cards)[elementNo]
	remove(cards, elementNo)
	return card
}

func updateIndicator(person *person) {
	isSameSymbol := true
	cards := (*person).cards
	for i := 1; i < len(cards); i++ {
		isSameSymbol = isSameSymbol && (cards[i-1]%4 == cards[i]%4)
	}
	if isSameSymbol {
		(*person).indicator = len(cards)
	}
}

func (person *person) firstShuffleWithScore(cards *[]int) {
	(*person).cards = append((*person).cards, shuffleCard(cards))
	(*person).cards = append((*person).cards, shuffleCard(cards))
	(*person).score = getCardsScore((*person).cards)
	if (*person).score <= 5 {
		(*person).cards = append((*person).cards, shuffleCard(cards))
		(*person).score = getCardsScore((*person).cards)
		updateIndicator(person)
		for i := 0; i < len((*person).cards); i++ {
			(*person).signCards = append((*person).signCards, getRealCardSign((*person).cards[i]))
		}
	} else {
		updateIndicator(person)
		for i := 0; i < len((*person).cards); i++ {
			(*person).signCards = append((*person).signCards, getRealCardSign((*person).cards[i]))
		}
	}
}

func compareScoreWithController(controller *person, person *person, competeMoney int) {

	fmt.Println("---------------------------------")

	if (*controller).score > (*person).score {
		deductMoney := (competeMoney * (*controller).indicator)
		(*person).money -= deductMoney
		(*controller).money += deductMoney
		fmt.Println("Controller cards: ", (*controller).signCards)
		fmt.Println((*person).name, " cards: ", (*person).signCards)
		fmt.Println("Controller won! x", (*controller).indicator)
		fmt.Println((*controller).name, ": ", (*controller).money, "  ", (*person).name, ": ", (*person).money)
	} else {
		deductMoney := (competeMoney * (*person).indicator)
		(*person).money += deductMoney
		(*controller).money -= deductMoney
		fmt.Println("Controller cards: ", (*controller).signCards)
		fmt.Println((*person).name, " cards: ", (*person).signCards)
		fmt.Println((*person).name, " won! x", (*person).indicator)
		fmt.Println((*controller).name, ": ", (*controller).money, "  ", (*person).name, ": ", (*person).money)
	}
}

func newCards() []int {
	var cards []int

	for i := 0; i < 52; i++ {
		cards = append(cards, i)
	}

	return cards
}

func newPlayers(number int) []person {
	var players []person
	for i := 0; i < number; i++ {
		newPlayers := person{
			name:      fmt.Sprintf("Player%d", i+1),
			indicator: 1,
			money:     100,
		}
		players = append(players, newPlayers)
	}
	return players
}

func newController() person {
	controller := person{
		name:      "Controller",
		indicator: 1,
		money:     100,
	}
	return controller
}

func getRealCardSign(num int) string {
	number := (num / 4) + 1
	symbolNo := num % 4
	var symbol string
	switch symbolNo {
	case 0:
		symbol = "♣"
	case 1:
		symbol = "♦"
	case 2:
		symbol = "♥"
	case 3:
		symbol = "♠"
	}
	if number > 10 {
		switch number {
		case 11:
			return "J" + symbol
		case 12:
			return "Q" + symbol
		case 13:
			return "K" + symbol
		}
	}
	i := strconv.Itoa(number)
	return i + symbol
}

func (person *person) emptyCards() {
	(*person).cards = nil
	(*person).signCards = nil
	(*person).indicator = 1
}

func main() {
	players := newPlayers(4)
	controller := newController()
	for {
		stop := false
		cards := newCards()
		controller.firstShuffleWithScore(&cards)
		for i := 0; i < len(players); i++ {
			players[i].firstShuffleWithScore(&cards)
			compareScoreWithController(&controller, &players[i], 10)
			players[i].emptyCards()
			if players[i].money <= 0 || controller.money <= 0 {
				stop = true
			}
		}
		controller.emptyCards()
		if stop {
			fmt.Println("---------------------------------")
			break
		}
	}
}
