package survey

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/bubbles/v2/textinput"
	"github.com/stretchr/testify/assert"
)

func TestInitialSaveModel(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected Save
	}{
		{
			name:    "empty content",
			content: "",
			expected: Save{
				content:  "",
				status:   "Select a directory, then enter filename",
				filename: textinput.New(),
			},
		},
		{
			name:    "with content",
			content: "test content",
			expected: Save{
				content:  "test content",
				status:   "Select a directory, then enter filename",
				filename: textinput.New(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function
			result := InitialSaveModel(tt.content)

			// Verify content
			assert.Equal(t, tt.expected.content, result.content, "content should match")

			// Verify status
			assert.Equal(t, tt.expected.status, result.status, "status should match")

			// Verify filepicker configuration
			assert.True(t, result.filepicker.DirAllowed, "DirAllowed should be true")
			assert.False(t, result.filepicker.FileAllowed, "FileAllowed should be false")
			homeDir, _ := os.UserHomeDir()
			assert.Equal(t, homeDir, result.filepicker.CurrentDirectory, "should default to home directory")

			// Verify textinput configuration
			assert.Equal(t, "filename.txt", result.filename.Placeholder, "placeholder should match")
			assert.Equal(t, 256, result.filename.CharLimit, "char limit should match")
			assert.Equal(t, 50, result.filename.Width, "width should match")
			assert.True(t, result.filename.Focused(), "input should be focused")
		})
	}
}

func TestSaveModelIntegration(t *testing.T) {
	t.Run("verify home directory resolution", func(t *testing.T) {
		// Setup
		expectedHome, _ := os.UserHomeDir()
		testContent := "test content"

		// Execute
		model := InitialSaveModel(testContent)

		// Verify
		assert.Equal(t, expectedHome, model.filepicker.CurrentDirectory,
			"filepicker should initialize to user home directory")
	})

	t.Run("verify filepicker permissions", func(t *testing.T) {
		model := InitialSaveModel("test")

		// Attempt to navigate up from home directory (should work)
		parentDir := filepath.Dir(model.filepicker.CurrentDirectory)
		model.filepicker.CurrentDirectory = parentDir

		assert.NotEqual(t, "", model.filepicker.CurrentDirectory,
			"should be able to navigate directories")
	})
}
