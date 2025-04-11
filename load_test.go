package survey

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Sample YAML content for testing
var sampleYAML = `
survey_questions:
  - prompt: "What is your favorite color?"
    name: "favorite_color"
    options: ["Red", "Blue", "Green"]
    default: "Blue"
    kind: "select"
  - prompt: "How old are you?"
    name: "age"
    type: "int"
    default: "25"
    kind: "ask"
    minLength: 2
    maxLength: 30
`

func TestLoadQuestionFile(t *testing.T) {
	filename := createTempYAMLFile(t, sampleYAML)
	defer func() {
		err := os.Remove(filename)
		assert.NoError(t, err)
	}()
	questions, err := LoadQuestionFile(filename, "survey_questions")
	assert.NoError(t, err)
	assert.NotEmpty(t, questions)
	assert.Equal(t, "favorite_color", questions[0].Name)
	assert.Equal(t, "Blue", questions[0].Default)

	// Test the age question with length constraints
	assert.Equal(t, "age", questions[1].Name)
	assert.Equal(t, 2, questions[1].MinLength)
	assert.Equal(t, 30, questions[1].MaxLength)
}

func createTempYAMLFile(t *testing.T, content string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "test_survey_*.yaml")
	assert.NoError(t, err)

	_, err = tmpFile.Write([]byte(content))
	assert.NoError(t, err)

	err = tmpFile.Close()
	assert.NoError(t, err)

	return tmpFile.Name()
}
