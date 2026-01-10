package base

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	coretel "github.com/brunojet/my-store-go/app/core/telemetry"
	"github.com/brunojet/my-store-go/app/infra/http/contracts"
)

type CRUDService[T any, CreateReq any, PatchReq any] interface {
	List(ctx context.Context) ([]T, error)
	Get(ctx context.Context, id uint) (*T, error)
	Create(ctx context.Context, req CreateReq) (*T, error)
	Patch(ctx context.Context, id uint, req PatchReq) (*T, error)
	Delete(ctx context.Context, id uint) error
}

type ToResponse[T any, Resp any] func(T) Resp

// CRUDHandler implements common HTTP plumbing for CRUD endpoints.
// It is adapter-level code (framework-agnostic) and should be reusable across resources.
//
// Error mapping:
// - if errors.Is(err, ValidationError) => 400
// - if errors.Is(err, NotFoundError) => 404
// - otherwise => 500
//
// Delete is assumed idempotent at the service level; HTTP returns 204 on success.
//
// This handler does not depend on Gin; the framework-specific adapter must provide a contracts.Context.
type CRUDHandler[T any, CreateReq any, PatchReq any, Resp any] struct {
	BaseHandler

	// Name identifies the resource (e.g. "apps"). Used for span/metric naming.
	Name string
	// Telemetry is optional; if nil, a noop provider is used.
	Telemetry coretel.Provider

	Service         CRUDService[T, CreateReq, PatchReq]
	ToResponse      ToResponse[T, Resp]
	NotFoundError   error
	ValidationError error
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) tel() coretel.Provider {
	if h.Telemetry == nil {
		return coretel.NoopProvider{}
	}
	return h.Telemetry
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) resource() string {
	r := strings.TrimSpace(h.Name)
	if r == "" {
		return "crud"
	}
	return r
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) spanName(action string) string {
	return h.resource() + "." + action
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) metricPrefix(action string) string {
	return "handler." + h.resource() + "." + strings.ToLower(action)
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) List(c contracts.Context) {
	tel := h.tel()
	ctx, span := tel.Tracer().Start(c.RequestContext(), h.spanName("List"))
	start := time.Now()

	items, err := h.Service.List(ctx)
	tel.Metrics().ObserveDuration(h.metricPrefix("List")+".duration", time.Since(start))
	if err != nil {
		tel.Metrics().Inc(h.metricPrefix("List")+".error", 1)
		tel.Logger().Error(ctx, "handler list failed", err)
		span.End(err)
		h.writeError(c, err)
		return
	}

	tel.Metrics().Inc(h.metricPrefix("List")+".ok", 1)
	span.End(nil, coretel.Field{Key: "count", Value: len(items)})

	resp := make([]Resp, 0, len(items))
	for _, it := range items {
		resp = append(resp, h.ToResponse(it))
	}
	h.JSON(c, http.StatusOK, resp)
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) Get(c contracts.Context) {
	id, ok := h.ParseUintParam(c, "id")
	if !ok {
		return
	}

	tel := h.tel()
	ctx, span := tel.Tracer().Start(c.RequestContext(), h.spanName("Get"), coretel.Field{Key: "id", Value: id})
	start := time.Now()

	item, err := h.Service.Get(ctx, id)
	tel.Metrics().ObserveDuration(h.metricPrefix("Get")+".duration", time.Since(start), coretel.Field{Key: "id", Value: id})
	if err != nil {
		tel.Metrics().Inc(h.metricPrefix("Get")+".error", 1, coretel.Field{Key: "id", Value: id})
		tel.Logger().Error(ctx, "handler get failed", err, coretel.Field{Key: "id", Value: id})
		span.End(err)
		h.writeError(c, err)
		return
	}

	tel.Metrics().Inc(h.metricPrefix("Get")+".ok", 1, coretel.Field{Key: "id", Value: id})
	span.End(nil, coretel.Field{Key: "id", Value: id})

	h.JSON(c, http.StatusOK, h.ToResponse(*item))
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) Create(c contracts.Context) {
	var req CreateReq
	if !h.BindJSON(c, &req) {
		return
	}

	tel := h.tel()
	ctx, span := tel.Tracer().Start(c.RequestContext(), h.spanName("Create"))
	start := time.Now()

	item, err := h.Service.Create(ctx, req)
	tel.Metrics().ObserveDuration(h.metricPrefix("Create")+".duration", time.Since(start))
	if err != nil {
		tel.Metrics().Inc(h.metricPrefix("Create")+".error", 1)
		tel.Logger().Error(ctx, "handler create failed", err)
		span.End(err)
		h.writeError(c, err)
		return
	}

	tel.Metrics().Inc(h.metricPrefix("Create")+".ok", 1)
	span.End(nil)

	h.JSON(c, http.StatusCreated, h.ToResponse(*item))
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) Patch(c contracts.Context) {
	id, ok := h.ParseUintParam(c, "id")
	if !ok {
		return
	}

	var req PatchReq
	if !h.BindJSON(c, &req) {
		return
	}

	tel := h.tel()
	ctx, span := tel.Tracer().Start(c.RequestContext(), h.spanName("Patch"), coretel.Field{Key: "id", Value: id})
	start := time.Now()

	item, err := h.Service.Patch(ctx, id, req)
	tel.Metrics().ObserveDuration(h.metricPrefix("Patch")+".duration", time.Since(start), coretel.Field{Key: "id", Value: id})
	if err != nil {
		tel.Metrics().Inc(h.metricPrefix("Patch")+".error", 1, coretel.Field{Key: "id", Value: id})
		tel.Logger().Error(ctx, "handler patch failed", err, coretel.Field{Key: "id", Value: id})
		span.End(err)
		h.writeError(c, err)
		return
	}

	tel.Metrics().Inc(h.metricPrefix("Patch")+".ok", 1, coretel.Field{Key: "id", Value: id})
	span.End(nil, coretel.Field{Key: "id", Value: id})

	h.JSON(c, http.StatusOK, h.ToResponse(*item))
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) Delete(c contracts.Context) {
	id, ok := h.ParseUintParam(c, "id")
	if !ok {
		return
	}

	tel := h.tel()
	ctx, span := tel.Tracer().Start(c.RequestContext(), h.spanName("Delete"), coretel.Field{Key: "id", Value: id})
	start := time.Now()

	if err := h.Service.Delete(ctx, id); err != nil {
		tel.Metrics().ObserveDuration(h.metricPrefix("Delete")+".duration", time.Since(start), coretel.Field{Key: "id", Value: id})
		tel.Metrics().Inc(h.metricPrefix("Delete")+".error", 1, coretel.Field{Key: "id", Value: id})
		tel.Logger().Error(ctx, "handler delete failed", err, coretel.Field{Key: "id", Value: id})
		span.End(err)
		h.writeError(c, err)
		return
	}

	tel.Metrics().ObserveDuration(h.metricPrefix("Delete")+".duration", time.Since(start), coretel.Field{Key: "id", Value: id})
	tel.Metrics().Inc(h.metricPrefix("Delete")+".ok", 1, coretel.Field{Key: "id", Value: id})
	span.End(nil, coretel.Field{Key: "id", Value: id})

	c.Status(http.StatusNoContent)
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) writeError(c contracts.Context, err error) {
	if h.ValidationError != nil && errors.Is(err, h.ValidationError) {
		h.Error(c, http.StatusBadRequest, "validation")
		return
	}
	if h.NotFoundError != nil && errors.Is(err, h.NotFoundError) {
		h.NotFound(c, "not found")
		return
	}
	h.InternalError(c)
}
