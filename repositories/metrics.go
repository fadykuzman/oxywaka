package repositories

import (
	"github.com/muety/wakapi/config"
	"gorm.io/gorm"
)

type MetricsRepository struct {
	BaseRepository
	config *config.Config
}

const sizeTplPostgres = `SELECT pg_database_size(?);`

func NewMetricsRepository(db *gorm.DB) *MetricsRepository {
	return &MetricsRepository{BaseRepository: NewBaseRepository(db), config: config.Get()}
}

func (srv *MetricsRepository) GetDatabaseSize() (size int64, err error) {
	cfg := srv.config.Db

	query := srv.db.Raw("SELECT 0")
	query = srv.db.Raw(sizeTplPostgres, cfg.Name)

	err = query.Scan(&size).Error
	return size, err
}
