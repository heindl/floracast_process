package name_usage

type CanonicalNameMatcher interface {
	MatchCanonicalNames(names ...string) ([]*NameUsageSource, error)
}
