package main

import (
	"backend-exercise/cmd"
	"backend-exercise/utils"
)

func init() {
	utils.LoadEnv()
}

func main() {
	cmd.Execute()
}
