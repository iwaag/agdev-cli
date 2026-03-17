package api

type AcceptedResponse struct {
	Status  string            `json:"status"`
	Kind    string            `json:"kind"`
	Message string            `json:"message"`
	Inputs  map[string]string `json:"inputs,omitempty"`
}
