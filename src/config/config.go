package config

type Flags struct {
	Path     string
	Name     string
	Format   bool
	Vet      bool
	Test     bool
	Generate bool
	Server   bool
	Open     bool
	Install  bool
	Build    bool
	Run      bool
	Legacy   bool
	NoConfig bool
}
