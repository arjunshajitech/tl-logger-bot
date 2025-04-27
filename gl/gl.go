package gl

type Gitlab struct {
	Token    string
	Projects []Project
}

type Project struct {
	ID     int
	Branch []string
	Author string
}

func NewGitLab(token string, projects []Project) *Gitlab {
	return &Gitlab{token, projects}
}
