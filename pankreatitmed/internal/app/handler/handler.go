package handler

import (
	"pankreatitmed/internal/app/authctx"
	"pankreatitmed/internal/app/services"

	"github.com/gin-gonic/gin"
)

type Handler struct{ svcs *services.Services }

func NewHandler(svcs *services.Services) *Handler { return &Handler{svcs: svcs} }

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		crit := api.Group("/criteria")
		{

			crit.GET("", h.CriteriaList)
			crit.GET("/:id", h.CriteriaGet)
			auth := crit.Group("")
			auth.Use(authctx.RequireAuth())
			{
				moder := auth.Group("")
				moder.Use(authctx.RequireModerator())
				{
					moder.POST("", h.CriteriaCreate)
					moder.PUT("/:id", h.CriteriaUpdate)
					moder.DELETE("/:id", h.CriteriaDelete)
					moder.POST("/:id/image", h.UploadCriterionImage)
				}
				auth.POST("/:id/add-to-draft", h.AddCriteriaToDraft)
			}

		}

		medord := api.Group("/pankreatitorders")
		{

			auth := medord.Group("")
			auth.Use(authctx.RequireAuth())
			{
				auth.GET("/cart", h.PankreatitOrderFromCart)
				auth.GET("", h.ListPankreatitOrders)

				moder := auth.Group("")
				moder.Use(authctx.RequireModerator())

				auth.GET(":id", h.PankreatitOrderGet)
				auth.PUT("/:id", h.PankreatitOrderUpdate)
				auth.PUT("/:id/form", h.PankreatitOrderForm)
				moder.PUT("/:id/set/:status", h.PankreatitOrderComplete)
				auth.DELETE("/:id", h.PankreatitOrderDelete)

				auth.DELETE("/items", h.DeletePankreatitOrderItem)
				auth.PUT("/items", h.UpdatePankreatitOrderItem)
			}
		}

		authmeduser := api.Group("/users")
		{
			authmeduser.POST("auth/register", h.MedUserRegistation)
			authmeduser.POST("auth/login", h.MedUserLogIn)
			auth := authmeduser.Group("")
			auth.Use(authctx.RequireAuth())
			{
				auth.GET("me", h.MedUserGetFields)
				auth.PUT("me", h.MedUserUpdateFields)
				moder := auth.Group("")
				moder.Use(authctx.RequireModerator())
				moder.POST("auth/logout/:token", h.MedUserLogOut)
			}
		}

		api.PUT("/setranson", h.PankreatitOrderSetRanson)
	}
}
