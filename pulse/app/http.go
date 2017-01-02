package app

import (
	"context"
	"log"
	"net/http"

	"github.com/bearyinnovative/radagast/config"
)

func Serve(ctx context.Context) {
	config := config.FromContext(ctx).Get("pulse.app").Config()

	bind := config.Get("bind").String()
	log.Printf("pulse.app listening on %s", bind)
	log.Fatal(http.ListenAndServe(bind, nil))
}
