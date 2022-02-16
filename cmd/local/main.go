package main

import (
	"context"

	"github.com/cyralinc/sidecar-failopen/internal/failopen"
)

// main could be a simple call `lambda.Start(<handler>)`. However, to avoid
// parsing the config every time the lambda is invoked, we build the config
// outside the lambda handler.
func main() {
	failopen.Run(context.Background())
}
