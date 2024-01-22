package main

import (
	"sync"

	_ "github.com/lib/pq"
	"github.com/vojdelenie/task0/internal/api"
	"github.com/vojdelenie/task0/internal/db"
	"github.com/vojdelenie/task0/internal/nats"
)

func main() {
	db.Run()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); go nats.Run() }()
	go func() { defer wg.Done(); go db.CacheInit() }()
	wg.Wait()
	api.Router()
}
