package services

import "Go-MCS/models"

type AuthService interface {
	SignUpUser(*models.SignUpInput) (*models.DBResponse, error)
	// SignInUser(*models.SignInInput) (*models.DBResponse, error)
}
