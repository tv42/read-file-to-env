# read-file-to-env -- Read files into environment variables and execute command

Example use:

```
read-file-to-env -one-line=HOST=/etc/hostname sh -c 'echo host=$HOST'
```

Use with systemd credentials, useful for example when a container image can *only* receive things via environment variables.
For example, `postgres:14` (see https://github.com/docker-library/docs/blob/master/postgres/README.md for background).
For container images you control, please support reading from files.

```
[Service]
# only relevant lines shown, see `podman generate systemd`
LoadCredential=pg-user:/etc/keys/postgres-admin-username
LoadCredential=pg-password:/etc/keys/postgres-admin-password
ExecStart=read-file-to-env \
    -one-line=POSTGRES_USER=${CREDENTIALS_DIRECTORY}/pg-user \ 
    -one-line=POSTGRES_PASSWORD=${CREDENTIALS_DIRECTORY}/pg-password \ 
    -- \
    podman run \
    ... \
    --env=POSTGRES_USER \
    --env=POSTGRES_PASSWORD \
    postgres:14
```
