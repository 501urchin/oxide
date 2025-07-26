package input

import (
	"fmt"
	"strings"
)

func GetConfirmationFromUser(format string, a ...any) bool {
	var answer string
	fmt.Printf(format, a...)
	_, err := fmt.Scanln(&answer)
	if err != nil {
		return false
	}

	answer = strings.ToLower(answer)

	if answer == "yes" || answer == "y" {
		return true
	}

	return false
}
