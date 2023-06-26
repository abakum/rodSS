package main

import (
	"encoding/json"
	"os"
)

type config struct {
	fn  string
	P   map[string][]string
	Ids []int
}

var (
	conf *config
)

func loader(fn string) (conf *config, err error) {
	obj := config{fn: fn}
	conf = &obj
	bytes, err := os.ReadFile(fn)
	if err != nil {
		letf.Println("loader")
		return
	}
	err = json.Unmarshal(bytes, conf)
	if err != nil {
		letf.Println("loader")
		return
	}
	ltf.Println("loader done")
	ltf.Println(obj.Ids)
	return
}

func (conf *config) saver() (err error) {
	ltf.Println(conf.Ids)
	bytes, err := json.Marshal(conf)
	if err != nil {
		letf.Println("saver")
		return
	}
	err = os.WriteFile(conf.fn, bytes, 0o644)
	if err != nil {
		letf.Println("saver")
		return
	}
	ltf.Println("saver done")
	return
}
