package entities

import "time"

type SecretCreateReport struct {
	ID string
}

type SecretCreateOptions struct {
	Driver string
}

type SecretListReport struct {
	ID        string
	Name      string
	Driver    string
	CreatedAt string
	UpdatedAt string
}

type SecretRmReport struct {
	ID  string
	Err error
}

type SecretInfoReport struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Spec      SecretSpec
}

type SecretSpec struct {
	Name   string
	Driver SecretDriverSpec
}

type SecretDriverSpec struct {
	Name    string
	Options map[string]string
}

type SecretCreateRequest struct {
	Name   string
	Data   string
	Driver SecretDriverSpec
}
