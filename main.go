package main

import (
	"blogagg/internal/config"
	"blogagg/internal/database"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const dbURL = "postgres://postgres:postgres@localhost:5432/gator"

func main(){
	var cfg config.Config

	db, err := sql.Open("postgres", dbURL)
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}

	dbQueries := database.New(db)


	cfg,_ = config.Read()
	//cfg.SetUser("Zexono")
	//cfg,_ = config.Read()
	//fmt.Println(cfg)
	s := &state{dbQueries,&cfg}
	c := commands{make(map[string]func(*state, command) error)}
	c.register("login",handlerLogin)
	c.register("register",handlerResgister)
	c.register("reset",handlerReset)
	c.register("users",handlerUsers)
	c.register("agg",handlerAgg)
	c.register("addfeed",handlerAddfeed)
	input := os.Args
	//if input[0] != "gator" {
	//	fmt.Println("plz use gator <command> <args>")
	//}
	if len(input) < 2 {
		fmt.Println("no command given")
		os.Exit(1)
	}
	//if len(input) < 3 { // move if check into func itself already
	//	fmt.Println("arg is required")
	//	os.Exit(1) 
	//}
	//fmt.Println(len(input))
	cmd := command{name: input[1],args: input[2:]}
	err = c.run(s,cmd)
	if err != nil {
		os.Exit(1)
	}
}