package models

type CapsuleCreateRequest struct {
	CapsuleTemplate CapsuleTemplate `json:"capsuleTemplate"`
	Name            string          `json:"name"`
	WorkingDir      string          `json:"workingDir"`
	UseLocalSSHKeys bool            `json:"useLocalSSHKeys"`
}
