package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/bradfitz/slice"
)

func main() {
	var lang string
	var pages int

	flag.StringVar(&lang, "l", "", "")
	flag.IntVar(&pages, "pages", 10, "")

	flag.Parse()

	var projects []Project

	if lang != "" {
		// download and save JSON from github
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
		// generate markdown
		// decode JSON
		if err := json.NewDecoder(os.Stdin).Decode(&projects); err != nil {
			fmt.Println(err)
			return
		}
		// group projects by topics
		topics := make(map[string][]Project)
		var topicsNames []string
		for _, project := range projects {
			for _, topic := range project.Topics {
				if topics[topic] == nil {
					topicsNames = append(topicsNames, topic)
				}
				topics[topic] = append(topics[topic], project)
			}
		}
		// sort topics
		sort.Strings(topicsNames)
		// generate markdown
		var topicProjects []Project
		for _, topicName := range topicsNames {
			topicProjects = topics[topicName]
			if len(topicProjects) > 1 {
				fmt.Printf("\n\n## %s\n\n", topicName)
				slice.Sort(topicProjects[:], func(i, j int) bool {
					return topicProjects[i].Stars > topicProjects[j].Stars
				})
				for _, project := range topicProjects {
					fmt.Println(project.getMarkdown())
				}
			}
		}
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
	return fmt.Sprintf("1. [%s](%s). %s", p.Name, p.getLink(), p.Descr)
}

func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
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
		project.Topics = filter(project.Topics, func(topic string) bool {
			return !strings.HasPrefix(topic, lang)
		})
		*projectsChan <- Project(project)
	}
}
