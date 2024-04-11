package common

import (
	"strings"
)

func CheckSqlInjection(text string) string {
	list := []string{"INSERT ", "UPDATE ", "DELETE ", "CREATE ", "EMPTY ", "DROP ", "ALTER ", "TRUNCATE "}
	if StringExistsInList(strings.ToUpper(text), list) {
		return "ERROR"
	}
	return text
}
