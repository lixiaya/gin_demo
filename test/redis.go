package main

import (
	"context"
	"fmt"
	"gin_demo/global"
)

func main() {
	rdb := global.Rdb.Get(context.Background(), "trhxnlove@yeah.net")
	fmt.Println(rdb)
}
