package main

import (
	"github.com/rtntubmt97/springprj/connector"
	"github.com/rtntubmt97/springprj/utils"
)

func main() {
	conn := connector.Connector{}
	utils.ReloadConfig("")
	conn.Init(1)
	conn.Listen(9090)
}
