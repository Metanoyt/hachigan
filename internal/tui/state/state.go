package state

type Screen int

const (
	ScreenOverview Screen = iota
	ScreenApplications
	ScreenApplicationDetail
	ScreenCluster
)
