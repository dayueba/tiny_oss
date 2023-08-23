package main

import "tiny_oss/internal"

func main() {
	app := internal.NewApp()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
