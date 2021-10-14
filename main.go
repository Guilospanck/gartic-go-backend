package main

import (
	_ "base/src/infrastructure/environments"
	"fmt"
	"os"
)

func main() {

	fmt.Println(os.Getenv("DB_DATABASE_NAME"))

}
