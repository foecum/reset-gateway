package main

import (
	"log"
	"os/exec"
	"strings"
	"time"
)

var networkErr = "connect: Network is unreachable"

func main() {
	exitChan := make(chan string, 1)
	go func() {
		for {
			b, _ := pingGoogle()

			if strings.Compare(networkErr, strings.Trim(string(b), "\n")) == 0 {
				log.Printf("Network gateway not set. Setting...\n")
				resetGW := exec.Command("sudo", "ip", "route", "add", "default", "via", "192.168.8.1")
				b, err := resetGW.CombinedOutput()

				if err != nil {
					log.Printf("%v\n", err)
					log.Printf("%v\n", string(b))
					exitChan <- "Exiting loop"

					return
				}
				log.Printf("Network gateway set\n")
			}

			time.Sleep(60 * time.Second)
		}

	}()

	<-exitChan
}

func pingGoogle() ([]byte, error) {
	ping := exec.Command("ping", "-c 1", "google.com")
	return ping.CombinedOutput()
}
