package abi

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/containers/common/pkg/secrets"
	"github.com/containers/podman/v2/pkg/domain/entities"
	"github.com/pkg/errors"
)

func (ic *ContainerEngine) SecretCreate(ctx context.Context, name string, reader io.Reader, options entities.SecretCreateOptions) (*entities.SecretCreateReport, error) {
	fmt.Println("created in abi")
	data, _ := ioutil.ReadAll(reader)
	secretsPath := filepath.Join(ic.Libpod.GetStore().GraphRoot(), "secrets")
	fmt.Printf("name: %s\ndata: %s\ndriver: %s\n", name, string(data), options.Driver)
	manager, err := secrets.NewManager(secretsPath) //should this be set in a flag? containers.conf?
	if err != nil {
		return nil, err
	}
	driverOptions := make(map[string]string)

	if options.Driver == "" {
		options.Driver = "file"
	}
	fmt.Printf("name: %s\ndata: %s\ndriver: %s\n", name, string(data), options.Driver)
	if options.Driver == "file" {
		driverOptions["path"] = secretsPath
	}
	secretID, err := manager.Store(name, data, options.Driver, driverOptions)
	if err != nil {
		return nil, err
	}
	return &entities.SecretCreateReport{
		ID: secretID,
	}, nil
}

func (ic *ContainerEngine) SecretInspect(ctx context.Context, nameOrIDs []string) ([]*entities.SecretInfoReport, []error, error) {
	fmt.Println("inspected")
	secretsPath := filepath.Join(ic.Libpod.GetStore().GraphRoot(), "secrets")
	manager, err := secrets.NewManager(secretsPath)
	if err != nil {
		return nil, nil, err
	}
	var errs []error
	var allInspect []secrets.Secret
	reports := make([]*entities.SecretInfoReport, 0, len(nameOrIDs))
	for _, nameOrID := range nameOrIDs {
		secret, err := manager.Lookup(nameOrID)
		if err != nil {
			if err.Error() == "no such secret" {
				errs = append(errs, errors.Errorf("no such secret %s", nameOrID))
				continue
			} else {
				return nil, nil, errors.Wrapf(err, "error inspecting volume %s", nameOrID)
			}
		}
		allInspect = append(allInspect, *secret)
		report := &entities.SecretInfoReport{
			ID:        secret.ID,
			CreatedAt: secret.CreatedAt,
			UpdatedAt: secret.CreatedAt,
			Spec: entities.SecretSpec{
				Name: secret.Name,
				Driver: entities.SecretDriverSpec{
					Name:    secret.Driver,
					Options: secret.DriverOptions,
				},
			},
		}
		reports = append(reports, report)

	}

	return reports, errs, nil
}
func (ic *ContainerEngine) SecretList(ctx context.Context) ([]*entities.SecretInfoReport, error) {
	secretsPath := filepath.Join(ic.Libpod.GetStore().GraphRoot(), "secrets")
	manager, err := secrets.NewManager(secretsPath)
	if err != nil {
		return nil, err
	}
	secretList, err := manager.List()
	if err != nil {
		return nil, err
	}
	fmt.Println(secretList)
	var report []*entities.SecretInfoReport
	for _, secret := range secretList {
		reportItem := entities.SecretInfoReport{
			ID:        secret.ID,
			CreatedAt: secret.CreatedAt,
			UpdatedAt: secret.CreatedAt,
			Spec: entities.SecretSpec{
				Name: secret.Name,
				Driver: entities.SecretDriverSpec{
					Name:    secret.Driver,
					Options: secret.DriverOptions,
				},
			},
		}
		report = append(report, &reportItem)
	}
	fmt.Println(report)
	return report, nil
}
func (ic *ContainerEngine) SecretRm(ctx context.Context, nameOrID string) (*entities.SecretRmReport, error) {
	fmt.Println("removed")
	secretsPath := filepath.Join(ic.Libpod.GetStore().GraphRoot(), "secrets")
	manager, err := secrets.NewManager(secretsPath)
	if err != nil {
		return nil, err
	}

	deletedID, err := manager.Delete(nameOrID)
	if err != nil {
		return nil, err
	}
	return &entities.SecretRmReport{
		ID: deletedID,
	}, nil
}
