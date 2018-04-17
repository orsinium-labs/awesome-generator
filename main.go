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
	var topic string
	var pages int

	flag.StringVar(&lang, "l", "", "")
	flag.StringVar(&topic, "t", "", "")
	flag.IntVar(&pages, "pages", 10, "")

	flag.Parse()

	if lang != "" || topic != "" {
		data, err := dumpProjects(lang, topic, pages)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			return
		}
		os.Stdout.Write(data)
	} else {
		var projects []Project
		// decode JSON
		if err := json.NewDecoder(os.Stdin).Decode(&projects); err != nil {
			fmt.Println(err)
			return
		}
		makeMarkdown(projects)
	}
}

// download and save JSON from github
func dumpProjects(lang string, topic string, pages int) ([]byte, error) {
	wg.Add(pages)
	var projectsChan = make(chan Project, pages*100)
	for page := 1; page <= pages; page++ {
		go getProjects(lang, topic, page, &projectsChan)
	}
	wg.Wait()
	close(projectsChan)

	var projects []Project
	for project := range projectsChan {
		projects = append(projects, project)
	}
	dump, err := json.Marshal(projects)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return dump, nil
}

// generate markdown from projects list
func makeMarkdown(projects []Project) {
	totalProjectsCount := 0
	// group projects by topics
	topics := make(map[string][]Project)
	var topicsNames []string
	for _, project := range projects {
		totalProjectsCount++
		for _, topic := range project.Topics {
			if topics[topic] == nil {
				topicsNames = append(topicsNames, topic)
			}
			topics[topic] = append(topics[topic], project)
		}
	}

	// sort topics
	sort.Strings(topicsNames)

	// filter topics
	topicsNames = filter(topicsNames, func(topicName string) bool {
		l := len(topics[topicName])
		return l > 1 && l <= totalProjectsCount/5
	})

	// generate TOC
	for _, topicName := range topicsNames {
		fmt.Printf("1. [%s](#%s)\n", topicName, topicName)
	}

	// generate projects list
	var topicProjects []Project
	for _, topicName := range topicsNames {
		topicProjects = topics[topicName]
		fmt.Printf("\n\n## %s\n\n", topicName)
		slice.Sort(topicProjects[:], func(i, j int) bool {
			return topicProjects[i].Stars > topicProjects[j].Stars
		})
		for _, project := range topicProjects {
			fmt.Println(project.getMarkdown())
		}
	}
}

var wg sync.WaitGroup

// endpoint for Github API search
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

// make project link
func (p *Project) getLink() string {
	return fmt.Sprintf(projectLinkTemplate, p.Author, p.Name)
}

// make markdown block with project info
func (p *Project) getMarkdown() string {
	return fmt.Sprintf("1. [%s](%s). %s", p.Name, p.getLink(), p.Descr)
}

// filter slice by key function
func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// retrieve and extract projects from Github API
func getProjects(lang string, topic string, page int, projectsChan *chan Project) {
	defer wg.Done()

	// make request
	values := url.Values{}
	if lang != "" {
		values.Set("q", "language:"+lang)
	} else if topic != "" {
		values.Set("q", "topic:"+topic)
	}
	values.Set("page", fmt.Sprintf("%d", page))
	url := fmt.Sprintf("%s?%s", searchAPI, values.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// demand topics list
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
	lang = topic // for filtering
	// get projects from response
	for _, project := range data.Items {
		project.Author = strings.Split(project.Author, "/")[0]
		project.Topics = filter(project.Topics, func(topic string) bool {
			// doesn't starts wirh language - ok
			if !strings.HasPrefix(topic, lang) {
				return true
			}
			// starts with "lang-" - bad
			if strings.HasPrefix(topic, lang+"-") {
				return false
			}
			// near to lang - bad
			if len(topic)-len(lang) < 3 {
				return false
			}
			// ok otherwise
			return true
		})
		*projectsChan <- Project(project)
	}
}
