package cmd

type Flags struct {
	AwsZone            FlagDetails
	InstanceIDs        FlagDetails
	BashScriptLocation FlagDetails
	LogLevel           FlagDetails
	OutputLocation     FlagDetails
	FreeFormCmd        FlagDetails
	AwsProfile         FlagDetails
	Mode               FlagDetails
	AnsiblePlaybook    FlagDetails
	AnsibleDryRun      FlagDetails
}

type FlagDetails struct {
	Name  string
	Usage string

	DefaultString string
	DefaultInt    int
	DefaultBool   bool

	ValueString    *string
	ValueStringArr []string
	ValueInt       *int
	ValueBool      *bool

	AllowedValuesStr []string
}

var UserFlags = Flags{
	AwsZone: FlagDetails{
		Name:          "aws-zone",
		Usage:         "aws zone where instances reside",
		DefaultString: "eu-central-1",
		ValueString:   new(string),
	},
	InstanceIDs: FlagDetails{
		Name:           "instances",
		Usage:          "instance IDs, separated by comma (,)",
		DefaultString:  "",
		ValueStringArr: make([]string, 0),
	},
	BashScriptLocation: FlagDetails{
		Name:          "script",
		Usage:         "the location of the script to run",
		DefaultString: "",
		ValueString:   new(string),
	},
	LogLevel: FlagDetails{
		Name:          "log-level",
		Usage:         "log output level",
		DefaultString: "info",
		ValueString:   new(string),
	},
	OutputLocation: FlagDetails{
		Name:          "output",
		Usage:         "the location of file to write json output (default: output to console)",
		DefaultString: "",
		ValueString:   new(string),
	},
	FreeFormCmd: FlagDetails{
		Name:          "cmd",
		Usage:         "freeform command, a single line bash command to be executed",
		DefaultString: "",
		ValueString:   new(string),
	},
	AwsProfile: FlagDetails{
		Name:          "aws-profile",
		Usage:         "aws credentials profile",
		DefaultString: "default",
		ValueString:   new(string),
	},
	Mode: FlagDetails{
		Name:             "mode",
		Usage:            "set command mode - bash script or ansible playbook",
		DefaultString:    "bash",
		ValueString:      new(string),
		AllowedValuesStr: []string{"bash", "ansible"},
	},
	AnsiblePlaybook: FlagDetails{
		Name:          "playbook",
		Usage:         "the location of Ansible playbook file",
		DefaultString: "",
		ValueString:   new(string),
	},
	AnsibleDryRun: FlagDetails{
		Name:        "dryrun",
		Usage:       "run Ansible script without changing any actual data",
		DefaultBool: false,
		ValueBool:   new(bool),
	},
}
