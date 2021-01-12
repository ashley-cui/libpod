package tunnel

import (
	"context"
	"fmt"
	"io"

	"github.com/containers/podman/v2/pkg/bindings/secrets"
	"github.com/containers/podman/v2/pkg/domain/entities"
)

func (ic *ContainerEngine) SecretCreate(ctx context.Context, name string, reader io.Reader, options entities.SecretCreateOptions) (*entities.SecretCreateReport, error) {
	fmt.Println("created tunnel")
	opts := new(secrets.CreateOptions).WithDriver(options.Driver).WithName(name)
	created, _ := secrets.Create(ic.ClientCtx, reader, opts)
	return created, nil
}

func (ic *ContainerEngine) SecretInspect(ctx context.Context, nameOrIDs []string) ([]*entities.SecretInfoReport, []error, error) {
	fmt.Println("inspected tunnel")
	var allInspect []*entities.SecretInfoReport
	for _, name := range nameOrIDs {
		inspected, _ := secrets.Inspect(ic.ClientCtx, name, nil)
		allInspect = append(allInspect, inspected)
		fmt.Println("thisistheinspected")
		fmt.Println(inspected)
	}
	return allInspect, nil, nil
}
func (ic *ContainerEngine) SecretList(ctx context.Context) ([]*entities.SecretInfoReport, error) {
	fmt.Println("listed tunnel")
	secrs, _ := secrets.List(ic.ClientCtx, nil)
	return secrs, nil
}
func (ic *ContainerEngine) SecretRm(ctx context.Context, nameOrID string) (*entities.SecretRmReport, error) {
	fmt.Println("removed tunnel")
	secrets.Remove(ic.ClientCtx, nameOrID, nil)
	return nil, nil
}
