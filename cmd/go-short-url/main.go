package main

import (
	"fmt"
	"go-short-url/internal/config"
)

func main() {

	// TODO init config: cleanenv
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO init logger: slog

	// TODO init storage: sqlite or postgresql

	// TODO init router: chi

	// TODO run server

}
