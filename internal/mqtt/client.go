package mqtt

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Callback func(string, []byte)

type Client struct {
	mqttClient mqtt.Client
}

func NewClient(clientId string, host string, port uint16) (*Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", host, port))
	opts.SetClientID(clientId)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		return nil, err
	}

	return &Client{
		mqttClient: client,
	}, nil
}

func (c *Client) Subscribe(topic string, callback Callback) {
	c.mqttClient.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		callback(msg.Topic(), msg.Payload())
	})
}
