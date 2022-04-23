package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	out, err := exec.Command("curl",  "https://get.k0s.sh", "-sSLf", "|", "sudo", "sh").Output()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(out))
}
