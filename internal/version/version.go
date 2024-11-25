package version

//nolint:gochecknoglobals
var (
	Version   = "no-version-tag"
	BuildDate string
	Commit    string
)

func String() string {
	return Version + ", " + Commit + ", build at " + BuildDate
}
