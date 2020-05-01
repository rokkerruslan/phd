package event

import (
	"fmt"
	"strings"
)

func validateForCreate(e Event) error {
	errors := e.baseValidation()

	if e.OwnerID == 0 {
		errors = append(errors, "`OwnerID` can't be empty")
	}

	if len(errors) == 0 {
		return nil
	}

	return fmt.Errorf("validateForCreate fails %v", strings.Join(errors, ", "))
}

func (e *Event) ValidateForUpdate() error {
	errors := e.baseValidation()

	if len(errors) == 0 {
		return nil
	}

	return fmt.Errorf("event.ValidateForUpdate fails: %v", strings.Join(errors, ", "))
}

func (e *Event) baseValidation() []string {
	var errors []string

	if e.Name == "" {
		errors = append(errors, "`Title` can't be empty")
	}
	if e.Description == "" {
		errors = append(errors, "`Description` can't be empty")
	}
	if e.Timelines == nil || len(e.Timelines) == 0 {
		errors = append(errors, "`Timelines` can't be empty")
	}
	for _, timeline := range e.Timelines {
		if tError := timeline.Validate(); tError != nil {
			errors = append(errors, tError.Error())
		}
	}

	return errors
}
