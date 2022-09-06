package model

type User struct {
	// Basic Info
	Name  string `json:"name"`
	Email string `json:"email"`
	// Git Info
	GitUser
}

type GitUser struct {
	GitID        int64  `json:"git_id"`
	GitLogin     string `json:"git_login"`
	GitAvatarURL string `json:"git_avatar_url"`
	GitName      string `json:"git_name"`
}
