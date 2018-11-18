package auth

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
)

func callBackHandler(cancel chan<- struct{}) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			cancel <- struct{}{}
		}()

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
func CreateAuthServer(ctx context.Context, cancel chan struct{}, address string) error {
	handler := http.NewServeMux()
	handler.HandleFunc("/callback", callBackHandler(cancel))

	// Generate a key pair from your pem-encoded cert and key ([]byte).
	cert, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err.Error())
		fmt.Println("unable to parse certificates")
		return err
	}
	// create new service
	srv := &http.Server{
		Addr: address,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true, // you really should not reuse for anything else
			Certificates:       []tls.Certificate{cert},
		},
		Handler: handler,
	}

	// idleConnClosed allows to keep idle connections for some period of time
	idleConnClosed := make(chan struct{})

	go func() {
		select {
		case <-ctx.Done():
		case <-cancel:
		}

		close(idleConnClosed)
	}()

	// start the server
	err = srv.ListenAndServeTLS("", "")
	if err != http.ErrServerClosed {
		fmt.Println("Error:")
		fmt.Println(err.Error())
		fmt.Println("unable to start server")
		return err
	}

	<-idleConnClosed

	return nil
}
