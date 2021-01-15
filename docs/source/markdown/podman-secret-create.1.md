% podman-secret-create(1)

## NAME
podman\-secret\-create - Create a new secret

## SYNOPSIS
**podman secret create** [*options*] *name* *file*

## DESCRIPTION

Creates a secret using standard input or from a file for the secret content.
Secrets are a blob of sensitive data which a container needs at runtime but
you don’t want to store in the image or in source control, such as Usernames and passwords,
TLS certificates and keys, SSH keys or other important generic strings or binary content (up to 500 kb in size)

## OPTIONS

#### **--driver**=*driver*

Specify the secret driver (default file, which is unencrypted)

#### **--help**

Print usage statement


## EXAMPLES

```
$ podman secret create my_secret ./secret.json
$ podman secret create --driver=file my_secret ./secret.json
$ printf <secret> | podman secret create my_secret -
```

## SEE ALSO
podman-secret (1)

## HISTORY
January 2021, Originally compiled by Ashley Cui <acui@redhat.com>