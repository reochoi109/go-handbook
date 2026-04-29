package domain

type User struct {
	ID   int
	Name string
}

type Post struct {
	ID     int
	UserID int
	Title  string
}
