package mqtt

import (
	"encoding/json"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-basic/uuid"
	"github.com/project-flogo/rules/configuration"
	"sync"
	"time"
)

const (
	modelName                 	= "rule"
	nodeId					  	= "hqqaass-edge"//todo :config
	pinNumberConfig           	= "GPIO-PIN-NUMBER"
	BleGateway				  	= "sys/cloud/8cd4950007da"
	DeviceETPrefix            	= "$hw/events/device/"
	NodeETPrefix			  	= "$hw/events/node/"
	MemberShipETGetSuffix     	= "/membership/get"
	MemberShipETGetResultSuffix = "/membership/get/result"
	MemberShipETUpdateSuffix    = "/membership/updated"
	DeviceETStateUpdateSuffix 	= "/state/update"
	TwinETUpdateSuffix        	= "/twin/update"
	TwinETUpdateDeltaSuffix     = "/twin/update/delta"
	TwinETCloudSyncSuffix     	= "/twin/cloud_updated"
	TwinETGetResultSuffix     	= "/twin/get/result"
	TwinETGetSuffix           	= "/twin/get"
)

var mqttClient mqtt.Client
var TokenClient Token
var wg sync.WaitGroup
var deviceTwinResult DeviceTwinUpdate
var configFile configuration.ReadConfigFile
var pinNumber float64

//Token interface to validate the MQTT connection.
type Token interface {
	Wait() bool
	WaitTimeout(time.Duration) bool
	Error() error
}

//DeviceStateUpdate is the structure used in updating the device state
type DeviceStateUpdate struct {
	State string `json:"state,omitempty"`
}

//BaseMessage the base struct of event message
type BaseMessage struct {
	EventID   string `json:"event_id"`
	Timestamp int64  `json:"timestamp"`
}

//TwinValue the struct of twin value
type TwinValue struct {
	Value    *string        `json:"value,omitempty"`
	Metadata *ValueMetadata `json:"metadata,omitempty"`
}

//ValueMetadata the meta of value
type ValueMetadata struct {
	Timestamp int64 `json:"timestamp,omitempty"`
}

//TypeMetadata the meta of value type
type TypeMetadata struct {
	Type string `json:"type,omitempty"`
}

//TwinVersion twin version
type TwinVersion struct {
	CloudVersion int64 `json:"cloud"`
	EdgeVersion  int64 `json:"edge"`
}


//MsgTwin the struct of device twin
type MsgTwin struct {
	Expected        *TwinValue    `json:"expected,omitempty"`
	Actual          *TwinValue    `json:"actual,omitempty"`
	Optional        *bool         `json:"optional,omitempty"`
	Metadata        *TypeMetadata `json:"metadata,omitempty"`
	ExpectedVersion *TwinVersion  `json:"expected_version,omitempty"`
	ActualVersion   *TwinVersion  `json:"actual_version,omitempty"`
}

//DeviceTwinUpdate the struct of device twin update
type DeviceTwinUpdate struct {
	BaseMessage
	Twin map[string]*MsgTwin `json:"twin"`
}

func Init() error {
	opts := mqtt.NewClientOptions().AddBroker("tcp://172.17.12.23:1883").SetClientID(uuid.New())
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.AutoReconnect = true

	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return errors.New("connect mqtt failed")
	}
	fmt.Printf("create mqtt client\n")

	topics := map[string]byte{
		BleGateway: 0,
		NodeETPrefix+nodeId+MemberShipETGetResultSuffix: 0,
		NodeETPrefix+nodeId+MemberShipETUpdateSuffix: 0,
	}
	if token := mqttClient.SubscribeMultiple(topics, msgHandler); token.Wait() && token.Error() != nil {
		return errors.New("subscribe topic failed")
	}
	fmt.Printf("subscibe topic succeed, topics: %v\n", topics)

	ListDevice()

	return nil
}

var msgHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("msg arrived, topic: %v\n", msg.Topic())
	switch msg.Topic() {
	case NodeETPrefix+nodeId+MemberShipETGetResultSuffix:
		rules := RuleList{}
		err := json.Unmarshal(msg.Payload(), &rules)
		if err != nil {
			fmt.Printf("error in unmarshal rules: %v\n", err)
		}
		for _, rule := range rules.Rules {
			fmt.Printf("rule: %v\n", rule)
		}
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Printf("connect succeed!\n")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("connect lost: %v\n", err)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("reconnect mqtt failed")
	}
	fmt.Printf("reconnect succeed")

	if token := mqttClient.Subscribe(BleGateway, 0, msgHandler); token.Wait() && token.Error() != nil {
		fmt.Printf("resubscribe topic failed")
	}
	fmt.Printf("resubscibe topic: %v", BleGateway)
}

func ListDevice() {
	getList := NodeETPrefix + nodeId + MemberShipETGetSuffix
	type event struct {
		EventId    string 	`json:"event_id"`
	}
	e := event{
		EventId: "",
	}
	eventId, err := json.Marshal(e)
	if err != nil {
		fmt.Printf("marshal eventId error: %v", err)
	}

	fmt.Printf("topic is: %v, body is: %v\n", getList, string(eventId))

	TokenClient = mqttClient.Publish(getList, 0, false, string(eventId))
	if TokenClient.Wait() && TokenClient.Error() != nil {
		fmt.Printf("client.publish() Error in device list get  is: %v ", TokenClient.Error())
	}
}

// OnSubMessageReceived callback function which is called when message is received
func OnSubMessageReceived(client mqtt.Client, message mqtt.Message) {
	err := json.Unmarshal(message.Payload(), &deviceTwinResult)
	if err != nil {
		fmt.Printf("Error in unmarshalling: %v ", err)
	}
}

//createActualUpdateMessage function is used to create the device twin update message
func createActualUpdateMessage(actualValue string) DeviceTwinUpdate {
	var deviceTwinUpdateMessage DeviceTwinUpdate
	//actualMap := map[string]*MsgTwin{powerStatus: {Actual: &TwinValue{Value: &actualValue}, Metadata: &TypeMetadata{Type: "Updated"}}}
	//deviceTwinUpdateMessage.Twin = actualMap
	return deviceTwinUpdateMessage
}