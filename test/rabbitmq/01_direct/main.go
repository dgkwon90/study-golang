package main

import (
	"direct/consumer"
	"direct/publisher"
	"fmt"
	"strconv"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const RabbitMqUrl = "amqp://dgkwon:test001@192.168.56.1:5672/"
const ExchangeName = "direct_test_exchange"

var consumerMsgs map[string]string
var mutex = &sync.Mutex{}

func reviceMsgHandler(name string, msg interface{}) {
	reviceMsg := msg.(amqp.Delivery)
	msgNum := reviceMsg.Headers["Msg"].(int32)
	mutex.Lock()
	if val, ok := consumerMsgs[name]; ok {
		consumerMsgs[name] = val + ", Msg" + strconv.Itoa(int(msgNum))
	} else {
		consumerMsgs[name] = "Msg" + strconv.Itoa(int(msgNum))
	}
	mutex.Unlock()
}

func StartConsumers() {
	fmt.Println("\n\nStart Consumers!!!!!!!!!!!!!!!!!!!!!!!")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		con1 := consumer.NewCon(RabbitMqUrl, "consumber:1", ExchangeName, "ucl", "user.created")
		defer con1.Close()
		con1.Connection()
		con1.OpenChannel()
		con1.Bind(reviceMsgHandler)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		con2 := consumer.NewCon(RabbitMqUrl, "consumber:2", ExchangeName, "uul", "user.updated")
		defer con2.Close()
		con2.Connection()
		con2.OpenChannel()
		con2.Bind(reviceMsgHandler)
	}()

	wg.Add(1)
	go func() {
		con3 := consumer.NewCon(RabbitMqUrl, "consumber:3", ExchangeName, "ucl.two", "user.created")
		defer con3.Close()
		con3.Connection()
		con3.OpenChannel()
		con3.Bind(reviceMsgHandler)
	}()

	wg.Add(1)
	go func() {
		con4 := consumer.NewCon(RabbitMqUrl, "consumber:4", ExchangeName, "ucl", "user.created")
		defer con4.Close()
		con4.Connection()
		con4.OpenChannel()
		con4.Bind(reviceMsgHandler)
	}()
	wg.Wait()
}

func StartPublisher() {
	fmt.Println("\n\nStart Publisher!!!!!!!!!!!!!!!!!!!!!!!")
	pub := publisher.NewPub(RabbitMqUrl, "publisher:1")
	defer pub.Close()
	pub.Connection()
	pub.OpenChannel()
	pub.Publish(
		ExchangeName,
		"user.created",
		map[string]interface{}{
			"Msg": 1,
		},
		[]byte("{\"username\":\"sysed\"}"),
	)
	pub.Publish(
		ExchangeName,
		"user.created",
		map[string]interface{}{
			"Msg": 2,
		},
		[]byte("{\"username\":\"sirajul\"}"),
	)
	pub.Publish(
		ExchangeName,
		"user.created",
		map[string]interface{}{
			"Msg": 3,
		},
		[]byte("{\"username\":\"islam\"}"),
	)
	pub.Publish(
		ExchangeName,
		"user.updated",
		map[string]interface{}{
			"Msg": 4,
		},
		[]byte("{\"username\":\"anik\", \"old\":\"syed\"}"),
	)
	pub.Publish(
		"",
		"ucl.two",
		map[string]interface{}{
			"Msg": 5,
		},
		[]byte("{\"username\":\"islam\"}"),
	)
}

func main() {
	consumerMsgs = make(map[string]string)
	go StartConsumers()
	time.Sleep(time.Second * 3)
	StartPublisher()

	fmt.Println("====== [result] ======")
	for con, msg := range consumerMsgs {
		fmt.Printf("%v: %v\n", con, msg)
	}
}
