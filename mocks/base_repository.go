// Package mocks
package mocks

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type BaseRepositoryMock struct {
	mock.Mock
}

func (m *BaseRepositoryMock) GetDialector() string {
	args := m.Called()
	return args.Get(0).(string)
}

func (m *BaseRepositoryMock) RunInTx(f func(db *gorm.DB) error) error {
	args := m.Called(f)
	return args.Error(0)
}

func (m *BaseRepositoryMock) VacuumOrOptimize() {
}
