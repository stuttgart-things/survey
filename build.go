package survey

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
)

// DefaultFunctions holds the registered functions
var DefaultFunctions = make(map[string]func(params map[string]interface{}) string)

// RegisterFunction adds a function to the registry
func RegisterFunction(name string, fn func(params map[string]interface{}) string) {
	DefaultFunctions[name] = fn
}

// BUILD THE SURVEY FUNCTION WITH THE NEW RANDOM SETUP
func BuildSurvey(questions []*Question) (*huh.Form, map[string]interface{}, error) {
	var groupFields []*huh.Group
	answers := make(map[string]interface{})
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, question := range questions {
		var field huh.Field

		// Set up default values for options if applicable
		if question.Default == "" && len(question.Options) > 0 {
			question.Default = question.Options[r.Intn(len(question.Options))]
		}

		switch question.Kind {
		case "function":
			if question.DefaultFunction != "" {
				if fn, ok := DefaultFunctions[question.DefaultFunction]; ok {
					question.Default = fn(question.DefaultParams)
				} else {
					return nil, nil, fmt.Errorf("DEFAULT FUNCTION %s NOT FOUND", question.DefaultFunction)
				}
			}

			field = huh.NewInput().
				Title(question.Prompt).
				Value(&question.Default).
				Validate(func(input string) error {
					if len(input) < question.MinLength {
						return fmt.Errorf("INPUT TOO SHORT, MINIMUM LENGTH IS %d", question.MinLength)
					}
					if len(input) > question.MaxLength {
						return fmt.Errorf("INPUT TOO LONG, MAXIMUM LENGTH IS %d", question.MaxLength)
					}
					return nil
				})

		case "ask":
			field = huh.NewInput().
				Title(question.Prompt).
				Value(&question.Default).
				Validate(func(input string) error {
					if len(input) < question.MinLength {
						return fmt.Errorf("INPUT TOO SHORT, MINIMUM LENGTH IS %d", question.MinLength)
					}
					if len(input) > question.MaxLength {
						return fmt.Errorf("INPUT TOO LONG, MAXIMUM LENGTH IS %d", question.MaxLength)
					}
					return nil
				})

			answers[question.Name] = ""

		case "list":
			var defaultValues []string
			if question.Default != "" {
				defaultValues = strings.Split(question.Default, ",")
			}

			options := make([]huh.Option[string], len(question.Options))
			for i, opt := range question.Options {
				options[i] = huh.NewOption(opt, opt)
			}

			field = huh.NewMultiSelect[string]().
				Title(question.Prompt).
				Options(options...).
				Value(&defaultValues)

			answers[question.Name] = defaultValues

		default:
			options := make([]huh.Option[string], len(question.Options))
			for i, opt := range question.Options {
				options[i] = huh.NewOption(opt, opt)
			}

			field = huh.NewSelect[string]().
				Title(question.Prompt).
				Options(options...).
				Value(&question.Default)

			answers[question.Name] = question.Default
		}

		group := huh.NewGroup(field)
		groupFields = append(groupFields, group)
	}

	return huh.NewForm(groupFields...), answers, nil
}

// RunSurveyWithRandomSelects runs the survey but generates random answers for select questions if runSurvey is false
func RunSurveyWithRandomSelects(profilePath, surveyKey string, runSurvey bool) map[string]interface{} {
	surveyValues := make(map[string]interface{})
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// READ PROFILE AND SURVEY BY KEY
	survey, _ := LoadQuestionFile(profilePath, surveyKey)

	// IF SURVEY EXISTS, RUN IT
	if len(survey) > 0 {
		if !runSurvey {
			// Generate random answers for select questions
			for _, question := range survey {
				if question.Kind == "select" && len(question.Options) > 0 {
					randomIndex := r.Intn(len(question.Options))
					question.Default = question.Options[randomIndex]
				}
				surveyValues[question.Name] = question.Default
			}
			return surveyValues
		}

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
