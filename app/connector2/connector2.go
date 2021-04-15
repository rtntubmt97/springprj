package main

import "github.com/rtntubmt97/springprj/connector"

func main() {
	conn := connector.Connector{}
	conn.Init(3)
	conn.Connect(1, 9090)
}
