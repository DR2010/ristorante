package main

import (
	"fmt"
	"strings"
)

// Curr does something
type Curr struct {
	Currency string
	Name     string
	Country  string
	Number   int
}

var currencies = []Curr{
	Curr{"DZD", "Algerian Dinar", "Algeria", 12},
	Curr{"AUD", "Australian Dollar", "Australia", 36},
	Curr{"HKD", "Hong Kong Dollar", "Hong Koong", 344},
	Curr{"USD", "US Dollar", "United States", 840},
}

func isDollar(curr Curr) bool {
	var result bool

	switch curr {
	default:
		result = false

	case Curr{"AUD", "Australian Dollar", "Australia", 36}:
		result = true

	case Curr{"HKD", "Hong Kong Dollar", "Hong Koong", 344}:
		result = true

	case Curr{"USD", "US Dollar", "United States", 840}:
		result = true

	case Curr{"DZD", "Algerian Dinar", "Algeria", 12}:
		result = false

	}

	return result
}

func isDollar2(curr Curr) bool {

	var ret bool

	dollars := []Curr{currencies[0], currencies[1], currencies[2], currencies[3]}

	switch curr {
	default:
		ret = false

	case dollars[0]:
		fallthrough

	case dollars[1]:
		ret = true
		var t = "dollars[1]"
		fmt.Println(" t is ", t, " dol= ", dollars[1])
		break

	case dollars[2]:
		var t = "dollars[2]"
		fmt.Println(" t is ", t, " dol= ", dollars[2])
		fallthrough

	case dollars[3]:
		ret = true

	}

	return ret
}

func main() {

	curr := Curr{"HKD", "Hong Kong Dollar", "Hong Koong", 344}

	if isDollar(curr) {
		fmt.Printf("%+v is Euro currency \n", curr)
	} else {
		fmt.Println("Currency is not dollar or Euro")
	}

	dol := Curr{"AUD", "Australian Dollar", "Australia", 36}

	if isDollar2(dol) {
		fmt.Println("Dollar currency found \n", dol)

	}
}

func find(name string) {
	for i := 0; i < 10; i++ {
		c := currencies[i]
		switch {
		case strings.Contains(c.Currency, name),
			strings.Contains(c.Name, name),
			strings.Contains(c.Country, name):
			fmt.Println("Found:", c)
		}
	}
}
