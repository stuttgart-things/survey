package main

import (
	"fmt"
	"log"

	"github.com/stuttgart-things/survey"
)

var allAnswers = make(map[string]interface{})
var runSurvey = false // Set to false to generate random answers

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

	// r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// In the main.go file, modify the random answers generation section:
	if !runSurvey {
		// Generate random answers

		allAnswers = survey.GetRandomAnswers(questions)

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
		allAnswers[q.Name] = survey.ConvertToType(q.Default, q.Type)
	}

	fmt.Println("Survey answers:")
	for k, v := range allAnswers {
		fmt.Printf("%-20s: %v\n", k, v)
	}
}
