
```
export DSN=postgres://postgres:12345@localhost:5432/godb?sslmode=disable

export DSN=postgres://postgres:12345@localhost:5432/d.talipovDB?sslmode=disable


migrate -path=./migrations -database="postgres://postgres:12345@localhost:5432/d.talipovDB?sslmode=disable" up

migrate -path=./migrations -database="postgres://postgres:12345@localhost:5432/d.talipovDB?sslmode=disable" down

migrate -path=./migrations -database="postgres://postgres:12345@localhost:5432/d.talipovDB?sslmode=disable" version

http://localhost:4000/v1/users?sort=-id

http://localhost:4000/v1/users?name=Damir

something
```