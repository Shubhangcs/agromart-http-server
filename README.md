# agromart-http-server

Go (chi + Postgres) API for South Canara Agro Mart. Migrations run automatically on start (goose, embedded).

## Run locally

```bash
createdb agromart_dev
DATABASE_HOST=localhost DATABASE_USER=$USER DATABASE_PASSWORD=dev DATABASE_NAME=agromart_dev \
DATABASE_SSL_MODE=disable PORT=8090 PUSH_DRY_RUN=true ADMIN_BOOTSTRAP_SECRET=dev-secret \
go run .
```
`.env` is auto-loaded; explicit environment variables win over it. Swagger UI: `/swagger/index.html`.

## Tests

```bash
createdb agromart_test
TEST_DATABASE_NAME=agromart_test TEST_DATABASE_USER=$USER TEST_DATABASE_PASSWORD=dev go test ./... -count=1
```
`internal/tests` exercises the real HTTP API against that database (schema is dropped and re-migrated each run). CI runs the same suite and blocks deploys on failure.

## Environment

| Variable | Required | Purpose |
|---|---|---|
| `DATABASE_HOST/PORT/USER/PASSWORD/NAME/SSL_MODE` | yes | Postgres |
| `JWT_SECRET_KEY`, `JWT_TOKEN_ISSUER` | yes | JWT signing |
| `ACCESS_KEY`, `SECRET_KEY`, `REGION`, `BUCKET_NAME` | yes | S3 uploads (presigned URLs) |
| `ADMIN_BOOTSTRAP_SECRET` | first deploy | Enables `POST /admin/bootstrap` while no admin exists |
| `RESEND_API_KEY`, `EMAIL_FROM` | prod | Transactional email via Resend (welcome, password reset). Unset → emails are logged only |
| `EXPO_ACCESS_TOKEN` | optional | Expo push API auth token |
| `PUSH_DRY_RUN` | dev | `true` logs pushes instead of sending |
| `AUTH_RATE_LIMIT_PER_MIN` | optional | Per-IP limit on login / signup / password reset (default 20) |
| `APP_MIN_VERSION`, `APP_LATEST_VERSION`, `APP_ANDROID_STORE_URL`, `APP_UPDATE_MESSAGE` | optional | Served by `GET /app/config` for the force-update gate |
| `PORT` | optional | default 8080 |

## Auth model

- `POST /user/login` / `POST /admin/login` return a JWT with `role` (`user` | `admin`) and, for sellers, `business_id`.
- Admin-only: `/admin/*`, category mutations, seller-application accept/reject, business verify/trust/block, banner management, user listing/blocking.
- Ownership is enforced server-side: a caller can only mutate their own business, products, RFQs, images and read their own leads. `user_id` on ratings, reviews and follows is always taken from the token.
- First admin: `POST /admin/bootstrap` with header `X-Bootstrap-Secret` — works only while the `admins` table is empty.
- WebSocket `/chat/ws`: send the token as `Sec-WebSocket-Protocol: bearer, <jwt>` (the `?token=` query form is deprecated).

## Notifications & email

- Devices register with `POST /user/push-token` (`ExponentPushToken[...]`). Pushes are sent for new chat messages (receiver offline), new enquiries, and seller-application decisions.
- Email: welcome on signup; `POST /user/forgot-password` emails a 6-digit code (15 min, 5 attempts) for `POST /user/reset-password`.

## Deploy

`.github/workflows/deploy.yml`: vet + tests on every push/PR; on `main` builds the binary, SCPs it to the VM and restarts pm2 (secrets `SSH_PRIVATE_KEY`, `VM_USER`, `VM_IP`).
`.github/workflows/backup.yml`: nightly `pg_dump` → S3 (secrets `BACKUP_DATABASE_URL`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `BACKUP_BUCKET`).
