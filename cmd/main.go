package main

import (
	"tiny_oss/internal/rpc"
)

func main() {
	app := rpc.NewApp()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
