package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
	curlOut := exec.Command("curl", "https://get.k0s.sh", "-sSLf")
	sudoCmd := exec.Command("sudo", "sh")

	r, w := io.Pipe()

	curlOut.Stdout = w
	sudoCmd.Stdin = r

	var b2 bytes.Buffer
	sudoCmd.Stdout = &b2

	err := curlOut.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = sudoCmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = curlOut.Wait()
	if err != nil {
		log.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = sudoCmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
	wr, err := io.Copy(os.Stdout, &b2)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(wr)



	// Create working directory
	err = os.Mkdir("/etc/k0s", 0755)

	if err != nil {
		log.Fatal(err)
	}
}