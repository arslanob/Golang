package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	fmt.Println("****** Starting generators, press enter to stop ******")
	//declaring a new channel
	file, err := os.Create("result.txt")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	var wg sync.WaitGroup
	wg.Add(3)

	c := make(chan string, 3)
	done := make(chan string)

	go generator1(c, done, &wg)
	go generator2(c, done, &wg)
	go write(c, done, file, &wg)

	go func() {
		defer close(done)
		fmt.Scanln()
	}()

	wg.Wait()

	fmt.Println("all done!")
}

func generator1(c chan string, done chan string, wg *sync.WaitGroup) {
	for i := 1; true; i++ {
		select {
		case <-done:
			return
		default:
			var message string = strconv.Itoa(i) + " sheep"
			fmt.Println("Generator1 = ", message)
			c <- message
			time.Sleep(1 * time.Second)
		}
	}
	defer wg.Done()
}

func generator2(c chan string, done chan string, wg *sync.WaitGroup) {
	for {
		select {
		case <-done:
			return
		default:
			resp, err := http.Get("https://random-data-api.com/api/v2/users?size=1&response_type=json")
			if err != nil {
				fmt.Printf("Cannot Get names: %v", err)
				return
			}
			defer resp.Body.Close()

			responseData, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Cannot read names: %v", err)
				return
			}

			type Response struct {
				Name string `json:"first_name"`
			}

			var responseObject Response
			err = json.Unmarshal(responseData, &responseObject)
			if err != nil {
				fmt.Printf("Cannot unmarshal names: %v", err)
				return
			}

			fmt.Println("Generator2 = ", responseObject.Name)
			c <- responseObject.Name

			// wait so we don't exhaust CPU or make too many GET calls
			time.Sleep(1 * time.Second)
		}
	}
	defer wg.Done()
}

func write(c chan string, done chan string, file *os.File, wg *sync.WaitGroup) {

	for {
		select {
		case <-done:
			return
		case msg := <-c:
			fmt.Println(msg, " - received from channel, it will be written to text file.")
			fmt.Fprintf(file, msg+"\n")
		}
	}

	defer wg.Done()
}
