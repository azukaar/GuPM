package utils

type GupmEntryPoint struct {
	Name              string
	Version           string
	WrapInstallFolder string
	Git               gupmEntryPointGit
	Publish           gupmEntryPointPublish
	Cli               gupmEntryPointCliList
	Config            gupmEntryPointConfigList
	Dependencies      gupmEntryPointDependenciesList
}

type gupmEntryPointCliList struct {
	DefaultProviders map[string]string
	Aliases          map[string]interface{}
}

type gupmEntryPointGit struct {
	Hooks gupmEntryPointGitHooks
}

type gupmEntryPointGitHooks struct {
	Precommit interface{}
	Prepush   interface{}
}

type gupmEntryPointDependenciesList struct {
	DefaultProvider string
	Default         map[string]string
}

type gupmEntryPointConfigList struct {
	Default gupmEntryPointConfig
}

type gupmEntryPointConfig struct {
	Entrypoint      string
	InstallPath     string
	DefaultProvider string
	OsProviders     map[string]string
}

type gupmEntryPointPublish struct {
	Source []string
	Dest   string
}
