package metama

import (
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
)

type ConsulBackend struct {
	client *api.KV
}

func NewConsulBackend() *ConsulBackend {
	client, _ := api.NewClient(api.DefaultConfig())
	kv := client.KV()

	return &ConsulBackend{
		client: kv,
	}
}

func (b *ConsulBackend) pathOf(serverID string, key string) string {
	const prefix = "metadata"
	return fmt.Sprintf("%s/%s/%s", prefix, serverID, key)
}

func (b *ConsulBackend) Save() error {
	fmt.Fprintf(os.Stderr, "Currently do nothing.")
	return nil
}

func (b *ConsulBackend) Get(serverID string, key string) (string, error) {
	p, _, err := b.client.Get(b.pathOf(serverID, key), nil)
	if err != nil {
		return "", err
	}

	return string(p.Value), nil
}

func (b *ConsulBackend) Put(serverID string, key string, value string) error {
	data := &api.KVPair{
		Key:   b.pathOf(serverID, key),
		Value: []byte(value),
	}
	_, err := b.client.Put(data, nil)

	return err
}

func (b *ConsulBackend) Delete(serverID string, key string) error {
	_, err := b.client.Delete(b.pathOf(serverID, key), nil)

	return err
}
