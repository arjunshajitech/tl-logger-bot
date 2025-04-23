package tl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	Layout = "2006-01-02T15:04:05-0700"
	Format = "2006-01-02"
)

type TL interface {
	LoggedToday() (bool, error)
	LogToday() error
}

type TLImpl struct {
	User, Token string
}

func NewTL(user, token string) *TLImpl {
	return &TLImpl{
		User:  user,
		Token: token,
	}
}

type Timesheet struct {
	Begin string `json:"begin"`
}

func (t *TLImpl) LoggedToday() (bool, error) {
	url := "https://tl.techgentsia.com/api/timesheets"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("X-AUTH-USER", t.User)
	req.Header.Add("X-AUTH-TOKEN", t.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return false, err
	}

	var timesheets []Timesheet
	if err := json.Unmarshal(body, &timesheets); err != nil {
		return false, errors.New("JSON parse error")
	}

	return loggedToday(timesheets[0].Begin)
}

func loggedToday(date string) (bool, error) {
	t, err := time.Parse(Layout, date)
	if err != nil {
		return false, errors.New("Error parsing timestamp")
	}
	d := t.Format(Format)
	t1 := time.Now()
	d1 := t1.Format(Format)
	return d == d1, nil
}

func (t *TLImpl) LogToday() error {
	url := "https://tl.techgentsia.com/api/timesheets"
	begin, end := generateDayTimes()

	data := map[string]interface{}{
		"begin":       begin,
		"end":         end,
		"project":     1,
		"activity":    2,
		"description": nil,
		"tags":        "vconsol sfu",
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("X-AUTH-USER", t.User)
	req.Header.Add("X-AUTH-TOKEN", t.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Response Body: %s\n", responseBody)
	return nil
}

func generateDayTimes() (string, string) {
	now := time.Now()
	// Create morning 10:00 AM and evening 6:00 PM times
	morning := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())
	evening := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())
	morningStr := morning.Format("2006-01-02T15:04:05")
	eveningStr := evening.Format("2006-01-02T15:04:05")
	return morningStr, eveningStr
}
