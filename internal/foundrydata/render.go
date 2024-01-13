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

func RenderIndex(saveDirPath string) {
	renderComponent(nil, "index.html", saveDirPath, "index")

}

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
		renderJournal(journal, saveDirPath)
	}
}

func renderJournal(journal Journal, saveDirPath string) {
	templateFilename := "journal.html"
	dirPath := filepath.Join(saveDirPath, "journals")
	id := journal.ID
	renderComponent(journal, templateFilename, dirPath, id)

	for _, page := range journal.Pages {
		renderJournalPage(page, saveDirPath)
	}
}

func renderComponent(renderData interface{}, templateFilename string, dirPath string, id string) {
	err, tmpl := getTemplate(templateFilename)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %s", err)
		panic(err)
	}

	savePath := filepath.Join(dirPath, fmt.Sprintf("%s.html", id))
	file, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %s", err)
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, renderData)
	if err != nil {
		fmt.Printf("Error rendering %s (%s): %s", templateFilename, id, err)
	}
}

func renderJournalPage(page JournalPage, saveDirPath string) {
	dirPath := filepath.Join(saveDirPath, "journal_pages")
	renderComponent(page, "journal_page.html", dirPath, page.ID)
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

			// Example: @Compendium[pf2e.actionspf2e.TiNDYUGlMmxzxBYU]{Search}
			re = regexp.MustCompile(`@Compendium\[pf2e.actionspf2e\.[A-Za-z0-9]+\]{([^}]+)}`)
			output = re.ReplaceAllStringFunc(output, func(match string) string {
				action := html.EscapeString(re.FindStringSubmatch(match)[1])
				return fmt.Sprintf(`<span class="action">%s</span>`, action)
			})

			// Example: @Check[type:perception|dc:18]{Perception}
			re = regexp.MustCompile(`@Check\[type:([^\|]+)\|dc:(\d+)\]{([^}]+)}`)
			output = re.ReplaceAllStringFunc(output, func(match string) string {
				checkType := html.EscapeString(re.FindStringSubmatch(match)[1])
				checkDc := html.EscapeString(re.FindStringSubmatch(match)[2])
				checkDesc := html.EscapeString(re.FindStringSubmatch(match)[3])
				return fmt.Sprintf(`<span class="check" data-type="%s">%s (DC %s)</span>`,
					checkType, checkDesc, checkDc)
				//return fmt.Sprintf(`<span class="action">%s</span>`, action)
			})

			return template.HTML(output)
		},
	}).Parse(templateString)
	if err != nil {
		panic(err)
	}
	return err, tmpl
}
