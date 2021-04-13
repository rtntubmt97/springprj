package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func main() {
	// New Buffer.
	b := new(bytes.Buffer)

	// Write strings to the Buffer.
	// b.WriteString("ABC")
	// b.WriteString("DEF")

	err := binary.Write(b, binary.LittleEndian, uint16(1225))
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	err = binary.Write(b, binary.LittleEndian, uint16(24))
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Println(b.Bytes())
	fmt.Println(b.Len())
	b.WriteString("foo")

	var outi16 int16
	err = binary.Read(b, binary.LittleEndian, &outi16)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	fmt.Println(outi16)

	err = binary.Read(b, binary.LittleEndian, &outi16)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	fmt.Println(outi16)

	out := b.Next(3)
	fmt.Println(string(out))

	fmt.Println(b.Bytes())
	fmt.Println(b.Len())
	// Convert to a string and print it.
	// b.Reset()
	// fmt.Println(b.String())
}
