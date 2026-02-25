package tui

import "sprayer/src/api/job"

// ScrapeProgressMsg reports incremental scraping progress.
type ScrapeProgressMsg struct {
	Source  string
	Current int
	Total   int
	Found   int
}

// ScrapeCompleteMsg signals scraping has finished.
type ScrapeCompleteMsg struct {
	Jobs []job.Job
	Err  error
}

// FilterApplyMsg carries updated filter values from the filter screen.
type FilterApplyMsg struct {
	Keywords   string
	Exclude    string
	Locations  string
	Companies  string
	MinScore   string
}

// FilterCancelMsg signals the user cancelled the filter screen.
type FilterCancelMsg struct{}

// ProfileSelectedMsg signals the user activated a profile.
type ProfileSelectedMsg struct {
	Index int
}

// ProfileNewMsg signals the user wants to create a new profile.
type ProfileNewMsg struct{}

// ProfileBackMsg signals returning from profiles.
type ProfileBackMsg struct{}

// HelpCloseMsg signals help was closed.
type HelpCloseMsg struct{}

// EmailOpenComposeMsg signals opening compose for an email.
type EmailOpenComposeMsg struct {
	Index int
}

// EmailBackMsg signals returning from emails.
type EmailBackMsg struct{}

// ComposeBackMsg signals returning from compose.
type ComposeBackMsg struct{}

// CVWizardDoneMsg signals the wizard completed.
type CVWizardDoneMsg struct {
	ProfileName string
}

// CVWizardCancelMsg signals the wizard was cancelled.
type CVWizardCancelMsg struct{}
