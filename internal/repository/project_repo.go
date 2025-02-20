package repository

import (
	"project-portfolio-api/internal/model"

	// "gorm.io/gorm"
	"database/sql"
	"log"

	"github.com/lib/pq"
)

type ProjectRepository interface {
	Create(project *model.Project) error
	GetAll() ([]model.Project, error)
	GetByID(id uint) (*model.Project, error)
	Update(project *model.Project) error
	Delete(id uint) error
}

type projectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(project *model.Project) error {
	query := `
        INSERT INTO projects (title, description, images, tools) 
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	// Log query dan values untuk debugging
	log.Printf("Executing query: %s", query)
	log.Printf("Values: title=%s, desc=%s, images=%v, tools=%v",
		project.Title,
		project.Description,
		project.Images,
		project.Tools,
	)

	err := r.db.QueryRow(
		query,
		project.Title,
		project.Description,
		pq.Array(project.Images),
		pq.Array(project.Tools),
	).Scan(&project.ID)

	if err != nil {
		log.Printf("Error executing query: %v", err)
		return err
	}

	log.Printf("Successfully inserted project with ID: %d", project.ID)
	return nil
}

func (r *projectRepository) GetAll() ([]model.Project, error) {
	query := `SELECT id, title, description, images, tools FROM projects`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []model.Project
	for rows.Next() {
		var project model.Project
		var images pq.StringArray
		var tools pq.StringArray

		err := rows.Scan(
			&project.ID,
			&project.Title,
			&project.Description,
			(*pq.StringArray)(&images),
			(*pq.StringArray)(&tools),
		)
		if err != nil {
			return nil, err
		}

		project.Images = images
		project.Tools = tools
		projects = append(projects, project)
	}

	return projects, nil
}

func (r *projectRepository) GetByID(id uint) (*model.Project, error) {
	query := `
        SELECT id, title, description, images, tools 
        FROM projects 
        WHERE id = $1`

	var project model.Project

	err := r.db.QueryRow(query, id).Scan(
		&project.ID,
		&project.Title,
		&project.Description,
		(*pq.StringArray)(&project.Images),
		(*pq.StringArray)(&project.Tools),
	)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *projectRepository) Update(project *model.Project) error {
	query := `
        UPDATE projects 
        SET title = $1, description = $2, images = $3, tools = $4
        WHERE id = $5`

	result, err := r.db.Exec(
		query,
		project.Title,
		project.Description,
		pq.Array(project.Images),
		pq.Array(project.Tools),
		project.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *projectRepository) Delete(id uint) error {
	query := `DELETE FROM projects WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
