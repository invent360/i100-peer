package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
)

func main() {
	Reset()

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


	if _, err := os.Stat("/etc/k0s"); !os.IsNotExist(err) {
		err := os.RemoveAll("/etc/k0s")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Create working directory
	err = os.Mkdir("/etc/k0s", 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Create working directory
	err = os.Chdir("/etc/k0s")
	if err != nil {
		log.Fatal(err)
	}


	// Create config file
	out, err := exec.Command("k0s", "config", "create").Output()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	err = os.WriteFile("/etc/k0s/k0s.yaml", out, 0755)
	if err != nil {
		log.Fatal(err)
	}

	//Install controller
	ctrOut, err := exec.Command("k0s", "install", "controller", "-c", "/etc/k0s/k0s.yaml").Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(ctrOut))

	//Start controller
	startOut, err := exec.Command("k0s", "start").Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(startOut))

	//Start controller
	statusOut, err := exec.Command("k0s", "status").Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(statusOut))
}


func Reset(){
	stopOut, err := exec.Command("k0s", "stop").Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(stopOut))

	resetOut, err := exec.Command("k0s", "reset").Output()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(resetOut))
}
