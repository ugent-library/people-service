package cli

type Version struct {
	Branch string `env:"SOURCE_BRANCH" json:"branch,omitempty"`
	Commit string `env:"SOURCE_COMMIT" json:"commit,omitempty"`
	Image  string `env:"IMAGE_NAME" json:"image,omitempty"`
}
