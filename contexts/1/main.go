// BEHAVIOR: Book an Hotel
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// init in main thread
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 6*time.Second)

	defer cancel()

	BookHotel(ctx)
}

func BookHotel(ctx context.Context) {
	select {
	case <-ctx.Done():
		fmt.Println("Hotel booking cancelled: Timeout reached.")
		return
	case <-time.After(5 * time.Second):
		fmt.Println("Hotel booking successful.")
		return
	}
}
