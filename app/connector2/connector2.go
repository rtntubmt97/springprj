package main

import (
	"github.com/rtntubmt97/springprj/connector"
	"github.com/rtntubmt97/springprj/utils"
)

func main() {
	conn := connector.Connector{}
	utils.UseLog = true
	conn.Init(3)
	conn.Connect(1, 9090)
}
