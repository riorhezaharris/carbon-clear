package handlers

import (
	"net/http"
	"project_service/models"
	"project_service/repositories"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ProjectHandler struct {
	repo  *repositories.ProjectRepository
	redis *redis.Client
}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{
		repo: repositories.NewProjectRepository(),
	}
}

// SetRedisClient sets the Redis client for caching
func (h *ProjectHandler) SetRedisClient(client *redis.Client) {
	h.redis = client
}

// CreateProject creates a new project
// @Summary Create a new project
// @Description Create a new carbon offset project (admin only)
// @Tags projects
// @Accept json
// @Produce json
// @Security AdminAuth
// @Param request body models.CreateProjectRequest true "Project details"
// @Success 201 {object} models.Project "Project created successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Failed to create project"
// @Router /api/v1/projects/admin [post]
func (h *ProjectHandler) CreateProject(c echo.Context) error {
	var req models.CreateProjectRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Convert request to project model
	project := &models.Project{
		Title:                req.Title,
		Description:          req.Description,
		Category:             req.Category,
		Region:               req.Region,
		Country:              req.Country,
		VerificationStandard: req.VerificationStandard,
		PricePerTonne:        req.PricePerTonne,
		TotalCapacity:        req.TotalCapacity,
		AvailableCapacity:    req.AvailableCapacity,
		ProjectDeveloper:     req.ProjectDeveloper,
		ProjectURL:           req.ProjectURL,
		ImageURL:             req.ImageURL,
		Status:               "active",
	}

	if err := h.repo.Create(project); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create project"})
	}

	// Invalidate cache
	h.invalidateProjectCache()

	return c.JSON(http.StatusCreated, project)
}

// GetProject retrieves a project by ID
// @Summary Get project by ID
// @Description Retrieve a specific project by its ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} models.Project "Project details"
// @Failure 400 {object} map[string]string "Invalid project ID"
// @Failure 404 {object} map[string]string "Project not found"
// @Failure 500 {object} map[string]string "Failed to retrieve project"
// @Router /api/v1/projects/{id} [get]
func (h *ProjectHandler) GetProject(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	// Try to get from cache first
	cacheKey := "project:" + idStr
	if h.redis != nil {
		if cached, err := h.getFromCache(cacheKey); err == nil {
			return c.JSON(http.StatusOK, cached)
		}
	}

	project, err := h.repo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Project not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve project"})
	}

	// Cache the result
	if h.redis != nil {
		h.setCache(cacheKey, project, 3600) // Cache for 1 hour
	}

	return c.JSON(http.StatusOK, project)
}

// GetAllProjects retrieves all projects with pagination
// @Summary Get all projects
// @Description Retrieve all carbon offset projects with pagination
// @Tags projects
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of results (max 100)" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} map[string]interface{} "List of projects"
// @Failure 500 {object} map[string]string "Failed to retrieve projects"
// @Router /api/v1/projects [get]
func (h *ProjectHandler) GetAllProjects(c echo.Context) error {
	// Parse pagination parameters
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := 10 // default limit
	offset := 0 // default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Try to get from cache first
	cacheKey := "projects:all:" + strconv.Itoa(limit) + ":" + strconv.Itoa(offset)
	if h.redis != nil {
		if cached, err := h.getFromCache(cacheKey); err == nil {
			return c.JSON(http.StatusOK, cached)
		}
	}

	projects, err := h.repo.GetAll(limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve projects"})
	}

	// Cache the result
	if h.redis != nil {
		h.setCache(cacheKey, projects, 1800) // Cache for 30 minutes
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"projects": projects,
		"limit":    limit,
		"offset":   offset,
		"count":    len(projects),
	})
}

// UpdateProject updates a project
// @Summary Update project
// @Description Update a project's information (admin only)
// @Tags projects
// @Accept json
// @Produce json
// @Security AdminAuth
// @Param id path int true "Project ID"
// @Param request body models.UpdateProjectRequest true "Project update details"
// @Success 200 {object} map[string]string "Project updated successfully"
// @Failure 400 {object} map[string]string "Invalid project ID or request body"
// @Failure 404 {object} map[string]string "Project not found"
// @Failure 500 {object} map[string]string "Failed to update project"
// @Router /api/v1/projects/admin/{id} [put]
func (h *ProjectHandler) UpdateProject(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	var req models.UpdateProjectRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	if err := h.repo.Update(uint(id), &req); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Project not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update project"})
	}

	// Invalidate cache
	h.invalidateProjectCache()
	h.invalidateProjectCacheByID(idStr)

	return c.JSON(http.StatusOK, map[string]string{"message": "Project updated successfully"})
}

// DeleteProject soft deletes a project
// @Summary Delete project
// @Description Delete a project (admin only)
// @Tags projects
// @Accept json
// @Produce json
// @Security AdminAuth
// @Param id path int true "Project ID"
// @Success 200 {object} map[string]string "Project deleted successfully"
// @Failure 400 {object} map[string]string "Invalid project ID"
// @Failure 404 {object} map[string]string "Project not found"
// @Failure 500 {object} map[string]string "Failed to delete project"
// @Router /api/v1/projects/admin/{id} [delete]
func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid project ID"})
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Project not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete project"})
	}

	// Invalidate cache
	h.invalidateProjectCache()
	h.invalidateProjectCacheByID(idStr)

	return c.JSON(http.StatusOK, map[string]string{"message": "Project deleted successfully"})
}

// SearchProjects searches projects with filters
// @Summary Search projects
// @Description Search and filter carbon offset projects
// @Tags projects
// @Accept json
// @Produce json
// @Param request body models.ProjectSearchRequest true "Search filters"
// @Success 200 {object} map[string]interface{} "Search results"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Search failed"
// @Router /api/v1/projects/search [post]
func (h *ProjectHandler) SearchProjects(c echo.Context) error {
	var req models.ProjectSearchRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Set defaults
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Try to get from cache first
	cacheKey := h.generateSearchCacheKey(&req)
	if h.redis != nil {
		if cached, err := h.getFromCache(cacheKey); err == nil {
			return c.JSON(http.StatusOK, cached)
		}
	}

	projects, err := h.repo.Search(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Search failed"})
	}

	result := map[string]interface{}{
		"projects": projects,
		"limit":    req.Limit,
		"offset":   req.Offset,
		"count":    len(projects),
	}

	// Cache the result
	if h.redis != nil {
		h.setCache(cacheKey, result, 900) // Cache for 15 minutes
	}

	return c.JSON(http.StatusOK, result)
}

// GetProjectCategories retrieves all available categories
// @Summary Get project categories
// @Description Retrieve all available project categories
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of categories"
// @Failure 500 {object} map[string]string "Failed to retrieve categories"
// @Router /api/v1/projects/categories [get]
func (h *ProjectHandler) GetProjectCategories(c echo.Context) error {
	// Try to get from cache first
	cacheKey := "project:categories"
	if h.redis != nil {
		if cached, err := h.getFromCache(cacheKey); err == nil {
			return c.JSON(http.StatusOK, cached)
		}
	}

	categories, err := h.repo.GetCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve categories"})
	}

	// Cache the result
	if h.redis != nil {
		h.setCache(cacheKey, categories, 3600) // Cache for 1 hour
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"categories": categories})
}

// GetProjectRegions retrieves all available regions
// @Summary Get project regions
// @Description Retrieve all available project regions
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of regions"
// @Failure 500 {object} map[string]string "Failed to retrieve regions"
// @Router /api/v1/projects/regions [get]
func (h *ProjectHandler) GetProjectRegions(c echo.Context) error {
	// Try to get from cache first
	cacheKey := "project:regions"
	if h.redis != nil {
		if cached, err := h.getFromCache(cacheKey); err == nil {
			return c.JSON(http.StatusOK, cached)
		}
	}

	regions, err := h.repo.GetRegions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve regions"})
	}

	// Cache the result
	if h.redis != nil {
		h.setCache(cacheKey, regions, 3600) // Cache for 1 hour
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"regions": regions})
}

// GetProjectCountries retrieves all available countries
// @Summary Get project countries
// @Description Retrieve all available project countries
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of countries"
// @Failure 500 {object} map[string]string "Failed to retrieve countries"
// @Router /api/v1/projects/countries [get]
func (h *ProjectHandler) GetProjectCountries(c echo.Context) error {
	// Try to get from cache first
	cacheKey := "project:countries"
	if h.redis != nil {
		if cached, err := h.getFromCache(cacheKey); err == nil {
			return c.JSON(http.StatusOK, cached)
		}
	}

	countries, err := h.repo.GetCountries()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve countries"})
	}

	// Cache the result
	if h.redis != nil {
		h.setCache(cacheKey, countries, 3600) // Cache for 1 hour
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"countries": countries})
}

// Cache helper methods
func (h *ProjectHandler) getFromCache(key string) (interface{}, error) {
	if h.redis == nil {
		return nil, redis.Nil
	}
	// Implementation would use h.redis.Get() and JSON unmarshaling
	return nil, redis.Nil
}

func (h *ProjectHandler) setCache(key string, value interface{}, expiration int) {
	if h.redis == nil {
		return
	}
	// Implementation would use h.redis.Set() with JSON marshaling
}

func (h *ProjectHandler) invalidateProjectCache() {
	if h.redis == nil {
		return
	}
	// Implementation would delete cache keys with pattern "projects:*"
}

func (h *ProjectHandler) invalidateProjectCacheByID(id string) {
	if h.redis == nil {
		return
	}
	// Implementation would delete cache key "project:" + id
}

func (h *ProjectHandler) generateSearchCacheKey(req *models.ProjectSearchRequest) string {
	// Generate a cache key based on search parameters
	// This is a simplified version - in production, you'd want a more sophisticated key generation
	return "search:" + req.Query + ":" + strconv.Itoa(req.Limit) + ":" + strconv.Itoa(req.Offset)
}
