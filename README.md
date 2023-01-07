## Kiteps

# Infrastructure

## Setup localy with docker

- Run the `docker-compose up` command
- Now you will have to check your network with `ifconfig | grep inet`
- Command above should give you a list of hostnames. Get the one that looks something like this '192.168.49.1'. This will be the host for your database.
- Now you will have a postgresql database with three databases inside it: books, users, chat. They all have different users asigned to them. books - books_user, users - users_user, chat - chat_user.

Here is an example of connecting with golang

```go
package main

import (
	"context"

	"github.com/sultanaliev-s/kiteps/pkg/db"
	"github.com/sultanaliev-s/kiteps/pkg/logging"
)

func main() {
	log, _ := logging.NewLogger("debug")

	pool, err := db.NewPGXPool(
		"postgres://books_user:password@192.168.49.1:5432/books?sslmode=disable",
		log, context.Background(),
	)
	if err != nil {
		log.Fatal("could not connect to database", logging.Error("err", err))
	}
	defer pool.Close()

	pool.Exec(context.Background(), "SELECT 1")
}
```

Good luck 😃
