package main

import "github.com/rtntubmt97/springprj/master"

func main() {
	master := new(master.Master)
	master.Listen(1)
}
