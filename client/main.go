package main

import "github.com/google/go-github/github"
import "code.google.com/p/goauth2/oauth"
import "fmt"
import "strings"
import "time"

import "database/sql"
import _ "github.com/lib/pq"

var db *sql.DB

func main() {
	//Channel configuration
	ch := make(chan *(github.Repository))

	//Authentication
	//apiKey := "0cc6bb4dd641f483323592eda52b8bcf6327e3ed"
	apiKey := "e0064750447ddbd49c3fb7a86bb21a411dff85fb"
	secret := &oauth.Transport {
		Token: &oauth.Token{AccessToken: apiKey},
	}

	//Credentials
	//username := "davidpabon"
	username := "icas"
	organization := "koombea"

	//Client
	client := github.NewClient(secret.Client())
	//var timeZone = map[string]Repo
	repos := make(map[string]*(github.Repository))
	
	start := time.Now().Unix()
	
	//Fetch User Events By Organizations
	events, err := getEvents(organization, username, 10, client)

	totalSuccess := 0

	if err != nil {
		fmt.Printf("Error listing events: %s\n", err)
	} else {	
		
		fmt.Printf("Total Events: %v\n", len(events))
		
		for _, event := range events {

			if _, ok := repos[*event.Repo.Name]; !ok {
				
				repoNameSplit := strings.Split(*event.Repo.Name, "/")
				
				if len(repoNameSplit[0])>0 && len(repoNameSplit[1])>0 {
					
					go func(repoOwner string, repoName string) {
						repo, _, errRepo := client.Repositories.Get(repoOwner, repoName)
						if errRepo != nil {
							fmt.Printf("Error getting error: %s\n", errRepo)
						} else {
							ch <- repo
							totalSuccess += 1
						}
					} (repoNameSplit[0], repoNameSplit[1])

					repos[*event.Repo.Name] = event.Repo
				}
			}
		}

		for i := 0; i < totalSuccess; i++ {
			select {
			case r := <-ch:
	        	repos[*r.Name] = r
	    	case <-time.After(30 * time.Second):
	    		fmt.Printf("Time out error")
			}
		}

		// if events!=nil && len(events)>0 {
		// 	saveActivities(events, repos)
		// }
	}

	end := time.Now().Unix()
	fmt.Printf("Total Repos: %v\n", len(repos))
	fmt.Printf("Time elapsed: %v - %v = %v\n", end, start, end - start)

	//titleArray := []string{*events[0].Actor.Login, " made a ", *events[0].Type, " to ", *events[0].Repo.Name} //title
	//title := strings.Join(titleArray, "") // title
}

func getEvents(org string, user string, pages int, client *github.Client) ([]github.Event, error) {
	
	var events []github.Event
	var err error
	ch := make(chan []github.Event)

	if pages <= 0 {pages = 1}
	
	for i := 0; i < pages; i++ {
		go func() {
			var temp []github.Event
			opts := &github.ListOptions{Page: i}
			temp, _, err = client.Activity.ListUserEventsForOrganization(org, user, opts)
			ch <- temp
		} ()
	}
	
	for i := 0; i < pages; i++ {
		select {
		case evs := <-ch:
        	events = append(events, evs...)
    	case <-time.After(30 * time.Second):
    		fmt.Printf("Time out error")
		}
	}

	return events, err 
}

func saveActivities(events []github.Event, repos map[string]*(github.Repository)) {

	//Connection to database
	var errConn error

	db, errConn = sql.Open("postgres", "user=davidpabon dbname=dashable sslmode=disable")
	if errConn != nil {
		fmt.Printf("Error connecting to database: %v\n", errConn)
	}
	
	defer db.Close()
	
	//
	query := "	INSERT INTO activities(	title, provider, url, created_at, updated_at) VALUES($1, $2, $3, $4, $5)"
	for _, ev := range events {
		repo := repos[*ev.Repo.Name]
		url := []string{"https://github.com", *repo.Name}
		_, errExec := db.Exec(query, *ev.Type, "github", strings.Join(url, "/"),  time.Now(), time.Now())
		if errExec != nil {
			fmt.Printf("Error executing query: %v\n", errExec)
		}
	}
}

/*

*/