package helpers

import (
	"regexp"

	"github.com/isaaskin/capsule-engine/models"

	"github.com/docker/docker/api/types"
)

func removeLeadingSlash(input string) string {
	re := regexp.MustCompile(`^/`)
	return re.ReplaceAllString(input, "")
}

func ConvertContainerToCapsule(containers []types.Container) []models.Capsule {
	capsules := []models.Capsule{}

	for _, container := range containers {
		capsules = append(capsules, models.Capsule{
			ID:     container.ID,
			Name:   removeLeadingSlash(container.Names[0]),
			Status: container.State,
			Image:  container.Image,
		})
	}

	return capsules
}
