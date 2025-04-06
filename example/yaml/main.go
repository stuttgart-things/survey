package main

import (
	"fmt"
	"log"

	"github.com/stuttgart-things/survey"
)

var allAnswers = make(map[string]interface{})

func main() {

	// LOAD THE QUESTIONS FROM A YAML FILE
	questions, err := survey.LoadQuestionFile("questions.yaml", "survey_questions")
	if err != nil {
		fmt.Println("Error loading questions:", err)
		return
	}
	fmt.Println("Loaded questions:", questions)

	// BUILD THE SURVEY FORM AND GET A MAP FOR ANSWERS
	surveyForm, _, err := survey.BuildSurvey(questions)
	if err != nil {
		log.Fatalf("Error building survey: %v", err)
	}

	// RUN THE INTERACTIVE SURVEY
	err = surveyForm.Run()
	if err != nil {
		log.Fatalf("Error running survey: %v", err)
	}

	// SET ANWERS TO ALL VALUES
	for _, question := range questions {
		allAnswers[question.Name] = question.Default
	}

	// OUTPUT ALL ANSWERS
	for key, value := range allAnswers {
		fmt.Printf("%s: %v\n", key, value)
	}

}
