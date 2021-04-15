package main

import "github.com/rtntubmt97/springprj/connector"

func main() {
	conn := connector.Connector{}
	conn.Init(1)
	conn.Listen(9090)
}
