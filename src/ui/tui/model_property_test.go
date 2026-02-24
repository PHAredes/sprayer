package tui

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"testing/quick"
	"time"

	"github.com/charmbracelet/bubbletea"

	"sprayer/src/api/job"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomString(r *rand.Rand, length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func genModelFromSeed(seed int64) Model {
	r := rand.New(rand.NewSource(seed))
	numJobs := r.Intn(101)
	jobs := make([]job.Job, numJobs)
	for i := 0; i < numJobs; i++ {
		jobs[i] = job.Job{
			ID:          fmt.Sprintf("job-%d", r.Intn(10000)),
			Title:       randomString(r, 10),
			Company:     randomString(r, 8),
			Location:    randomString(r, 6),
			Description: randomString(r, 20),
			URL:         "https://example.com",
			Source:      randomString(r, 5),
			PostedDate:  time.Now(),
			Salary:      "",
			JobType:     "",
			Email:       "",
			Score:       r.Intn(100),
			HasTraps:    r.Float32() > 0.5,
			Traps:       nil,
			Applied:     false,
		}
	}

	selectedIndex := 0
	if numJobs > 0 {
		selectedIndex = r.Intn(numJobs)
	}

	profileName := randomString(r, 8)
	if profileName == "" {
		profileName = "Default"
	}

	viewState := JobList
	if numJobs == 0 {
		viewState = EmptyState
	}

	width := r.Intn(100) + 10
	height := r.Intn(50) + 5

	return Model{
		jobs:          jobs,
		selectedIndex: selectedIndex,
		profileName:   profileName,
		viewState:     viewState,
		width:         width,
		height:        height,
	}
}

func TestProperty_NavigationStaysInBounds(t *testing.T) {
	f := func(seed int64) bool {
		m := genModelFromSeed(seed)
		originalIndex := m.selectedIndex

		r := rand.New(rand.NewSource(seed))
		numDownPresses := r.Intn(200)
		for i := 0; i < numDownPresses; i++ {
			newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
			m = newModel.(Model)
		}

		if len(m.jobs) > 0 && m.selectedIndex > len(m.jobs)-1 {
			t.Errorf("selectedIndex %d exceeds len(jobs)-1 %d after down navigation", m.selectedIndex, len(m.jobs)-1)
			return false
		}

		m.selectedIndex = originalIndex

		numUpPresses := r.Intn(200)
		for i := 0; i < numUpPresses; i++ {
			newModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
			m = newModel.(Model)
		}

		if m.selectedIndex < 0 {
			t.Errorf("selectedIndex %d is below 0 after up navigation", m.selectedIndex)
			return false
		}

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Errorf("Navigation bounds property failed: %v", err)
	}
}

func TestProperty_ViewAlwaysRenders(t *testing.T) {
	f := func(seed int64) bool {
		m := genModelFromSeed(seed)

		if m.width < 10 {
			m.width = 10
		}
		if m.height < 5 {
			m.height = 5
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("View() panicked: %v", r)
				}
			}()
			view := m.View()
			if view == "" {
				t.Errorf("View() returned empty string")
			}
		}()

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Errorf("View always renders property failed: %v", err)
	}
}

func TestProperty_JobListRendersSelectedJob(t *testing.T) {
	f := func(seed int64) bool {
		m := genModelFromSeed(seed)

		if len(m.jobs) == 0 {
			return true
		}

		if m.selectedIndex < 0 || m.selectedIndex >= len(m.jobs) {
			m.selectedIndex = len(m.jobs) - 1
		}

		view := m.View()
		selectedJob := m.jobs[m.selectedIndex]

		if !strings.Contains(view, selectedJob.Company) {
			t.Errorf("View does not contain company of job at selectedIndex %d: %q", m.selectedIndex, selectedJob.Company)
			return false
		}

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Errorf("Job list rendering property failed: %v", err)
	}
}

func TestProperty_EmptyStateRenders(t *testing.T) {
	f := func(seed int64) bool {
		m := genModelFromSeed(seed)
		m.jobs = []job.Job{}
		m.width = 80
		m.height = 24

		view := m.View()

		if !strings.Contains(view, "No jobs found") {
			t.Errorf("Empty state does not contain 'No jobs found'")
			return false
		}

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Errorf("Empty state rendering property failed: %v", err)
	}
}

func TestProperty_TopBarShowsProfileAndJobCount(t *testing.T) {
	f := func(seed int64) bool {
		m := genModelFromSeed(seed)
		m.width = 80
		m.height = 24

		view := m.View()

		if !strings.Contains(view, m.profileName) {
			t.Errorf("View does not contain profileName: %q", m.profileName)
			return false
		}

		if !strings.Contains(view, "Jobs:") {
			t.Errorf("View does not contain 'Jobs:' label")
			return false
		}

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Errorf("TopBar property failed: %v", err)
	}
}
