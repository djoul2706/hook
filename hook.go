// File: main.go
package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    //"encoding/json"
    "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

/* Globals variables */
var topic_name string = "test-from-ws"
var kafka_config *kafka.ConfigMap = &kafka.ConfigMap{"bootstrap.servers": "localhost:9092"}

/*
func logDoc(w http.ResponseWriter, r *http.Request, producerChan chan *kafka.Message) {
    var jsonData map[string]interface{}

    byteData, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Println(err)
    }
    if err := json.Unmarshal(byteData, &jsonData); err != nil {
        log.Println(err)
        log.Println("byteData")
        fmt.Println(string(byteData))
    } else {
        log.Println("jsonData")
        fmt.Println(jsonData)
    }
    go produce_record(producerChan)
}
*/

func formatRecord(byteData []byte) *kafka.Message {
    record := &kafka.Message{
        	TopicPartition: kafka.TopicPartition{
        		Topic:     &topic_name,
        		Partition: kafka.PartitionAny,
        	},
        	Value: byteData,
        }
    return record
}

func syncProduce(m *kafka.Message, p *kafka.Producer, ch chan kafka.Event) {
   	p.Produce(m, ch)
	e := <-ch
	s := e.(*kafka.Message)
    if s.TopicPartition.Error != nil {
    	fmt.Printf("Delivery failed: %v\n", s.TopicPartition.Error)
    } else {
    	fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
    		*s.TopicPartition.Topic, s.TopicPartition.Partition, s.TopicPartition.Offset)
    }
}

func produce_record(byteData []byte, producerChan chan *kafka.Message)  {
    // create record
	record := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic_name,
			Partition: kafka.PartitionAny,
		},
		Value: byteData,
	}
	producerChan <- record
}

func log_errors(eventChan chan kafka.Event) {
    for e := range eventChan {
    	switch ev := e.(type) {
    	case *kafka.Message:
    		m := ev
    		if m.TopicPartition.Error != nil {
    			fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
    		} else {
    			fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
    				*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
    		}
    		return
    	default:
    		fmt.Printf("Ignored event: %s\n", ev)
    	}
    }
}

func main() {

    // init kafka producer
    producer, err := kafka.NewProducer(kafka_config)
    if err != nil {
    	log.Fatalf("Error creating producer: %s\n", err)
    }
    fmt.Printf("Created Producer %v\n", producer)
    producerChan := producer.ProduceChannel()
    eventChan := producer.Events()
    synchChan := producer.Events()

    // async endpoint
    mux := http.NewServeMux()
    mux.HandleFunc("/async", func(w http.ResponseWriter, r *http.Request) {
        byteData, err := ioutil.ReadAll(r.Body)
        if err != nil {
            fmt.Println(err)
        }
        go produce_record(byteData, producerChan)
        go log_errors(eventChan)
    })

    mux.HandleFunc("/sync", func(w http.ResponseWriter, r *http.Request) {
        byteData, err := ioutil.ReadAll(r.Body)
        if err != nil {
            fmt.Println(err)
        }
        syncProduce(formatRecord(byteData), producer, synchChan)
    })

    // start webserver
    log.Println("Starting server on :4000...")
    ws_err := http.ListenAndServe(":4000", mux)
    log.Fatal(ws_err)
}
