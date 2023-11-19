# AWS Commander

A tool used for easier automation of the AWS EC2 instances, leveraging AWS Systems Manager - Run Command feature.   
Supported scripts:
* One-liner bash command
* Bash script loaded from a local filesystem
* Ansible playbook loaded from a local filesystem

The command/script/playbook will run across all EC2 instances simultaneously.     
EC2 instances can be targeted by instance IDs or tags.


## Prerequisites

* The **AmazonSSMManagedInstanceCore** IAM role, attached on all EC2 instances.    
* Authenticated AWS CLI session 

## Usage

### AWS credentials
AWS access must be authenticated via `aws cli`.

### General Parameters
* `log-level` - the level of logging output (`info`, `debug`, `error`). Default: `info`
* `mode` - commands running mode (`bash`, `ansible`) Default: `bash`
* `profile` - AWS profile as defined in *aws credentials* file.
* `region` - AWS region in which EC2 instances reside.
* `ids` - instance IDs, separated by comma (`,`). This is a mandatory flag.
* `tags` - instance tags. Tags are semicolon delimited key - multiple value pairs (example: `Name=bar,baz;Role=foo,faz`)
* `max-wait` - maximum wait time in seconds to run the command Default: `30`
* `max-exec` - maximum wait time in seconds to get command result Default: `300`

### Running Bash scripts
* `cmd` - one-liner bash command that will be executed on EC2 instances.
* `script` - the location of bash script file that will run on EC2 instances.
* `mode` - for running Bash script or oneliner `mode` can be omitted or set to `bash`

#### Example

```bash
# AWS authentication
aws sso login --profile test-account

# oneliner using instance ids
aws-commander -instances i-0bf9c273c67f684a0,i-011c9b3e3607a63b5,i-0e53e37f7b34517f5,i-0f02ca10faf8f349e -cmd "cd /tmp && ls -lah" -aws-profile test-account

# or bash script using instance ids
aws-commander -instances i-0bf9c273c67f684a0,i-011c9b3e3607a63b5,i-0e53e37f7b34517f5,i-0f02ca10faf8f349e -script ./script.sh -aws-profile test-account

# or oneliner using tags
aws-commander -tags "Name=Test,Test2,Test3;Role=test" -cmd "cd /tmp && ls -lah" -aws-profile test-account

# or bash script using tags
aws-commander -tags "Name=Test,Test2,Test3;Role=test" -script ./script.sh -aws-profile test-account
```

### Running Ansible Playbook
* `playbook` - the location of Ansible playbook that will be executed on EC2 instances.
* `ansible-url` - the URL locaction of the Ansible playbook
* `extra-vars` - comma delimited, key value pairs of Ansible variables
* `dryrun` - when set to true, Ansible playbook will run and the output will be shown, but 
  no data will be changed. Default: `false`
* `mode` - for running Ansible playbook `mode` must be set to `ansible`

#### Example
```bash
# AWS authentication
aws sso login

# run local playbook using instance ids
aws-commander -instances i-0bf9c273c67f684a0,i-011c9b3e3607a63b5,i-0e53e37f7b34517f5,i-0f02ca10faf8f349e -mode ansible -playbook scripts/init.yaml -extra-vars foo=bar,faz=baz

# or from url using instance ids
aws-commander -instances i-0bf9c273c67f684a0,i-011c9b3e3607a63b5,i-0e53e37f7b34517f5,i-0f02ca10faf8f349e -mode ansible -ansible-url https://example.com/init.yaml -extra-vars foo=bar,faz=baz

# run local playbook using tags
aws-commander -tags "Name=Test,Test2,Test3;Role=test" -mode ansible -playbook scripts/init.yaml -extra-vars foo=bar,faz=baz

# or from url using tags
aws-commander -tags "Name=Test,Test2,Test3;Role=test" -mode ansible -ansible-url https://example.com/init.yaml -extra-vars foo=bar,faz=baz
```