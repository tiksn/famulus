package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
)

func ConfigDir() string {
	dir, err := homedir.Expand("~/.config/famulus")
	if err != nil {
		log.Fatalln(err)
	}
	return dir
}

func ConfigFile() string {
	return path.Join(ConfigDir(), "config.yml")
}

func ParseDefaultConfig() (Config, error) {
	return ParseConfig(ConfigFile())
}

func ParseConfig(filename string) (Config, error) {
	f, ferr := os.Open(filename)
	if ferr != nil {
		return nil, ferr
	}
	defer f.Close()

	data, rerr := ioutil.ReadAll(f)
	if rerr != nil {
		return nil, rerr
	}

	var root yaml.Node
	yerr := yaml.Unmarshal(data, &root)
	if yerr != nil {
		return nil, yerr
	}

	return &fileConfig{
		documentRoot: &root,
	}, nil
}
