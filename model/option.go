package model

// OptionData actual data within an option
type OptionData struct {
	Value string `bson:"value" json:"value" example:"ACTURSCI"`
	Text  string `bson:"text" json:"text"  example:"Actuarial Sciences"`
}

// Option struct to define a form option to be stored in database
type Option struct {
	Source SourceInfo `bson:"source" json:"source"`
	Time   TimeInfo   `bson:"time" json:"time"`
	Data   OptionData `bson:"data" json:"data"`
}
