package main

import (
	"fmt"
	"log"
	"os"

	"github.com/devforma/ultimate/internal/config"
	"github.com/devforma/ultimate/internal/database"
	"github.com/streadway/amqp"
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

	case "mqconsume":
		if len(os.Args) != 3 {
			log.Fatalln("no consume target specified")
		}

		mqConsume(os.Args[2])
	}
}

func mqConsume(queueName string) {
	mqConn, err := InitMQConn("amqp://guest:guest@127.0.0.1:5672/")
	if err != nil {
		log.Fatalf("init mq connection failed: %v", err)
	}
	defer mqConn.Close()

	channel, err := mqConn.Channel()
	if err != nil {
		log.Fatalf("get mq channel failed: %v", err)
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatalf("declare queue failed: %v", err)
	}

	msgs, err := channel.Consume(queue.Name, "adminConsume", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("register consumer failed: %v", err)
	}

	for msg := range msgs {
		log.Printf("consume: %s\n", msg.Body)
		msg.Ack(false)
	}
}

// InitMQConn get mq connection
func InitMQConn(mqUrl string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(mqUrl)
	if err != nil {
		return nil, err
	}

	return conn, nil
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
