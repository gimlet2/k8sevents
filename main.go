package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	apis "github.com/ericchiang/k8s/apis/events/v1beta1"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.Print("Start")
	client := NewClient(os.Getenv("KUBE_CONFIG"))
	client.WatchEvents(func(eventType string, event *apis.Event) {
		log.WithFields(log.Fields{
			"eventType": eventType,
			"event":     event,
		}).Info()
	})
}
