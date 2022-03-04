package client

import (
	"time"
)

type AuthRequest struct {
	Name     string `json:"Name" validate:"required,gt=0,lt=15"`
	Password string `json:"Password" validate:"required,gt=0,lt=30"`
}

type AuthResponse struct {
	Message string `json:"Message"`
	Name    string `json:"Name"`
	UUID    string `json:"UUID"`
}

type CenterRequest struct {
	UUID            string       `json:"UUID" validate:"max=50"`    // user id
	User            string       `json:"User" validate:"max=15"`    // user name
	BatchID         string       `json:"BatchID" validate:"max=50"` // unique id for request
	Timestamp       int64        `json:"Timestamp" validate:"required"`
	TaskType        string       `json:"TaskType" validate:"required,oneof='task' 'group' 'chord' 'chain'"`
	Signatures      []*Signature `json:"Signatures" validate:"required,gt=0"`
	TimeoutDuration int          `json:"TimeoutDuration" validate:"required,min=100,max=5000"`
	SleepDuration   int          `json:"SleepDuration" validate:"required,min=50,max=500"`
	SendConcurrency int          `json:"SendConcurrency" validate:"min=0,max=10"`
	CallBack        *Signature   `json:"CallBack" validate:"required_if=TaskType chord"`
}

type CenterResponse struct {
	UUID          string          `json:"UUID"` // user id
	User          string          `json:"User"`
	BatchID       string          `json:"BatchID"` // unique id for request
	Timestamp     int64           `json:"Timestamp"`
	TaskType      string          `json:"TaskType"`
	TaskResponses []*TaskResponse `json:"TaskResponses"`
}

type TaskResponse struct {
	Results    []interface{} `json:"Results"`
	Signatures []*Signature  `json:"Signatures"`
	CallBack   *Signature    `json:"CallBack"`
}

// Signature represents a single task invocation
type Signature struct {
	UUID           string
	Name           string
	RoutingKey     string
	ETA            *time.Time
	GroupUUID      string
	GroupTaskCount int
	Args           []Arg
	Headers        Headers
	Priority       uint8
	Immutable      bool
	RetryCount     int
	RetryTimeout   int
	OnSuccess      []*Signature
	OnError        []*Signature
	ChordCallback  *Signature
	//MessageGroupId for Broker, e.g. SQS
	BrokerMessageGroupId string
	//ReceiptHandle of SQS Message
	SQSReceiptHandle string
	// StopTaskDeletionOnError used with sqs when we want to send failed messages to dlq,
	// and don't want aurora to delete from source queue
	StopTaskDeletionOnError bool
	// IgnoreWhenTaskNotRegistered auto removes the request when there is no handeler available
	// When this is true a task with no handler will be ignored and not placed back in the queue
	IgnoreWhenTaskNotRegistered bool
}

type Arg struct {
	Name  string      `bson:"name"`
	Type  string      `bson:"type"`
	Value interface{} `bson:"value"`
}

// Headers represents the headers which should be used to direct the task
type Headers map[string]interface{}
