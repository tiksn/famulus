package config

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
)

func ConfigDir() (string, error) {
	return homedir.Expand("~/.config/famulus")
}

func ConfigFile() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(dir, "config.yml"), nil
}

func ParseDefaultConfig() (Config, error) {
	file, err := ConfigFile()
	if err != nil {
		return nil, err
	}
	return ParseConfig(file)
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
