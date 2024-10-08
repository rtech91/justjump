package local

import (
	"strings"

	"github.com/manifoldco/promptui"
)

func promptSearcher(input string, index int, jumpPointPaths []map[string]string) bool {
	jumpPoint := jumpPointPaths[index]
	name := strings.Replace(strings.ToLower(jumpPoint["jumpPoint"]), " ", "", -1)
	input = strings.Replace(strings.ToLower(input), " ", "", -1)

	return strings.Contains(name, input)
}

func promptTemplates() *promptui.SelectTemplates {
	return &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .jumpPoint | cyan }} ({{ .fullPath | faint }})",
		Inactive: "  {{ .jumpPoint | cyan }}",
		Selected: "\U0001F449 {{ .jumpPoint | red }}",
	}
}

func PromptSelector(jumpPointPaths []map[string]string) *promptui.Select {
	return &promptui.Select{
		Label:     "Select a jump point",
		Items:     jumpPointPaths,
		Templates: promptTemplates(),
		Searcher: func(input string, index int) bool {
			return promptSearcher(input, index, jumpPointPaths)
		},
	}
}
