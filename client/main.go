package main

import "github.com/google/go-github/github"
import "code.google.com/p/goauth2/oauth"
import "fmt"
import "strings"
import "time"

/*
func main() {
	//Channel configuration
	ch := make(chan *(github.Repository))

	//Authentication
	apiKey := "0cc6bb4dd641f483323592eda52b8bcf6327e3ed"
	secret := &oauth.Transport{
		Token: &oauth.Token{AccessToken: apiKey},
	}

	//Client
	client := github.NewClient(secret.Client())
	//var timeZone = map[string]Repo
	repos := make(map[string]*(github.Repository))
	
	start := time.Now().Unix()
	
	//Fetch User Events By Organizations
	opts := &github.ListOptions{Page: 0, PerPage: 10}
	events, _, err := client.Activity.ListUserEventsForOrganization("koombea", "davidpabon", opts)

	if err != nil {
		fmt.Printf("Error listing events: %s\n", err)
	} else {	
		
		fmt.Printf("Total Events : %v\n", len(events))
		
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
						}
					} (repoNameSplit[0], repoNameSplit[1])

					repos[*event.Repo.Name] = event.Repo
				}
			}
		}

		for key, _ := range repos {
			select {
			case r := <-ch:
	        	repos[key] = r
	    	case <-time.After(30 * time.Second):
	    		fmt.Printf("Time out error")
			}
		}

		
	}

	end := time.Now().Unix()
	fmt.Printf("Total Repos: %v\n", len(repos))
	fmt.Printf("Total Repos: %v\n", repos)
	fmt.Printf("Time elapsed: %v - %v = %v\n", end, start, end - start)

}
*/

/*
func asyncHttpGets(urls []string) []*HttpResponse {
  ch := make(chan *HttpResponse)
  responses := []*HttpResponse{}
  for _, url := range urls {
      go func(url string) {
          fmt.Printf("Fetching %s \n", url)
          resp, err := http.Get(url)
          ch <- &HttpResponse{url, resp, err}
      }(url)
  }

  for {
      select {
      case r := <-ch:
          fmt.Printf("%s was fetched\n", r.url)
          responses = append(responses, r)
          if len(responses) == len(urls) {
              return responses
          }
      case <-time.After(50 * time.Millisecond):
          fmt.Printf(".")
      }
  }
  return responses
}
*/

func main() {

	//Authentication
	apiKey := "0cc6bb4dd641f483323592eda52b8bcf6327e3ed"
	secret := &oauth.Transport{
		Token: &oauth.Token{AccessToken: apiKey},
	}

	//Client
	client := github.NewClient(secret.Client())
	//var timeZone = map[string]Repo
	repos := make(map[string]*(github.Repository))
	
	start := time.Now().Unix()
	//Fetch User Events By Organizations
	opts := &github.ListOptions{Page: 0, PerPage: 10}
	events, _, err := client.Activity.ListUserEventsForOrganization("koombea", "davidpabon", opts)

	if err != nil {
		fmt.Printf("Error listing events: %s\n", err)
	} else {	
		fmt.Printf("Total Events : %v\n", len(events))
		for _, event := range events {

			if _, ok := repos[*event.Repo.Name]; !ok {
				repoNameSplit := strings.Split(*event.Repo.Name, "/")
				if len(repoNameSplit[0])>0 && len(repoNameSplit[1])>0 {
					repo, _, errRepo := client.Repositories.Get(repoNameSplit[0], repoNameSplit[1])
					if errRepo != nil {
						fmt.Printf("Error getting error: %s\n", errRepo)
					} else {
						repos[*event.Repo.Name] = repo
					}
				}
			}
		}

		for name, repository := range repos {
			fmt.Printf("Repo: %v, Created By: %v\n", name, *repository.Owner.Login)
		}
	}
	end := time.Now().Unix()
	fmt.Printf("Time elapsed: %v - %v = %v\n", end, start, end - start)

}
