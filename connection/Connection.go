package connection

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"regexp"
	"strings"
)

var hosts []string
var params = ""

// get a db to execute queries
func Connect() (*gorm.DB, error) {
	var err error
	if hosts == nil && params == "" {
		hosts, params, err = initializeHosts()
		if err != nil {
			return nil, err
		}
	}

	// Attempt to connect to each host
	for _, host := range hosts {
		currentDSN := fmt.Sprintf("postgresql://%s/%s", host, params)
		// log.Printf("Attempting to connect to: %s", host)
		db, err := gorm.Open(postgres.New(postgres.Config{
			DriverName: "pgx",
			DSN:        currentDSN,
		}), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   "",   // table name prefix
				SingularTable: true, // use singular table names
			},
		})
		if err != nil {
			// log.Printf("Failed to connect to: %s, error: %v", host, err)
			// try next host
			continue
		}
		// log.Printf("Successfully connected to: %s", host)
		return db, nil
	}
	return nil, fmt.Errorf("error enstablishing connection to db: all hosts are unreachable: %w", err)
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

