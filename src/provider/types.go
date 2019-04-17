package provider

type GupmEntryPoint struct {
	Name string
	Version string
	Config gupmEntryPointConfigList
}

type gupmEntryPointConfigList  struct {
	Default gupmEntryPointConfig
}

type gupmEntryPointConfig struct {
	Entrypoint string
	InstallPath string
}