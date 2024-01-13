package foundrydata

import (
	"embed"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
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

func RenderItemList(data DBData, saveDirPath string) {

	sort.Slice(data.Items, func(i, j int) bool {
		return data.Items[i].Name < data.Items[j].Name

	})
	err, tmpl := getTemplate("items.html")
	if err != nil {
		panic(err)
	}
	listSavePath := filepath.Join(saveDirPath, "items.html")
	file, err := os.OpenFile(listSavePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %s", err)
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		fmt.Printf("Error writing journal %s", err)
	}

	for _, item := range data.Items {
		renderItem(item, saveDirPath)
	}
}

func renderItem(item Item, saveDirPath string) {
	dirPath := filepath.Join(saveDirPath, "items")
	renderComponent(item, "item.html", dirPath, item.ID)
}

func RenderActors(data DBData, saveDirPath string) {

	sort.Slice(data.Actors, func(i, j int) bool {
		return data.Actors[i].Name < data.Actors[j].Name
	})
	err, tmpl := getTemplate("actor_list.html")
	if err != nil {
		panic(err)
	}

	savePath := filepath.Join(saveDirPath, "actors.html")
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
	for _, actor := range data.Actors {
		renderActor(actor, saveDirPath)
	}
}

func renderActor(actor Actor, saveDirPath string) {
	const templateFile = "actor_page.html"
	dirPath := filepath.Join(saveDirPath, "actors")
	id := actor.ID
	renderComponent(actor, templateFile, dirPath, id)
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

func foundryTagProcessor(text string) template.HTML {
	titler := cases.Title(language.English)
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

	re = regexp.MustCompile(`@Check\[([^\]]+)\]({[^}]+})?`)
	output = re.ReplaceAllStringFunc(output, func(match string) string {
		stats := re.FindStringSubmatch(match)[1]
		description := re.FindStringSubmatch(match)[2]
		parts := strings.Split(stats, "|")
		checkType := ""
		dc := ""
		name := ""
		for _, part := range parts {
			details := strings.Split(part, ":")
			if len(details) > 1 {
				if details[0] == "type" {
					checkType = html.EscapeString(titler.String(details[1]))
				}
				if details[0] == "dc" {
					dc = html.EscapeString(details[1])
				}
				if details[0] == "name" {
					name = html.EscapeString(details[1])
				}
			}
		}
		if checkType != "" && dc != "" && name != "" {
			return fmt.Sprintf(`<span class="check">%s (%s DC %s)</span>`,
				name,
				checkType,
				dc)
		}
		if len(description) > 2 {
			// not empty, and longer than `{}`
			description = description[1 : len(description)-1]
			if description == checkType {
				description = ""
			}
		} else {
			description = ""
		}
		if checkType != "" && dc != "" && description != "" {

			return fmt.Sprintf(`<span class="check">%s (%s DC %s)</span>`,
				description, checkType, dc)
		}
		if checkType != "" && dc != "" {
			return fmt.Sprintf(`<span class="check">%s (DC %s)</span>`,
				checkType, dc)
		}
		return stats
	})

	// generic unhandled stuff
	re = regexp.MustCompile(`@Actor\[([a-zA-Z0-9]+)\]{([^}]+)}`)
	output = re.ReplaceAllStringFunc(output, func(output string) string {
		id := html.EscapeString(re.FindStringSubmatch(output)[1])
		desc := html.EscapeString(re.FindStringSubmatch(output)[2])
		return fmt.Sprintf(`<a href="/actors/%s.html">%s</a>`,
			id, desc)
	})

	//re = regexp.MustCompile(`@[^\]]+\[[^\]]+\]{([^}]+)}`)
	//output = re.ReplaceAllStringFunc(output, func(output string) string {
	//	desc := re.FindStringSubmatch(output)[1]
	//	return desc
	//
	//})

	return template.HTML(output)
}

func getTemplate(templateFilename string) (error, *template.Template) {
	p := fmt.Sprintf("templates/%s", templateFilename)
	templateContent, err := templates.ReadFile(p)
	if err != nil {
		panic(err)
	}
	templateString := string(templateContent)
	tmpl, err := template.New("template").Funcs(template.FuncMap{
		"safe": foundryTagProcessor,
	}).Parse(templateString)
	if err != nil {
		panic(err)
	}
	return err, tmpl
}
