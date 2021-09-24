package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/devforma/ultimate/internal/util"
	"github.com/streadway/amqp"
)

type MQAPI struct {
	Logger *log.Logger
	Conn *amqp.Connection
}

func (m MQAPI) WorkQueue(w http.ResponseWriter, r *http.Request) {
	sendMessage(m.Conn, m.Logger, r.URL.String() + time.Now().Format("2006-01-02 15:04:05"))
}


// sendMessage knows how to send message to mq
func sendMessage(mqConn *amqp.Connection, logger *log.Logger, message string) {
	if mqConn.IsClosed() {
		logger.Println("mq connection closed")
		return
	}

	channel, err := mqConn.Channel()
	if err != nil {
		logger.Printf("get channel failed: %v", err)
		return
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare("workqueue_1", false, false, false, false, nil)
	if err != nil {
		logger.Printf("declare queue failed: %v", err)
		return
	}

	err = channel.Publish("", queue.Name, false, false, amqp.Publishing{
		Body:        util.StringToBytes(message),
		ContentType: "text/plain",
	})
	if err != nil {
		logger.Printf("publish message failed: %v", err)
	}

	time.Sleep(5*time.Second)
}
