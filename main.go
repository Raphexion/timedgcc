package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	args := os.Args[1:]
	fmt.Println(args)

	gcc := "/usr/bin/g++"

	t0 := time.Now()
	cmd := exec.Command(gcc, args...)
	t1 := time.Now()

	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
	fmt.Println(t1.Sub(t0))
}
