package base

import "net/http"

type CRUDService[T any] interface {
	List(r *http.Request) ([]T, error)
	Get(r *http.Request) (*T, error)
	Create(r *http.Request) (*T, error)
	Update(r *http.Request) (*T, error)
	Delete(r *http.Request) error
}

type CRUDHandler[T any] struct {
	BaseHandler
	Service CRUDService[T]
}

func (h CRUDHandler[T]) Register(mux *http.ServeMux, basePath string) {
	// Intencionalmente vazio (chassis): o mapeamento de rotas vai ser definido
	// quando os contratos HTTP (paths, payloads, ids) forem acordados.
	_ = mux
	_ = basePath
}
