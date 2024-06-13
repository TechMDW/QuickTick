# QuickTick

QuickTick is a Go package that provides a customizable clock implementation, allowing you to create clocks that can run faster or slower than real time based on a specified multiplier. This is useful for simulations and testing.

## Features

- **Customizable Time Multiplier:** Adjust the clock to run faster or slower than real time.
- **Custom Update Intervals:** Set custom intervals for how often the clock updates.
- **Context Integration:** Use contexts to control the lifecycle of the clock.
- **Concurrency Safe:** Safely access and manipulate the clock from multiple goroutines.

## Installation

To install QuickTick, use `go get`:

```sh
go get github.com/TechMDW/QuickTick@latest
```

## Usage

### Basic Usage

Create a new QuickTick clock with a multiplier. By default, the clock updates every millisecond.

```go
package main

import (
	"fmt"
	"time"
	"github.com/TechMDW/QuickTick"
)

func main() {
	multiplier := 2.0
	clock := quicktick.New(multiplier)
	defer clock.Stop()

	time.Sleep(2 * time.Second)
	fmt.Println("Accelerated Time:", clock.Now())
}
```

### Custom Usage

Create a QuickTick clock with a custom start time and update interval.

```go
package main

import (
	"fmt"
	"time"
	"github.com/TechMDW/QuickTick"
)

func main() {
	startTime := time.Now().Add(-1 * time.Hour)
	multiplier := 1.5
	updateInterval := 500 * time.Millisecond
	clock := quicktick.NewCustom(startTime, multiplier, updateInterval)
	defer clock.Stop()

	time.Sleep(2 * time.Second)
	fmt.Println("Accelerated Time:", clock.Now())
}
```

### Context Usage

Create a QuickTick clock that stops when the context is done.

```go
package main

import (
	"context"
	"fmt"
	"time"
	"github.com/TechMDW/QuickTick"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	multiplier := 2.0
	clock := quicktick.NewCtx(ctx, multiplier)
	defer clock.Stop()

	time.Sleep(2 * time.Second)
	fmt.Println("Accelerated Time:", clock.Now())
	cancel()
}
```

### Custom Context Usage

Create a QuickTick clock with a custom start time, multiplier, and update interval, controlled by a context.

```go
package main

import (
	"context"
	"fmt"
	"time"
	"github.com/TechMDW/QuickTick"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startTime := time.Now().Add(-1 * time.Hour)
	multiplier := 1.5
	updateInterval := 500 * time.Millisecond
	clock := quicktick.NewCustomCtx(ctx, startTime, multiplier, updateInterval)
	defer clock.Stop()

	time.Sleep(2 * time.Second)
	fmt.Println("Accelerated Time:", clock.Now())
	cancel()
}
```

## Concurrency

QuickTick is designed to be concurrency safe. You can safely access and manipulate the clock from multiple goroutines.

```go
package main

import (
	"fmt"
	"sync"
	"time"
	"github.com/TechMDW/QuickTick"
)

func main() {
	clock := quicktick.New(1.5)
	defer clock.Stop()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			fmt.Println("Goroutine 1:", clock.Now())
			time.Sleep(100 * time.Millisecond)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			fmt.Println("Goroutine 2:", clock.Now())
			time.Sleep(100 * time.Millisecond)
		}
	}()

	wg.Wait()
}
```

## Testing

QuickTick includes tests to ensure its functionality. Use the following command to run the tests:

```sh
go test ./...
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

This README provides an overview of the QuickTick package, installation instructions, usage examples, concurrency considerations, and testing information.
