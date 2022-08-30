# AWS Commander

A tool used for running bash scripts on AWS EC2 instances, leveraging AWS Systems Manager > Run Command feature.   
User can load a bash script or define a single command, that will execute on all instances with defined instance ID.

## Prerequisites

* The **AmazonSSMManagedInstanceCore** must be placed on all instances that need to be managed via this tool.    
* AWS API credentials defined in *aws credentials* file or as environment variables

## Usage

### AWS credentials
AWS credentials can be pulled from environment variables or from aws credentials file.   
To define a which profile from credentials file should be used, set `aws-profile` flag. By default, it is set to `default`.   
Environment variables with credentials that can be set:
* `AWS_ACCESS_KEY_ID` - the aws access key id
* `AWS_SECRET_ACCESS_KEY` - the access key secret
* `AWS_SESSION_TOKEN` - the session token (optional)


### General Parameters
* `aws-profile` - AWS profile as defined in *aws credentials* file. Default: `default`
* `aws-zone` - AWS zone in which EC2 instances reside. Default: `eu-central-1`
* `instances` - instance IDs, separated by comma (,). This is a mandatory flag.
* `log-level` - the level of logging output (info, debug, error). Default: `info`
* `output` - a file name to write the output result of a command/script. Default: `console output`
* `mode` - switch between modes - Bash script or Ansible playbook. Default: `bash`

### Running Bash scripts
* `cmd` - one-liner bash command that will be executed on EC2 instances.
* `script` - the location of bash script file that will run on EC2 instances.
* `mode` - for running bash scripts `mode` can be omitted as the default value is `bash`

If both `cmd` and `script` flags are defined, `script` will take precedence, and `cmd` will be disregarded.

#### Example

```bash
aws-commander -instances i-0bf9c273c67f684a0,i-011c9b3e3607a63b5,i-0e53e37f7b34517f5,i-0f02ca10faf8f349e -cmd "cd /tmp && ls -lah" -aws-profile test-account
```

### Running Ansible Playbook
* `playbook` - the location of Ansible playbook that will be executed on EC2 instances.
* `dryrun` - when set to true, Ansible playbook will run and the output will be shown, but 
  no data will be changed. 
* `mode` - for running Ansible playbook `mode` must be set to `ansible`

#### Ansible prerequisites
Every EC2 instance, that should run Ansible playbook, must have Ansible already installed.   
If Ansible is not installed, the deployment will fail.    
You can use `bash` mode to simply install Ansible from your OS package manager before running the playbook.

#### Example
```bash
## if Ansible is not installed on host - install Ansible
aws-commander -instances i-0bf9c273c67f684a0,i-011c9b3e3607a63b5,i-0e53e37f7b34517f5,i-0f02ca10faf8f349e -cmd "sudo apt install -y ansible" -aws-profile test-account -aws-zone us-west-2
## run playbook
aws-commander -instances i-0bf9c273c67f684a0,i-011c9b3e3607a63b5,i-0e53e37f7b34517f5,i-0f02ca10faf8f349e -mode ansible -playbook scripts/nodes-restart.yaml -aws-profile test-account -aws-zone us-west-2
```

#### Missing features
Currently, running the Ansible playbook from a remote location via URL / S3 is not supported.    
It will be supported in the future release.

## License

Copyright 2022 Trapesys

Licensed under the Apache License, Version 2.0 (the “License”); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

### http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an “
AS IS” BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
