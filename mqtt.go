package sdm630

import (
	"fmt"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
	"log"
)

type MQTTSubmitter struct {
	mqtt       *client.Client
	devicename string
	datastream ReadingChannel
	control    ControlChannel
}

func NewMQTTSubmitter(ds ReadingChannel, cc ControlChannel,
	brokerurl string, username string, password string, devicename string) (*MQTTSubmitter, error) {
	mqttclient := client.New(&client.Options{
		ErrorHandler: func(err error) {
			log.Printf("MQTT error occured: %s\n", err)
		},
	})
	err := mqttclient.Connect(&client.ConnectOptions{
		Network:      "tcp",
		Address:      brokerurl,
		ClientID:     []byte("gosdm360_submitter:" + devicename),
		CleanSession: true,
		KeepAlive:    30,
		WillQoS:      mqtt.QoS0,
	})
	if err != nil {
		return nil, err
	} else {
		log.Println("Connected to broker.")
		return &MQTTSubmitter{mqtt: mqttclient, devicename: devicename, datastream: ds, control: cc}, nil
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
	}
}

func (ms *MQTTSubmitter) ConsumeData() {
	basechannel := fmt.Sprintf("%s/readings", ms.devicename)
	for {
		// TODO: Read on control, terminate goroutine when
		readings := <-ms.datastream
		ms.submitReading(basechannel, "L1/Voltage", readings.L1Voltage)
		ms.submitReading(basechannel, "L2/Voltage", readings.L2Voltage)
		ms.submitReading(basechannel, "L3/Voltage", readings.L3Voltage)
		ms.submitReading(basechannel, "L1/Current", readings.L1Current)
		ms.submitReading(basechannel, "L2/Current", readings.L2Current)
		ms.submitReading(basechannel, "L3/Current", readings.L3Current)
		ms.submitReading(basechannel, "L1/Power", readings.L1Power)
		ms.submitReading(basechannel, "L2/Power", readings.L2Power)
		ms.submitReading(basechannel, "L3/Power", readings.L3Power)
		ms.submitReading(basechannel, "L1/CosPhi", readings.L1CosPhi)
		ms.submitReading(basechannel, "L2/CosPhi", readings.L2CosPhi)
		ms.submitReading(basechannel, "L3/CosPhi", readings.L3CosPhi)

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
		},
	})
	err := mqttclient.Connect(&client.ConnectOptions{
		Network:      "tcp",
		Address:      brokerurl,
		ClientID:     []byte("gosdm360_receiver:" + devicename),
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

// TODO: Inject handler as a function parameter
func (mq *MQTTSource) Subscribe(topic string) error {
	topicfilter := fmt.Sprintf("%s/#", topic)
	log.Println("Subscribing to", topicfilter)
	return mq.mqtt.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				TopicFilter: []byte(topicfilter),
				QoS:         mqtt.QoS0,
				// Define the processing of the message handler.
				Handler: func(topicName, message []byte) {
					log.Println(string(topicName), string(message))
				},
			},
		},
	})
}
