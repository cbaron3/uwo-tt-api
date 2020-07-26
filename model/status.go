package model

// ScraperStatus track
type ScraperStatus struct {
	Active bool `bson:"active" json:"active"`
	// Need time and also x/documents
	// Options or courses
}

// Status track
type Status struct {
	Scraper ScraperStatus `bson:"scraper" json:"scraper"`
}
