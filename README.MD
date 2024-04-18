
```
export DSN=postgres://postgres:12345@localhost:5432/godb?sslmode=disable

migrate -path=./migrations -database="postgres://postgres:12345@localhost:5432/godb?sslmode=disable" up

migrate -path=./migrations -database="postgres://postgres:12345@localhost:5432/godb?sslmode=disable" down

migrate -path=./migrations -database="postgres://postgres:12345@localhost:5432/godb?sslmode=disable" version

```