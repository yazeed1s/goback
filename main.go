package main

import "goback/cmd"

func main() {
	err := cmd.Execute()
	if err!= nil {
        panic(err)
    }
}
