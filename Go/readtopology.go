package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func (conf *ConfType) readTopologyFile(fileName string) error {
	yfile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(yfile, &conf.Topology); err != nil {
		return err
	}

	return nil

}
