package sdm630

import (
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"log"
)

type MQTTSubmitter struct {
	mqtt       *MQTT.Client
	devicename string
	datastream ReadingChannel
	control    ControlChannel
}

//define a function for the default message handler
var f MQTT.MessageHandler = func(client *MQTT.Client, msg MQTT.Message) {
	log.Printf("TOPIC: %s - MSG:%s\r\n", msg.Topic(), msg.Payload())
}

//define a function for the connection lost handler
var l MQTT.ConnectionLostHandler = func(client *MQTT.Client, err error) {
	log.Printf("Lost broker connection: %s\r\n", err.Error())
}

func NewMQTTSubmitter(ds ReadingChannel, cc ControlChannel,
	brokerurl string, username string, password string, devicename string) (*MQTTSubmitter, error) {
	opts := MQTT.NewClientOptions().AddBroker(brokerurl)
	opts.SetClientID("gosdm360")
	opts.SetDefaultPublishHandler(f)
	opts.SetConnectionLostHandler(l)
	opts.SetPassword(password)
	opts.SetUsername(username)
	opts.SetAutoReconnect(false)
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	} else {
		return &MQTTSubmitter{mqtt: c, devicename: devicename, datastream: ds, control: cc}, nil
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
	ms.mqtt.Disconnect(250)
}
