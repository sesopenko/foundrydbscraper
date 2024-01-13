package foundrydata

import (
	"embed"
	"fmt"
	"html"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
)

//go:embed templates
var templates embed.FS

func RenderJournalList(data DBData, saveDirPath string) {
	templateFilename := "journal_list.html"
	err, tmpl := getTemplate(templateFilename)
	if err != nil {
		panic(err)
	}

	savePath := filepath.Join(saveDirPath, "journals.html")
	file, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %s", err)
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		fmt.Printf("Error writing journal %s", err)
	}

	for _, journal := range data.Journals {
		RenderJournal(journal, saveDirPath)
	}
}

func RenderJournal(journal Journal, saveDirPath string) {
	templateFilename := "journal.html"
	err, tmpl := getTemplate(templateFilename)
	if err != nil {
		panic(err)
	}

	dirPath := filepath.Join(saveDirPath, "journals")
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %s", err)
		panic(err)
	}

	savePath := filepath.Join(dirPath, fmt.Sprintf("%s.html", journal.ID))
	file, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %s", err)
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, journal)
	if err != nil {
		fmt.Printf("Error writing journal %s: %s", journal.ID, err)
	}

	for _, page := range journal.Pages {
		renderJournalPage(page, saveDirPath)
	}
}

func renderJournalPage(page JournalPage, saveDirPath string) {
	err, tmpl := getTemplate("journal_page.html")
	if err != nil {
		panic(err)
	}
	dirPath := filepath.Join(saveDirPath, "journal_pages")
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %s", err)
		panic(err)
	}
	savePath := filepath.Join(dirPath, fmt.Sprintf("%s.html", page.ID))
	file, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %s", err)
		panic(err)
	}
	defer file.Close()
	err = tmpl.Execute(file, page)
	if err != nil {
		fmt.Printf("Error writing journal %s: %s", page.ID, err)
	}
}

func getTemplate(templateFilename string) (error, *template.Template) {
	p := fmt.Sprintf("templates/%s", templateFilename)
	templateContent, err := templates.ReadFile(p)
	if err != nil {
		panic(err)
	}
	templateString := string(templateContent)
	tmpl, err := template.New("template").Funcs(template.FuncMap{
		"safe": func(text string) template.HTML {
			// Example:
			// @UUID[JournalEntry.3T1M395V6J75OsEp.JournalEntryPage.8Jl0TWoH4iUJzJxK]{Beginning the Adventure}
			re := regexp.MustCompile(`@UUID\[JournalEntry\.([a-zA-Z0-9]+)\.JournalEntryPage\.([a-zA-Z0-9]+)\]{([^}]+)}`)
			output := re.ReplaceAllStringFunc(text, func(match string) string {
				//journalId := html.EscapeString(re.FindStringSubmatch((match))[1])
				pageId := html.EscapeString(re.FindStringSubmatch(match)[2])
				linkTitle := html.EscapeString(re.FindStringSubmatch(match)[3])
				return fmt.Sprintf(`<a href="/journal_pages/%s.html">%s</a>`, pageId, linkTitle)

			})

			re = regexp.MustCompile(`@UUID\[JournalEntry\.([a-zA-Z0-9]+)\]{([^}]+)}`)
			output = re.ReplaceAllStringFunc(output, func(match string) string {
				journalId := html.EscapeString(re.FindStringSubmatch((match))[1])
				linkTitle := html.EscapeString(re.FindStringSubmatch(match)[2])
				return fmt.Sprintf(`<a href="/journals/%s.html">%s</a>`, journalId, linkTitle)
			})

			re = regexp.MustCompile(`@UUID\[\.([a-zA-Z0-9]+)\]{([^}]+)}`)
			output = re.ReplaceAllStringFunc(output, func(match string) string {
				journalPageId := html.EscapeString(re.FindStringSubmatch((match))[1])
				linkTitle := html.EscapeString(re.FindStringSubmatch(match)[2])
				return fmt.Sprintf(`<a href="/journal_pages/%s.html">%s</a>`, journalPageId, linkTitle)
			})
			return template.HTML(output)
		},
	}).Parse(templateString)
	if err != nil {
		panic(err)
	}
	return err, tmpl
}
