package main

import "github.com/abs3ntdev/haunt/cmd"

func main() {
	// names := haunt.AssetNames()
	// fmt.Println("Assets in bindata:")
	// fmt.Println(names)
	//
	// wd := "www"
	// for i, v := range names {
	// 	fmt.Printf("Restoring asset [%v] [%s]\n", i, v)
	// 	haunt.RestoreAsset(wd, v)
	// }
	cmd.Execute()
}
