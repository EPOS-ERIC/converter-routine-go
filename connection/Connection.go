package connection

import (
	"context"
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
	log = logging.Get("database")
	dnsRegex  = regexp.MustCompile(`(&?(targetServerType|loadBalanceHosts)=[^&]+)`)
)

// dbPools is a map from an environment variable (that holds a DSN)
// to a slice of *gorm.DB connections. This allows you to have multiple
// distinct DSNs/connection sets under different environment variables
var dbPools = make(map[string][]*gorm.DB)

// Protect dbPools with a mutex if multiple goroutines might race to init
var mu sync.Mutex

// ConnectConverter is a thin wrapper that uses the manager to fetch
// a connection for CONVERTER_CATALOGUE_CONNECTION_STRING
func ConnectConverter() (*gorm.DB, error) {
	return connectManager("CONVERTER_CATALOGUE_CONNECTION_STRING")
}

// connectManager checks if we have a pool of *gorm.DB for the given
// environment variable. If not, it initializes it, then returns a *gorm.DB
func connectManager(envVar string) (*gorm.DB, error) {
	mu.Lock()
	defer mu.Unlock()

	log.Debug("requesting database connection", "env_var", envVar)

	// check if we already have a pool
	if _, exists := dbPools[envVar]; !exists || len(dbPools[envVar]) == 0 {
		log.Info("initializing new database pool", "env_var", envVar)

		// initialize a new pool of connections
		err := initializePool(envVar)
		if err != nil {
			log.Error("failed to initialize database pool", "env_var", envVar, "error", err)
			return nil, fmt.Errorf("initialization error: %w", err)
		}

		// check if that succeeded in creating any connections
		if len(dbPools[envVar]) == 0 {
			log.Error("no database connections available after initialization", "env_var", envVar)
			return nil, fmt.Errorf("no database connections available for %s", envVar)
		}

		log.Info("database pool initialized successfully", "env_var", envVar, "connection_count", len(dbPools[envVar]))
	}

	// At this point, dbPools[envVar] should have at least 1 *gorm.DB
	// Try each one in turn and return the first that is reachable
	for i, db := range dbPools[envVar] {
		sqlDB, err := db.DB()
		if err != nil {
			log.Error("failed to get underlying sql.DB", "env_var", envVar, "connection_index", i, "error", err)
			continue
		}

		// Use a 2 sec timeout for the ping
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		err = sqlDB.PingContext(ctx)
		cancel()

		if err != nil {
			log.Warn("database ping failed, trying next connection", "env_var", envVar, "connection_index", i, "error", err)
			continue
		}

		log.Debug("database connection established successfully", "env_var", envVar, "connection_index", i)

		// Return the first connection that works
		return db, nil
	}

	log.Error("all database hosts unreachable", "env_var", envVar, "total_connections", len(dbPools[envVar]))

	return nil, fmt.Errorf("all database hosts are unreachable for %s", envVar)
}

// initializePool reads the DSN from envVar, parses out the hosts, sets up
// multiple connections (one per host) and stores them in dbPools[envVar]
func initializePool(envVar string) error {
	hosts, params, err := parseMultiHostDSN(envVar)
	if err != nil {
		return fmt.Errorf("failed to parse DSN for %s: %w", envVar, err)
	}

	log.Info("initializing database connections", "env_var", envVar, "host_count", len(hosts), "hosts", hosts)

	// Make a slice to hold the *gorm.DB for each host
	newDbs := make([]*gorm.DB, 0, len(hosts))

	for i, host := range hosts {
		currentDSN := fmt.Sprintf("postgresql://%s/%s", host, params)

		log.Debug("attempting to connect to database host", "env_var", envVar, "host_index", i, "host", host)

		db, err := gorm.Open(postgres.New(postgres.Config{
			DriverName: "pgx",
			DSN:        currentDSN,
		}), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "",
				SingularTable: true,
			},
			DisableAutomaticPing: true, // We manually do the ping
		})
		if err != nil {
			log.Error("failed to connect to database host", "env_var", envVar, "host", host, "error", err)
			continue
		}

		log.Info("successfully connected to database host", "env_var", envVar, "host", host)

		newDbs = append(newDbs, db)
	}

	if len(newDbs) == 0 {
		return fmt.Errorf("failed to initialize any DB connections for %s", envVar)
	}

	// Store in the global map
	dbPools[envVar] = newDbs

	log.Info("database pool initialization completed", "env_var", envVar, "successful_connections", len(newDbs), "total_hosts", len(hosts))

	return nil
}

// parseMultiHostDSN fetches the DSN from the given envVar and
// splits it into (hosts, params)
func parseMultiHostDSN(envVar string) ([]string, string, error) {
	dsn, ok := os.LookupEnv(envVar)
	if !ok {
		log.Error("environment variable not set", "env_var", envVar)
		return nil, "", fmt.Errorf("%s is not set", envVar)
	}

	log.Debug("parsing DSN from environment variable", "env_var", envVar, "original_dsn", dsn)

	// Remove "jdbc:" prefix if present
	if strings.HasPrefix(dsn, "jdbc:") {
		dsn = strings.Replace(dsn, "jdbc:", "", 1)
		log.Debug("removed jdbc prefix from DSN", "cleaned_dsn", dsn)
	}

	// Remove unsupported parameters like targetServerType & loadBalanceHosts
	if dnsRegex.MatchString(dsn) {
		dsn = dnsRegex.ReplaceAllString(dsn, "")
		log.Debug("removed unsupported parameters from DSN", "cleaned_dsn", dsn)
	}

	// Clean up trailing "?" or "&"
	trailingRe := regexp.MustCompile(`[?&]$`)
	if trailingRe.MatchString(dsn) {
		dsn = trailingRe.ReplaceAllString(dsn, "")
		log.Debug("removed trailing characters from DSN", "cleaned_dsn", dsn)
	}

	// Must contain "//"
	hostStart := strings.Index(dsn, "//")
	if hostStart == -1 {
		log.Error("invalid DSN format: missing '//'", "env_var", envVar, "dsn", dsn)
		return nil, "", fmt.Errorf("invalid connection string format: missing '//'")
	}

	// Extract everything after `//` (hosts and params)
	hostsAndParams := dsn[hostStart+2:]
	splitIndex := strings.Index(hostsAndParams, "/")
	if splitIndex == -1 {
		log.Error("invalid DSN format: missing '/' after hosts", "env_var", envVar, "hosts_and_params", hostsAndParams)
		return nil, "", fmt.Errorf("invalid connection string format: missing '/' after hosts")
	}

	hosts := hostsAndParams[:splitIndex]
	params := hostsAndParams[splitIndex+1:]

	hostList := strings.Split(hosts, ",")

	log.Debug("DSN parsing completed successfully", "env_var", envVar, "hosts", hostList, "params", params)

	return hostList, params, nil
}
