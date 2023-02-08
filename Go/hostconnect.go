package main

import (
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

func (conf *ConfType) connectHost(host string, port uint, user string, pw string) error {

	// Unfortunately we can't use goph.New() as the port is fixed
	// To allow port input, we need to use goph.NewConn()
	client, err := goph.NewConn(&goph.Config{
		User:     user,
		Addr:     host,
		Port:     port,
		Auth:     goph.Password(pw),
		Callback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return err
	}

	conf.Client = client

	return nil

}
