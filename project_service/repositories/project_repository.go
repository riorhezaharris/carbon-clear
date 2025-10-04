package repositories

import (
	"context"
	"fmt"
	"project_service/config"
	"project_service/models"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db    *gorm.DB
	es    *elasticsearch.Client
	index string
}

func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{
		db:    config.DB,
		index: "projects",
	}
}

// SetElasticsearchClient sets the Elasticsearch client
func (r *ProjectRepository) SetElasticsearchClient(client *elasticsearch.Client) {
	r.es = client
}

// Create creates a new project
func (r *ProjectRepository) Create(project *models.Project) error {
	if err := r.db.Create(project).Error; err != nil {
		return err
	}

	// Index in Elasticsearch
	if r.es != nil {
		go r.indexProject(project)
	}

	return nil
}

// GetByID retrieves a project by ID
func (r *ProjectRepository) GetByID(id uint) (*models.Project, error) {
	var project models.Project
	if err := r.db.First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// GetAll retrieves all projects with pagination
func (r *ProjectRepository) GetAll(limit, offset int) ([]models.Project, error) {
	var projects []models.Project
	if err := r.db.Where("status = ?", "active").
		Limit(limit).Offset(offset).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

// Update updates a project
func (r *ProjectRepository) Update(id uint, project *models.UpdateProjectRequest) error {
	if err := r.db.Model(&models.Project{}).Where("id = ?", id).Updates(project).Error; err != nil {
		return err
	}

	// Update in Elasticsearch
	if r.es != nil {
		go r.indexProjectByID(id)
	}

	return nil
}

// Delete soft deletes a project (sets status to inactive)
func (r *ProjectRepository) Delete(id uint) error {
	if err := r.db.Model(&models.Project{}).Where("id = ?", id).
		Update("status", "inactive").Error; err != nil {
		return err
	}

	// Remove from Elasticsearch
	if r.es != nil {
		go r.deleteFromElasticsearch(id)
	}

	return nil
}

// Search searches projects using Elasticsearch
func (r *ProjectRepository) Search(searchReq *models.ProjectSearchRequest) ([]models.Project, error) {
	if r.es == nil {
		// Fallback to database search if Elasticsearch is not available
		return r.databaseSearch(searchReq)
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{},
				"filter": []interface{}{
					map[string]interface{}{
						"term": map[string]interface{}{
							"status": "active",
						},
					},
				},
			},
		},
		"from": searchReq.Offset,
		"size": searchReq.Limit,
	}

	// Add text search
	if searchReq.Query != "" {
		textQuery := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  searchReq.Query,
				"fields": []string{"title^2", "description", "category", "region", "country"},
			},
		}
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]interface{}),
			textQuery,
		)
	}

	// Add category filter
	if len(searchReq.Category) > 0 {
		categoryFilter := map[string]interface{}{
			"terms": map[string]interface{}{
				"category": searchReq.Category,
			},
		}
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{}),
			categoryFilter,
		)
	}

	// Add region filter
	if len(searchReq.Region) > 0 {
		regionFilter := map[string]interface{}{
			"terms": map[string]interface{}{
				"region": searchReq.Region,
			},
		}
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{}),
			regionFilter,
		)
	}

	// Add country filter
	if len(searchReq.Country) > 0 {
		countryFilter := map[string]interface{}{
			"terms": map[string]interface{}{
				"country": searchReq.Country,
			},
		}
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{}),
			countryFilter,
		)
	}

	// Add price range filter
	if searchReq.MinPrice > 0 || searchReq.MaxPrice > 0 {
		priceRange := map[string]interface{}{}
		if searchReq.MinPrice > 0 {
			priceRange["gte"] = searchReq.MinPrice
		}
		if searchReq.MaxPrice > 0 {
			priceRange["lte"] = searchReq.MaxPrice
		}
		priceFilter := map[string]interface{}{
			"range": map[string]interface{}{
				"price_per_tonne": priceRange,
			},
		}
		query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"] = append(
			query["query"].(map[string]interface{})["bool"].(map[string]interface{})["filter"].([]interface{}),
			priceFilter,
		)
	}

	// Execute search
	req := esapi.SearchRequest{
		Index: []string{r.index},
		Body:  strings.NewReader(fmt.Sprintf("%v", query)),
	}

	res, err := req.Do(context.Background(), r.es)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch search error: %s", res.String())
	}

	// Parse response and convert to projects
	// This is a simplified version - in production, you'd want proper JSON parsing
	return r.parseSearchResponse(res)
}

// databaseSearch is a fallback search method using database
func (r *ProjectRepository) databaseSearch(searchReq *models.ProjectSearchRequest) ([]models.Project, error) {
	query := r.db.Where("status = ?", "active")

	if searchReq.Query != "" {
		searchPattern := "%" + searchReq.Query + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ? OR category ILIKE ? OR region ILIKE ? OR country ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	if len(searchReq.Category) > 0 {
		query = query.Where("category IN ?", searchReq.Category)
	}

	if len(searchReq.Region) > 0 {
		query = query.Where("region IN ?", searchReq.Region)
	}

	if len(searchReq.Country) > 0 {
		query = query.Where("country IN ?", searchReq.Country)
	}

	if searchReq.MinPrice > 0 {
		query = query.Where("price_per_tonne >= ?", searchReq.MinPrice)
	}

	if searchReq.MaxPrice > 0 {
		query = query.Where("price_per_tonne <= ?", searchReq.MaxPrice)
	}

	var projects []models.Project
	if err := query.Limit(searchReq.Limit).Offset(searchReq.Offset).Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

// indexProject indexes a project in Elasticsearch
func (r *ProjectRepository) indexProject(project *models.Project) {
	// Implementation for indexing project in Elasticsearch
	// This would involve creating/updating the document in the index
}

// indexProjectByID indexes a project by ID in Elasticsearch
func (r *ProjectRepository) indexProjectByID(id uint) {
	// Fetch project from database and index it
	project, err := r.GetByID(id)
	if err != nil {
		return
	}
	r.indexProject(project)
}

// deleteFromElasticsearch removes a project from Elasticsearch index
func (r *ProjectRepository) deleteFromElasticsearch(id uint) {
	// Implementation for deleting project from Elasticsearch
}

// parseSearchResponse parses Elasticsearch search response
func (r *ProjectRepository) parseSearchResponse(res *esapi.Response) ([]models.Project, error) {
	// Implementation for parsing Elasticsearch response
	// This is a placeholder - proper implementation would parse JSON response
	return []models.Project{}, nil
}

// GetCategories retrieves all unique categories
func (r *ProjectRepository) GetCategories() ([]string, error) {
	var categories []string
	if err := r.db.Model(&models.Project{}).
		Where("status = ?", "active").
		Distinct("category").
		Pluck("category", &categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetRegions retrieves all unique regions
func (r *ProjectRepository) GetRegions() ([]string, error) {
	var regions []string
	if err := r.db.Model(&models.Project{}).
		Where("status = ?", "active").
		Distinct("region").
		Pluck("region", &regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

// GetCountries retrieves all unique countries
func (r *ProjectRepository) GetCountries() ([]string, error) {
	var countries []string
	if err := r.db.Model(&models.Project{}).
		Where("status = ?", "active").
		Distinct("country").
		Pluck("country", &countries).Error; err != nil {
		return nil, err
	}
	return countries, nil
}
