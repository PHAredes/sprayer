package cv

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"job-scraper/pkg/models"
	"job-scraper/pkg/tracking"
)

// CVTemplate represents a CV template
type CVTemplate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Template    string `json:"template"`
}

// GenerateCV generates a CV with tracking integration
func GenerateCV(cvData models.CVData, job models.Job, templateName string) (string, error) {
	// Get the appropriate template
	templateContent := getTemplate(templateName)
	
	// Parse the template
	tmpl, err := template.New("cv").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	// Prepare data for template
	data := struct {
		CVData models.CVData
		Job    models.Job
		Tracking struct {
			PixelURL string
			LinkURL  string
		}
	}{
		CVData: cvData,
		Job:    job,
	}

	// Generate tracking URLs
	if cvData.Email != "" {
		obfuscatedEmail := tracking.ObfuscatedEmail(cvData.Email)
		data.Tracking.PixelURL = tracking.GenerateTrackingPixelURL(job.ID, obfuscatedEmail)
		data.Tracking.LinkURL = tracking.GenerateTrackingLink(job.ID, obfuscatedEmail, job.URL)
	}

	// Execute template
	var result bytes.Buffer
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return result.String(), nil
}

// getTemplate returns the template content based on template name
func getTemplate(templateName string) string {
	switch templateName {
	case "modern":
		return modernTemplate()
	case "professional":
		return professionalTemplate()
	case "minimal":
		return minimalTemplate()
	default:
		return modernTemplate()
	}
}

// modernTemplate returns a modern CV template
func modernTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.CVData.Name}} - CV</title>
    <style>
        body { font-family: 'Segoe UI', Arial, sans-serif; line-height: 1.6; margin: 0; padding: 20px; background: #f8f9fa; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 40px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; border-bottom: 2px solid #4ECDC4; padding-bottom: 20px; }
        .name { font-size: 32px; font-weight: bold; color: #2c3e50; margin: 0; }
        .contact { color: #7f8c8d; margin: 10px 0; }
        .section { margin-bottom: 25px; }
        .section-title { font-size: 20px; font-weight: bold; color: #4ECDC4; border-bottom: 1px solid #ecf0f1; padding-bottom: 5px; margin-bottom: 15px; }
        .job-title { font-weight: bold; color: #2c3e50; }
        .job-company { color: #7f8c8d; font-style: italic; }
        .skills { display: flex; flex-wrap: wrap; gap: 10px; }
        .skill { background: #4ECDC4; color: white; padding: 5px 15px; border-radius: 20px; font-size: 14px; }
        .footer { text-align: center; margin-top: 30px; padding-top: 20px; border-top: 1px solid #ecf0f1; color: #7f8c8d; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1 class="name">{{.CVData.Name}}</h1>
            <div class="contact">
                {{.CVData.Email}} • {{.CVData.Phone}} • {{.CVData.Location}}
            </div>
        </div>

        <div class="section">
            <h2 class="section-title">Professional Summary</h2>
            <p>{{.CVData.Summary}}</p>
        </div>

        {{if .CVData.Experience}}
        <div class="section">
            <h2 class="section-title">Experience</h2>
            {{range .CVData.Experience}}
            <div style="margin-bottom: 15px;">
                <div class="job-title">{{.}}</div>
            </div>
            {{end}}
        </div>
        {{end}}

        {{if .CVData.Skills}}
        <div class="section">
            <h2 class="section-title">Skills</h2>
            <div class="skills">
                {{range .CVData.Skills}}
                <span class="skill">{{.}}</span>
                {{end}}
            </div>
        </div>
        {{end}}

        {{if .CVData.Education}}
        <div class="section">
            <h2 class="section-title">Education</h2>
            {{range .CVData.Education}}
            <div style="margin-bottom: 10px;">{{.}}</div>
            {{end}}
        </div>
        {{end}}

        {{if .CVData.Projects}}
        <div class="section">
            <h2 class="section-title">Projects</h2>
            {{range .CVData.Projects}}
            <div style="margin-bottom: 10px;">{{.}}</div>
            {{end}}
        </div>
        {{end}}

        {{if .CVData.Certifications}}
        <div class="section">
            <h2 class="section-title">Certifications</h2>
            {{range .CVData.Certifications}}
            <div style="margin-bottom: 10px;">{{.}}</div>
            {{end}}
        </div>
        {{end}}

        <div class="footer">
            <p>Generated by Sprayer CV Generator • {{now.Format "January 2, 2006"}}</p>
            {{if .Tracking.PixelURL}}
            <!-- Tracking pixel for email opens -->
            <img src="{{.Tracking.PixelURL}}" width="1" height="1" style="display:none;" alt="">
            {{end}}
        </div>
    </div>
</body>
</html>`
}

// professionalTemplate returns a professional CV template
func professionalTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.CVData.Name}} - Curriculum Vitae</title>
    <style>
        body { font-family: 'Times New Roman', serif; line-height: 1.4; margin: 0; padding: 40px; }
        .container { max-width: 700px; margin: 0 auto; }
        .header { text-align: center; margin-bottom: 30px; border-bottom: 2px solid #333; padding-bottom: 20px; }
        .name { font-size: 24px; font-weight: bold; text-transform: uppercase; letter-spacing: 2px; margin: 0; }
        .contact { font-size: 14px; margin: 10px 0; }
        .section { margin-bottom: 20px; }
        .section-title { font-size: 18px; font-weight: bold; text-transform: uppercase; border-bottom: 1px solid #333; padding-bottom: 5px; margin-bottom: 10px; }
        .item { margin-bottom: 10px; }
        .item-title { font-weight: bold; }
        .footer { text-align: center; margin-top: 30px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1 class="name">{{.CVData.Name}}</h1>
            <div class="contact">
                {{.CVData.Email}} | {{.CVData.Phone}} | {{.CVData.Location}}
            </div>
        </div>

        <div class="section">
            <h2 class="section-title">Professional Summary</h2>
            <p>{{.CVData.Summary}}</p>
        </div>

        {{if .CVData.Experience}}
        <div class="section">
            <h2 class="section-title">Professional Experience</h2>
            {{range .CVData.Experience}}
            <div class="item">
                <div class="item-title">{{.}}</div>
            </div>
            {{end}}
        </div>
        {{end}}

        {{if .CVData.Skills}}
        <div class="section">
            <h2 class="section-title">Technical Skills</h2>
            <p>{{join .CVData.Skills ", "}}</p>
        </div>
        {{end}}

        {{if .CVData.Education}}
        <div class="section">
            <h2 class="section-title">Education</h2>
            {{range .CVData.Education}}
            <div class="item">{{.}}</div>
            {{end}}
        </div>
        {{end}}

        <div class="footer">
            <p>Generated on {{now.Format "January 2, 2006"}}</p>
            {{if .Tracking.PixelURL}}
            <!-- Email tracking pixel -->
            <img src="{{.Tracking.PixelURL}}" width="1" height="1" style="display:none;" alt="">
            {{end}}
        </div>
    </div>
</body>
</html>`
}

// minimalTemplate returns a minimal CV template
func minimalTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.CVData.Name}}</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, sans-serif; line-height: 1.6; margin: 40px; max-width: 600px; }
        .name { font-size: 24px; font-weight: 600; margin-bottom: 10px; }
        .contact { color: #666; margin-bottom: 20px; }
        .section { margin-bottom: 20px; }
        .section-title { font-weight: 600; margin-bottom: 10px; }
        .footer { margin-top: 30px; font-size: 12px; color: #999; }
    </style>
</head>
<body>
    <div class="name">{{.CVData.Name}}</div>
    <div class="contact">
        {{.CVData.Email}} • {{.CVData.Phone}} • {{.CVData.Location}}
    </div>

    <div class="section">
        <div class="section-title">Summary</div>
        <div>{{.CVData.Summary}}</div>
    </div>

    {{if .CVData.Experience}}
    <div class="section">
        <div class="section-title">Experience</div>
        {{range .CVData.Experience}}
        <div style="margin-bottom: 8px;">{{.}}</div>
        {{end}}
    </div>
    {{end}}

    {{if .CVData.Skills}}
    <div class="section">
        <div class="section-title">Skills</div>
        <div>{{join .CVData.Skills ", "}}</div>
    </div>
    {{end}}

    <div class="footer">
        Generated {{now.Format "2006-01-02"}}
        {{if .Tracking.PixelURL}}
        <img src="{{.Tracking.PixelURL}}" width="1" height="1" style="display:none;">
        {{end}}
    </div>
</body>
</html>`
}

// join is a helper function to join strings with a separator
func join(items []string, separator string) string {
	return strings.Join(items, separator)
}