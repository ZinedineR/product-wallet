package api

import "product-wallet/internal/delivery/http"

type Middleware struct {
	http.Handler
}
