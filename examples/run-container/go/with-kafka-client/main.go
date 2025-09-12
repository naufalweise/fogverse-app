package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {

	bootstrapServers := os.Getenv("BOOTSTRAP_SERVERS")
	if len(bootstrapServers) == 0 {
		fmt.Fprintf(os.Stderr, "Environment variable BOOTSTRAP_SERVERS must be set\n")
		os.Exit(1)
	}

	inputTopic := os.Getenv("INPUT_TOPIC")
	if inputTopic == "" {
		fmt.Fprintf(os.Stderr, "Please set INPUT_TOPIC environment variable\n")
		os.Exit(1)
	}
	outputTopic := os.Getenv("OUTPUT_TOPIC")
	if outputTopic == "" {
		fmt.Fprintf(os.Stderr, "Please set OUTPUT_TOPIC environment variable\n")
		os.Exit(1)
	}

	group := os.Getenv("GROUP")
	if group == "" {
		group = "default-group"
		fmt.Println("GROUP environment variable is not set. Using default group: ", group)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  bootstrapServers,
		"group.id":           group,
		"session.timeout.ms": 6000,
		// Start reading from the first message of each assigned
		// partition if there are no previously committed offsets
		// for this group.
		"auto.offset.reset": "earliest",
		// Whether we store offsets automatically.
		"enable.auto.offset.store": false,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	err = consumer.SubscribeTopics([]string{inputTopic}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to subscribe to topic: %s\n", err)
		os.Exit(1)
	}

	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServers})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", producer)

	// Listen to all the events on the default events channel
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				// The message delivery report, indicating success or
				// permanent failure after retries have been exhausted.
				// Application level retries won't help since the client
				// is already configured to do that.
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
						*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
			case kafka.Error:
				// Generic client instance-level errors, such as
				// broker connection failures, authentication issues, etc.
				//
				// These errors should generally be considered informational
				// as the underlying client will automatically try to
				// recover from any errors encountered, the application
				// does not need to take action on them.
				fmt.Printf("Error: %v\n", ev)
			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				// Process the message received.
				output := processMsg(e.Value)
				writeToTopic(producer, outputTopic, output)

				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
			}
		}
	}

	// Flush and close the producer and the events channel
	for producer.Flush(10000) > 0 {
		fmt.Print("Still waiting to flush outstanding messages\n")
	}
	producer.Close()
}

// processMsg does operations to the message and returns the operations result
// the operations results, in raw byte array, will be written to the output topic
// in this case, the operation is summing the numbers contained in the message
func processMsg(msg []byte) []byte {
	msgStr := string(msg)
	// convert output to byte array
	out := fmt.Sprintf("Msg %s, Sum: %d\n", msgStr, sumNumbers(msgStr))
	return []byte(out)
}

// sumNumbers returns the sums of numbers contained in a string
func sumNumbers(line string) int {
	sum := 0
	numberStrList := strings.Split(line, " ")
	for _, numberStr := range numberStrList {
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			continue
		}
		sum += number
	}
	return sum
}

func writeToTopic(p *kafka.Producer, topic string, msg []byte) {
	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
		Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, nil)

	if err != nil {
		fmt.Printf("Failed to produce message: %v\n", err)
	}
}
