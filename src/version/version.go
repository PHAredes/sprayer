package version

var Version = "dev"

func WithPrefix() string {
	return "v" + Version
}
