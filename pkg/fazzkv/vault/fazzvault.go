package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/payfazz/go-apt/pkg/fazzcommon/httpError"
	"github.com/payfazz/go-apt/pkg/fazzkv"
)

type Interface interface {
	fazzkv.Store
	ReadPath(path string) (Interface, error)
}

type vault struct {
	client     *api.Client
	collection map[string]string
}

func (v *vault) ReadPath(path string) (Interface, error) {
	var result map[string]string
	data, err := v.client.Logical().Read(path)
	if err != nil {
		return nil, err
	}

	b, _ := json.Marshal(data.Data)
	_ = json.Unmarshal(b, &result)
	v.collection = result
	return v, nil
}

func (v *vault) Set(key string, value interface{}) error {
	return nil
}

func (v *vault) Get(key string) (string, error) {
	if len(v.collection) < 1 {
		return "", httpError.NotFound("collection still empty, please use vault ReadPath() function first.")
	}
	return v.collection[key], nil
}

func (v *vault) Delete(key string) error {
	return nil
}

func (v *vault) Truncate() error {
	return nil
}

func authenticateUser(client *api.Client, username string, password string) (*api.Client, error) {
	options := map[string]interface{}{
		"password": password,
	}
	path := fmt.Sprintf("auth/userpass/login/%s", username)
	secret, err := client.Logical().Write(path, options)
	if err != nil {
		return nil, err
	}
	client.SetToken(secret.Auth.ClientToken)
	return client, nil
}

// NewFazzVault is a function that used to get new vault client
func NewFazzVault(url string, username string, password string) (Interface, error) {
	client, err := api.NewClient(&api.Config{Address: url, HttpClient: &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}})
	if err != nil {
		return nil, httpError.InternalServer("cannot connect to vault")
	}
	client, err = authenticateUser(client, username, password)
	if err != nil {
		return nil, httpError.InternalServer("cannot authenticate user to vault")
	}

	return &vault{
		client: client,
	}, nil
}
