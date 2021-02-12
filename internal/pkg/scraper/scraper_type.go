package scraper

type contact struct {
	numbers     string
	website     string
	title       string
	description string
}

type Contact interface {
	GetNumbers() string
	GetWebsite() string
	GetTitle() string
	GetDescription() string
}

func (c *contact) GetNumbers() string {
	return c.numbers
}

func (c *contact) GetWebsite() string {
	return c.website
}

func (c *contact) GetTitle() string {
	return c.title
}

func (c *contact) GetDescription() string {
	return c.description
}
