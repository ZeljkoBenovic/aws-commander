package conf

import (
	"errors"
	"flag"
	"log"
)

var (
	ErrNoBashCMDOrScriptProvided         = errors.New("bash cmd or script not provided")
	ErrAnsiblePlaybookNotProvided        = errors.New("ansible playbook not provided")
	ErrEC2TagsOrIDsNotSpecified          = errors.New("ec2 instance ids or tags not specified")
	ErrEC2TagsAndIDsAreMutuallyExclusive = errors.New("ec2 instance ids and tags are mutually exclusive")
)

type Config struct {
	LogLevel string
	Mode     string

	BashOneLiner string
	BashFile     string

	AnsiblePlaybook  string
	AnsibleURL       string
	AnsibleExtraVars string
	AnsibleDryRun    bool

	AWSProfile string
	AWSRegion  string

	AWSInstanceIDs  string
	AWSInstanceTags string

	CommandResultMaxWait int
	CommandExecMaxWait   int64
}

func New() Config {
	conf := DefaultConfig()

	conf.processFlags()

	if err := conf.validateFlags(); err != nil {
		log.Fatalln(err)
	}

	return conf
}

func DefaultConfig() Config {
	return Config{
		LogLevel:             "info",
		Mode:                 "bash",
		BashOneLiner:         "",
		BashFile:             "",
		AnsiblePlaybook:      "",
		AnsibleURL:           "",
		AWSProfile:           "",
		AWSRegion:            "",
		AWSInstanceIDs:       "",
		AWSInstanceTags:      "",
		AnsibleExtraVars:     "",
		CommandResultMaxWait: 30,
		CommandExecMaxWait:   300,
		AnsibleDryRun:        false,
	}
}

func (c *Config) processFlags() {
	flag.StringVar(&c.LogLevel, "log-level", c.LogLevel,
		"log output level",
	)
	flag.StringVar(&c.Mode, "mode", c.Mode,
		"running mode",
	)
	flag.StringVar(&c.BashOneLiner, "cmd", c.BashOneLiner,
		"bash command to run",
	)
	flag.StringVar(&c.BashFile, "script", c.BashFile,
		"bash script to run",
	)
	flag.StringVar(&c.AnsiblePlaybook, "playbook", c.AnsiblePlaybook,
		"ansible playbook to run",
	)
	flag.StringVar(&c.AnsibleURL, "ansible-url", c.AnsibleURL,
		"ansible url where the playbook can be read from",
	)
	flag.StringVar(&c.AWSProfile, "profile", c.AWSProfile,
		"aws profile",
	)
	flag.StringVar(&c.AWSRegion, "region", c.AWSRegion,
		"aws region",
	)
	flag.StringVar(&c.AWSInstanceIDs, "ids", c.AWSInstanceIDs,
		"comma delimited list of aws ec2 ids",
	)
	flag.StringVar(&c.AWSInstanceTags, "tags", c.AWSInstanceTags,
		"semi-column delimited list of ec2 tags (\"foo=bar,baz;faz=baz,bar\")",
	)
	flag.IntVar(&c.CommandResultMaxWait, "max-wait", c.CommandResultMaxWait,
		"maximum wait time in seconds for command execution",
	)
	flag.Int64Var(&c.CommandExecMaxWait, "max-exec", c.CommandExecMaxWait,
		"maximum command execution time in seconds",
	)
	flag.BoolVar(&c.AnsibleDryRun, "dryrun", c.AnsibleDryRun,
		"run ansible in dry-run mode",
	)
	flag.StringVar(&c.AnsibleExtraVars, "extra-vars", c.AnsibleExtraVars,
		"comma separated key value pairs for extra vars (foo=bar,fus=baz)",
	)
	flag.Parse()
}

func (c *Config) validateFlags() error {
	if c.Mode == "bash" && c.BashFile == "" && c.BashOneLiner == "" {
		return ErrNoBashCMDOrScriptProvided
	}

	if c.Mode == "ansible" && c.AnsiblePlaybook == "" {
		return ErrAnsiblePlaybookNotProvided
	}

	if c.AWSInstanceIDs == "" && c.AWSInstanceTags == "" {
		return ErrEC2TagsOrIDsNotSpecified
	}

	if c.AWSInstanceTags != "" && c.AWSInstanceIDs != "" {
		return ErrEC2TagsAndIDsAreMutuallyExclusive
	}

	return nil
}
