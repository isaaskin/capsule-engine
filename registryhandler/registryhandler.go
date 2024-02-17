package registryhandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/isaaskin/capsule-engine/models"
)

const REGISTRY = "http://hub.docker.com/v2"
const NAMESPACE = "isaaskin"

func GetRepositoryList() ([]models.CapsuleTemplate, error) {
	url := REGISTRY + "/repositories/" + NAMESPACE

	capsuleTemplates := []models.CapsuleTemplate{}

	// Make the GET request
	resp, err := http.Get(url)
	if err != nil {
		return capsuleTemplates, errors.New(fmt.Sprintln("Failed to fetch data: ", err))
	}
	defer resp.Body.Close()

	// Decode the JSON response
	var repositoryList RepositoryList
	err = json.NewDecoder(resp.Body).Decode(&repositoryList)
	if err != nil {
		return capsuleTemplates, errors.New(fmt.Sprintln("Failed to decode JSON: ", err))
	}

	// Create Capsule templates
	for _, repository := range repositoryList.Results {
		capsuleTemplates = append(capsuleTemplates, models.CapsuleTemplate{
			Name:        repository.Name,
			Namespace:   repository.Namespace,
			Description: repository.Description,
		})
	}

	return capsuleTemplates, nil
}
