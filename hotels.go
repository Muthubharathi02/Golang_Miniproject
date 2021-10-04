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
type hotel struct {
	Name  string
	price int
}
type hotel_list []hotel

func (hotels hotel_list) Len() int {
	return len(hotels)
}
func (hotels hotel_list) Swap(i, j int) {
	hotels[i], hotels[j] = hotels[j], hotels[i]
}

//function to sort data based on price
func (hotels hotel_list) Less(i, j int) bool {
	return hotels[i].price < hotels[j].price
}

func main() {
	hotels := make(hotel_list, 20)
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
	if err := wd.Get("https://www.kayak.co.in/hotels/Bengaluru,Karnataka,India-c14559/2021-10-13/2021-10-20/2adults?sort=rank_a"); err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)
	//getting hotel names
	wes, err := wd.FindElements(selenium.ByCSSSelector, "div.FLpo-big-name")
	if err != nil {
		panic(err)
	}
	//getting hotel price
	wep, err := wd.FindElements(selenium.ByCSSSelector, ".zV27-price-section")
	if err != nil {
		panic(err)
	}

	//Loop to get information for each element
	for i, we := range wes {
		text, err := we.Text()
		text1, err1 := wep[i].Text()
		text1 = text1[4:9]
		text1 = strings.ReplaceAll(text1, ",", "")
		text_, _ := strconv.Atoi(text1)

		if err != nil {
			panic(err)
		}
		if err1 != nil {
			panic(err)
		}

		hotels[i] = hotel{
			Name:  text,
			price: text_,
		}

	}
	//Delaying service shutdown
	defer service.Stop()
	//sorting data
	sort.Sort(hotels)
	//printing result
	fmt.Println("The cheapest Hotel available in Bangalore :")
	fmt.Printf("Hotel Name :%v  Price :%d", hotels[0].Name, hotels[0].price)
}
