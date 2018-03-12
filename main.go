package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func get_next_timeout(last_timeout time.Duration) time.Duration {
	var backoff_factor int64 = 2
	return time.Duration(backoff_factor * int64(last_timeout))
}

func perform_get_request_on_service(client http.Client, service_address string) *http.Response {
	rs, err := client.Get(service_address)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.Println(err)
	}
	return rs
}

func request_get_request_on_service_with_backoff(service_address string, max_num_backoffs int) *http.Response {
	var successful bool = false
	var cur_timeout time.Duration = time.Second
	var response *http.Response = nil
	var client http.Client = http.Client{
		Timeout: cur_timeout,
	}
	for i := 0; i < max_num_backoffs || successful; i++ {
		log.Printf("Performing GET request on %s\n", service_address)
		response = perform_get_request_on_service(client, service_address)
		if response != nil && (response.StatusCode == 200) {
			log.SetOutput(os.Stdout)
			log.Println("Successful Response.")
			successful = true
			return response
		} else {
			log.SetOutput(os.Stdout)
			log.Printf("Timed out after %d seconds.", int(cur_timeout.Seconds()))
			cur_timeout = get_next_timeout(cur_timeout)
			log.Printf("Trying again with a timeout of %d seconds.\n", int(cur_timeout.Seconds()))
			client = http.Client{
				Timeout: cur_timeout,
			}
		}
	}
	return response
}

func main() {
	request_get_request_on_service_with_backoff("https://httpbin.org/delay/3", 3)
}
