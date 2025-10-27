package db

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/epos-eu/converter-routine/logging"
	sloggorm "github.com/orandin/slog-gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	log      = logging.Get("database")
	dnsRegex = regexp.MustCompile(`(&?(targetServerType|loadBalanceHosts|readOnly)=[^&]+)`)
)

var (
	converterDB *gorm.DB
	once        sync.Once
	initErr     error
)

func Connect() (*gorm.DB, error) {
	once.Do(func() {
		converterDB, initErr = initializeConnection("POSTGRESQL_CONNECTION_STRING", "CONVERTER_CATALOGUE_CONNECTION_STRING")
	})
	return converterDB, initErr
}

func initializeConnection(envVars ...string) (*gorm.DB, error) {
	for _, envVar := range envVars {

		dsn, err := parseAndCleanDSN(envVar)
		if err != nil {
			log.Error("failed to parse DSN", "env_var", envVar, "error", err)
			continue
		}

		gormLogger := sloggorm.New(
			sloggorm.WithHandler(logging.Get("gorm").Handler()),
			sloggorm.WithSlowThreshold(200*time.Millisecond),
			sloggorm.WithRecordNotFoundError(),
		)

		const maxRetries = 3
		for attempt := range maxRetries {
			if attempt > 0 {
				backoff := time.Duration(1<<uint(attempt-1)) * time.Second
				log.Info("retrying database connection", "attempt", attempt+1, "backoff", backoff)
				time.Sleep(backoff)
			}

			log.Info("connecting to database", "env_var", envVar, "dsn", sanitizeDSN(dsn))

			db, err := gorm.Open(postgres.New(postgres.Config{
				DriverName: "pgx",
				DSN:        dsn,
			}), &gorm.Config{
				Logger: gormLogger,
				NamingStrategy: schema.NamingStrategy{
					TablePrefix:   "",
					SingularTable: true,
				},
			})
			if err != nil {
				log.Error("failed to connect to database", "env_var", envVar, "error", err)
				if attempt == maxRetries-1 {
					log.Error("all retries failed for DSN", "env_var", envVar, "error", err)
					break
				}
				continue
			}

			log.Info("successfully connected to database", "env_var", envVar)
			return db, nil
		}
	}
	return nil, fmt.Errorf("failed to connect to any database: all connection strings exhausted")
}

func parseAndCleanDSN(envVar string) (string, error) {
	dsn, ok := os.LookupEnv(envVar)
	if !ok {
		log.Error("environment variable not set", "env_var", envVar)
		return "", fmt.Errorf("%s is not set", envVar)
	}

	log.Debug("parsing DSN from environment variable", "env_var", envVar)

	dsn = strings.TrimPrefix(dsn, "jdbc:")
	dsn = dnsRegex.ReplaceAllString(dsn, "")
	dsn = strings.TrimRight(dsn, "?&")

	separator := "?"
	if strings.Contains(dsn, "?") {
		separator = "&"
	}

	dsn = dsn + separator + "sslmode=disable&target_session_attrs=read-write"

	log.Debug("DSN parsed and configured", "env_var", envVar)
	return dsn, nil
}

func sanitizeDSN(dsn string) string {
	passwordRe := regexp.MustCompile(`password=[^&\s]+`)
	return passwordRe.ReplaceAllString(dsn, "password=***")
}
