package mqtt

type RuleList struct {
	BaseMessage
	Rules 			[]Rule	`json:"devices"`
}

type Rule struct {
	Id 				string		`json:"id"`
	Name 			string		`json:"name"`
	Description 	string		`json:"description"`
}
