package scraper

type contact struct {
	numbers string
}

type Contact interface {
	GetNumbers() string
}

func (c *contact) GetNumbers() string {
	return c.numbers
}
