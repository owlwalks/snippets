package snippets

import (
	"encoding/json"
)

func publicStruct() ([]byte, error) {
	type omit *struct{}
	type priv struct {
		Password string `json:"password"`
	}
	type pub struct {
		*priv
		Password omit `json:"password,omitempty"`
	}
	return json.Marshal(pub{priv: &priv{"password"}})
}
