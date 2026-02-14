package llm

import (
	"strings"
)

// RenderPrompt replaces {{key}} placeholders in a template with values from vars
func RenderPrompt(template string, vars map[string]string) string {
	replacements := make([]string, 0, len(vars)*2)
	for key, value := range vars {
		placeholder := "{{" + key + "}}"
		replacements = append(replacements, placeholder, value)
	}

	r := strings.NewReplacer(replacements...)
	return r.Replace(template)
}
