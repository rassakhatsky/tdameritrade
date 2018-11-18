package auth

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
)

func callBackHandler(wg *sync.WaitGroup) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()

		// Get code

		oauthCode = r.URL.Query().Get("code")
		if len(oauthCode) == 0 {
			http.Error(w, "failed to parse oauth code", http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, "Auth code received. Check CLI")
	}

	return fn
}

// CreateAuthServer creates new server instance for Authenticated Request
func CreateAuthServer(ctx context.Context, wg *sync.WaitGroup, address string) error {
	handler := http.NewServeMux()
	handler.HandleFunc("/callback", callBackHandler(wg))

	// create new service
	srv := &http.Server{
		Addr:      address,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
		Handler:   handler,
	}

	// idleConnClosed allows to keep idle connections for some period of time
	idleConnClosed := make(chan struct{})

	go func() {
		wg.Wait()

		srv.Shutdown(ctx)
		close(idleConnClosed)
	}()

	// start the server
	if err := srv.ListenAndServeTLS("auth/server.crt", "auth/server.key"); err != http.ErrServerClosed {
		fmt.Println("Error:")
		fmt.Println(err.Error())
		fmt.Println("unable to start server")
		return err
	}

	<-idleConnClosed

	return nil
}
