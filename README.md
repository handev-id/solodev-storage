# Storage API (Go + Fiber + S3)

Simple API to upload and fetch files from S3.

## Project Structure

```text
.
├── main.go
├── utils/
│   ├── auth.go
│   ├── config.go
│   ├── key.go
│   ├── storage.go
│   └── url.go
├── Dockerfile
├── compose.yaml
├── .dockerignore
├── .env.example
├── go.mod
└── go.sum
```

## Environment Variables

Copy `.env.example` to `.env`.

Required:

- `UPLOAD_SECRET_KEY`
- `PUBLIC_BASE_URL` (example: `https://your-public-domain.example/`)
- `AWS_REGION`
- `AWS_S3_BUCKET`
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`

Optional:

- `APP_PORT` (default: `3000`)
- `AWS_SESSION_TOKEN`
- `AWS_ENDPOINT_URL`
- `AWS_S3_USE_PATH_STYLE` (`true/false`)

## Run (Local)

```bash
go run .
```

## Docker (Production)

```bash
docker build -t storage-api:latest .
docker run -d --name storage-api \
  --env-file .env \
  -p ${APP_PORT:-3000}:${APP_PORT:-3000} \
  --restart unless-stopped \
  storage-api:latest
```

Or with compose:

```bash
docker compose up -d --build
```

## Endpoints

### `POST /`

- Auth header: `Authorization: Bearer <UPLOAD_SECRET_KEY>`
- Body: `multipart/form-data`
- Required field: `file`
- Optional `key`: query `?key=...` or form field `key`
- Optional `folder`: query `?folder=...` or form field `folder`
- Priority: query value is used first, then form value
- If `folder` is provided, final key becomes: `<folder>/<key-or-generated-name>`

Success response (`201`):

```json
{
  "message": "file uploaded",
  "key": "docs/test.txt",
  "public_url": "https://your-public-domain.example/docs/test.txt"
}
```

Notes:

- Returned URL is a normal public URL (not a signed URL).
- Upload uses `public-read` ACL.
- Bucket/public access policy must allow public reads for this to work externally.

### `GET /*`

- Streams file from S3.
- URL format: `/<key>` (supports nested path)
- Example: `/my-file.png` or `/docs/test.txt`

### `GET /healthz`

```json
{ "status": "ok" }
```

## Quick cURL

Upload:

```bash
curl -X POST 'http://localhost:3000/?folder=docs&key=test.txt' \
  -H 'Authorization: Bearer your-static-upload-secret' \
  -F 'file=@/absolute/path/to/test.txt'
```

Download:

```bash
curl 'http://localhost:3000/docs/test.txt' --output test.txt
```
