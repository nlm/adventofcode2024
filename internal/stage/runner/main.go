package main

import (
	"fmt"
	"io"
	"plugin"

	"github.com/nlm/adventofcode2023/internal/stage"
	"github.com/nlm/adventofcode2023/internal/utils"
)

// var flagPlugin = flag.String("plugin", "", "plugin to load")

func main() {
	// flag.Parse()
	// fns := make([]stage.StageFunc, 0)
	for i := 0; i < 15; i++ {
		fmt.Printf("----- DAY %02d -----\n", i+1)
		for j := 0; j < 2; j++ {
			fn := stage.StageFunc(utils.Must(utils.Must(plugin.Open(fmt.Sprintf("%02d.so", i+1))).Lookup(fmt.Sprintf("Stage%d", j+1))).(func(io.Reader) (any, error)))
			if fn == nil {
				fmt.Println("no func for", i+1, j+1)
			}
			stage.RunCLI(nil, fn)
		}
	}
}
