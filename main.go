package main

import (
	"blogagg/internal/config"
	"fmt"
)

func main(){
	var cfg config.Config
	//fmt.Print("2nd go \n")
	cfg,_ = config.Read()
	cfg.SetUser("Zexono")
	cfg,_ = config.Read()
	fmt.Println(cfg)
}