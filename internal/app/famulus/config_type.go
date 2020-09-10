package config

import (
	"errors"

	"gopkg.in/yaml.v3"
)

type Config interface {
	ListAddress() ([]string, error)
	GetAddress(string) (AddressConfig, error)
	GetContactsCsvFilePath() (string, error)
}

type AddressConfig interface {
	GetAddress() (string, error)
	GetPhoneAddress() (string, error)
}

type fileConfig struct {
	documentRoot *yaml.Node
}

type addressConfig struct {
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

func getAddressMap(c *fileConfig) (map[string]AddressConfig, error) {
	_, addressesValueNode, err1 := findEntry(c.documentRoot.Content[0], "Addresses")
	if err1 != nil {
		return nil, err1
	}
	addressMap, err2 := getMapEntry(addressesValueNode)
	if err2 != nil {
		return nil, err2
	}

	result := make(map[string]AddressConfig)

	/* Copy Content from Map1 to Map2*/
	for index, element := range addressMap {
		result[index] = &addressConfig{
			rootNode: element,
		}
	}

	return result, nil
}

func (c *fileConfig) ListAddress() ([]string, error) {
	addressMap, err := getAddressMap(c)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(addressMap))

	for k := range addressMap {
		keys = append(keys, k)
	}
	return keys, nil
}

func (c *fileConfig) GetAddress(name string) (AddressConfig, error) {
	addressMap, err := getAddressMap(c)
	if err != nil {
		return nil, err
	}
	return addressMap[name], nil
}

func (c *fileConfig) GetContactsCsvFilePath() (string, error) {
	_, contactsValueNode, err1 := findEntry(c.documentRoot.Content[0], "Contacts")
	if err1 != nil {
		return "", err1
	}

	return contactsValueNode.Value, nil
}

func (c *addressConfig) GetAddress() (string, error) {
	_, addressValueNode, err1 := findEntry(c.rootNode, "Address")
	if err1 != nil {
		return "", err1
	}

	return addressValueNode.Value, nil
}

func (c *addressConfig) GetPhoneAddress() (string, error) {
	_, phoneAddressValueNode, err1 := findEntry(c.rootNode, "PhoneAddress")
	if err1 != nil {
		return "", err1
	}

	return phoneAddressValueNode.Value, nil
}
