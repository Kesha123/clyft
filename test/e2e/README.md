# E2E tests

End-to-end tests, fixtures and services for the `clyft` CLI.

## Zot OCI registry

Local registry backed by [zot](https://github.com/project-zot/zot), used as the
OCI registry target for `clyft`. Run via podman kube from `test/Makefile`:

```sh
make start-zot   # creates network + plays zot-pod.yaml
make stop-zot    # tears down pod + network
```

- **Pod**: `zot/zot-pod.yaml` — image `ghcr.io/project-zot/zot:v2.1.18`,
  `args: [serve, /etc/zot/config.yaml]`.
- **Network**: bridge `zot-net` (`172.29.0.0/24`), created by the Makefile target.
- **Port**: `127.0.0.1:5000` → container `5000` (HTTP; loopback only).
- **Volumes**: `./e2e/zot/config` → `/etc/zot` (read-only); PVC `zot-data` →
  `/var/lib/registry` (persisted registry storage).
- **Config**: `zot/config/config.yaml` — dist-spec `1.1.1`, HTTP on `0.0.0.0:5000`,
  htpasswd auth at `/etc/zot/htpasswd`, UI + search extensions enabled.
- **Auth**: `zot/config/htpasswd` — test user `clyft` / `test` (bcrypt).

## Fixtures

Fixtures contain all possible scenarios for `clyft` stack deployment.

Nginx endpoint container on a fixed bridge network
(`172.28.0.0/16`, gateway `172.28.0.1`, nginx at `172.28.0.10`, host `8080` → `80`):

- `fixtures/compose/compose.yaml` — Docker Compose: bridge network with IPAM, static IPv4.
- `fixtures/quadlet/nginx.network` + `fixtures/quadlet/nginx.container` — Podman
  Quadlet equivalents.
- `fixtures/pod/nginx-pod.yaml` — podman kube Pod.

Image: `nginx:1.27-alpine`, serving default static content.
