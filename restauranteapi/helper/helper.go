package helper

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/go-redis/redis"
)

// DatabaseX is a struct
type DatabaseX struct {
	Location   string // location of the database localhost, something.com, etc
	Database   string // database name
	Collection string // collection name
}

// Resultado is a struct
type Resultado struct {
	ErrorCode        string // error code
	ErrorDescription string // description
	IsSuccessful     string // Y or N
	ReturnedValue    string // Any string
}

// GetRedisPointer returns
func GetRedisPointer() *redis.Client {

	redisclient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return redisclient
}

// RestEnvVariables = restaurante environment variables
//
type RestEnvVariables struct {
	APIMongoDBLocation    string // location of the database localhost, something.com, etc
	APIMongoDBDatabase    string // database name
	APIAPIServerPort      string // collection name
	APIAPIServerIPAddress string // apiserver name
	WEBDebug              string // debug
	CollectionOrders      string // Collection Names
	CollectionSecurity    string // Collection Names
	CollectionDishes      string // Collection Names
	CollectionEvents      string // Collection Names
}

// Readfileintostruct is
func Readfileintostruct() RestEnvVariables {
	dat, err := ioutil.ReadFile("restauranteapi.ini")
	check(err)
	fmt.Print(string(dat))

	var restenv RestEnvVariables

	json.Unmarshal(dat, &restenv)

	return restenv
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type PlayerRegistrationFile struct {
	FFA  string
	Name string
	DOB  string
}

// Capitalfootball is
func Capitalfootball(redisclient *redis.Client) []PlayerRegistrationFile {

	file, err := os.Open("capitalfootball.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var playerlist []PlayerRegistrationFile

	scanner := bufio.NewScanner(file)

	playerlist = make([]PlayerRegistrationFile, 52)

	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(scanner.Text())

		tmp := strings.Split(line, ",")

		i++
		playerlist[i] = PlayerRegistrationFile{}
		playerlist[i].FFA = strings.Trim(tmp[0], " ")
		playerlist[i].Name = strings.Trim(tmp[1], " ")
		playerlist[i].DOB = strings.Trim(tmp[2], " ")

		fmt.Println(playerlist[i].FFA)

		err = redisclient.Set(playerlist[i].FFA, playerlist[i].Name, 0).Err()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return playerlist
}
