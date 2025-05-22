package main

type Event struct {
	OrderType  string `json:"orderType"`
	SessionID  string `json:"sessionId"`
	Card       string `json:"card"`
	EventDate  string `json:"eventDate"`
	WebsiteURL string `json:"websiteUrl"`
}
