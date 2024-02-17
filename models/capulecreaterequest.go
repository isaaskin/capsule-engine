package models

type CapsuleCreateRequest struct {
	CapsuleTemplate CapsuleTemplate
	Name            string
	WorkingDir      string
}
