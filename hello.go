package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	//declaring a new channel
	file, err := os.Create("result.txt")

	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	c := make(chan string, 2)

	go generator1(c)
	go generator2(c)

	for msg := range c {
		fmt.Println(msg, " - received from channel, it will be written to text file.")
		write(msg, file)
	}

	//this is to close the channel when the sender is done.
	close(c)

	//press enter to exit
	fmt.Scanln()
}

func generator1(c chan string) {
	for i := 1; true; i++ {
		var message string = strconv.Itoa(i) + " sheep"
		fmt.Println("Generator1 = ", message)
		c <- message
		time.Sleep(1 * time.Second)
	}
}

func generator2(c chan string) {
	for i := 1; true; i++ {
		resp, err := http.Get("https://random-data-api.com/api/v2/users?size=1&response_type=json")
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		type Response struct {
			Name string `json:"first_name"`
		}

		var responseObject Response
		json.Unmarshal(responseData, &responseObject)

		fmt.Println("Generator2 = ", responseObject.Name)
		time.Sleep(1 * time.Second)

		c <- responseObject.Name
	}
}

func write(data string, file *os.File) {
	data = data + "\n"
	fmt.Fprintf(file, data)
}
