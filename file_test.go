package mime

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var detectFileTypeTestcases = []struct {
	filename string
	mimetype string
}{
	{"test.json", "application/json"},
	{"test.xml", "application/xml"},
	{"test.html", "text/html"},
	{"test.css", "text/css"},
	{"test.js", "text/javascript"},
	{"test.doesnotexist", ""},
}

func TestDetectFileType(t *testing.T) {
	for _, testcase := range detectFileTypeTestcases {
		t.Run("", func(t *testing.T) {
			mimetype := DetectFileType(testcase.filename)
			require.Equal(t, testcase.mimetype, mimetype)
		})
	}
}
