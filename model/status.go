package model

// ScraperStatus track
type ScraperStatus struct {
	Active bool `bson:"active" json:"active"`
}

// Status track
type Status struct {
	Scraper ScraperStatus `bson:"scraper" json:"scraper"`
}
