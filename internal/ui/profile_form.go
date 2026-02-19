package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"sprayer/internal/profile"
)

type ProfileFormState int

const (
	ProfileFormIdle ProfileFormState = iota
	ProfileFormCreate
	ProfileFormEdit
	ProfileFormDeleteConfirm
	ProfileFormImport
)

type ProfileForm struct {
	state        ProfileFormState
	form         *huh.Form
	profile      profile.Profile
	profileStore *profile.Store
	width        int
	height       int
	err          error

	deleteConfirm bool
	importPath    string

	name               string
	keywords           string
	cvPath             string
	coverPath          string
	contactEmail       string
	locations          string
	preferRemote       bool
	minScore           string
	maxScore           string
	excludeTraps       bool
	mustHaveEmail      bool
	preferredTech      string
	avoidTech          string
	preferredCompanies string
	avoidCompanies     string
}

func NewProfileForm(store *profile.Store) *ProfileForm {
	return &ProfileForm{
		state:        ProfileFormIdle,
		profileStore: store,
		minScore:     "0",
		maxScore:     "100",
	}
}

func (pf *ProfileForm) StartCreate() tea.Cmd {
	pf.state = ProfileFormCreate
	pf.profile = profile.NewDefaultProfile()
	pf.profile.ID = fmt.Sprintf("profile-%d", time.Now().Unix())

	pf.resetFields()
	pf.form = pf.buildForm("Create Profile")
	return pf.form.Init()
}

func (pf *ProfileForm) StartEdit(p profile.Profile) tea.Cmd {
	pf.state = ProfileFormEdit
	pf.profile = p

	pf.name = p.Name
	pf.keywords = strings.Join(p.Keywords, ", ")
	pf.cvPath = p.CVPath
	pf.coverPath = p.CoverPath
	pf.contactEmail = p.ContactEmail
	pf.locations = strings.Join(p.Locations, ", ")
	pf.preferRemote = p.PreferRemote
	pf.minScore = strconv.Itoa(p.MinScore)
	pf.maxScore = strconv.Itoa(p.MaxScore)
	pf.excludeTraps = p.ExcludeTraps
	pf.mustHaveEmail = p.MustHaveEmail
	pf.preferredTech = strings.Join(p.PreferredTech, ", ")
	pf.avoidTech = strings.Join(p.AvoidTech, ", ")
	pf.preferredCompanies = strings.Join(p.PreferredCompanies, ", ")
	pf.avoidCompanies = strings.Join(p.AvoidCompanies, ", ")

	pf.form = pf.buildForm("Edit Profile")
	return pf.form.Init()
}

func (pf *ProfileForm) StartDeleteConfirm(p profile.Profile) tea.Cmd {
	pf.state = ProfileFormDeleteConfirm
	pf.profile = p
	pf.deleteConfirm = false

	pf.form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Delete profile '%s'?", p.Name)).
				Description("This action cannot be undone.").
				Value(&pf.deleteConfirm).
				Affirmative("Yes, delete").
				Negative("No, cancel"),
		),
	)

	return pf.form.Init()
}

func (pf *ProfileForm) StartImport() tea.Cmd {
	pf.state = ProfileFormImport
	pf.importPath = ""

	pf.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Import Profile from JSON").
				Description("Enter the path to the profile JSON file").
				Value(&pf.importPath).
				Placeholder("/path/to/profile.json"),
		),
	)

	return pf.form.Init()
}

func (pf *ProfileForm) resetFields() {
	pf.name = ""
	pf.keywords = ""
	pf.cvPath = ""
	pf.coverPath = ""
	pf.contactEmail = ""
	pf.locations = ""
	pf.preferRemote = false
	pf.minScore = "0"
	pf.maxScore = "100"
	pf.excludeTraps = true
	pf.mustHaveEmail = false
	pf.preferredTech = ""
	pf.avoidTech = ""
	pf.preferredCompanies = ""
	pf.avoidCompanies = ""
}

func (pf *ProfileForm) buildForm(title string) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Profile Name").
				Value(&pf.name).
				Placeholder("My Profile").
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("profile name is required")
					}
					return nil
				}),

			huh.NewInput().
				Title("Keywords (comma-separated)").
				Value(&pf.keywords).
				Placeholder("golang, rust, remote"),

			huh.NewInput().
				Title("CV Path").
				Value(&pf.cvPath).
				Placeholder("/path/to/cv.pdf"),

			huh.NewInput().
				Title("Cover Letter Path").
				Value(&pf.coverPath).
				Placeholder("/path/to/cover.pdf"),

			huh.NewInput().
				Title("Contact Email").
				Value(&pf.contactEmail).
				Placeholder("me@example.com"),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Preferred Locations (comma-separated)").
				Value(&pf.locations).
				Placeholder("remote, san francisco, new york"),

			huh.NewConfirm().
				Title("Prefer Remote?").
				Value(&pf.preferRemote),

			huh.NewInput().
				Title("Minimum Score (0-100)").
				Value(&pf.minScore).
				Placeholder("0"),

			huh.NewInput().
				Title("Maximum Score (0-100)").
				Value(&pf.maxScore).
				Placeholder("100"),

			huh.NewConfirm().
				Title("Exclude Traps?").
				Value(&pf.excludeTraps),

			huh.NewConfirm().
				Title("Must Have Email?").
				Value(&pf.mustHaveEmail),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Preferred Tech (comma-separated)").
				Value(&pf.preferredTech).
				Placeholder("go, rust, kubernetes"),

			huh.NewInput().
				Title("Avoid Tech (comma-separated)").
				Value(&pf.avoidTech).
				Placeholder("cobol, fortran"),

			huh.NewInput().
				Title("Preferred Companies (comma-separated)").
				Value(&pf.preferredCompanies).
				Placeholder("google, meta, apple"),

			huh.NewInput().
				Title("Avoid Companies (comma-separated)").
				Value(&pf.avoidCompanies).
				Placeholder("palantir, theranos"),
		),
	).WithWidth(pf.width).WithHeight(pf.height)
}

func (pf *ProfileForm) SetSize(width, height int) {
	pf.width = width
	pf.height = height
}

func (pf *ProfileForm) Update(msg tea.Msg) (tea.Cmd, tea.Cmd) {
	if pf.form == nil {
		return nil, nil
	}

	model, cmd := pf.form.Update(msg)
	pf.form = model.(*huh.Form)

	if pf.form.State == huh.StateCompleted {
		switch pf.state {
		case ProfileFormCreate, ProfileFormEdit:
			return cmd, pf.saveProfile()
		case ProfileFormDeleteConfirm:
			if pf.deleteConfirm {
				return cmd, pf.deleteProfile()
			}
			return cmd, func() tea.Msg { return ProfileOperationCancelledMsg{} }
		case ProfileFormImport:
			return cmd, pf.importProfile()
		}
	}

	return cmd, nil
}

func (pf *ProfileForm) saveProfile() tea.Cmd {
	return func() tea.Msg {
		pf.profile.Name = pf.name
		pf.profile.Keywords = parseCommaList(pf.keywords)
		pf.profile.CVPath = pf.cvPath
		pf.profile.CoverPath = pf.coverPath
		pf.profile.ContactEmail = pf.contactEmail
		pf.profile.Locations = parseCommaList(pf.locations)
		pf.profile.PreferRemote = pf.preferRemote
		pf.profile.MinScore, _ = strconv.Atoi(pf.minScore)
		pf.profile.MaxScore, _ = strconv.Atoi(pf.maxScore)
		pf.profile.ExcludeTraps = pf.excludeTraps
		pf.profile.MustHaveEmail = pf.mustHaveEmail
		pf.profile.PreferredTech = parseCommaList(pf.preferredTech)
		pf.profile.AvoidTech = parseCommaList(pf.avoidTech)
		pf.profile.PreferredCompanies = parseCommaList(pf.preferredCompanies)
		pf.profile.AvoidCompanies = parseCommaList(pf.avoidCompanies)

		if pf.profile.ID == "" {
			pf.profile.ID = fmt.Sprintf("profile-%d", time.Now().Unix())
		}

		if pf.profileStore == nil {
			return ProfileErrorMsg{Err: fmt.Errorf("profile store not initialized")}
		}

		if err := pf.profileStore.Save(pf.profile); err != nil {
			return ProfileErrorMsg{Err: err}
		}

		if pf.state == ProfileFormCreate {
			return ProfileCreatedMsg{Profile: pf.profile}
		}
		return ProfileUpdatedMsg{Profile: pf.profile}
	}
}

func (pf *ProfileForm) deleteProfile() tea.Cmd {
	return func() tea.Msg {
		if pf.profileStore == nil {
			return ProfileErrorMsg{Err: fmt.Errorf("profile store not initialized")}
		}

		if err := pf.profileStore.Delete(pf.profile.ID); err != nil {
			return ProfileErrorMsg{Err: err}
		}

		return ProfileDeletedMsg{Profile: pf.profile}
	}
}

func (pf *ProfileForm) importProfile() tea.Cmd {
	return func() tea.Msg {
		if pf.importPath == "" {
			return ProfileErrorMsg{Err: fmt.Errorf("no import path specified")}
		}

		data, err := os.ReadFile(pf.importPath)
		if err != nil {
			return ProfileErrorMsg{Err: fmt.Errorf("failed to read file: %w", err)}
		}

		var importedProfile profile.Profile
		if err := json.Unmarshal(data, &importedProfile); err != nil {
			return ProfileErrorMsg{Err: fmt.Errorf("failed to parse profile JSON: %w", err)}
		}

		if importedProfile.ID == "" {
			importedProfile.ID = fmt.Sprintf("imported-%d", time.Now().Unix())
		}

		if pf.profileStore == nil {
			return ProfileErrorMsg{Err: fmt.Errorf("profile store not initialized")}
		}

		if err := pf.profileStore.Save(importedProfile); err != nil {
			return ProfileErrorMsg{Err: err}
		}

		return ProfileImportedMsg{Profile: importedProfile}
	}
}

func (pf *ProfileForm) View() string {
	if pf.form == nil {
		return ""
	}
	return pf.form.View()
}

func (pf *ProfileForm) State() ProfileFormState {
	return pf.state
}

func (pf *ProfileForm) Reset() {
	pf.state = ProfileFormIdle
	pf.form = nil
	pf.err = nil
	pf.deleteConfirm = false
	pf.importPath = ""
}

func parseCommaList(s string) []string {
	if s == "" {
		return nil
	}
	items := strings.Split(s, ",")
	result := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

type ProfileCreatedMsg struct {
	Profile profile.Profile
}

type ProfileUpdatedMsg struct {
	Profile profile.Profile
}

type ProfileDeletedMsg struct {
	Profile profile.Profile
}

type ProfileImportedMsg struct {
	Profile profile.Profile
}

type ProfileErrorMsg struct {
	Err error
}

type ProfileOperationCancelledMsg struct{}
