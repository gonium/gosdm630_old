package sdm630

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
	"log"
	"os"
)

func GenUUID(prefix string) string {
	return fmt.Sprintf("%s-%s", prefix, uuid.NewV4())
}

//////////////////////////////////////////////////////////////////////////////////

type MQTTSubmitter struct {
	mqtt       *client.Client
	topic      string
	datastream ReadingChannel
	control    ControlChannel
}

func NewMQTTSubmitter(ds ReadingChannel, cc ControlChannel,
	brokerurl string, username string, password string, topic string) (*MQTTSubmitter, error) {
	mqttclient := client.New(&client.Options{
		ErrorHandler: func(err error) {
			log.Printf("MQTT error occured: %s\n", err)
		},
	})
	log.Printf("Connecting as user %s to broker %s", username, brokerurl)
	err := mqttclient.Connect(&client.ConnectOptions{
		Network:      "tcp",
		Address:      brokerurl,
		ClientID:     []byte(GenUUID("gosdm360-submitter")),
		UserName:     []byte(username),
		Password:     []byte(password),
		CleanSession: true,
		KeepAlive:    30,
		WillQoS:      mqtt.QoS0,
	})
	if err != nil {
		return nil, err
	} else {
		log.Println("Connected to broker.")
		return &MQTTSubmitter{mqtt: mqttclient, topic: topic, datastream: ds, control: cc}, nil
	}
}

func (ms *MQTTSubmitter) Close() {
	ms.mqtt.Terminate()
}

func (ms *MQTTSubmitter) submitReading(basechannel string,
	subchannel string, reading float32) {
	payload := fmt.Sprintf("%f", reading)
	channel := fmt.Sprintf("%s/%s", basechannel, subchannel)
	err := ms.mqtt.Publish(&client.PublishOptions{
		QoS:       mqtt.QoS0,
		TopicName: []byte(channel),
		Message:   []byte(payload),
	})
	if err != nil {
		log.Printf("Error: >%s< while submitting %s\r\n", err, payload)
		ms.Close()
		os.Exit(1)
	}
}

func (ms *MQTTSubmitter) ConsumeData() {
	for {
		// TODO: Read on control, terminate goroutine when
		readings := <-ms.datastream
		ms.submitReading(ms.topic, "Voltage/L1", readings.L1Voltage)
		ms.submitReading(ms.topic, "Voltage/L2", readings.L2Voltage)
		ms.submitReading(ms.topic, "Voltage/L3", readings.L3Voltage)
		ms.submitReading(ms.topic, "Current/L1", readings.L1Current)
		ms.submitReading(ms.topic, "Current/L2", readings.L2Current)
		ms.submitReading(ms.topic, "Current/L3", readings.L3Current)
		ms.submitReading(ms.topic, "Power/L1", readings.L1Power)
		ms.submitReading(ms.topic, "Power/L2", readings.L2Power)
		ms.submitReading(ms.topic, "Power/L3", readings.L3Power)
		ms.submitReading(ms.topic, "CosPhi/L1", readings.L1CosPhi)
		ms.submitReading(ms.topic, "CosPhi/L2", readings.L2CosPhi)
		ms.submitReading(ms.topic, "CosPhi/L3", readings.L3CosPhi)

	}
	ms.mqtt.Terminate()
}

///////////////////////////////////////////////////////////////////////////

type MQTTSource struct {
	mqtt       *client.Client
	devicename string
	datastream ReadingChannel
	control    ControlChannel
}

func NewMQTTSource(ds ReadingChannel, cc ControlChannel,
	brokerurl string, username string, password string, devicename string) (*MQTTSource, error) {
	mqttclient := client.New(&client.Options{
		ErrorHandler: func(err error) {
			log.Printf("MQTT error occured: %s\n", err)
			os.Exit(1)
		},
	})
	log.Printf("Connecting as user %s to broker %s", username, brokerurl)
	err := mqttclient.Connect(&client.ConnectOptions{
		Network:      "tcp",
		Address:      brokerurl,
		ClientID:     []byte(GenUUID("gosdm360-receiver")),
		UserName:     []byte(username),
		Password:     []byte(password),
		CleanSession: true,
		KeepAlive:    30,
		WillQoS:      mqtt.QoS0,
	})
	if err != nil {
		return nil, err
	} else {
		log.Println("Connected to broker.")
		return &MQTTSource{mqtt: mqttclient, devicename: devicename, datastream: ds, control: cc}, nil
	}
}

func (ms *MQTTSource) Close() {
	ms.mqtt.Terminate()
}

func (mq *MQTTSource) Subscribe(topic string, handler func(topicname, message []byte)) error {
	topicfilter := fmt.Sprintf("%s/#", topic)
	log.Println("Subscribing to", topicfilter)
	return mq.mqtt.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte(topicfilter),
				QoS:         mqtt.QoS0,
				Handler:     handler,
			},
		},
	})
}
