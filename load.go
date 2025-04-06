package survey

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadQuestionFile(filename, yamlKey string) ([]*Question, error) {
	var questions []*Question

	// READ THE YAML FILE
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// ATTEMPT TO UNMARSHAL AS A LIST DIRECTLY (FOR YAML WITHOUT `yamlKey` KEY)
	if err := yaml.Unmarshal(data, &questions); err == nil {
		return questions, nil
	}

	// IF UNMARSHALING DIRECTLY FAILS, UNMARSHAL INTO A MAP AND EXTRACT BY `yamlKey`
	var genericMap map[string]interface{}
	if err := yaml.Unmarshal(data, &genericMap); err != nil {
		return nil, err
	}

	// EXTRACT THE DATA ASSOCIATED WITH `yamlKey`
	if rawQuestions, found := genericMap[yamlKey]; found {
		rawData, err := yaml.Marshal(rawQuestions) // Marshal back into YAML for unmarshaling into []*Question
		if err != nil {
			return nil, err
		}

		if err := yaml.Unmarshal(rawData, &questions); err != nil {
			return nil, err
		}
		return questions, nil
	}

	// RETURN AN ERROR IF `yamlKey` IS NOT FOUND
	return nil, fmt.Errorf("key '%s' not found in YAML file", yamlKey)
}
