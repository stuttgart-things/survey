package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/stuttgart-things/survey"
)

var allAnswers = make(map[string]interface{})
var runSurvey = true // Set to false to generate random answers

func main() {

	// Register functions before loading questions
	survey.RegisterFunction("getDefaultFavoriteFood", func(params map[string]interface{}) string {
		if spiceLevel, ok := params["spiceLevel"].(string); ok && spiceLevel != "" {
			return fmt.Sprintf("spicy %s", spiceLevel)
		}
		return "steak"
	})

	survey.RegisterFunction("getDefaultDrink", func(params map[string]interface{}) string {
		if temp, ok := params["temperature"].(string); ok && temp != "" {
			return fmt.Sprintf("%s water", temp)
		}
		return "water"
	})

	// LOAD THE QUESTIONS FROM YAML
	questions, err := survey.LoadQuestionFile("questions.yaml", "survey_questions")
	if err != nil {
		log.Fatalf("Error loading questions: %v", err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// In the main.go file, modify the random answers generation section:
	if !runSurvey {
		// Generate random answers
		for _, q := range questions {
			switch q.Kind {
			case "select":
				if len(q.Options) > 0 {
					q.Default = q.Options[rand.Intn(len(q.Options))]
				}
			case "ask":
				if q.Default == "" {
					q.Default = generateRandomValue(q, r)
				}
			case "function":
				if q.DefaultFunction != "" {
					if fn, ok := survey.DefaultFunctions[q.DefaultFunction]; ok {
						q.Default = fn(q.DefaultParams)
					} else {
						log.Printf("Function %s not found, using empty string", q.DefaultFunction)
						q.Default = ""
					}
				}
				// If no default was generated, use a random string
				if q.Default == "" {
					q.Default = generateRandomValue(q, r)
				}
			}

			// Convert to proper type
			allAnswers[q.Name] = convertToType(q.Default, q.Type)
		}

		// Print results without type information
		fmt.Println("Generated answers:")
		for k, v := range allAnswers {
			fmt.Printf("%-20s: %v\n", k, v)
		}
		return
	}

	// Original survey flow
	surveyForm, _, err := survey.BuildSurvey(questions)
	if err != nil {
		log.Fatalf("Error building survey: %v", err)
	}

	if err := surveyForm.Run(); err != nil {
		log.Fatalf("Error running survey: %v", err)
	}

	for _, q := range questions {
		allAnswers[q.Name] = convertToType(q.Default, q.Type)
	}

	fmt.Println("Survey answers:")
	for k, v := range allAnswers {
		fmt.Printf("%-20s: %v\n", k, v)
	}
}

func generateRandomValue(q *survey.Question, r *rand.Rand) string {
	minLen := q.MinLength
	if minLen < 0 {
		minLen = 0
	}

	maxLen := q.MaxLength
	if maxLen <= minLen {
		maxLen = minLen + 10
	}

	switch q.Type {
	case "int":
		min := 0
		max := 9999
		if minLen > 0 {
			min = pow10(minLen - 1)
		}
		if maxLen > 0 {
			max = pow10(maxLen) - 1
			if max < min {
				max = min * 10
			}
		}
		return fmt.Sprintf("%d", r.Intn(max-min+1)+min)

	case "boolean":
		if r.Intn(2) == 0 {
			return "false"
		}
		return "true"

	default: // string
		length := r.Intn(maxLen-minLen+1) + minLen
		var sb strings.Builder
		for i := 0; i < length; i++ {
			sb.WriteByte(byte(r.Intn(26) + 'a'))
		}
		return sb.String()
	}
}

func convertToType(value string, typ string) interface{} {
	switch typ {
	case "int":
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
		return 0
	case "boolean":
		return strings.ToLower(value) == "true" || value == "Yes"
	default:
		return value
	}
}

func pow10(n int) int {
	if n <= 0 {
		return 1
	}
	return 10 * pow10(n-1)

}
