package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Create database tables
	CreateTable()

	// Insert Record
	InsertWallet(phoneNumber, pin, publicKey, keystorePath string) error

	// Select Wallet by phone number
	SelectWalletByPhone(phoneNumber string) (*WalletRecord, error)

	// Update keystore path
	UpdateKeystorePathByID(path string, id uint64)
	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type service struct {
	db *sqlx.DB
}

var (
	dburl      = os.Getenv("BLUEPRINT_DB_URL")
	dbInstance *service
)

func (s *service) CreateTable() {
	schema := `
	CREATE TABLE IF NOT EXISTS wallets (
		id integer primary key autoincrement,
		phone_number text not null unique,
		public_key text not null unique,
		pin text not null,
		keystore_path text not null
	);`
	s.db.MustExec(schema)
}

func (s *service) InsertWallet(phoneNumber, pin, publicKey, keystorePath string) error {
	stmt := `
	insert into wallets (
		phone_number,
		pin,
		public_key,
		keystore_path
	) values (
		?, ?, ?,?
	);
	`
	_, err := s.db.Exec(stmt, phoneNumber, pin, publicKey, keystorePath)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SelectWalletByPhone(phoneNumber string) (*WalletRecord, error) {
	stmt := `
	select *
	from wallets
	where phone_number = ?;
	`

	record := WalletRecord{}

	err := s.db.Get(&record, stmt, phoneNumber)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (s *service) UpdateKeystorePathByID(path string, id uint64) {
	stmt := `
	update wallets 
	set keystore_path = ?
	where id = ?;
	`

	s.db.MustExec(stmt, path, id)
}

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	db, err := sqlx.Open("sqlite3", dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dburl)
	return s.db.Close()
}
