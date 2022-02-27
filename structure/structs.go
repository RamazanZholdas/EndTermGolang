package structure

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type GamePage struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	Rating    string    `json:"rating"`
	Developer string    `json:"developer"`
	Genre     string    `json:"genre"`
	VideoUrl  string    `json:"url"`
	PhotoUrl  string    `json:"photoUrl"`
	Comments  []Comment `json:"comments:"`
}

type Comment struct {
	Username   string `json:"username"`
	UserRating string `json:"userRating"`
	Review     string `json:"review"`
}

type Auth struct {
	IsAuth   bool
	Username string
}
