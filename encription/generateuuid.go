package encription

import "github.com/google/uuid"

func Generateuudi() string {

	createuuid := uuid.New()

	return createuuid.String()
}
