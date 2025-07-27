package types

type Change struct {
	Type    string
	Summary string
}

type CommitMessage struct {
	Changes     []Change          `yaml:"changes"`
	Issue       string            `yaml:"issue"`
	Subject     string            `yaml:"subject"`
	Description []string          `yaml:"description"`
	Footer      map[string]string `yaml:"footer"`
}
