package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"project_service/models"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/joho/godotenv"
)

var ES *elasticsearch.Client

func InitElasticsearch() (*elasticsearch.Client, error) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Elasticsearch configuration
	esConfig := elasticsearch.Config{
		Addresses: []string{
			os.Getenv("ELASTICSEARCH_URL"),
		},
		Username: os.Getenv("ELASTICSEARCH_USERNAME"),
		Password: os.Getenv("ELASTICSEARCH_PASSWORD"),
	}

	// If no URL is provided, use default localhost
	if esConfig.Addresses[0] == "" {
		esConfig.Addresses = []string{"http://localhost:9200"}
	}

	client, err := elasticsearch.NewClient(esConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating elasticsearch client: %v", err)
	}

	// Test connection
	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("error connecting to elasticsearch: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	ES = client

	// Create index if it doesn't exist
	err = createProjectIndex(client)
	if err != nil {
		return nil, fmt.Errorf("error creating project index: %v", err)
	}

	return client, nil
}

func createProjectIndex(client *elasticsearch.Client) error {
	indexName := "projects"

	// Check if index exists
	req := esapi.IndicesExistsRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// If index exists, return
	if res.StatusCode == 200 {
		return nil
	}

	// Create index with mapping
	mapping := `{
		"mappings": {
			"properties": {
				"id": {
					"type": "integer"
				},
				"title": {
					"type": "text",
					"analyzer": "standard"
				},
				"description": {
					"type": "text",
					"analyzer": "standard"
				},
				"category": {
					"type": "keyword"
				},
				"region": {
					"type": "keyword"
				},
				"country": {
					"type": "keyword"
				},
				"verification_standard": {
					"type": "keyword"
				},
				"price_per_tonne": {
					"type": "float"
				},
				"total_capacity": {
					"type": "float"
				},
				"available_capacity": {
					"type": "float"
				},
				"project_developer": {
					"type": "keyword"
				},
				"project_url": {
					"type": "keyword"
				},
				"image_url": {
					"type": "keyword"
				},
				"status": {
					"type": "keyword"
				},
				"created_at": {
					"type": "date"
				},
				"updated_at": {
					"type": "date"
				}
			}
		},
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		}
	}`

	createReq := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(mapping),
	}

	createRes, err := createReq.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		return fmt.Errorf("error creating index: %s", createRes.String())
	}

	log.Printf("Successfully created Elasticsearch index: %s", indexName)
	return nil
}

// IndexProject indexes a project document in Elasticsearch
func IndexProject(client *elasticsearch.Client, project *models.Project) error {
	if client == nil {
		return nil // Skip if Elasticsearch is not configured
	}

	// Convert project to JSON
	doc := map[string]interface{}{
		"id":                    project.ID,
		"title":                 project.Title,
		"description":           project.Description,
		"category":              project.Category,
		"region":                project.Region,
		"country":               project.Country,
		"verification_standard": project.VerificationStandard,
		"price_per_tonne":       project.PricePerTonne,
		"total_capacity":        project.TotalCapacity,
		"available_capacity":    project.AvailableCapacity,
		"project_developer":     project.ProjectDeveloper,
		"project_url":           project.ProjectURL,
		"image_url":             project.ImageURL,
		"status":                project.Status,
		"created_at":            project.CreatedAt,
		"updated_at":            project.UpdatedAt,
	}

	// Index the document
	req := esapi.IndexRequest{
		Index:      "projects",
		DocumentID: fmt.Sprintf("%d", project.ID),
		Body:       strings.NewReader(mustJSON(doc)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

// DeleteProjectFromIndex removes a project document from Elasticsearch
func DeleteProjectFromIndex(client *elasticsearch.Client, projectID uint) error {
	if client == nil {
		return nil // Skip if Elasticsearch is not configured
	}

	req := esapi.DeleteRequest{
		Index:      "projects",
		DocumentID: fmt.Sprintf("%d", projectID),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 { // 404 is OK, document doesn't exist
		return fmt.Errorf("error deleting document: %s", res.String())
	}

	return nil
}

// Helper function to convert struct to JSON
func mustJSON(v interface{}) string {
	// This is a simplified version - in production, use proper JSON marshaling
	return fmt.Sprintf(`{"id": %d, "title": "%s", "description": "%s", "category": "%s", "region": "%s", "country": "%s", "verification_standard": "%s", "price_per_tonne": %.2f, "total_capacity": %.2f, "available_capacity": %.2f, "project_developer": "%s", "project_url": "%s", "image_url": "%s", "status": "%s", "created_at": "%s", "updated_at": "%s"}`,
		v.(map[string]interface{})["id"],
		v.(map[string]interface{})["title"],
		v.(map[string]interface{})["description"],
		v.(map[string]interface{})["category"],
		v.(map[string]interface{})["region"],
		v.(map[string]interface{})["country"],
		v.(map[string]interface{})["verification_standard"],
		v.(map[string]interface{})["price_per_tonne"],
		v.(map[string]interface{})["total_capacity"],
		v.(map[string]interface{})["available_capacity"],
		v.(map[string]interface{})["project_developer"],
		v.(map[string]interface{})["project_url"],
		v.(map[string]interface{})["image_url"],
		v.(map[string]interface{})["status"],
		v.(map[string]interface{})["created_at"],
		v.(map[string]interface{})["updated_at"],
	)
}
