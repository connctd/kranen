package main

type PushData struct {
	Images   []string
	PushedAt float64 `json:"pushed_at"`
	Pusher   string
	Tag      string `json:"tag"`
}

type Repository struct {
	CommentCount     int     `json:"comment_count"`
	DateCreated      float64 `json:"date_created"`
	Description      string
	FulleDescription string `json:"full_description, omitempty"`
	Dockerfile       string `json:"_, omitempty"`
	Official         bool   `json:"is_official"`
	Private          bool   `json:"is_private"`
	Trusted          bool   `json:"is_trusted"`
	Name             string
	Namespace        string
	Owner            string
	RepoName         string `json:"repo_name"`
	RepoUrl          string `json:"repo_url"`
	StarCount        int    `json:"star_count"`
	Status           string
}

type Payload struct {
	CallbackUrl string      `json:"callback_url"`
	PushData    *PushData   `json:"push_data"`
	Repo        *Repository `json:"repository"`
}
