package service

import (
	"project-portfolio-api/internal/model"
	"project-portfolio-api/internal/model/request"
	"project-portfolio-api/internal/repository"
	"project-portfolio-api/pkg/custom_error"
)

type ProjectService interface {
	Create(req *request.CreateProjectRequest, images []string) error
	GetAll() ([]model.Project, error)
	GetByID(id uint) (*model.Project, error)
	Update(id uint, req *request.UpdateProjectRequest) error
	Delete(id uint) error
}

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &projectService{repo: repo}
}

func (s *projectService) Create(req *request.CreateProjectRequest, images []string) error {
	project := &model.Project{
		Title:       req.Title,
		Description: req.Description,
		Tools:       req.Tools,
		Images:      images,
	}
	return s.repo.Create(project)
}

func (s *projectService) GetAll() ([]model.Project, error) {
	return s.repo.GetAll()
}

func (s *projectService) GetByID(id uint) (*model.Project, error) {
	project, err := s.repo.GetByID(id)
	if err != nil {
		return nil, custom_error.NewAppError(404, "Project not found")
	}
	return project, nil
}

func (s *projectService) Update(id uint, req *request.UpdateProjectRequest) error {
	project, err := s.repo.GetByID(id)
	if err != nil {
		return custom_error.NewAppError(404, "Project not found")
	}

	if req.Title != "" {
		project.Title = req.Title
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if len(req.Tools) > 0 {
		project.Tools = req.Tools
	}

	return s.repo.Update(project)
}

func (s *projectService) Delete(id uint) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return custom_error.NewAppError(404, "Project not found")
	}
	return s.repo.Delete(id)
}
