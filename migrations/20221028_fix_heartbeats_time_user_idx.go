package migrations

import (
	"log/slog"

	"github.com/muety/wakapi/config"
	"github.com/muety/wakapi/models"
	"gorm.io/gorm"
)

// due to an error in the model definition, idx_time_user used to only cover 'user_id', but not time column
// if that's the case in the current state of the database, drop the index and let it be recreated by auto migration afterwards
func init() {
	const name = "20221028-fix_heartbeats_time_user_idx"
	f := migrationFunc{
		name: name,
		f: func(db *gorm.DB, cfg *config.Config) error {
			migrator := db.Migrator()

			if !migrator.HasTable(&models.Heartbeat{}) {
				return nil
			}

			var drop bool
			indexes, err := migrator.GetIndexes(&models.Heartbeat{})
			if err != nil {
				return err
			}

			for _, idx := range indexes {
				if idx.Table() == "heartbeats" && idx.Name() == "idx_time_user" && len(idx.Columns()) == 1 {
					drop = true
					break
				}
			}

			if !drop {
				return nil
			}

			if err := migrator.DropIndex(&models.Heartbeat{}, "idx_time_user"); err != nil {
				return err
			}
			slog.Info("index 'idx_time_user' needs to be recreated, this may take a while")

			return nil
		},
	}

	registerPreMigration(f)
}
