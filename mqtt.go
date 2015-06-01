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

func (ms *MQTTSubmitter) submitReading(basechannel string,
	subchannel string, reading float32) {
	payload := fmt.Sprintf("%f", reading)
	channel := fmt.Sprintf("%s/%s", basechannel, subchannel)
	token := ms.mqtt.Publish(channel, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		fmt.Printf("Error: >%s< while submitting %s\r\n", token.Error().Error(), payload)
	}
}

func (ms *MQTTSubmitter) ConsumeData() {
	for {
		// TODO: Read on control, terminate goroutine when
		readings := <-ms.datastream
		basechannel := "SDM630/foo"
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
	ms.mqtt.Disconnect(250)
}
