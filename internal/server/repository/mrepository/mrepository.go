package mrepository

import (
	"fmt"

	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/config/db"
	"github.com/dmitryDevGoMid/go-service-collect-metrics/internal/server/repository"
)

type ManagerRepository interface {
	GetRepositoryActive() repository.MetricsRepository
}

type managerRepository struct {
	repositoryDB     repository.MetricsRepository
	repositoryMemory repository.MetricsRepository
	db               db.Connection
}

func NewMamangerRepository(repositoryDB repository.MetricsRepository,
	repositoryMemory repository.MetricsRepository,
	db db.Connection) ManagerRepository {
	return &managerRepository{repositoryDB: repositoryDB, repositoryMemory: repositoryMemory, db: db}

}

func (m *managerRepository) GetRepositoryActive() repository.MetricsRepository {
	err := m.db.Ping()

	if err != nil {
		fmt.Println("Active Local")
		return m.repositoryMemory
	}

	fmt.Println("Active DB")
	return m.repositoryDB
}
