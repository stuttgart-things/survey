package survey

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	allAnswers = make(map[string]interface{})
)

func GetRandomAnswers(questions []*Question) map[string]interface{} {

	for _, q := range questions {

		r := rand.New(rand.NewSource(time.Now().UnixNano()))

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
				if fn, ok := DefaultFunctions[q.DefaultFunction]; ok {
					q.Default = fn(q.DefaultParams)
				} else {
					log.Printf("FUNCTION %s NOT FOUND, USING EMPTY STRING", q.DefaultFunction)
					q.Default = ""
				}
			}
			// IF NO DEFAULT WAS GENERATED, USE A RANDOM STRING
			if q.Default == "" {
				q.Default = generateRandomValue(q, r)
			}
		}

		// CONVERT TO PROPER TYPE
		allAnswers[q.Name] = ConvertToType(q.Default, q.Type)
	}

	return allAnswers
}

func generateRandomValue(q *Question, r *rand.Rand) string {
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

func ConvertToType(value string, typ string) interface{} {
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
