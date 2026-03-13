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

func ScanHTTPNetwork(inspector Inspector) (shutdown func(ctx context.Context) error, err error) {
    ctx, cancel := context.WithCancel(context.Background())

    handler := &ProxyHandler{Inspector: inspector,
                             cache: map[string]time.Time{},
                             cacheTTL: 30*time.Minute.Abs(),
                             max_tries: 5,
                             mu: &sync.RWMutex{},
                            }

    srv := &http.Server{
        Addr:    "127.0.0.1:4444",
        Handler: handler,
        BaseContext: func(l net.Listener) context.Context {
            return context.WithValue(ctx, KeyServerAddr, l.Addr().String())
        },
    }

    go handler.startCleanup(ctx, 45*time.Second)

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
    return func(shutdownCtx context.Context) error {
        err := srv.Shutdown(shutdownCtx)
        cancel()
        return err
    }, nil
}