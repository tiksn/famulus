package config

import (
	"errors"
	"log"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
)

type Config interface {
	ListSources() ([]string, error)
	GetSource(string) (SourceConfig, error)
	GetContactsCsvFilePath() (string, error)
}

type SourceConfig interface {
	GetAddress() (string, error)
	GetPhoneAddress() (string, error)
	GetDefaultRegion() (string, error)
}

type fileConfig struct {
	documentRoot *yaml.Node
}

type sourceConfig struct {
	rootNode *yaml.Node
}

type NotFoundError struct {
	error
}

func findEntry(node *yaml.Node, key string) (keyNode *yaml.Node, valueNode *yaml.Node, err error) {
	topLevelKeys := node.Content
	for i, v := range topLevelKeys {
		if v.Value == key {
			if i+1 < len(topLevelKeys) {
				return v, topLevelKeys[i+1], nil
			} else {
				return v, nil, nil
			}
		}
	}

	return nil, nil, &NotFoundError{errors.New("not found")}
}

func getMapEntry(node *yaml.Node) (map[string]*yaml.Node, error) {
	if node.Kind != yaml.MappingNode {
		return nil, errors.New("Not a mapping node")
	}

	result := make(map[string]*yaml.Node)

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		result[keyNode.Value] = valueNode
	}
	return result, nil
}

func getSourceMap(c *fileConfig) (map[string]SourceConfig, error) {
	_, sourcesValueNode, err1 := findEntry(c.documentRoot.Content[0], "Sources")
	if err1 != nil {
		return nil, err1
	}
	sourceMap, err2 := getMapEntry(sourcesValueNode)
	if err2 != nil {
		return nil, err2
	}

	result := make(map[string]SourceConfig)

	/* Copy Content from Map1 to Map2*/
	for index, element := range sourceMap {
		result[index] = &sourceConfig{
			rootNode: element,
		}
	}

	return result, nil
}

func (c *fileConfig) ListSources() ([]string, error) {
	sourceMap, err := getSourceMap(c)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(sourceMap))

	for k := range sourceMap {
		keys = append(keys, k)
	}
	return keys, nil
}

func (c *fileConfig) GetSource(name string) (SourceConfig, error) {
	sourceMap, err := getSourceMap(c)
	if err != nil {
		return nil, err
	}
	return sourceMap[name], nil
}

func (c *fileConfig) GetContactsCsvFilePath() (string, error) {
	_, contactsValueNode, err1 := findEntry(c.documentRoot.Content[0], "Contacts")
	if err1 != nil {
		return "", err1
	}

	path, err := homedir.Expand(contactsValueNode.Value)
	if err != nil {
		log.Fatalln(err)
	}

	return path, nil
}

func (c *sourceConfig) GetAddress() (string, error) {
	_, sourceValueNode, err1 := findEntry(c.rootNode, "Address")
	if err1 != nil {
		return "", err1
	}

	return sourceValueNode.Value, nil
}

func (c *sourceConfig) GetPhoneAddress() (string, error) {
	_, phoneAddressValueNode, err1 := findEntry(c.rootNode, "PhoneAddress")
	if err1 != nil {
		return "", err1
	}

	return phoneAddressValueNode.Value, nil
}

func (c *sourceConfig) GetDefaultRegion() (string, error) {
	_, regionValueNode, err1 := findEntry(c.rootNode, "Region")
	if err1 != nil {
		return "", err1
	}

	return regionValueNode.Value, nil
}
