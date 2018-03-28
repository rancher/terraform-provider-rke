package rke

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func validateIntegerInRange(min, max int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if value < min {
			errors = append(errors, fmt.Errorf(
				"%q cannot be lower than %d: %d", k, min, value))
		}
		if value > max {
			errors = append(errors, fmt.Errorf(
				"%q cannot be higher than %d: %d", k, max, value))
		}
		return
	}
}

func validateStringInWord(allowWords []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		var found bool
		for _, t := range allowWords {
			if v.(string) == t {
				found = true
			}
		}
		if !found {
			errors = append(errors, fmt.Errorf("%q must be one of [%s]", k, strings.Join(allowWords, "/")))

		}
		return
	}
}
