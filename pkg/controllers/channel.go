package controllers

import (
	"fmt"
	"net/http"
	"sync"
)

func Channels(w http.ResponseWriter, r *http.Request) {
	urls := []string{
		"https://reqres.in/api/users?page=1",
		"https://reqres.in/api/users?page=2",
		"https://reqres.in/api/users?page=3",
	}

	var wg sync.WaitGroup
	responses := make(chan string)

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			response, err := http.Get(url)
			if err != nil {
				responses <- fmt.Sprintf("Error fetching %s: %s", url, err.Error())
				return
			}
			defer response.Body.Close()

			// Gunakan response.Body untuk memproses data response jika diperlukan

			responses <- fmt.Sprintf("Response from %s", url)
		}(url)
	}

	go func() {
		wg.Wait()
		close(responses)
	}()

	for response := range responses {
		fmt.Println(response)
	}
}
