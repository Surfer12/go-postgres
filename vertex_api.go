package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/aiplatform/v1"
	"google.golang.org/api/option"
)

// VertexAIClient is a wrapper for the Vertex AI client
type VertexAIClient struct {
	client *aiplatform.Service
}

// NewVertexAIClient creates a new Vertex AI client
func NewVertexAIClient(projectID, location, endpointID string) (*VertexAIClient, error) {
	ctx := context.Background()
	client, err := aiplatform.NewService(ctx, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
	if err != nil {
		return nil, fmt.Errorf("error creating Vertex AI client: %v", err)
	}
	return &VertexAIClient{client: client}, nil
}

// Predict sends a prediction request to Vertex AI
func (c *VertexAIClient) Predict(instances [][]float64) ([]interface{}, error) {
	req := &aiplatform.GoogleCloudAiplatformV1PredictRequest{
		Instances: instances,
	}
	resp, err := c.client.Projects.Locations.Endpoints.Predict("projects/"+os.Getenv("PROJECT_ID")+"/locations/"+os.Getenv("REGION")+"/endpoints/"+os.Getenv("ENDPOINT_ID"), req).Do()
	if err != nil {
		return nil, fmt.Errorf("error making prediction: %v", err)
	}
	return resp.Predictions, nil
}

// InputSchema defines the structure for input validation
type InputSchema struct {
	Instances [][]float64 `json:"instances"`
}

func main() {
	http.HandleFunc("/predict", func(w http.ResponseWriter, r *http.Request) {
		var input InputSchema
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, fmt.Sprintf("{\"error\": \"Request must be in JSON format\", \"code\": \"INVALID_FORMAT\"}"), http.StatusBadRequest)
			return
		}

		if len(input.Instances) == 0 {
			http.Error(w, fmt.Sprintf("{\"error\": \"JSON must contain 'instances' key\", \"code\": \"MISSING_KEY\"}"), http.StatusBadRequest)
			return
		}

		client, err := NewVertexAIClient(os.Getenv("PROJECT_ID"), os.Getenv("REGION"), os.Getenv("ENDPOINT_ID"))
		if err != nil {
			http.Error(w, fmt.Sprintf("{\"error\": \"Internal server error\", \"code\": \"INTERNAL_ERROR\", \"details\": \"%s\"}", err.Error()), http.StatusInternalServerError)
			return
		}

		predictions, err := client.Predict(input.Instances)
		if err != nil {
			log.Printf("Vertex AI error: %v", err)
			http.Error(w, fmt.Sprintf("{\"error\": \"Vertex AI prediction failed\", \"code\": \"VERTEX_AI_ERROR\", \"details\": \"%s\"}", err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(predictions)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
