package pkg

import (
	"log/slog"
	"os"
)

func GetPackageList() Dependencies {
	packages := Dependencies{
		InitialDependencies: []string{"sudo"},
		BasicDependencies:   []string{"curl", "ca-certificates", "vim"},
		Docker:              []string{"docker-ce", "docker-ce-cli", "containerd.io", "docker-buildx-plugin", "docker-compose-plugin"},
		DebianDev:           []string{"build-essential", "checkinstall", "libz-dev", "dh-make", "libssl-dev", "devscripts"},
	}
	return packages
}

func GenerateTemplates(c *Config) error {
	if err := generateBashRCTemplate(c.homeDiretory + "/" + ".bashrc"); err != nil {
		return err
	}
	return nil
}

func generateBashRCTemplate(fileKey string) error {
	contents :=
		`# Custom PS1
export PS1='\u@\h \W $ '

# User Aliases
alias ls='ls -al'

# Terraform
alias tf="terraform"
alias tfv="terraform validate"
alias tfp="terraform plan"
alias tfi="terraform init"
alias tff="terraform fmt -recursive"
alias tfaa="terraform apply -auto-approve"
alias tfdaa="terraform destroy -auto-approve"

# Github
alias shbr='git branch --show-current'

# AWS
alias awsc="aws configure"
alias awswho="aws sts get-caller-identity"

# Extend PATH
export PATH=$PATH:$HOME/.scriptbin:/.tfenv/bin

# User Stuff

# Extend path
export PATH=$PATH:/usr/local/go/bin:$HOME/tools/bin

# Google
alias g='gcloud'`
	file, err := os.OpenFile(fileKey, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("error closing file stream", "error", err)
		}
	}()
	if _, err := file.WriteString(contents); err != nil {
		return err
	}
	return nil
}
