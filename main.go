package main

import (
	"blogagg/internal/config"
	"fmt"
	"os"
)

func main(){
	var cfg config.Config
	
	cfg,_ = config.Read()
	//cfg.SetUser("Zexono")
	//cfg,_ = config.Read()
	//fmt.Println(cfg)
	s := &state{&cfg}
	c := commands{make(map[string]func(*state, command) error)}
	c.register("login",handlerLogin)
	input := os.Args
	//if input[0] != "gator" {
	//	fmt.Println("plz use gator <command> <args>")
	//}
	if len(input) < 2 {
		fmt.Println("no command given")
		os.Exit(1)
	}
	if len(input) < 3 && input[1] == "login" {
		fmt.Println("username is required")
		os.Exit(1)
	}
	//fmt.Println(len(input))
	cmd := command{name: input[1],args: input[2:]}
	c.run(s,cmd)
}