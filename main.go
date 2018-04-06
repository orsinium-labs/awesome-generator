package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

func main() {
	var lang string
	var toJSON bool
	var pages int

	flag.StringVar(&lang, "l", "go", "")
	flag.BoolVar(&toJSON, "json", true, "")
	flag.IntVar(&pages, "pages", 1, "")

	flag.Parse()

	var projects []Project

	if toJSON {
		wg.Add(pages)
		for page := 1; page <= pages; page++ {
			fmt.Println(page)
			go getProjects(lang, page)
		}
		wg.Wait()
		close(projectsChan)
		for project := range projectsChan {
			copy(append(projects, project), projects)
		}
		fmt.Println(json.Marshal(projects))
	}
}

var wg sync.WaitGroup

const searchAPI = "https://api.github.com/search/repositories"
const projectLinkTemplate = "https://github.com/%s/%s"

// Project is type for one github project
type Project struct {
	name   string
	author string
	descr  string
	stars  int32
	topics []string
}

func (p *Project) getLink() string {
	return fmt.Sprintf(projectLinkTemplate, p.author, p.name)
}

func (p *Project) getMarkdown() string {
	return fmt.Sprintf("[%s](%s). %s", p.name, p.getLink(), p.descr)
}

var projectsChan chan Project

func getProjects(lang string, page int) {
	defer wg.Done()

	// make request
	values := url.Values{}
	values.Set("q", "language:"+lang)
	values.Set("page", fmt.Sprintf("%d", page))
	url := fmt.Sprintf("%s?%s", searchAPI, values.Encode())
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// parse JSON response
	var data struct {
		items             []map[string]interface{}
		totalCount        int64
		incompleteResults bool
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println(err)
		return
	}

	// get projects from response
	fmt.Println(data)
	for _, info := range data.items {
		projectsChan <- Project{
			name:   info["name"].(string),
			author: info["owner"].(map[string]interface{})["login"].(string),
			descr:  info["description"].(string),
			stars:  info["stargazers_count"].(int32),
			topics: info["topics"].([]string),
		}
	}
}
