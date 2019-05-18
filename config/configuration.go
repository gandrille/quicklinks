package config

import "github.com/gandrille/go-commons/result"

type Configuration struct {
	key         string
	runner      func() result.Result
}

func (c Configuration) Key() string {
	return c.key
}

func (c Configuration) Run() result.Result {
	return c.runner()
}
