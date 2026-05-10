package pkg

type DebianAptSource struct {
	Types         string
	URIs          string
	Suites        string
	Components    string
	Architectures string
	SignedBy      string
}

func GetPackageList() Dependencies {
	packages := Dependencies{
		BasicDependencies: []string{"curl", "ca-certificates"},
		Docker:            []string{"docker-ce", "docker-ce-cli", "containerd.io", "docker-buildx-plugin", "docker-compose-plugin"},
		DebianDev:         []string{"build-essential", "checkinstall", "libz-dev", "dh-make", "libssl-dev", "devscripts"},
	}
	return packages
}
