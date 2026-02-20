package apply

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"sprayer/internal/job"
	"sprayer/internal/llm"
	"sprayer/internal/profile"
)

type LatexCV struct {
	Name            string               `json:"name"`
	Title           string               `json:"title"`
	SkillsHighlight []string             `json:"skills_highlight"`
	Contact         LatexContact         `json:"contact"`
	Summary         string               `json:"summary"`
	Experience      []LatexExperience    `json:"experience"`
	Projects        []LatexProject       `json:"projects"`
	TechnicalSkills LatexTechnicalSkills `json:"technical_skills"`
	Education       []LatexEducation     `json:"education"`
}

type LatexContact struct {
	Email    string `json:"email"`
	LinkedIn string `json:"linkedin"`
	GitHub   string `json:"github"`
	Location string `json:"location"`
}

type LatexExperience struct {
	Role         string   `json:"role"`
	Company      string   `json:"company"`
	Duration     string   `json:"duration"`
	Achievements []string `json:"achievements"`
	TechStack    []string `json:"tech_stack"`
}

type LatexProject struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
}

type LatexTechnicalSkills struct {
	Languages  []string `json:"languages"`
	Frameworks []string `json:"frameworks"`
	Tools      []string `json:"tools"`
	Databases  []string `json:"databases"`
}

type LatexEducation struct {
	Degree      string `json:"degree"`
	Field       string `json:"field"`
	Institution string `json:"institution"`
	Year        string `json:"year"`
}

func (g *CVGenerator) GenerateLatexCV(j *job.Job, p *profile.Profile) (*LatexCV, error) {
	if g.client == nil {
		return nil, fmt.Errorf("LLM client not available")
	}

	cvData := p.CVData
	if cvData == nil && p.CVPath != "" {
		parser := profile.NewCVParser()
		var err error
		cvData, err = parser.ParseCVFromFile(p.CVPath)
		if err != nil {
			return nil, fmt.Errorf("parse CV: %w", err)
		}
	}

	if cvData == nil {
		return nil, fmt.Errorf("no CV data available for profile")
	}

	vars := map[string]string{
		"job_title":       j.Title,
		"company":         j.Company,
		"job_description": truncate(j.Description, 3000),
		"applicant_name":  cvData.Name,
		"applicant_email": cvData.Email,
		"applicant_phone": cvData.Phone,
		"applicant_title": cvData.Title,
		"linkedin":        p.ContactEmail,
		"github":          "",
		"location":        cvData.Location,
		"summary":         cvData.Summary,
		"technologies":    strings.Join(cvData.Technologies, ", "),
		"skills":          strings.Join(cvData.Skills, ", "),
		"experience":      formatExperience(cvData.Experience),
		"education":       formatEducation(cvData.Education),
		"projects":        "",
	}

	prompt, err := llm.LoadPrompt("cv_latex", vars)
	if err != nil {
		return nil, fmt.Errorf("load prompt: %w", err)
	}

	response, err := g.client.Complete(
		"You are a CV optimization expert. Generate a tailored CV in valid JSON format. Output only valid JSON, no markdown.",
		prompt,
	)
	if err != nil {
		return nil, fmt.Errorf("LLM generation: %w", err)
	}

	var latexCV LatexCV
	if err := json.Unmarshal([]byte(cleanJSONResponse(response)), &latexCV); err != nil {
		return nil, fmt.Errorf("parse CV JSON: %w", err)
	}

	return &latexCV, nil
}

func (g *CVGenerator) GenerateLatexDocument(j *job.Job, p *profile.Profile) (string, error) {
	cv, err := g.GenerateLatexCV(j, p)
	if err != nil {
		return "", err
	}

	return cv.ToLatex()
}

func (cv *LatexCV) ToLatex() (string, error) {
	tmpl := `\documentclass[10pt,a4paper]{article}
\usepackage[utf8]{inputenc}
\usepackage[T1]{fontenc}
\usepackage[margin=1.2cm]{geometry}
\usepackage{lmodern}
\usepackage{xcolor}
\usepackage{hyperref}
\usepackage{enumitem}
\usepackage{titlesec}
\usepackage{parskip}
\usepackage{microtype}

\pagestyle{empty}

\definecolor{linkblue}{HTML}{0066cc}
\definecolor{textblack}{HTML}{000000}
\definecolor{headergray}{HTML}{333333}
\definecolor{subtitlegray}{HTML}{555555}

\hypersetup{
    colorlinks=true,
    linkcolor=linkblue,
    urlcolor=linkblue,
    pdftitle={{{.Name}} - {{.Title}}},
    pdfauthor={{{.Name}}}
}

\titleformat{\section}
  {\normalfont\large\bfseries\uppercase}{\thesection}{1em}{}[\titlerule]
\titlespacing*{\section}{0pt}{10pt}{6pt}

\newcommand{\name}[1]{
    \begin{center}
        {\LARGE \textbf{\uppercase{#1}}}\\[2pt]
    \end{center}
}

\newcommand{\jobtitle}[1]{
    \begin{center}
        {\large \textcolor{headergray}{#1}}\\[2pt]
    \end{center}
}

\newcommand{\jobsubtitle}[1]{
    \begin{center}
        {\small \textcolor{subtitlegray}{#1}}\\[2pt]
    \end{center}
}

\newcommand{\contact}[1]{
    \begin{center}
        {\small #1}
    \end{center}
}

\newcommand{\tech}[1]{\textbf{#1}}

\newcommand{\entry}[4]{
    \noindent\textbf{#1} --- \textit{#2} \hfill #3\par
}

\begin{document}

\name{` + escapeLatex(cv.Name) + `}
\jobtitle{` + escapeLatex(cv.Title) + `}
\jobsubtitle{` + cv.formatSkillsHighlight() + `}

\contact{` + cv.formatContact() + `}

\vspace{10pt}

\section*{Professional Summary}
` + escapeLatex(cv.Summary) + `

\section*{Work Experience}
` + cv.formatExperience() + `

` + cv.formatProjects() + `

\section*{Technical Skills}
\begin{itemize}[leftmargin=*, noitemsep, topsep=0pt]
` + cv.formatTechnicalSkills() + `
\end{itemize}

\section*{Education}
\begin{itemize}[leftmargin=*, noitemsep, topsep=0pt]
` + cv.formatEducation() + `
\end{itemize}

\end{document}
`

	return tmpl, nil
}

func (cv *LatexCV) formatSkillsHighlight() string {
	if len(cv.SkillsHighlight) == 0 {
		return ""
	}
	var parts []string
	for _, s := range cv.SkillsHighlight {
		parts = append(parts, `\tech{`+escapeLatex(s)+`}`)
	}
	return strings.Join(parts, ` $\cdot$ `)
}

func (cv *LatexCV) formatContact() string {
	var parts []string
	if cv.Contact.Email != "" {
		parts = append(parts, `\href{mailto:`+cv.Contact.Email+`}{`+cv.Contact.Email+`}`)
	}
	if cv.Contact.LinkedIn != "" {
		parts = append(parts, `\href{https://`+cv.Contact.LinkedIn+`}{`+cv.Contact.LinkedIn+`}`)
	}
	if cv.Contact.GitHub != "" {
		parts = append(parts, `\href{https://`+cv.Contact.GitHub+`}{`+cv.Contact.GitHub+`}`)
	}
	if cv.Contact.Location != "" {
		parts = append(parts, escapeLatex(cv.Contact.Location))
	}
	return strings.Join(parts, ` $\mid$ `)
}

func (cv *LatexCV) formatExperience() string {
	var sb strings.Builder
	for i, exp := range cv.Experience {
		sb.WriteString(`\entry{` + escapeLatex(exp.Role) + `}{` + escapeLatex(exp.Company) + `}{` + escapeLatex(exp.Duration) + `}{}` + "\n")
		sb.WriteString(`\begin{itemize}[leftmargin=*, noitemsep, topsep=0pt]` + "\n")
		for _, ach := range exp.Achievements {
			sb.WriteString(`    \item ` + escapeLatex(ach) + `.\n`)
		}
		if len(exp.TechStack) > 0 {
			sb.WriteString(`    \item \tech{Stack:} ` + escapeLatex(strings.Join(exp.TechStack, ", ")) + `.\n`)
		}
		sb.WriteString(`\end{itemize}` + "\n")
		if i < len(cv.Experience)-1 {
			sb.WriteString("\n\\vspace{8pt}\n\n")
		}
	}
	return sb.String()
}

func (cv *LatexCV) formatProjects() string {
	if len(cv.Projects) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(`\section*{Projects}` + "\n\n")
	for _, proj := range cv.Projects {
		sb.WriteString(`\noindent\textbf{` + escapeLatex(proj.Name) + `}`)
		if proj.URL != "" {
			sb.WriteString(` \hfill \href{` + proj.URL + `}{GitHub}`)
		}
		sb.WriteString("\n")
		sb.WriteString(`\vspace{-\parskip}` + "\n")
		sb.WriteString(`\vspace{2pt}` + "\n")
		sb.WriteString(`\begin{itemize}[leftmargin=12pt, noitemsep, topsep=0pt, labelsep=4pt]` + "\n")
		sb.WriteString(`    \item ` + escapeLatex(proj.Description) + `.\n`)
		sb.WriteString(`\end{itemize}` + "\n\n")
	}
	return sb.String()
}

func (cv *LatexCV) formatTechnicalSkills() string {
	var lines []string
	if len(cv.TechnicalSkills.Languages) > 0 {
		lines = append(lines, `    \item \textbf{Languages:} `+escapeLatex(strings.Join(cv.TechnicalSkills.Languages, ", "))+`.`)
	}
	if len(cv.TechnicalSkills.Frameworks) > 0 {
		lines = append(lines, `    \item \textbf{Frameworks:} `+escapeLatex(strings.Join(cv.TechnicalSkills.Frameworks, ", "))+`.`)
	}
	if len(cv.TechnicalSkills.Tools) > 0 {
		lines = append(lines, `    \item \textbf{Tools:} `+escapeLatex(strings.Join(cv.TechnicalSkills.Tools, ", "))+`.`)
	}
	if len(cv.TechnicalSkills.Databases) > 0 {
		lines = append(lines, `    \item \textbf{Databases:} `+escapeLatex(strings.Join(cv.TechnicalSkills.Databases, ", "))+`.`)
	}
	return strings.Join(lines, "\n")
}

func (cv *LatexCV) formatEducation() string {
	var lines []string
	for _, edu := range cv.Education {
		line := `    \item `
		if edu.Degree != "" {
			line += escapeLatex(edu.Degree)
		}
		if edu.Field != "" {
			line += `, ` + escapeLatex(edu.Field)
		}
		if edu.Institution != "" {
			line += `, ` + escapeLatex(edu.Institution)
		}
		if edu.Year != "" {
			line += ` (` + escapeLatex(edu.Year) + `)`
		}
		line += `.`
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func escapeLatex(s string) string {
	replacements := []struct {
		from string
		to   string
	}{
		{`\`, `\textbackslash{}`},
		{`&`, `\&`},
		{`%`, `\%`},
		{`$`, `\$`},
		{`#`, `\#`},
		{`_`, `\_`},
		{`{`, `\{`},
		{`}`, `\}`},
		{`~`, `\textasciitilde{}`},
		{`^`, `\textasciicircum{}`},
	}
	result := s
	for _, r := range replacements {
		result = strings.ReplaceAll(result, r.from, r.to)
	}
	return result
}

func cleanJSONResponse(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}

func SaveLatexCV(content, jobID, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("create output dir: %w", err)
	}

	filename := fmt.Sprintf("cv_%s_%d.tex", sanitize(jobID), time.Now().Unix())
	filePath := filepath.Join(outputDir, filename)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("write CV file: %w", err)
	}

	return filePath, nil
}

func CompileLatexToPDF(texPath string) (string, error) {
	dir := filepath.Dir(texPath)
	base := strings.TrimSuffix(filepath.Base(texPath), ".tex")

	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", "-output-directory="+dir, base+".tex")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("pdflatex failed: %w\nOutput: %s", err, string(output))
	}

	pdfPath := filepath.Join(dir, base+".pdf")
	if _, err := os.Stat(pdfPath); err != nil {
		return "", fmt.Errorf("PDF not generated: %w", err)
	}

	return pdfPath, nil
}
