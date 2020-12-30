package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

var url string

/*
For load balancing, run fabioLoadBalancing();
For simple service discovery, run serviceDiscoveryWithConsul();
*/
func main() {
	fabioLoadBalancing()
	fmt.Println("Starting Client.")
	var client = &http.Client{
		Timeout: time.Second * 30,
	}
	callServerEvery(10*time.Second, client)
}

/* Load balancing with Fabio */
func fabioLoadBalancing() {
	address := os.Getenv("FABIO_HTTP_ADDR")
	url = fmt.Sprintf("http://%s/helloworld", address)
}

/* Service Discovery with Consul */
func serviceDiscoveryWithConsul() {
	config := consulapi.DefaultConfig()
	consul, error := consulapi.NewClient(config)
	if error != nil {
		fmt.Println(error)
	}
	services, error := consul.Agent().Services()
	if error != nil {
		fmt.Println(error)
	}

	service := services["helloworld-server"]
	address := service.Address
	port := service.Port
	url = fmt.Sprintf("http://%s:%v/helloworld", address, port)
}

func hello(t time.Time, client *http.Client) {
	response, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("%s. Time is %v\n", body, t)
}

func callServerEvery(d time.Duration, client *http.Client) {
	for x := range time.Tick(d) {
		hello(x, client)
	}
}
