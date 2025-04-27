package gl

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Commits struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type MyProjects struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Branch struct {
	Name string `json:"name"`
}

func (g *Gitlab) Branches(projectId string) []Branch {

	var branches []Branch
	baseURL := "https://gitlab.techgentsia.com/api/v4/projects/" + projectId + "/repository/branches"
	params := url.Values{}
	//params.Add("pagination", "keyset")
	//params.Add("per_page", "100")
	//params.Add("order_by", "id")
	//params.Add("sort", "asc")

	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Println(err.Error())
		return branches
	}
	req.Header.Add("PRIVATE-TOKEN", g.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return branches
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println(readErr.Error())
		return branches
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("GitLab API error: status %d, body: %s\n", res.StatusCode, string(body))
		return branches
	}

	if err := json.Unmarshal(body, &branches); err != nil {
		log.Println(err.Error())
		return branches
	}

	return branches
}

func (g *Gitlab) MyProjects() []MyProjects {

	var projects []MyProjects
	baseURL := "https://gitlab.techgentsia.com/api/v4/projects"
	params := url.Values{}
	params.Add("pagination", "keyset")
	params.Add("per_page", "100")
	params.Add("order_by", "id")
	params.Add("sort", "asc")

	fullURL := baseURL + "?" + params.Encode()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Println(err.Error())
		return projects
	}
	req.Header.Add("PRIVATE-TOKEN", g.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err.Error())
		return projects
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println(readErr.Error())
		return projects
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("GitLab API error: status %d, body: %s\n", res.StatusCode, string(body))
		return projects
	}

	if err := json.Unmarshal(body, &projects); err != nil {
		log.Println(err.Error())
		return projects
	}

	return projects
}

func (g *Gitlab) Commits(since, until time.Time) ([]Commits, error) {

	getCommitMessages := func(projectID int, branch string, author string) []Commits {

		var commits []Commits
		baseURL := "https://gitlab.techgentsia.com/api/v4/projects/" + strconv.Itoa(projectID) + "/repository/commits"
		params := url.Values{}
		params.Add("since", since.Format("2006-01-02T15:04:05Z"))
		params.Add("until", until.Format("2006-01-02T15:04:05Z"))
		params.Add("author", author)
		params.Add("ref_name", branch)

		fullURL := baseURL + "?" + params.Encode()

		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			log.Println(err.Error())
			return commits
		}
		req.Header.Add("PRIVATE-TOKEN", g.Token)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err.Error())
			return commits
		}
		defer res.Body.Close()

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Println(readErr.Error())
			return commits
		}

		if res.StatusCode != http.StatusOK {
			log.Printf("GitLab API error: status %d, body: %s\n", res.StatusCode, string(body))
			return commits
		}

		if err := json.Unmarshal(body, &commits); err != nil {
			log.Println(err.Error())
			return commits
		}

		return commits
	}

	var allCommits []Commits
	for _, p := range g.Projects {
		for _, b := range p.Branch {
			commits := getCommitMessages(p.ID, b, p.Author)
			allCommits = append(allCommits, commits...)
		}
	}

	return allCommits, nil
}
