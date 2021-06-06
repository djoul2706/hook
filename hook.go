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
    "time"
    "io/ioutil"
    //"encoding/json"
)

/* Globals variables */

var (
    topicName = flag.String("topic", "default-topic", "topic name")
    brokerList = flag.String("brokers", "localhost:9092", "bootstrap URL")
    listenAddr = flag.String("listen", "localhost:4000", "ip:port to bind service")
    validate = flag.Bool("validate", false, "if set a msg will be sent to topic at start")
)

func init() {
    flag.Parse()
    log.SetFormatter(&log.JSONFormatter{})
    log.SetOutput(os.Stdout)
    log.SetLevel(log.DebugLevel)
}

func getKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		//Balancer: &kafka.LeastBytes{}, default balancer is round-robin distribution
		RequiredAcks: kafka.RequireAll,
		MaxAttempts: 30,    // default 30
		BatchSize:  100,    // default 100
		BatchTimeout:   time.Duration(100)*time.Millisecond, // equals linger.ms, set to 100ms
		Compression:    kafka.Snappy,
	}
}

func formatRecord(byteData []byte, key string) kafka.Message {
    msg := kafka.Message{
    	Key:   []byte(key),
    	Value: byteData,
    }
    return msg
}

func sendRecord(msg kafka.Message, kafkaWriter *kafka.Writer) error {
    err := kafkaWriter.WriteMessages(context.Background(), msg)
	if err != nil {
	    return err
	} else {
		log.Debug("msg sent")
		return nil
	}
}

func main() {
    kafkaWriter := getKafkaWriter(*brokerList, *topicName)
    log.Info(fmt.Sprint("Producer created for ", *brokerList, *topicName))
    defer kafkaWriter.Close()

    if *validate == true {
        err := sendRecord(formatRecord([]byte("random data"), "random key"), kafkaWriter)
        if err != nil {
            log.Fatal(err)
        } else {
            log.Info(fmt.Sprint("Test message correctly sent into "), *topicName)
        }
    }

    // sync endpoint
    mux := http.NewServeMux()
    mux.HandleFunc("/sync", func(w http.ResponseWriter, r *http.Request) {
        byteData, err := ioutil.ReadAll(r.Body)
        if err != nil {
            log.Error(err)
        }
        err = sendRecord(formatRecord(byteData, fmt.Sprintf("address-%s", r.RemoteAddr)), kafkaWriter)
        if err != nil {
            log.Error(err)
            w.WriteHeader(http.StatusServiceUnavailable)
            fmt.Fprintf(w, "ERROR - message not sent to kafka\n")
        } else {
            w.WriteHeader(http.StatusAccepted)
            fmt.Fprintf(w, "OK - message sent\n")
        }
    })

    // start webserver
    log.Info(fmt.Sprint("Webserver starting on ", *listenAddr))
    ws_err := http.ListenAndServe(*listenAddr, mux)
    log.Fatal(ws_err)
}
