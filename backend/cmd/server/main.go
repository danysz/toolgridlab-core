package main

import (
	"github.com/toolgridlab/toolgridlab-core/backend/pkg/engine"
)

func main() {
	// Start the open-source core engine (commercial modules are injected in toolgridlab-dc)
	engine.Run()
}
