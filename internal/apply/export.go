package apply

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"sprayer/internal/job"
)

// ExportJSON writes jobs to a JSON file.
func ExportJSON(jobs []job.Job, path string) error {
	data, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// ExportCSV writes jobs to a CSV file.
func ExportCSV(jobs []job.Job, path string) error {
	var b strings.Builder
	b.WriteString("Title,Company,Location,Source,Score,Posted Date,URL\n")
	for _, j := range jobs {
		b.WriteString(fmt.Sprintf("%q,%q,%q,%q,%d,%q,%q\n",
			j.Title, j.Company, j.Location, j.Source, j.Score,
			j.PostedDate.Format("2006-01-02"), j.URL))
	}
	return os.WriteFile(path, []byte(b.String()), 0644)
}
