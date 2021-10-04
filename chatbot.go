package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//struct of type pairList
type pairList struct {
	Pairs []pair `json:"pairs"`
}

//struct of type pair
type pair struct {
	Question string `json: "question"`
	Response string `json: "response"`
}

func main() {
	//reading json file having chat pairs
	file, _ := ioutil.ReadFile("conversation.json")

	data := pairList{}
	//coverting json file to struct type
	_ = json.Unmarshal([]byte(file), &data)

	for {
		//reading chat input
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("YOU :")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "Quit" {
			break
		}
		for i := 0; i < len(data.Pairs); i++ {
			if strings.Contains(input, data.Pairs[i].Question) {
				fmt.Print("BOT_KMB :")
				//printing response for the input chat
				fmt.Println(data.Pairs[i].Response)
			}
		}
	}

}
