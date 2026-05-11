# RDIPs Dashboard

Small Go + htmx gateway for browsing RDIPs devices and viewing electric capacity performance data.

## Run locally

```bash
go run .
```

The server listens on `PORT` or `3000` by default.

## Configuration

Configuration is loaded with Viper from `.env`, then overridden by process environment variables.

- `API_BASE_URL`: upstream RDIPs API URL. Defaults to `http://sunflower_api:8080`.
- `API_KEY`: upstream API key.
- `PORT`: dashboard HTTP port. Defaults to `3000`.
- `HTTP_TIMEOUT_SECONDS`: upstream HTTP timeout. Defaults to `20`.

## Structure

- `internal/controllers`: HTTP route handlers.
- `internal/services`: upstream RDIPs API access.
- `internal/models`: shared request, response, and view models.
- `internal/views`: template rendering.
- `internal/chart`: chart view-model construction.
- `internal/scenario`: scenario payload generation.

## Build

```bash
go build -o rdips-dashboard .
```

Docker uses a multi-stage Go build and keeps the runtime image on Alpine, matching the RDIPs backend container style.
