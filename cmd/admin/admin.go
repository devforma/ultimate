package main

import (
	"fmt"
	"log"
	"os"

	"github.com/devforma/ultimate/internal/config"
	"github.com/devforma/ultimate/internal/database"
)

const help = `Commands You Can Use
================================================================
[init] connect to mysql base on the config.yaml and seeding data 
`

const seedSql = `
INSERT INTO canal
(title, duration)
VALUES
("达瓦大无", 131),
("dc", 11),
("我晚点哇", 89),
("房管局让你a", 9)
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(help)
		return
	}

	switch os.Args[1] {
	case "dbinit":
		if err := InitDB(); err != nil {
			log.Fatalf("init db failed: %v", err)
		}

		fmt.Println("seeding finished!")
	}
}

// InitDB init the mysql database
func InitDB() error {
	cfg, err := getConfig("config.yaml")
	if err != nil {
		return err
	}

	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Seed(seedSql); err != nil {
		return err
	}

	return nil
}

// getConfig parse config file and return dbconfig
func getConfig(path string) (*database.Config, error) {
	var cfg database.Config
	if err := config.Parse(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
