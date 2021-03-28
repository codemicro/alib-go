package alib

import (
	"os"
	"strings"
)

func OsPathJoin(secs ...string) string {
	return strings.Join(secs, string(os.PathSeparator))
}
