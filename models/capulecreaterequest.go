package models

type CapsuleCreateRequest struct {
	CapsuleTemplate CapsuleTemplate
	Name            string `json:"name"`
	WorkingDir      string `json:"workingDir"`
	UseLocalSSHKeys bool   `json:"useLocalSSHKeys"`
}
