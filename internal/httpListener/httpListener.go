package httplistener

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

const KeyServerAddr = "ServerAdress"

type Inspector interface {
    InspectRequest(req *http.Request) (res bool, reason string)
}

func ScanHTTPNetwork(inspector Inspector) (shutdown func(ctx context.Context) error, manageCache func(cmd CacheCommand), err error) {
    ctx, cancel := context.WithCancel(context.Background())

    handler := &ProxyHandler{Inspector: inspector,
                             cache: map[string]time.Time{},
                             cacheTTL: 30*time.Minute.Abs(),
                             max_tries: 5,
                             mu: sync.RWMutex{},
                             cacheCmds: make(chan CacheCommand, 8),
                            }

    srv := &http.Server{
        Addr:    "127.0.0.1:4444",
        Handler: handler,
        BaseContext: func(l net.Listener) context.Context {
            return context.WithValue(ctx, KeyServerAddr, l.Addr().String())
        },
    }

    go handler.startCacheRoutine(ctx, 45*time.Second)

    go func() {
        e := srv.ListenAndServe()
        if errors.Is(e, http.ErrServerClosed) {
            fmt.Println("Server is closed")
        } else if e != nil {
            fmt.Printf("Error while starting server: %s\n", e)
        }
        cancel()
    }()

    // return a shutdown function to the caller (controller/main)
    shutdown = func(shutdownCtx context.Context) error {
        err := srv.Shutdown(shutdownCtx)
        cancel()
        return err
    }
    manageCache = func(cmd CacheCommand) {
        select {
        case handler.cacheCmds <- cmd:
        case <-ctx.Done():
        }
    }
    
    return shutdown, manageCache, nil
}