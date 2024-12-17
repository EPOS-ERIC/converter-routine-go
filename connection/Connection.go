package connection

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/go-pg/pg/v10"
)

var dbHost1 = ""
var dbHost2 = ""

// get a db to execute queries
func Connect() (*pg.DB, error) {
	if dbHost1 == "" && dbHost2 == "" {
		conn, ok := os.LookupEnv("POSTGRESQL_CONNECTION_STRING")
		if !ok {
			return nil, fmt.Errorf("POSTGRESQL_CONNECTION_STRING is not set")
		}

		// Check if single host or multiple hosts
		singleHostConn, multiHostsConns, err := parsePostgresURL(conn)
		if err != nil {
			return nil, err
		}

		if singleHostConn != "" {
			dbHost1 = singleHostConn
		} else {
			switch len(multiHostsConns) {
			case 0:
				return nil, fmt.Errorf("no hosts found in connection string")
			case 1:
				dbHost1 = multiHostsConns[0]
			default:
				dbHost1 = multiHostsConns[0]
				dbHost2 = multiHostsConns[1]
			}
		}
	}

	// if we only have one host, just connect directly
	if dbHost2 == "" {
		return connect(dbHost1)
	}

	// if we have two or more hosts try the first, if it fails try the second
	db, err := connect(dbHost1)
	if err != nil {
		db, err = connect(dbHost2)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

// try to connect to a db, ping it to see if it fails
func connect(conn string) (*pg.DB, error) {
	opt, err := pg.ParseURL(conn)
	if err != nil {
		return nil, err
	}

	db := pg.Connect(opt)
	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return db, nil
}

// parsePostgresURL handles both single and multi-host URLs.
// If it's a single host URL (no commas), it returns singleHostConn as the original URL and an empty multiHostsConns slice.
// If it's a multi-host URL (commas in the host), it returns singleHostConn as "" and multiHostsConns as the parsed multiple host URLs.
func parsePostgresURL(connStr string) (singleHostConn string, multiHostsConns []string, err error) {
	// if the string contains jdbc: just remove it
	connStr = strings.Replace(connStr, "jdbc:", "", 1)

	parsed, err := url.Parse(connStr)
	if err != nil {
		return "", nil, err
	}

	if strings.Contains(parsed.Host, ",") {
		multiHostsConns, err := parseMultiHostPostgresURL(connStr)
		return "", multiHostsConns, err
	} else {
		return connStr, nil, nil
	}
}

// postgresql://IP1:PORT1,IP2:PORT2/cerif?user=USER&password=PASSWORD&targetServerType=master&loadBalanceHosts=true
// to
// postgresql://${POSTGRESQL_USER}:${POSTGRESQL_PASSWORD}@${PREFIX}${POSTGRESQL_HOST}:${POSTGRESQL_PORT}/${POSTGRES_DB}
// removing the extra parameters
func parseMultiHostPostgresURL(connStr string) ([]string, error) {
	parsed, err := url.Parse(connStr)
	if err != nil {
		return nil, err
	}

	dbPath := parsed.Path
	q := parsed.Query()

	user := q.Get("user")
	password := q.Get("password")

	q.Del("user")
	q.Del("password")

	allHosts := parsed.Host
	hosts := strings.Split(allHosts, ",")
	if len(hosts) == 0 {
		return nil, fmt.Errorf("no hosts found in connection string")
	}

	var results []string
	for _, host := range hosts {
		newURL := &url.URL{
			Scheme: parsed.Scheme,
			User:   url.UserPassword(user, password),
			Host:   host,
			Path:   dbPath,
		}

		results = append(results, newURL.String())
	}

	return results, nil
}

