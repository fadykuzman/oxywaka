package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
A quick note to myself including some clarifications about time zones.

For Postgres:
- There are `timestamp` and `timestamptz` columns, whereby the former seems to behave very strangely.
- Apparently, when storing dates, they'll just chop off zone information entirely, while for retrieval, all dates are interpreted as UTC
- If Wakapi runs in CEST, '2025-03-25T14:00:00 +02:00' will end up as '2025-03-25T14:00:00' in the database and become '2025-03-25T16:00:00 +02:00' when retrieved back (at least in case of our annoying models.CustomTime, because of https://github.com/muety/wakapi/blob/bc2096f4117275d110a84f5b367aa8fdb4bd87ba/models/shared.go#L95)
- For `timestamptz`, columns don't actually store tz information either (https://github.com/jackc/pgx/issues/520#issuecomment-479692198), but at least allow for correct retrieval. The driver will return dates already in the application's TZ
- Here (https://github.com/go-gorm/gorm/issues/4834) is an interesting discussion on the issue
- According to https://www.cybertec-postgresql.com/en/time-zone-management-in-postgresql/, good practice is to always use timestamptz, leaving conversions to the database itself
*/

func (c *dbConfig) GetDialector() gorm.Dialector {
	connectionString := connectionString(c)
	log.Println("Connecting to Postgres with connection string:", connectionString)
	return postgres.New(postgres.Config{
		DSN: connectionString,
	})
}

func connectionString(config *dbConfig) string {
	if len(config.DSN) > 0 {
		return config.DSN
	}

	sslmode := "disable"
	if config.Ssl {
		sslmode = "require"
	}

	// note: passing a `timezone` param here doesn't seem to have any effect, neither with `timestamp`, not for `timestamptz` columns
	// possibly related to https://github.com/go-gorm/postgres/issues/199 ?

	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Name,
		config.Password,
		sslmode,
	)
}
