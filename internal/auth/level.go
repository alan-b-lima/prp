package auth

type Level int

const (
	Unlogged Level = iota

	valid_start

	Admin
	User

	valid_end
)

var levelStrings = map[Level]string{
	Admin:    "admin",
	User:     "user",
	Unlogged: "unlogged",
}

func (l Level) IsValid() bool {
	return valid_start < l && l < valid_end
}

func (l Level) IsValidOrUnlogged() bool {
	return l == Unlogged || valid_start < l && l < valid_end
}

func (l Level) String() string {
	return levelStrings[l]
}
