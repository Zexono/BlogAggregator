package main

import (
	"blogagg/internal/config"
	"fmt"
)

func handlerLogin(s *state, cmd command) error{
	if len(cmd.args) == 0{
		return fmt.Errorf("command empty missing username")
	}
	s.ptrconfig.SetUser(cmd.args[0])
	fmt.Println("username has been set")

	return nil
}

type state struct {
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