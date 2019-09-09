//+build e2e

package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/satori/go.uuid"

	"github.com/gomodule/redigo/redis"
	"github.com/nsqio/go-nsq"
)

// e2e test connecting to the different microservices. Using same parameters as docker-compose.yaml file
func TestE2E(t *testing.T) {
	assertThat := assert.New(t)
	topic := "locations"

	tearDown := func() {
		nsq.UnRegister(topic, "")

		redisAddr := os.Getenv("REDIS")
		if redisAddr == "" {
			redisAddr = "localhost:6379"
		}

		conn, err := redis.Dial("tcp", redisAddr)
		if err != nil {
			println(err)
		}
		_, err = conn.Do("FLUSHALL")
		assertThat.NoError(err)
	}

	syncChan := make(chan bool, 0)
	secondSyncChan := make(chan bool, 0)
	numMessages := 0

	setup := func() {
		config := nsq.NewConfig()
		consumer, err := nsq.NewConsumer(topic, "te", config)
		assertThat.NoError(err)
		handler := nsq.HandlerFunc(func(message *nsq.Message) error {
			numMessages++
			if numMessages == 2 {
				time.Sleep(1 * time.Second)
				syncChan <- true
			}
			if numMessages == 4 {
				time.Sleep(1 * time.Second)
				secondSyncChan <- true
			}
			return nil
		})
		consumer.AddHandler(handler)
		consumer.ConnectToNSQD("localhost:4150")

		resp, err := http.Get("http://localhost:8080/healthz") // driver-location
		assertThat.NoError(err)
		assertThat.Equal(http.StatusOK, resp.StatusCode)

		resp, err = http.Get("http://localhost:7070/healthz") // zombie-driver
		assertThat.NoError(err)
		assertThat.Equal(http.StatusOK, resp.StatusCode)

		resp, err = http.Get("http://localhost:1138/healthz") // gateway
		assertThat.NoError(err)
		assertThat.Equal(http.StatusOK, resp.StatusCode)

		body := map[string]interface{}{"max_meters": 500, "last_minutes": 5}
		bodyJson, _ := json.Marshal(body)
		req, err := http.NewRequest(http.MethodPatch, "http://localhost:7070/config", bytes.NewBuffer(bodyJson))
		assertThat.NoError(err)
		resp, err = http.DefaultClient.Do(req)
		assertThat.NoError(err)
		assertThat.Equal(http.StatusOK, resp.StatusCode)
	}

	t.Run("Given all the services running and the default configuration set", func(t *testing.T) {
		setup()
		defer tearDown()

		t.Run("When creating near locations for a driver", func(t *testing.T) {
			driverID := uuid.NewV4().String()

			// create location
			body := map[string]interface{}{"latitude": 48.864193, "longitude": 2.350498}
			bodyJson, _ := json.Marshal(body)
			req, err := http.NewRequest(http.MethodPatch, "http://localhost:1138/drivers/"+driverID+"/locations", bytes.NewBuffer(bodyJson))
			assertThat.NoError(err)
			resp, err := http.DefaultClient.Do(req)
			assertThat.NoError(err)
			assertThat.Equal(http.StatusOK, resp.StatusCode)

			// create near location
			body = map[string]interface{}{"latitude": 48.863921, "longitude": 2.349211}
			bodyJson, _ = json.Marshal(body)
			req, err = http.NewRequest(http.MethodPatch, "http://localhost:1138/drivers/"+driverID+"/locations", bytes.NewBuffer(bodyJson))
			assertThat.NoError(err)
			resp, err = http.DefaultClient.Do(req)
			assertThat.NoError(err)
			assertThat.Equal(http.StatusOK, resp.StatusCode)

			t.Run("Then the driver is a zombie", func(t *testing.T) {
				<-syncChan
				resp, err = http.Get("http://localhost:1138/drivers/" + driverID)
				assertThat.NoError(err)
				assertThat.Equal(http.StatusOK, resp.StatusCode)

				body, err := ioutil.ReadAll(resp.Body)
				assertThat.NoError(err)

				var response ZombieResponse
				err = json.Unmarshal(body, &response)
				assertThat.NoError(err)

				assertThat.Equal(true, response.Zombie)
				assertThat.Equal(driverID, response.Id)
			})
		})
	})
	t.Run("Given all the services running and the default configuration set", func(t *testing.T) {
		setup()
		defer tearDown()

		t.Run("When creating far locations for a driver", func(t *testing.T) {
			driverID := uuid.NewV4().String()

			// create location
			body := map[string]interface{}{"latitude": 48.864193, "longitude": 2.350498}
			bodyJson, _ := json.Marshal(body)
			req, err := http.NewRequest(http.MethodPatch, "http://localhost:1138/drivers/"+driverID+"/locations", bytes.NewBuffer(bodyJson))
			assertThat.NoError(err)
			resp, err := http.DefaultClient.Do(req)
			assertThat.NoError(err)
			assertThat.Equal(http.StatusOK, resp.StatusCode)

			//create far location
			body = map[string]interface{}{"latitude": 48.858093, "longitude": 2.294694}
			bodyJson, _ = json.Marshal(body)
			req, err = http.NewRequest(http.MethodPatch, "http://localhost:1138/drivers/"+driverID+"/locations", bytes.NewBuffer(bodyJson))
			assertThat.NoError(err)
			resp, err = http.DefaultClient.Do(req)
			assertThat.NoError(err)
			assertThat.Equal(http.StatusOK, resp.StatusCode)

			t.Run("Then the driver is not a zombie", func(t *testing.T) {
				<-secondSyncChan
				resp, err := http.Get("http://localhost:1138/drivers/" + driverID)
				assertThat.NoError(err)
				assertThat.Equal(http.StatusOK, resp.StatusCode)

				body, err := ioutil.ReadAll(resp.Body)
				assertThat.NoError(err)

				var response ZombieResponse
				err = json.Unmarshal(body, &response)
				assertThat.NoError(err)

				assertThat.Equal(false, response.Zombie)
				assertThat.Equal(driverID, response.Id)
			})
		})
	})
}

type ZombieResponse struct {
	Id     string `json:"id"`
	Zombie bool   `json:"zombie"`
}
