package main

import (
	"flag"
	"log"

	"thomas.vn/hr_recruitment/internal/config"
	"thomas.vn/hr_recruitment/internal/migration"
	xmigration "thomas.vn/hr_recruitment/pkg/migration"
)

var (
	// Name is the name of the application
	Name = "app"
	// Env is the environment of the application
	Env = "dev"
	// Command is the migration command
	Command = ""
	// Steps is the number of migration steps
	Steps = 0
	// DB is the database driver
	DB = ""
)

const (
	MySQL         string = "mysql"
	MongoDB       string = "mongo"
	Elasticsearch string = "es"

	UP      string = "up"
	DOWN    string = "down"
	FORCE   string = "force"
	VERSION string = "version"
)

func init() {
	flag.StringVar(&Name, "name", Name, "Name")
	flag.StringVar(&Env, "env", Env, "Environment")
	flag.StringVar(&Command, "command", Command, "migration command (up/down/version)")
	flag.IntVar(&Steps, "steps", Steps, "number of migration steps")
	flag.StringVar(&DB, "db", DB, "database driver (mysql/mongo/es)")
}

func main() {
	flag.Parse()

	// Validate inputs
	if DB != MySQL && DB != MongoDB && DB != Elasticsearch {
		log.Fatalf("Invalid database driver: %s", DB)
	}
	if Command != UP && Command != DOWN && Command != FORCE && Command != VERSION {
		log.Fatalf("Invalid migration command: %s", Command)
	}

	// Load configuration for the specified environment
	cfg, err := config.LoadConfig(config.Environment(Env), "./config")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Execute migrations based on database type
	switch DB {
	case MySQL:
		executeMySQLMigrations(cfg, Command, Steps)
	case Elasticsearch:
		executeESMigrations(cfg, Command, Steps)
	}
}

func executeMySQLMigrations(cfg *config.Config, command string, steps int) {
	client, err := cfg.InitMySQLDB()
	if err != nil {
		log.Fatalf("Failed to initialize MySQL client: %v", err)
	}

	migrator := migration.NewMySQLMigrator(Name, client.DB)

	if err := migrator.Init(); err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}

	executeCommand(command, steps, migrator)
}

func executeESMigrations(cfg *config.Config, command string, steps int) {
	client, err := cfg.InitElasticSearch()
	if err != nil {
		log.Fatalf("Failed to initialize Elasticsearch client: %v", err)
	}

	migrator := migration.NewESMigrator(Name, client.Client)

	if err := migrator.Init(); err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}

	executeCommand(command, steps, migrator)
}

func executeCommand(command string, steps int, executor xmigration.Migrator) {
	switch command {
	case VERSION:
		version, err := executor.Version()
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		log.Printf("Current version: %d", version)
	case UP:
		if err := executor.Up(steps); err != nil {
			log.Fatalf("Failed to migrate up: %v", err)
		}
		log.Printf("Migration up successful")
	case DOWN:
		if err := executor.Down(steps); err != nil {
			log.Fatalf("Failed to migrate down: %v", err)
		}
		log.Printf("Migration down successful")
	case FORCE:
		if err := executor.Force(steps); err != nil {
			log.Fatalf("Failed to force migrate: %v", err)
		}
		log.Printf("Force migration successful")
	default:
		log.Fatalf("Invalid migration command: %s", command)
	}
}
