package main

import "fmt"

var name, desc string
var radius int32
var mass float64
var active bool
var satellites []string

type currency struct {
	Name    string
	Country string
	Number  int
}

func main() {

	if "a" == "ab" {
		fmt.Println("Hello")
	}

	const SUNRADIUS int32 = 685800
	const SUNMASS float64 = 1.989E+30

	var minhavariable = " New variable mine "
	var one, two = "one", "two"

	var anothermass = SUNRADIUS

	name = "Sun"
	desc = "Star"
	radius = 685800
	mass = 1.989E+30
	active = true
	satellites = []string{
		"Mercury",
		"Venus",
		"Earth",
		"Mars",
		"Jupiter",
		"Saturn",
		"Uranus",
		"Neptune",
	}

	fmt.Println(name)
	fmt.Println(desc)
	fmt.Println("Radius (km)", radius)
	fmt.Println("Mass (kg)", mass)
	fmt.Println("Satellites", satellites)
	fmt.Println("minha variable", minhavariable)
	fmt.Println("one:", one, "two:", two)
	fmt.Println("Sun Radius (const):", SUNRADIUS)
	fmt.Println("Sun Mass (const):", SUNMASS)
	fmt.Println("Sun Mass (const):", anothermass)

	const (
		StarHyperGiant = 2 * iota
		StarSuperGiant
		StarBrightGiant
		StarGiant
		StarSubGiant
		StarDwarf
		StarSubDwarf
		StarWhiteDwarf
	)

	fmt.Println("StarBrightGiant: ", StarHyperGiant)
	fmt.Println("StarBrightGiant: ", StarWhiteDwarf)
	fmt.Println("StarBrightGiant: ", StarGiant)

	var CAD = currency{
		Name:    "Cannadian Dollar",
		Country: "Canada",
		Number:  124}

	var AUD = currency{
		Name:    "Australian Dollar",
		Country: "Canada",
		Number:  125}

	var BRL = currency{
		Name:    "Brazilian Real",
		Country: "Brazil",
		Number:  127}

	fmt.Println("Currency Canada: ", CAD)
	fmt.Println("Currency Australia: ", AUD)
	fmt.Println("Currency Brazil: ", BRL)

	var currencies = []currency{
		currency{"DZD", "Algerian Dinar", 123},
		currency{"AUD", "Australian Dolar", 124},
		currency{"USD", "United States Dolar", 125},
	}

	fmt.Println("Currencies: ", currencies)

}

func isDollar(currency Curr) bool {
	
}