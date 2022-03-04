package client_test

import (
	client "aurora/client/go"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	loginURl        = "http://localhost/auth"
	tasksUrl        = "http://localhost/tasks/send"
	touchUrl        = "http://localhost/tasks/touch"
	username        = "admin"
	password        = "password"
	timeoutDuration = 100 + rand.Intn(4900)
	sleepDuration   = 50 + rand.Intn(450)
)

func TestSendSyncWithTask(t *testing.T) {
	t.Parallel()
	initTasks()
	connector := client.NewAuroraConnector(loginURl, tasksUrl, touchUrl)
	defer connector.Close()
	connector.Init(username, password)
	requestOBJ := &client.CenterRequest{
		TaskType:  "task",
		Timestamp: time.Now().Local().Unix(),
		Signatures: []*client.Signature{
			&addTask0,
		},
		TimeoutDuration: timeoutDuration,
		SleepDuration:   sleepDuration,
	}
	responseOBJ, err := connector.SendSync(requestOBJ)
	assert.NoError(t, err)
	b, err := json.Marshal(responseOBJ)
	assert.NoError(t, err)
	// (1+1)=?
	t.Logf("responseOBJ: %s", string(b))
}

func TestSendSyncWithGroup(t *testing.T) {
	t.Parallel()
	initTasks()
	connector := client.NewAuroraConnector(loginURl, tasksUrl, touchUrl)
	defer connector.Close()
	connector.Init(username, password)
	requestOBJ := &client.CenterRequest{
		TaskType:  "group",
		Timestamp: time.Now().Local().Unix(),
		Signatures: []*client.Signature{
			&addTask0,
			&addTask1,
			&addTask2,
		},
		TimeoutDuration: timeoutDuration,
		SleepDuration:   sleepDuration,
		SendConcurrency: 3, // max counts per send subtask
	}
	responseOBJ, err := connector.SendSync(requestOBJ)
	assert.NoError(t, err)
	b, err := json.Marshal(responseOBJ)
	assert.NoError(t, err)
	// (1+1)=? (2+2)=? (5+6)=?
	t.Logf("responseOBJ: %s", string(b))
}

func TestSendSyncWithChain(t *testing.T) {
	t.Parallel()
	initTasks()
	connector := client.NewAuroraConnector(loginURl, tasksUrl, touchUrl)
	defer connector.Close()
	connector.Init(username, password)
	requestOBJ := &client.CenterRequest{
		TaskType:  "chain",
		Timestamp: time.Now().Local().Unix(),
		Signatures: []*client.Signature{
			&addTask0,
			&addTask1,
			&addTask2,
			&multiplyTask0,
		},
		TimeoutDuration: timeoutDuration,
		SleepDuration:   sleepDuration,
	}
	responseOBJ, err := connector.SendSync(requestOBJ)
	assert.NoError(t, err)
	b, err := json.Marshal(responseOBJ)
	assert.NoError(t, err)
	// ((((1 + 1) + (2 + 2)) + (5 + 6)) * 4) = ?
	t.Logf("responseOBJ: %s", string(b))
}

func TestSendSyncWithChord(t *testing.T) {
	t.Parallel()
	initTasks()
	connector := client.NewAuroraConnector(loginURl, tasksUrl, touchUrl)
	defer connector.Close()
	connector.Init(username, password)
	requestOBJ := &client.CenterRequest{
		TaskType:  "chord",
		Timestamp: time.Now().Local().Unix(),
		Signatures: []*client.Signature{
			&addTask0,
			&addTask1,
		},
		TimeoutDuration: timeoutDuration,
		SleepDuration:   sleepDuration,
		SendConcurrency: 2,
		CallBack:        &multiplyTask1,
	}
	responseOBJ, err := connector.SendSync(requestOBJ)
	assert.NoError(t, err)
	b, err := json.Marshal(responseOBJ)
	assert.NoError(t, err)
	t.Logf("responseOBJ: %s", string(b))
}
