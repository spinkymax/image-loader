package model

type User struct {
	ID          int
	Name        string
	Login       string
	Password    string
	Description string
	ImageUrls   []string
}
