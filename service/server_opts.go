package service

import "github.com/frantjc/sequence/conf"

type serverOpts struct {
	conf *conf.Config
}

type ServerOpt func(so *serverOpts) error

func WithConfig(c *conf.Config) ServerOpt {
	return func(so *serverOpts) error {
		so.conf = c
		return nil
	}
}
