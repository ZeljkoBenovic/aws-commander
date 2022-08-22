package ssm

type Instances struct {
	Instance []Instance `json:"instance"`
}

type Instance struct {
	ID            string `json:"id"`
	CommandOutput string `json:"commandOutput"`
	ErrorOutput   string `json:"errorOutput"`
}
