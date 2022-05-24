package models

type Version struct {
	Semver string
}

func (v Version) String() string {
	return v.Semver
}

func (v Version) IsNewer(o Version) bool {
	return true
}
