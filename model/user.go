package model

const STATUS_WAPPR = "WAITING_FOR_APPROVAL"
const STATUS_APPROVED = "APPROVED"
const STATUS_REJECTED = "REJECTED"

type RegisterRequest struct {
	Email    string `json:"email" example:"test@email.net"`
	Fullname string `json:"fullname" example:"test user"`
	Password string `json:"password" example:"mysecret123"`
	Stage    string `swaggerignore:"true"`
}

type ApprovalRequest struct {
	ActivityId    string `json:"activity_id"`
	WorkflowId    string `json:"workflow_id"`
	WorkflowRunId string `json:"workflow_run_id"`
	IsApproved    bool   `json:"is_approved"`
}

type WorkflowResponse struct {
	WorkflowId    string `json:"workflow_id"`
	WorkflowRunId string `json:"workflow_run_id"`
}
