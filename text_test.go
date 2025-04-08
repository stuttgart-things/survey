package survey

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadYAMLFile(t *testing.T) {
	// Setup - create a temporary YAML file
	content := `key: value
list:
  - item1
  - item2`
	tmpfile, err := os.CreateTemp("", "test*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	tests := []struct {
		name     string
		filePath string
		want     string
		wantErr  bool
	}{
		{
			name:     "valid yaml file",
			filePath: tmpfile.Name(),
			want:     content,
			wantErr:  false,
		},
		{
			name:     "non-existent file",
			filePath: "nonexistent.yaml",
			want:     "",
			wantErr:  true,
		},
		{
			name:     "empty path",
			filePath: "",
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadYAMLFile(tt.filePath)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestValidateYAML(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name: "valid yaml",
			content: `key: value
list:
  - item1
  - item2`,
			wantErr: false,
		},
		{
			name:    "empty content",
			content: "",
			wantErr: false, // Empty is technically valid YAML
		},
		{
			name: "invalid yaml - bad indentation",
			content: `key: value
 list:
  - item1`,
			wantErr: true,
		},
		{
			name:    "invalid yaml - malformed",
			content: `key: value: anothervalue`,
			wantErr: true,
		},
		{
			name: "complex valid yaml",
			content: `---
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  containers:
  - name: nginx
    image: nginx:1.14.2
    ports:
    - containerPort: 80`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateYAML(tt.content)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReadAndValidateIntegration(t *testing.T) {
	// Setup - create valid temporary YAML file
	validContent := `key: value
valid: true`
	tmpfile, err := os.CreateTemp("", "valid*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(validContent); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	t.Run("read and validate valid yaml", func(t *testing.T) {
		content, err := ReadYAMLFile(tmpfile.Name())
		assert.NoError(t, err)
		assert.Equal(t, validContent, content)

		err = validateYAML(content)
		assert.NoError(t, err)
	})

	// Setup - create invalid temporary YAML file
	invalidContent := `key: value: wrong`
	tmpfile2, err := os.CreateTemp("", "invalid*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile2.Name())

	if _, err := tmpfile2.WriteString(invalidContent); err != nil {
		t.Fatal(err)
	}
	tmpfile2.Close()

	t.Run("read and validate invalid yaml", func(t *testing.T) {
		content, err := ReadYAMLFile(tmpfile2.Name())
		assert.NoError(t, err)
		assert.Equal(t, invalidContent, content)

		err = validateYAML(content)
		assert.Error(t, err)
	})
}
