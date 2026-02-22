package explanations

import "fmt"

// prepositionCaseMap defines which case each preposition requires
var prepositionCaseMap = map[string]string{
	"σε":     "accusative", // to, at, in
	"από":    "genitive",   // from
	"για":    "genitive",   // for
	"με":     "accusative", // with
	"χωρίς":  "accusative", // without
	"μετά":   "accusative", // after
	"πριν":   "accusative", // before
}

// SyntacticRoleTemplate returns the rule explanation for a given context
func SyntacticRoleTemplate(contextType string, caseType string, prep *string) string {
	switch contextType {
	case "direct_object":
		return "Direct objects use accusative case"

	case "possession":
		return "Possession requires genitive case"

	case "preposition":
		if prep != nil && *prep != "" {
			return fmt.Sprintf("The preposition '%s' requires %s case", *prep, caseType)
		}
		return fmt.Sprintf("This preposition requires %s case", caseType)

	default:
		return fmt.Sprintf("This context uses %s case", caseType)
	}
}
