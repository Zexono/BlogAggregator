package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Zexono/blogagg/internal/config"
	"github.com/Zexono/blogagg/internal/database"

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

func handlerAgg(s *state, cmd command) error{
	if len(cmd.args) <1{
		return fmt.Errorf("command need time like 1s 1m 1h")
	}
	time_between_reqs := cmd.args[0] 

	timeBetweenRequests,err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	//if len(cmd.args) != 0{
	//	return fmt.Errorf("command do not need args")
	//}
	//feed , err := fetchFeed(context.Background(),"https://www.wagslane.dev/index.xml")
	//if err != nil {
	//	return err
	//}
	//fmt.Println(feed)


	//return nil
}

func handlerAddfeed(s *state, cmd command, user database.User) error{
	if len(cmd.args) == 0{
		return fmt.Errorf("command empty missing feed name and url")
	}

	if len(cmd.args) == 1{
		return fmt.Errorf("command empty missing url")
	}

	name := cmd.args[0]
	url := cmd.args[1]

	userid := user.ID
	
	
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
    ID:        uuid.New(),
    CreatedAt: time.Now().UTC(),
    UpdatedAt: time.Now().UTC(),
    Name:      name,
    Url:       url,
    UserID:    userid})

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

func handlerFollow(s *state, cmd command, user database.User) error{
	if len(cmd.args) == 0{
		return fmt.Errorf("command empty missing feed url")
	}

	url := cmd.args[0]

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

func handlerFollowing(s *state, cmd command, user database.User) error{
	if len(cmd.args) != 0{
		return fmt.Errorf("command do not need args")
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
		fmt.Printf("Feed: %+v By %s\n", feed_name.Name,user.Name)
	}
	//fmt.Printf("Feed: %+v By %s\n", feed,user_cur)
	
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error{
	if len(cmd.args) == 0{
		return fmt.Errorf("command empty missing feed url")
	}

	url := cmd.args[0]


	feed_url ,err:= s.db.GetFeedFromURL(context.Background(),url)
	if err != nil {
		return err
	}

	userid := user.ID
	feedid := feed_url.ID
	
	err = s.db.DeleteFeedFollows(context.Background(),database.DeleteFeedFollowsParams{UserID: userid,FeedID: feedid})
	if err != nil {
    	return err
	}

	fmt.Printf("unfollow feed: %v\n",url)
	
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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error{
	
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.ptrconfig.CurrentUserName)
		if err != nil {
			return err
		}
		err = handler(s,cmd,user)
		if err != nil {
			return err
		}
		
		return nil
	}
}

func scrapeFeeds(s *state) {
	nextfeed ,err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("Error occurred: %v\n", err)
	}
	err = s.db.MarkFeedFetched(context.Background(),database.MarkFeedFetchedParams{
	LastFetchedAt: sql.NullTime{
    Time:  time.Now().UTC(),
    Valid: true},
	UpdatedAt: time.Now().UTC(),ID: nextfeed.ID})
	if err != nil {
		log.Printf("Error occurred: %v\n", err)
	}

	rssFeed, err := fetchFeed(context.Background(),nextfeed.Url)

	if err != nil  {
		log.Printf("Error occurred: %v\n", err)
	}
	
	for _, v := range rssFeed.Channel.Item {
		fmt.Printf("Feed title: %s",v.Title)
	}
	//layout := "2006-01-02 15:04:05"
	for _, v := range rssFeed.Channel.Item {
		pubDate,err := time.Parse(time.RFC1123Z,v.PubDate)
		if err != nil  {
			log.Printf("Error occurred: %v\n", err)
		}
		_,err = s.db.CreatePost(context.Background(),database.CreatePostParams{
		ID: uuid.New(),
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
		Title: v.Title,
		Url: v.Link,
		Description: v.Description,
		PublishedAt: pubDate,
		FeedID: nextfeed.ID, 
		})
		if err != nil {
			msg := err.Error() 
			if strings.Contains(msg,"duplicate key value")  {
				
			}else{
			log.Printf("Error occurred: %v\n", err)	
			}
			
		}
	}
	
	

}

func handlerBrowse(s *state, cmd command,user database.User) error{
	var limit_show int32
	if len(cmd.args) <1{
		//return fmt.Errorf("command need time like 1s 1m 1h")
		limit_show = 2
	}else{
		v, err := strconv.Atoi(cmd.args[0])

		if err != nil {
			fmt.Println("need to input number")
			return err
		}
		limit_show = int32(v)
	}
	
	post,err := s.db.GetPostForUser(context.Background(),database.GetPostForUserParams{
		UserID: user.ID,
		Limit: limit_show })
	if err != nil {
		return err
	}

	for _, v := range post {
		fmt.Printf("post %v\n", v)
	}
	
	return nil

}

