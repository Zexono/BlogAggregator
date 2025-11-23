package main

import (
	"blogagg/internal/config"
	"blogagg/internal/database"
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error{
	if len(cmd.args) == 0{
		return fmt.Errorf("command empty missing username")
	}

	name := cmd.args[0]
	if _ , err := s.db.GetUser(context.Background(),name); err != nil{
		
		return fmt.Errorf("username doesn't exist in the database")
	}

	s.ptrconfig.SetUser(name)
	fmt.Println("username has been set")

	return nil
}

func handlerResgister(s *state, cmd command) error{
	if len(cmd.args) == 0{
		return fmt.Errorf("command empty missing username arg")
	}
	
	name := cmd.args[0]
	if _ , err := s.db.GetUser(context.Background(),name); err == nil{
		
		return fmt.Errorf("user already exists")
	}
	s.db.CreateUser(context.Background(),database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
		Name: name })

	s.ptrconfig.SetUser(name)
	fmt.Println("user was created")

	user , _ := s.db.GetUser(context.Background(),name)
	fmt.Printf("created user: %+v\n", user)

	return nil
}

func handlerReset(s *state, cmd command) error{
	if len(cmd.args) != 0{
		return fmt.Errorf("command do not need args")
	}

	err := s.db.DeleteAllUser(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("All user deleted")
	return nil
}

func handlerUsers(s *state, cmd command) error{
	if len(cmd.args) != 0{
		return fmt.Errorf("command do not need args")
	}

	users,err := s.db.GetAllUser(context.Background())
	user_cur := s.ptrconfig.CurrentUserName
	if err != nil {
		return err
	}
	//name := users{}
	fmt.Println("List all users")
	for _, v := range users {
		if v.Name == user_cur {
			fmt.Printf("%s (current) \n",v.Name)
		}else{
			fmt.Printf("%s \n",v.Name)
		}
	}
	return nil
}

func handlerAgg(_ *state, cmd command) error{
	if len(cmd.args) != 0{
		return fmt.Errorf("command do not need args")
	}
	feed , err := fetchFeed(context.Background(),"https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(feed)


	return nil
}

func handlerAddfeed(s *state, cmd command) error{
	if len(cmd.args) == 0{
		return fmt.Errorf("command empty missing feed name and url")
	}

	if len(cmd.args) == 1{
		return fmt.Errorf("command empty missing url")
	}

	name := cmd.args[0]
	url := cmd.args[1]
	user_cur := s.ptrconfig.CurrentUserName

	user,err := s.db.GetUser(context.Background(),user_cur)
	if err != nil {
		return err
	}

	userid := user.ID
	
	
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
    ID:        uuid.New(),
    CreatedAt: time.Now().UTC(),
    UpdatedAt: time.Now().UTC(),
    Name:      name,
    Url:       url,
    UserID:    userid,})

	if err != nil {
    	return err
	}



	fmt.Printf("created feed: %+v\n", feed)

	feedfollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
	ID:        uuid.New(),
    CreatedAt: time.Now().UTC(),
    UpdatedAt: time.Now().UTC(),
    UserID:	userid,
	FeedID:	feed.ID})

	if err != nil {
    	return err
	}

	fmt.Printf("followed: %+v\n", feedfollow)

	return nil
}

func handlerFeeds(s *state, cmd command) error{
	if len(cmd.args) != 0{
		return fmt.Errorf("command do not need args")
	}

	feeds,err := s.db.GetAllFeed(context.Background())
	
	if err != nil {
		return err
	}
	
	fmt.Println("List all feeds")
	for _, v := range feeds {
		cur_user_name ,err := s.db.GetUserById(context.Background(),v.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("%s URL: %s Create By: %s \n",v.Name,v.Url,cur_user_name)
	}
	return nil
}

func handlerFollow(s *state, cmd command) error{
	if len(cmd.args) == 0{
		return fmt.Errorf("command empty missing feed url")
	}

	url := cmd.args[0]

	user_cur := s.ptrconfig.CurrentUserName
	user,err := s.db.GetUser(context.Background(),user_cur)
	if err != nil {
		println("user error")
		return err
	}
	feed_url ,err:= s.db.GetFeedFromURL(context.Background(),url)
	if err != nil {
		println("feed error")
		return err
	}

	userid := user.ID
	feedid := feed_url.ID
	
	
	feedfollow, err  := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
	ID:        uuid.New(),
    CreatedAt: time.Now().UTC(),
    UpdatedAt: time.Now().UTC(),
    UserID:	userid,
	FeedID:	feedid})

	if err != nil {
		println("follow error")
    	return err
	}

	fmt.Printf("created feed follow: %+v\n", feedfollow)
	
	return nil
}

func handlerFollowing(s *state, cmd command) error{
	if len(cmd.args) != 0{
		return fmt.Errorf("command do not need args")
	}

	user_cur := s.ptrconfig.CurrentUserName
	user,err := s.db.GetUser(context.Background(),user_cur)
	if err != nil {
		return err
	}

	userid := user.ID
	feed,err := s.db.GetFeedFollowsForUser(context.Background(),userid)

	if err != nil {
    	return err
	}

	for _, v := range feed {
		feed_name,err := s.db.GetFeedFromID(context.Background(),v.FeedID)
		if err != nil {
    	return err
		}
		fmt.Printf("Feed: %+v By %s\n", feed_name.Name,user_cur)
	}
	//fmt.Printf("Feed: %+v By %s\n", feed,user_cur)
	
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

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}


func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error){
	req, err := http.NewRequestWithContext(ctx,"GET",feedURL,nil)
	req.Header.Set("User-Agent", "gator")
	if err != nil {
		//fmt.Printf("Error creating request \n")
		
        return nil , fmt.Errorf("error creating request %v ",err)
	}
	
	//body, err := io.ReadAll(resp.Body)
	client := &http.Client{}

	resp, err := client.Do(req)
    if err != nil {
        //fmt.Printf("Error making request \n")
        return nil , fmt.Errorf("error making request %v ",err)
    }
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		//fmt.Printf("Error Reading Body \n")
        return nil , fmt.Errorf("error Reading Body %v ",err)
	}

	var feed RSSFeed
	xml.Unmarshal(body,&feed)
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	
	return &feed,nil
}


