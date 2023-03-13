package model

type User struct {
	// Basic Info
	Name         string `json:"name,omitempty"`
	Email        string `json:"email,omitempty"`
	IsActive     bool   `json:"active,omitempty"`
	JobNumber    string `json:"job_number,omitempty"`
	IsEmployee   bool   `json:"is_employee,omitempty"`
	HrEmployeeID string `json:"hr_employee_id,omitempty"`

	// Git Info
	GitUser
}

type GitUser struct {
	GitID        int64  `json:"git_id,omitempty"`
	GitLogin     string `json:"git_login,omitempty"`
	GitAvatarURL string `json:"git_avatar_url,omitempty"`
	GitName      string `json:"git_name,omitempty"`
}
