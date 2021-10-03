package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	// TODO
	ROM, err := ioutil.ReadFile("./ROMs/Pokemon - Blue Version (USA, Europe) (SGB Enhanced).gb")
	if err != nil {
		log.Fatal(err.Error())
	}
	for i := 0; i < 255; i += 2 {
		fmt.Printf("%d %02X\n", i, (ROM[i]<<4)|ROM[i+1])
	}
	var mem = [0xFF]uint8{}
	test(mem)
	fmt.Println(mem)
}

func test(mem [0xFF]uint8) {
	mem[0] = 100
}
