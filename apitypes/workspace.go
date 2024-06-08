package apitypes

type Workspace struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

type CreateWorkspaceBody struct {
	Name string `json:"name"`
}
