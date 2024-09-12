package global

import (
	"strings"

	"github.com/manifoldco/promptui"
)

func promptSearcher(input string, index int, jumpRootPaths []map[string]string) bool {
	jumpRoot := jumpRootPaths[index]
	name := strings.Replace(strings.ToLower(jumpRoot["jumpRoot"]), " ", "", -1)
	input = strings.Replace(strings.ToLower(input), " ", "", -1)

	return strings.Contains(name, input)
}

func promptTemplates() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .jumpRoot | cyan }} ({{ .fullPath | faint }})",
		Inactive: "  {{ .jumpRoot | cyan }}",
		Selected: "\U0001F449 {{ .jumpRoot | red }}",
	}
}

func PromptSelector(jumpRootPaths []map[string]string) *promptui.Select {
	return &promptui.Select{
		Label:     "Select a jump root",
		Items:     jumpRootPaths,
		Templates: promptTemplates(),
		Searcher: func(input string, index int) bool {
			return promptSearcher(input, index, jumpRootPaths)
		},
	}
}
