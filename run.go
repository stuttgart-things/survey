package survey

import (
	log "github.com/sirupsen/logrus"
)

func RunSurvey(profilePath, surveyKey string) (surveyValues map[string]interface{}) {
	surveyValues = make(map[string]interface{})

	// READ PROFILE AND SURVEY BY KEY
	survey, _ := LoadQuestionFile(profilePath, surveyKey)

	// IF SURVEY EXISTS, RUN IT
	if len(survey) > 0 {
		surveyQuestions, _, err := BuildSurvey(survey)

		if err != nil {
			log.Fatalf("ERROR BUILDING SURVEY: %v", err)
		}
		log.Info("SURVEY FOUND")

		// RUN THE INTERACTIVE SURVEY
		err = surveyQuestions.Run()
		if err != nil {
			log.Fatalf("ERROR RUNNING SURVEY: %v", err)
		}

		// SET ANWERS TO ALL VALUES
		for _, question := range survey {
			surveyValues[question.Name] = question.Default
		}

	} else {
		log.Info("NO SURVEY FOUND")
	}

	return surveyValues
}
