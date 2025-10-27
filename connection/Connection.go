package connection

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/epos-eu/converter-routine/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	log      = logging.Get("database")
	dnsRegex = regexp.MustCompile(`(&?(targetServerType|loadBalanceHosts)=[^&]+)`)
)

var (
	converterDB *gorm.DB
	once        sync.Once
	initErr     error
)

func ConnectConverter() (*gorm.DB, error) {
	once.Do(func() {
		converterDB, initErr = initializeConnection("CONVERTER_CATALOGUE_CONNECTION_STRING")
	})
	return converterDB, initErr
}

func initializeConnection(envVar string) (*gorm.DB, error) {
	dsn, err := parseAndCleanDSN(envVar)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN for %s: %w", envVar, err)
	}

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
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "",
				SingularTable: true,
			},
		})
		if err != nil {
			log.Error("failed to connect to database", "env_var", envVar, "error", err)
			if attempt == maxRetries-1 {
				return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
			}
			continue
		}

		log.Info("successfully connected to database", "env_var", envVar)
		return db, nil
	}
	return nil, fmt.Errorf("unexpected error in connection retry")
}

func parseAndCleanDSN(envVar string) (string, error) {
	dsn, ok := os.LookupEnv(envVar)
	if !ok {
		log.Error("environment variable not set", "env_var", envVar)
		return "", fmt.Errorf("%s is not set", envVar)
	}

	log.Debug("parsing DSN from environment variable", "env_var", envVar)

	// Remove "jdbc:" prefix if present
	dsn = strings.TrimPrefix(dsn, "jdbc:")

	// Remove unsupported parameters
	dsn = dnsRegex.ReplaceAllString(dsn, "")

	// Clean up trailing "?" or "&"
	dsn = strings.TrimRight(dsn, "?&")

	// Determine separator
	separator := "?"
	if strings.Contains(dsn, "?") {
		separator = "&"
	}

	// Add pgx configuration
	dsn = dsn + separator + "target_session_attrs=read-write"

	log.Debug("DSN parsed and configured", "env_var", envVar)
	return dsn, nil
}

func sanitizeDSN(dsn string) string {
	passwordRe := regexp.MustCompile(`password=[^&\s]+`)
	return passwordRe.ReplaceAllString(dsn, "password=***")
}
