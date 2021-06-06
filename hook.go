// File: main.go
package main

import (
    "fmt"
    log "github.com/sirupsen/logrus"
    kafka "github.com/segmentio/kafka-go"
    "context"
    "flag"
    "os"
    "net/http"
    "io/ioutil"
    //"encoding/json"
)

/* Globals variables */

var (
    topicName = flag.String("topic", "default-topic", "topic name")
    brokerList = flag.String("brokers", "localhost:9092", "bootstrap URL")
    listenAddr = flag.String("listen", "localhost:4000", "ip:port to bind service")
)

func init() {
    flag.Parse()
    log.SetFormatter(&log.JSONFormatter{})
    log.SetOutput(os.Stdout)
    log.SetLevel(log.InfoLevel)
}

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func formatRecord(byteData []byte, key string) kafka.Message {
    msg := kafka.Message{
    	Key:   []byte(key),
    	Value: byteData,
    }
    return msg
}

func sendRecord(msg kafka.Message, kafkaWriter *kafka.Writer) {
    err := kafkaWriter.WriteMessages(context.Background(), msg)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("message produced")
	}
}

func main() {
    kafkaWriter := getKafkaWriter(*brokerList, *topicName)
    log.Info(fmt.Sprint("Producer created for ", *brokerList, *topicName))

    defer kafkaWriter.Close()

    // sync endpoint
    mux := http.NewServeMux()
    mux.HandleFunc("/sync", func(w http.ResponseWriter, r *http.Request) {
        byteData, err := ioutil.ReadAll(r.Body)
        if err != nil {
            fmt.Println(err)
        }
        sendRecord(formatRecord(byteData, fmt.Sprintf("address-%s", r.RemoteAddr)), kafkaWriter)
    })

    // start webserver
    log.Info(fmt.Sprint("Webserver starting on ", *listenAddr))
    ws_err := http.ListenAndServe(*listenAddr, mux)

    log.Fatal(ws_err)
}
