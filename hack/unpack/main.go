package main

import (
	"fmt"

	"github.com/abs3ntdev/haunt/src/haunt"
)

func main() {
	names := haunt.AssetNames()
	fmt.Println("Assets in bindata:")
	fmt.Println(names)

	wd := "assets"
	for i, v := range names {
		fmt.Printf("Restoring asset [%v] [%s]\n", i, v)
		err := haunt.RestoreAsset(wd, v)
		if err != nil {
			fmt.Println("Failed to restore", v)
		}
	}
}
