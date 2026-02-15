package httplistener

import (
    "context"
    "errors"
    "fmt"
    "net"
    "net/http"
)

const KeyServerAddr = "ServerAdress"

type Inspector interface {
    InspectGET(r *http.Request) (block bool, reason string)
    InspectPOST(r *http.Request, bodyPreview []byte) (block bool, reason string)
    SeenRecently(host string) bool
    MarkSeen(host string)
}


func ScanHTTPNetwork(inspector Inspector) (shutdown func(ctx context.Context) error, err error) {
    ctx, cancel := context.WithCancel(context.Background())

    handler := &ProxyHandler{Inspector: inspector}

    srv := &http.Server{
        Addr:    "127.0.0.1:4444",
        Handler: handler,
        BaseContext: func(l net.Listener) context.Context {
            return context.WithValue(ctx, KeyServerAddr, l.Addr().String())
        },
    }

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