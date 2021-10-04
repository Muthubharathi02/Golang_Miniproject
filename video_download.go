package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/kkdai/youtube/v2"
	"github.com/tebeka/selenium"
)

//connecting to chrome driver
const (
	seleniumPath = "/Users/muthubharathi/Downloads/chromedriver"
	port         = 4444
)

func main() {
	//Enable selenium service
	ops := []selenium.ServiceOption{}
	service, err := selenium.NewChromeDriverService(seleniumPath, port, ops...)
	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v", err)
	}
	//Delay service shutdown
	defer service.Stop()

	//Call browser
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		panic(err)
	}
	//Delay exiting chrome
	defer wd.Quit()
	//loading website
	if err := wd.Get("https://medium.datadriveninvestor.com/build-your-own-chat-bot-using-python-95fdaaed620f"); err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Second)

	//Finding <a> tag webelement objects
	element, err := wd.FindElements(selenium.ByCSSSelector, "a.bv.jj")
	//getting youtube link from element object "href" attribute
	link, _ := element[1].GetAttribute("href")

	if err != nil {
		panic(err)
	}

	videoID := link
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
	}

	stream, _, err := client.GetStream(video, &video.Formats[0])
	fmt.Println(stream)
	if err != nil {
		panic(err)
	}
	//creating mp4 file
	file, err := os.Create("Chatbot_video.mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//downloading video content
	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}
	fmt.Println("Downloaded /Chatbot_video.mp4")
	time.Sleep(10 * time.Second)
}
