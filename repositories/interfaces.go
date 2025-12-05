package repositories

import (
	"time"

	"github.com/muety/wakapi/models"
)

type UserRepository interface {
	FindOne(attributes models.User) (*models.User, error)
	GetByIds(userIds []string) ([]*models.User, error)
	GetAll() ([]*models.User, error)
	GetMany(ids []string) ([]*models.User, error)
	GetAllByReports(reportsEnabled bool) ([]*models.User, error)
	GetAllByLeaderboard(leaderboardEnabled bool) ([]*models.User, error)
	GetByLoggedInAfter(t time.Time) ([]*models.User, error)
	GetByLoggedInBefore(t time.Time) ([]*models.User, error)
	GetByLastActiveAfter(t time.Time) ([]*models.User, error)
	Count() (int64, error)
	InsertOrGet(user *models.User) (*models.User, bool, error)
	Update(user *models.User) (*models.User, error)
	UpdateField(user *models.User, key string, value interface{}) (*models.User, error)
	Delete(user *models.User) error
}
