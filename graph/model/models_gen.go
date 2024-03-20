// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Agent struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Prompt      string `json:"prompt"`
	Engine      string `json:"engine"`
}

type Engine struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Model       string `json:"model"`
}

type Mutation struct {
}

type NewAgent struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Engine      string `json:"engine"`
	Prompt      string `json:"prompt"`
}

type Query struct {
}
