package server

import (
	"net/http"

	"github.com/beozel/go-microservices/internal/dberrors"
	"github.com/beozel/go-microservices/internal/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (s *EchoServer) GetAllVendors(ctx echo.Context) error {
	vendors, err := s.DB.GetAllVendors(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, vendors)
}

func (s *EchoServer) AddVendor(ctx echo.Context) error {
	vendor := new(models.Vendor)
	if err := ctx.Bind(vendor); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}
	vendor, err := s.DB.AddVendor(ctx.Request().Context(), vendor)
	if err != nil {
		switch err.(type) {
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}
	return ctx.JSON(http.StatusCreated, vendor)
}

func (s *EchoServer) GetVendorById(ctx echo.Context) error {
	id := ctx.Param("id")

	// Optional: Validate UUID format early
	if _, err := uuid.Parse(id); err != nil {
		return ctx.JSON(http.StatusNotFound, echo.Map{
			"error": "Vendor not found",
		})
	}

	vendor, err := s.DB.GetVendorById(ctx.Request().Context(), id)
	if err != nil {
		switch err.(type) {
		case *dberrors.NotFoundError:
			return ctx.JSON(http.StatusNotFound, echo.Map{
				"error": "Vendor not found",
			})
		default:
			return ctx.JSON(http.StatusInternalServerError, echo.Map{
				"error": "Internal server error",
			})
		}
	}

	return ctx.JSON(http.StatusOK, vendor)
}

func (s *EchoServer) UpdateVendor(ctx echo.Context) error {
	ID := ctx.Param("id")
	vendor := new(models.Vendor)
	if err := ctx.Bind(vendor); err != nil {
		return ctx.JSON(http.StatusUnsupportedMediaType, err)
	}
	if ID != vendor.VendorID {
		return ctx.JSON(http.StatusBadRequest, "id on path doesn't match id on body")
	}
	vendor, err := s.DB.UpdateVendor(ctx.Request().Context(), vendor)
	if err != nil {
		switch err.(type) {
		case *dberrors.NotFoundError:
			return ctx.JSON(http.StatusNotFound, err)
		case *dberrors.ConflictError:
			return ctx.JSON(http.StatusConflict, err)
		default:
			return ctx.JSON(http.StatusInternalServerError, err)
		}
	}
	return ctx.JSON(http.StatusOK, vendor)
}

func (s *EchoServer) DeleteVendor(ctx echo.Context) error {
	ID := ctx.Param("id")
	err := s.DB.DeleteVendor(ctx.Request().Context(), ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.NoContent(http.StatusResetContent)
}