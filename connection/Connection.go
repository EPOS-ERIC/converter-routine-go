package connection

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var dbs []*gorm.DB

// Connect returns an available database connection from the pool
func Connect() (*gorm.DB, error) {
	if len(dbs) == 0 {
		err := initializeDbs()
		if err != nil {
			return nil, fmt.Errorf("initialization error: %w", err)
		}
		if len(dbs) == 0 {
			return nil, fmt.Errorf("no database connections available")
		}
	}

	// try each connection in order
	for _, db := range dbs {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error getting underlying *sql.DB: %v", err)
			continue
		}

		// Use a 2 sec timeout for the ping
		ctx, ctxCancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
		// Check connectivity
		err = sqlDB.PingContext(ctx)
		// Cancel immediately after use
		ctxCancelFunc()

		if err != nil {
			// log.Printf("Failed to ping database: %v", err)
			continue
		}

		return db, nil
	}

	return nil, fmt.Errorf("all database hosts are unreachable")
}

// initialize the dbs from the hosts
func initializeDbs() error {
	hosts, params, err := initializeHosts()
	if err != nil {
		return fmt.Errorf("failed to initialize hosts: %w", err)
	}

	// GORM logger
	logConfig := logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Error,
		IgnoreRecordNotFoundError: false,
	}

	// clear any existing connections
	dbs = make([]*gorm.DB, 0, len(hosts))

	// create a db for each host
	for _, host := range hosts {
		currentDSN := fmt.Sprintf("postgresql://%s/%s", host, params)

		db, err := gorm.Open(postgres.New(postgres.Config{
			DriverName: "pgx",
			DSN:        currentDSN,
		}), &gorm.Config{
			Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logConfig),
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "",
				SingularTable: true,
			},
			DisableAutomaticPing: true, // Needed because the base timeout might be very high, we manually do the ping anyway
		})
		if err != nil {
			log.Printf("Failed to connect to host %s: %v", host, err)
			continue
		}

		// add to the initialized databases
		dbs = append(dbs, db)
	}

	if len(dbs) == 0 {
		return fmt.Errorf("failed to initialize any database connections")
	}

	return nil
}

func initializeHosts() ([]string, string, error) {
	dsn, ok := os.LookupEnv("POSTGRESQL_CONNECTION_STRING")
	log.Println("POSTGRESQL_CONNECTION_STRING:", dsn)
	if !ok {
		return nil, "", fmt.Errorf("POSTGRESQL_CONNECTION_STRING is not set")
	}

	// Remove the "jdbc:" prefix if it exists
	dsn = strings.Replace(dsn, "jdbc:", "", 1)

	log.Println("Cleaned DSN (jdbc prefix removed):", dsn)

	// Remove unsupported parameters like targetServerType and loadBalanceHosts
	re := regexp.MustCompile(`(&?(targetServerType|loadBalanceHosts)=[^&]+)`)
	dsn = re.ReplaceAllString(dsn, "")

	log.Println("Cleaned DSN (unsupported parameters removed):", dsn)

	// Clean up trailing "?" or "&"
	dsn = regexp.MustCompile(`[?&]$`).ReplaceAllString(dsn, "")

	log.Println("Cleaned DSN (multi-host supported):", dsn)

	// Parse hosts and connection parameters correctly
	hostStart := strings.Index(dsn, "//")
	if hostStart == -1 {
		return nil, "", fmt.Errorf("invalid connection string format: missing '//'")
	}

	// Extract everything after `//` (hosts and parameters)
	hostsAndParams := dsn[hostStart+2:]
	splitIndex := strings.Index(hostsAndParams, "/")
	if splitIndex == -1 {
		return nil, "", fmt.Errorf("invalid connection string format: missing '/' after hosts")
	}

	hosts := hostsAndParams[:splitIndex]
	params := hostsAndParams[splitIndex+1:]

	hostList := strings.Split(hosts, ",")

	log.Printf("Parsed Hosts: %v", hostList)
	log.Printf("Connection Parameters: %s", params)

	return hostList, params, nil
}
