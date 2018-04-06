package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
		var projectsChan = make(chan Project, pages*100)
		for page := 1; page <= pages; page++ {
			go getProjects(lang, page, &projectsChan)
		}
		wg.Wait()
		close(projectsChan)
		for project := range projectsChan {
			copy(append(projects, project), projects)
		}
		fmt.Println(projects)
		// fmt.Println(json.Marshal(projects))
	}
}

var wg sync.WaitGroup

const searchAPI = "https://api.github.com/search/repositories"
const projectLinkTemplate = "https://github.com/%s/%s"

// Project is type for one github project
type Project struct {
	Name   string   `json:"name"`
	Author string   `json:"full_name"`
	Descr  string   `json:"description"`
	Stars  int32    `json:"stargazers_count"`
	Topics []string `json:"topics"`
}

func (p *Project) getLink() string {
	return fmt.Sprintf(projectLinkTemplate, p.Author, p.Name)
}

func (p *Project) getMarkdown() string {
	return fmt.Sprintf("[%s](%s). %s", p.Name, p.getLink(), p.Descr)
}

func getProjects(lang string, page int, projectsChan *chan Project) {
	defer wg.Done()

	// make request
	values := url.Values{}
	values.Set("q", "language:"+lang)
	values.Set("page", fmt.Sprintf("%d", page))
	url := fmt.Sprintf("%s?%s", searchAPI, values.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/vnd.github.mercy-preview+json")

	// get response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// parse JSON response
	var data struct {
		Items []Project `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println(err)
		return
	}
	// get projects from response
	for _, project := range data.Items {
		project.Author = strings.Split(project.Author, "/")[0]
		*projectsChan <- project
	}
}
