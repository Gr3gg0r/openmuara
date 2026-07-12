package testsdk

// ScenarioResponse is returned by the admin scenario simulation endpoint.
type ScenarioResponse struct {
	Reference string `json:"reference"`
	Status    string `json:"status"`
	Scenario  string `json:"scenario"`
}
