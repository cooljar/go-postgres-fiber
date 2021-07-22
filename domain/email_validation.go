package domain

import "encoding/json"

type EmailValidation struct {
	Email string `json:"email"`
	SecretCode string `json:"secret_code"`
}

// FromJSON decode json to user struct
func (e *EmailValidation) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, e)
}

// ToJSON encode user struct to json
func (e *EmailValidation) ToJSON() []byte {
	str, _ := json.Marshal(e)
	return str
}
