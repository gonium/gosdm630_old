package sdm630

import (
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

type MQTTSubmitter struct {
	mqtt       *MQTT.Client
	datastream ReadingChannel
	control    ControlChannel
}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client *MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func NewMQTTSubmitter(ds ReadingChannel, cc ControlChannel) (*MQTTSubmitter, error) {
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("SDM630")
	opts.SetDefaultPublishHandler(f)
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	} else {
		return &MQTTSubmitter{mqtt: c, datastream: ds, control: cc}, nil
	}
}

func (ms *MQTTSubmitter) ConsumeData() {
	for {
		// TODO: Read on control, terminate goroutine when
		readings := <-ms.datastream
		payload := fmt.Sprintf("%s", &readings)
		channel := "SDM630/foo"
		// TODO: Fan out and publish individual sensor readings on separate
		// channels.
		token := ms.mqtt.Publish(channel, 0, false, payload)
		token.Wait()
		if token.Error() != nil {
			fmt.Printf("Error: >%s< while submitting %s\r\n", token.Error().Error(), payload)
		}
	}
	ms.mqtt.Disconnect(250)
}
