package client

type Installer interface {
	Name() string
	Detect() (found bool, where string)
	Install(token, baseURL, model string) error
	Uninstall() error
}

func All() []Installer {
	return []Installer{
		&codexCLI{},
		&codexVSCode{},
		&codexApp{},
		&openCode{},
	}
}
