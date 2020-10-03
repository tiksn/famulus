package scraper

type contact struct {
	numbers string
	website string
}

type Contact interface {
	GetNumbers() string
	GetWebsite() string
}

func (c *contact) GetNumbers() string {
	return c.numbers
}

func (c *contact) GetWebsite() string {
	return c.website
}
