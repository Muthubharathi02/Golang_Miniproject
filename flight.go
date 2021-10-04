package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

//connecting to chrome driver
const (
	seleniumPath = `/Users/muthubharathi/Downloads/chromedriver`
	port         = 4444
)

//struct of type flight
type flight struct {
	Name  string
	price int
}

type flight_list []flight

func (flights flight_list) Len() int {
	return len(flights)
}
func (flights flight_list) Swap(i, j int) {
	flights[i], flights[j] = flights[j], flights[i]
}

//function to sort data based on price
func (flights flight_list) Less(i, j int) bool {
	return flights[i].price < flights[j].price
}

func main() {
	flights := make(flight_list, 15)

	ops := []selenium.ServiceOption{}
	//Enabling selenium service
	service, err := selenium.NewChromeDriverService(seleniumPath, port, ops...)
	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v", err)
	}

	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	wd, err := selenium.NewRemote(caps, "")
	//Delaying chrome exit
	defer wd.Quit()
	if err != nil {
		panic(err)
	}
	//loading the website
	if err := wd.Get("https://www.kayak.co.in/flights/IXM-BLR/2021-10-13/2021-10-20?sort=price_a&attempt=1&lastms=1633026564376&force=true"); err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
	//getting flight names
	wes, err := wd.FindElements(selenium.ByCSSSelector, ".keel .Flights-Results-ExtraFlightInfo .section.codeshares")
	if err != nil {
		panic(err)
	}
	//getting flight price
	wep, err := wd.FindElements(selenium.ByCSSSelector, ".Flights-Results-FlightPriceSection.sleek .Theme-featured-large .price")
	if err != nil {
		panic(err)
	}

	//Loop to get information for each element
	for i, we := range wes {
		text, err := we.Text()
		text1, err1 := wep[i].Text()
		text1 = text1[4:]
		text1 = strings.ReplaceAll(text1, ",", "")
		text_, _ := strconv.Atoi(text1)
		if err != nil {
			panic(err)
		}
		if err1 != nil {
			panic(err)
		}

		flights[i] = flight{
			Name:  text,
			price: text_,
		}

	}
	//Delaying service shutdown
	defer service.Stop()
	//sorting data
	sort.Sort(flights)
	//printing result
	fmt.Println("The cheapest Flight available from Madurai to Bangalore :")
	fmt.Printf("Airline Name :%v  Price :%d", flights[0].Name, flights[0].price)
}
