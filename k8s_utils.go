package main

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/ericchiang/k8s"
	apis "github.com/ericchiang/k8s/apis/events/v1beta1"
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
)

// Client is wrapper on top of k8s.Client
type Client interface {
	WatchEvents(fn func(string, *apis.Event))
}

// ClientImpl implementation of Client interface
type ClientImpl struct {
	*k8s.Client
}

// NewClient constructor of Client
func NewClient(kubeConfigPath string) Client {
	var client *k8s.Client
	var err error
	if kubeConfigPath != "" {
		client, err = loadClientWithConfig(kubeConfigPath)
	} else {
		client, err = k8s.NewInClusterClient()
	}
	if err != nil {
		log.Fatal(err)
	}
	return &ClientImpl{client}
}

func loadClientWithConfig(kubeconfigPath string) (*k8s.Client, error) {
	data, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("read kubeconfig: %v", err)
	}

	// Unmarshal YAML into a Kubernetes config object.
	var config k8s.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unmarshal kubeconfig: %v", err)
	}
	return k8s.NewClient(&config)
}

// WatchEvents is watching events with lambda
func (c ClientImpl) WatchEvents(fn func(string, *apis.Event)) {
	for {
		var event apis.Event
		watcher, err := c.Watch(context.Background(), k8s.AllNamespaces, &event)
		if err != nil {
			log.Printf("Failed to watch events - %v", err)
			break
		}
		defer watcher.Close() // Always close the returned watcher.

		for {
			ev := new(apis.Event)
			eventType, err := watcher.Next(ev)
			if err != nil {
				log.Printf("Failed to get next event, rewatch - %v", err)
				break
			}
			go fn(eventType, ev)
		}
	}
}
