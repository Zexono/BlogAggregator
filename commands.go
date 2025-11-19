package main

import (
	"blogagg/internal/config"
	"blogagg/internal/database"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error{
	if len(cmd.args) == 0{
		os.Exit(1)
		return fmt.Errorf("command empty missing username")
	}

	name := cmd.args[0]
	if _ , err := s.db.GetUser(context.Background(),name); err != nil{
		
		os.Exit(1)
		return fmt.Errorf("username doesn't exist in the database")
	}

	s.ptrconfig.SetUser(name)
	fmt.Println("username has been set")

	return nil
}

func handlerResgister(s *state, cmd command) error{
	if len(cmd.args) == 0{
		os.Exit(1)
		return fmt.Errorf("command empty missing username arg")
	}
	
	name := cmd.args[0]
	if _ , err := s.db.GetUser(context.Background(),name); err == nil{
		
		os.Exit(1)
		return fmt.Errorf("user already exists")
	}
	s.db.CreateUser(context.Background(),database.CreateUserParams{ID: uuid.New(),CreatedAt: time.Now().Local(),UpdatedAt: time.Now().Local(),Name: name })

	s.ptrconfig.SetUser(name)
	fmt.Println("user was created")

	user , _ := s.db.GetUser(context.Background(),name)
	fmt.Printf("created user: %+v\n", user)

	return nil
}

func reset(s *state, cmd command) error{
	if len(cmd.args) != 0{
		return fmt.Errorf("command do not need args")
	}

	err := s.db.DeleteAllUser(context.Background())
	if err != nil {
		os.Exit(1)
		return err
	}
	fmt.Println("All user deleted")
	return nil
}

type state struct {
	db  *database.Queries
	ptrconfig *config.Config
}

type command struct {
	name string
	args  []string
}

type commands struct {
	handlers 	map[string]func(*state, command) error

}
func (c *commands) run(s *state, cmd command) error{
	//checkmap := c.handlers
	//val,have := checkmap[cmd.name]

	val, ok := c.handlers[cmd.name]
	
	if !ok {
		return fmt.Errorf("command not found: %s",cmd.name)
	}
	
	return val(s,cmd)
}

func (c *commands) register(name string, f func(*state, command) error){
	//setval := c.handlers 
	//setval[name] = f
	if c.handlers == nil {
        c.handlers = make(map[string]func(*state, command) error)
    }
	c.handlers[name] = f
}

