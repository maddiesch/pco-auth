package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	callbackPageHTML = `
	<html>
		<body>
			<h2 style="text-align: center;">You can close this page</h2>
		</body>
	</html>
	`
)

func getCallbackURL(ctx context.Context, logger *log.Logger, port int) (*url.URL, error) {
	errs := make(chan error)
	complete := make(chan *url.URL)

	m := http.NewServeMux()
	s := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: m,
	}

	m.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		logger.Print("Received callback URL")

		w.Write([]byte(callbackPageHTML))

		go func() {
			complete <- r.URL
		}()
	})

	go func() {
		logger.Printf("Starting server on :%d", port)

		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errs <- err
		}
	}()

	shutdown := func() {
		logger.Printf("Shutdown callback server on :%d", port)

		err := s.Shutdown(ctx)
		if err != nil {
			_fail(ctx, err)
		}
	}

	for {
		select {
		case err := <-errs:
			shutdown()
			return nil, err
		case uri := <-complete:
			shutdown()
			return uri, nil
		}
	}
}
