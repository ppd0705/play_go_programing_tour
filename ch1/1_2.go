package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var name string

func main() {
	fmt.Println(len(os.Args), os.Args)
	flag.Parse()
	goCmd := flag.NewFlagSet("go", flag.ExitOnError)
	goCmd.StringVar(&name, "name", "default go", "default help msg")
	pyCmd := flag.NewFlagSet("py", flag.ExitOnError)
	pyCmd.StringVar(&name, "name", "default py", "default help msg")

	args := flag.Args()
	fmt.Println("cmd:", args[0])
	switch args[0] {
	case "go":
		_ = goCmd.Parse(args[1:])
	case "py":
		_ = pyCmd.Parse(args[1:])
	default:
		panic(fmt.Sprintf("unknown cmd: %s", args[0]))
	}
	log.Printf("name: %s\n", name)
}
