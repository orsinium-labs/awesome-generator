package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func main() {
	var lang string
	var pages int

	flag.StringVar(&lang, "l", "", "")
	flag.IntVar(&pages, "pages", 1, "")

	flag.Parse()

	var projects []Project

	if lang != "" {
		wg.Add(pages)
		var projectsChan = make(chan Project, pages*100)
		for page := 1; page <= pages; page++ {
			go getProjects(lang, page, &projectsChan)
		}
		wg.Wait()
		close(projectsChan)
		for project := range projectsChan {
			projects = append(projects, project)
		}
		b, err := json.Marshal(projects)
		if err != nil {
			fmt.Println(err)
			return
		}
		os.Stdout.Write(b)
	} else {

	}
}

var wg sync.WaitGroup

const searchAPI = "https://api.github.com/search/repositories"
const projectLinkTemplate = "https://github.com/%s/%s"

// Project is type for one github project
type Project struct {
	Name   string   `json:"name"`
	Author string   `json:"author"`
	Descr  string   `json:"description"`
	Stars  int32    `json:"stars"`
	Topics []string `json:"topics"`
}

// ProjectImport is type for one github project importing
type ProjectImport struct {
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
		Items []ProjectImport `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println(err)
		return
	}
	// get projects from response
	for _, project := range data.Items {
		project.Author = strings.Split(project.Author, "/")[0]
		*projectsChan <- Project(project)
	}
}
