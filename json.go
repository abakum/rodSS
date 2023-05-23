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
		stdo.Println("loader")
		return
	}
	err = json.Unmarshal(bytes, conf)
	if err != nil {
		stdo.Println("loader")
		return
	}
	stdo.Println("loader done")
	stdo.Println(obj.Ids)
	return
}

func (conf *config) saver() (err error) {
	stdo.Println(conf.Ids)
	bytes, err := json.Marshal(conf)
	if err != nil {
		stdo.Println("saver")
		return
	}
	err = os.WriteFile(conf.fn, bytes, 0o644)
	if err != nil {
		stdo.Println("saver")
		return
	}
	stdo.Println("saver done")
	return
}
