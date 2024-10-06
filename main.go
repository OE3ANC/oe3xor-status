package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type MessageMmdvm struct{ MMDVM MMDVM }
type MessageText struct{ Text Text }
type MessageBer struct{ BER BER }
type MessageM17 struct{ M17 M17 }
type MessageFm struct{ FM FM }

type MMDVM struct {
	Mode string
}

type M17 struct {
	Action        string
	Ber           float32
	Duration      float32
	DestinationCs string `json:"destination_cs"`
	SourceCs      string `json:"source_cs"`
	Source        string
	TrafficType   string `json:"traffic_type"`
}

type FM struct {
	State string
}

type BER struct {
	Mode  string
	Value float32
}

type Text struct {
	Mode  string
	Value string
}

var app = pocketbase.New()

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

	if msg.Topic() == "mmdvm/json" {
		fmt.Printf("MSG: %s\n", msg.Payload())

		var mmdvm MessageMmdvm
		_ = json.Unmarshal(msg.Payload(), &mmdvm)
		if mmdvm.MMDVM.Mode != "" {
			handleMmdvmMessage(mmdvm)
		}

		var text MessageText
		_ = json.Unmarshal(msg.Payload(), &text)

		if text.Text.Mode != "" {
			handleTextMessage(text)
		}

		var ber MessageBer
		_ = json.Unmarshal(msg.Payload(), &ber)

		if ber.BER.Mode != "" {
			handleBerMessage(ber)
		}

		var m17 MessageM17
		_ = json.Unmarshal(msg.Payload(), &m17)

		if m17.M17.Action != "" {
			handleM17Message(m17)
		}

		var fm MessageFm
		_ = json.Unmarshal(msg.Payload(), &fm)

		if fm.FM.State != "" {
			handleFmMessage(fm)
		}
	} else {
		handleRawLog(string(msg.Payload()))
	}

}

func handleMmdvmMessage(message MessageMmdvm) {
	fmt.Printf("Repeater mode is now %s\n", message.MMDVM.Mode)

	collection, err := app.Dao().FindCollectionByNameOrId("mmdvm_mode")
	if err != nil {
		fmt.Printf("Can't find collection: %s", err)
		return
	}

	record := models.NewRecord(collection)
	record.Set("mode", message.MMDVM.Mode)

	if err := app.Dao().SaveRecord(record); err != nil {
		fmt.Printf("Can't save record: %s", err)
		return
	}
}

func handleTextMessage(message MessageText) {
	fmt.Printf("Text Message %s\n", message.Text.Value)

	collection, err := app.Dao().FindCollectionByNameOrId("mmdvm_logs_text")
	if err != nil {
		fmt.Printf("Can't find collection: %s", err)
		return
	}

	record := models.NewRecord(collection)
	record.Set("mode", message.Text.Mode)
	record.Set("value", message.Text.Value)

	if err := app.Dao().SaveRecord(record); err != nil {
		fmt.Printf("Can't save record: %s", err)
		return
	}
}

func handleBerMessage(message MessageBer) {
	fmt.Printf("Last BER was %f\n", message.BER.Value)

	collection, err := app.Dao().FindCollectionByNameOrId("mmdvm_logs_ber")
	if err != nil {
		fmt.Printf("Can't find collection: %s", err)
		return
	}

	record := models.NewRecord(collection)
	record.Set("mode", message.BER.Mode)
	record.Set("value", message.BER.Value)

	if err := app.Dao().SaveRecord(record); err != nil {
		fmt.Printf("Can't save record: %s", err)
		return
	}

}

func handleM17Message(message MessageM17) {
	if message.M17.Action == "end" || message.M17.Action == "lost" {
		// Update last open QSO in DB
		fmt.Printf("%s of M17 Transmission. BER was %f. Duration was %f\n", message.M17.Action, message.M17.Ber, message.M17.Duration)

		record := &models.Record{}

		err := app.Dao().RecordQuery("mmdvm_qso_m17").
			OrderBy("created DESC").
			AndWhere(dbx.HashExp{inflector.Columnify("duration"): 0}).
			Limit(1).
			One(record)

		if err != nil {
			fmt.Printf("Can't find record: %s", err)
			return
		}

		record.Set("ber", message.M17.Ber)
		record.Set("duration", message.M17.Duration)
		record.Set("action", message.M17.Action)

		if err := app.Dao().SaveRecord(record); err != nil {
			fmt.Printf("Can't save record: %s", err)
			return
		}

	} else {
		fmt.Printf("New M17 %s-Transmission from %s to %s via %s \n", message.M17.TrafficType, message.M17.SourceCs, message.M17.DestinationCs, message.M17.Source)

		collection, err := app.Dao().FindCollectionByNameOrId("mmdvm_qso_m17")
		if err != nil {
			fmt.Printf("Can't find collection: %s", err)
			return
		}

		record := models.NewRecord(collection)
		record.Set("action", message.M17.Action)
		record.Set("destination_callsign", message.M17.DestinationCs)
		record.Set("source_callsign", message.M17.SourceCs)
		record.Set("source", message.M17.Source)
		record.Set("traffic_type", message.M17.TrafficType)

		if err := app.Dao().SaveRecord(record); err != nil {
			fmt.Printf("Can't save record: %s", err)
			return
		}
	}
}

func handleFmMessage(message MessageFm) {
	fmt.Printf("FM State is now %s\n", message.FM.State)

	collection, err := app.Dao().FindCollectionByNameOrId("mmdvm_logs_fm")
	if err != nil {
		fmt.Printf("Can't find collection: %s", err)
		return
	}

	record := models.NewRecord(collection)
	record.Set("state", message.FM.State)

	if err := app.Dao().SaveRecord(record); err != nil {
		fmt.Printf("Can't save record: %s", err)
		return
	}
}

func handleRawLog(message string) {
	fmt.Printf("Raw Log: %s\n", message)

	collection, err := app.Dao().FindCollectionByNameOrId("mmdvm_logs_raw")
	if err != nil {
		fmt.Printf("Can't find collection: %s", err)
		return
	}

	record := models.NewRecord(collection)
	record.Set("value", message)

	if err := app.Dao().SaveRecord(record); err != nil {
		fmt.Printf("Can't save record: %s", err)
		return
	}
}

func main() {

	broker := os.Getenv("OE3XOR_STATUS_BROKER")

	if broker == "" {
		broker = "localhost:1883"
	}

	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID("oe3xor_status_client")

	opts.SetKeepAlive(60 * time.Second)

	// Set the message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe to a topic
	if token := c.Subscribe("mmdvm/#", 2, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())

	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Admin UI
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("./pb_public"), false))
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
