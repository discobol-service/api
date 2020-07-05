package main

import (
	"github.com/discobol-service/api/daemon"
)


func main() {
	daemon.New(":3000").Run()
}