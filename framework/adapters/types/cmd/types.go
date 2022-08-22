package cmd

type Flags struct {
	AwsZone            *string
	InstanceIDs        []string
	BashScriptLocation *string
	LogLevel           *string
	OutputLocation     *string
	FreeFormCmd        *string
	AwsProfile         *string
}
