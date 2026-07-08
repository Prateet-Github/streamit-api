package models

type ViewHeartbeatRequest struct {
	Elapsed  int    `json:"elapsed"`
	ViewerID string `json:"viewerId,omitempty"`
}
