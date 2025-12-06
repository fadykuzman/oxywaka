package repositories

import (
	"database/sql"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	conf "github.com/muety/wakapi/config"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const chunkSize = 4096

type BaseRepository struct {
	db *gorm.DB
}

func NewBaseRepository(db *gorm.DB) BaseRepository {
	return BaseRepository{db: db}
}

func (r *BaseRepository) GetDialector() string {
	return r.db.Dialector.Name()
}

func (r *BaseRepository) RunInTx(f func(tx *gorm.DB) error) error {
	return r.db.Transaction(f)
}

func (r *BaseRepository) VacuumOrOptimize() {
	// postgres require manual vacuuming regularly to reclaim free storage from deleted records
	// see https://www.postgresql.org/docs/current/sql-vacuum.html
	// also see https://github.com/muety/wakapi/issues/785
	t0 := time.Now()

	if err := r.db.Exec("vacuum").Error; err != nil {
		conf.Log().Error("vacuuming failed", "error", err.Error())
		return
	}
	conf.Log().Info("vacuuming done", "time_elapsed", time.Since(t0))
}

func InsertBatchChunked[T any](data []T, model T, db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		chunks := slice.Chunk[T](data, chunkSize)
		for _, chunk := range chunks {
			if err := insertBatch[T](chunk, model, tx); err != nil {
				return err
			}
		}
		return nil
	})
}

func insertBatch[T any](data []T, model T, db *gorm.DB) error {
	if err := db.
		Clauses(clause.OnConflict{
			DoNothing: true,
		}).
		Model(model).
		Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func streamRows[T any](rows *sql.Rows, channel chan *T, db *gorm.DB, onErr func(error)) {
	defer close(channel)
	defer rows.Close()
	for rows.Next() {
		var item T
		if err := db.ScanRows(rows, &item); err != nil {
			onErr(err)
			continue
		}
		channel <- &item
	}
}

func streamRowsBatched[T any](rows *sql.Rows, channel chan []*T, db *gorm.DB, batchSize int, onErr func(error)) {
	defer close(channel)
	defer rows.Close()

	buffer := make([]*T, 0, batchSize)

	for rows.Next() {
		var item T
		if err := db.ScanRows(rows, &item); err != nil {
			onErr(err)
			continue
		}

		buffer = append(buffer, &item)

		if len(buffer) == batchSize {
			channel <- buffer
			buffer = make([]*T, 0, batchSize)
		}
	}

	if len(buffer) > 0 {
		channel <- buffer
	}
}

func filteredQuery(q *gorm.DB, filterMap map[string][]string) *gorm.DB {
	for col, vals := range filterMap {
		q = q.Where(col+" in ?", slice.Map[string, string](vals, func(i int, val string) string {
			// query for "unknown" projects, languages, etc.
			if val == "-" {
				return ""
			}
			return val
		}))
	}
	return q
}
