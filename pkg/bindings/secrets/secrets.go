package secrets

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/containers/podman/v2/pkg/bindings"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/pkg/errors"
)

// List returns the configurations for existing volumes in the form of a slice.  Optionally, filters
// can be used to refine the list of volumes.
func List(ctx context.Context, options *ListOptions) ([]*entities.SecretInfoReport, error) {
	fmt.Println("in binding")
	var (
		secrs []*entities.SecretInfoReport
	)
	if options == nil {
		options = new(ListOptions)
	}
	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	response, err := conn.DoRequest(nil, http.MethodGet, "/secrets/json", nil, nil)
	if err != nil {
		return secrs, err
	}
	return secrs, response.Process(&secrs)
}

// Inspect returns low-level information about a volume.
func Inspect(ctx context.Context, nameOrID string, options *InspectOptions) (*entities.SecretInfoReport, error) {
	fmt.Println("inspect in bindings")
	var (
		inspect *entities.SecretInfoReport
	)
	if options == nil {
		options = new(InspectOptions)
	}
	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	response, err := conn.DoRequest(nil, http.MethodGet, "/secrets/%s/json", nil, nil, nameOrID)
	if err != nil {
		return inspect, err
	}
	return inspect, response.Process(&inspect)
}

func Remove(ctx context.Context, nameOrID string, options *RemoveOptions) error {
	fmt.Println("remove in bindings")

	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return err
	}
	params, err := options.ToParams()
	if err != nil {
		return err
	}
	response, err := conn.DoRequest(nil, http.MethodDelete, "/secrets/%s", params, nil, nameOrID)
	if err != nil {
		return err
	}
	return response.Process(nil)
}

func Create(ctx context.Context, reader io.Reader, options *CreateOptions) (*entities.SecretCreateReport, error) {
	var (
		create *entities.SecretCreateReport
	)
	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	params, err := options.ToParams()
	if err != nil {
		return nil, err
	}

	response, err := conn.DoRequest(reader, http.MethodPost, "/secrets/create", params, nil)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}
	return create, response.Process(&create)
}
