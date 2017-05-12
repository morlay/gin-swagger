package http_error_code

import (
	"strings"
	"strconv"
)

func ParseHttpCodeDesc(str string) (msg string, desc string, canBeErrTalk bool) {
	lines := strings.Split(str, "\n")
	firstLine := strings.Split(lines[0], "@errTalk")

	if (len(firstLine) > 1) {
		canBeErrTalk = true
		msg = strings.TrimSpace(firstLine[1])
	} else {
		canBeErrTalk = false
		msg = strings.TrimSpace(firstLine[0])
	}

	if (len(lines) > 1) {
		desc = strings.TrimSpace(strings.Join(lines[1:], "\n"));
	}

	return
}

func CodeToStatus(s string) int {
	status, _ := strconv.Atoi(s[:3])
	return status
}
