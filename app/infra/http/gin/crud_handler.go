package base

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
// It is adapter-level code (depends on Gin) and should be reusable across resources.
//
// Error mapping:
// - if errors.Is(err, ValidationError) => 400
// - if errors.Is(err, NotFoundError) => 404
// - otherwise => 500
//
// Delete is assumed idempotent at the service level; HTTP returns 204 on success.
type CRUDHandler[T any, CreateReq any, PatchReq any, Resp any] struct {
	BaseHandler

	Service         CRUDService[T, CreateReq, PatchReq]
	ToResponse      ToResponse[T, Resp]
	NotFoundError   error
	ValidationError error
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) List(c *gin.Context) {
	items, err := h.Service.List(c.Request.Context())
	if err != nil {
		h.writeError(c, err)
		return
	}

	resp := make([]Resp, 0, len(items))
	for _, it := range items {
		resp = append(resp, h.ToResponse(it))
	}
	h.JSON(c, http.StatusOK, resp)
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) Get(c *gin.Context) {
	id, ok := h.ParseUintParam(c, "id")
	if !ok {
		return
	}

	item, err := h.Service.Get(c.Request.Context(), id)
	if err != nil {
		h.writeError(c, err)
		return
	}

	h.JSON(c, http.StatusOK, h.ToResponse(*item))
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) Create(c *gin.Context) {
	var req CreateReq
	if !h.BindJSON(c, &req) {
		return
	}

	item, err := h.Service.Create(c.Request.Context(), req)
	if err != nil {
		h.writeError(c, err)
		return
	}

	h.JSON(c, http.StatusCreated, h.ToResponse(*item))
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) Patch(c *gin.Context) {
	id, ok := h.ParseUintParam(c, "id")
	if !ok {
		return
	}

	var req PatchReq
	if !h.BindJSON(c, &req) {
		return
	}

	item, err := h.Service.Patch(c.Request.Context(), id, req)
	if err != nil {
		h.writeError(c, err)
		return
	}

	h.JSON(c, http.StatusOK, h.ToResponse(*item))
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) Delete(c *gin.Context) {
	id, ok := h.ParseUintParam(c, "id")
	if !ok {
		return
	}

	if err := h.Service.Delete(c.Request.Context(), id); err != nil {
		h.writeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *CRUDHandler[T, CreateReq, PatchReq, Resp]) writeError(c *gin.Context, err error) {
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
