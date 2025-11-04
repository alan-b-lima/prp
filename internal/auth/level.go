package auth

type Level int

const (
	Unlogged Level = iota

	Admin
	User
)

var levelStrings = map[Level]string{
	Admin:    "admin",
	User:     "user",
	Unlogged: "unlogged",
}

func (l Level) String() string {
	return levelStrings[l]
}
