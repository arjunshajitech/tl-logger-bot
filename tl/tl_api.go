package tl

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Timesheet struct {
	Begin string `json:"begin"`
}

func (t *CustomLogger) processLogging(desc string) error {
	url := "https://tl.techgentsia.com/api/timesheets"
	begin, end := generateDayTimes()

	data := map[string]interface{}{
		"begin":       begin,
		"end":         end,
		"project":     1,
		"activity":    2,
		"description": desc,
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

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return nil
}

func (t *CustomLogger) loggedToday() (bool, error) {
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
