package provider

type GupmEntryPoint struct {
	Name string
	Version string
	Cli gupmEntryPointCliList
	Config gupmEntryPointConfigList
}

type gupmEntryPointCliList  struct {
	Aliases map[string]interface {}
}

type gupmEntryPointConfigList  struct {
	Default gupmEntryPointConfig
}

type gupmEntryPointConfig struct {
	Entrypoint string
	InstallPath string
}