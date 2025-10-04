package models

import "time"

type Project struct {
	ID                   uint      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Title                string    `json:"title" gorm:"column:title;not null"`
	Description          string    `json:"description" gorm:"column:description;type:text"`
	Category             string    `json:"category" gorm:"column:category;not null"`
	Region               string    `json:"region" gorm:"column:region;not null"`
	Country              string    `json:"country" gorm:"column:country;not null"`
	VerificationStandard string    `json:"verification_standard" gorm:"column:verification_standard;not null"`
	PricePerTonne        float64   `json:"price_per_tonne" gorm:"column:price_per_tonne;not null"`
	TotalCapacity        float64   `json:"total_capacity" gorm:"column:total_capacity"`
	AvailableCapacity    float64   `json:"available_capacity" gorm:"column:available_capacity"`
	ProjectDeveloper     string    `json:"project_developer" gorm:"column:project_developer"`
	ProjectURL           string    `json:"project_url" gorm:"column:project_url"`
	ImageURL             string    `json:"image_url" gorm:"column:image_url"`
	Status               string    `json:"status" gorm:"column:status;default:'active'"`
	CreatedAt            time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt            time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type CreateProjectRequest struct {
	Title                string  `json:"title" validate:"required"`
	Description          string  `json:"description" validate:"required"`
	Category             string  `json:"category" validate:"required"`
	Region               string  `json:"region" validate:"required"`
	Country              string  `json:"country" validate:"required"`
	VerificationStandard string  `json:"verification_standard" validate:"required"`
	PricePerTonne        float64 `json:"price_per_tonne" validate:"required,gt=0"`
	TotalCapacity        float64 `json:"total_capacity"`
	AvailableCapacity    float64 `json:"available_capacity"`
	ProjectDeveloper     string  `json:"project_developer"`
	ProjectURL           string  `json:"project_url"`
	ImageURL             string  `json:"image_url"`
}

type UpdateProjectRequest struct {
	Title                string  `json:"title"`
	Description          string  `json:"description"`
	Category             string  `json:"category"`
	Region               string  `json:"region"`
	Country              string  `json:"country"`
	VerificationStandard string  `json:"verification_standard"`
	PricePerTonne        float64 `json:"price_per_tonne"`
	TotalCapacity        float64 `json:"total_capacity"`
	AvailableCapacity    float64 `json:"available_capacity"`
	ProjectDeveloper     string  `json:"project_developer"`
	ProjectURL           string  `json:"project_url"`
	ImageURL             string  `json:"image_url"`
	Status               string  `json:"status"`
}

type ProjectSearchRequest struct {
	Query    string   `json:"query"`
	Category []string `json:"category"`
	Region   []string `json:"region"`
	Country  []string `json:"country"`
	MinPrice float64  `json:"min_price"`
	MaxPrice float64  `json:"max_price"`
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
}

type ProjectResponse struct {
	ID                   uint      `json:"id"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	Category             string    `json:"category"`
	Region               string    `json:"region"`
	Country              string    `json:"country"`
	VerificationStandard string    `json:"verification_standard"`
	PricePerTonne        float64   `json:"price_per_tonne"`
	TotalCapacity        float64   `json:"total_capacity"`
	AvailableCapacity    float64   `json:"available_capacity"`
	ProjectDeveloper     string    `json:"project_developer"`
	ProjectURL           string    `json:"project_url"`
	ImageURL             string    `json:"image_url"`
	Status               string    `json:"status"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}
