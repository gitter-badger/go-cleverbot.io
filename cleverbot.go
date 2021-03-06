// Package cleverbot implements wrapper for the cleverbot.io API.
package cleverbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// baseURL for the cleverbot.io API.
const baseURL = "https://cleverbot.io/1.0/"

// New bot instance.
// "nick" is optional if you did not specify it, a random one is generated for you.
// A successful call returns err == nil.
func New(user, key, nick string) (s *Session, err error) {
	s = &Session{
		User: user,
		Key:  key,
		Nick: nick,
	}

	params, err := json.Marshal(s)
	if err != nil {
		return
	}

	response, err := http.Post(baseURL+"create", "application/json", bytes.NewBuffer(params))
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	m := map[string]string{}
	err = json.Unmarshal([]byte(body), &m)
	if err != nil {
		return
	}

	if m["status"] == "success" {
		s.Nick = m["nick"]
	} else {
		err = fmt.Errorf(m["status"])
		return
	}

	return
}

// Ask Cleverbot a question, returns Cleverbots response.
// A successful call returns err == nil.
func (s *Session) Ask(text string) (output string, err error) {

	s.Text = text

	params, err := json.Marshal(s)
	if err != nil {
		return
	}

	response, err := http.Post(baseURL+"ask", "application/json", bytes.NewBuffer(params))
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	m := map[string]string{}
	err = json.Unmarshal([]byte(body), &m)
	if err != nil {
		return
	}

	if m["status"] != "success" {
		err = fmt.Errorf(m["status"])
		return
	}

	// return the bots asnwer.
	return m["response"], nil
}
