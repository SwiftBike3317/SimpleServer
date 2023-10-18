package main

import api "simpleServer/api_server"

func main() {

	server := api.NewAPIServer(":3001")
	server.Run()
}
