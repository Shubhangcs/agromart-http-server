package routes

import (
	"github.com/go-chi/chi/v5"
	_ "github.com/shubhangcs/agromart-server/docs"
	"github.com/shubhangcs/agromart-server/internal/app"
	"github.com/shubhangcs/agromart-server/internal/middlewares"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.CORSMiddleware)
	r.Get("/health", app.HealthCheck)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	usersRoutes(app, r)
	businessRoutes(app, r)
	categoryRoutes(app, r)
	followerRoutes(app, r)
	rfqRoutes(app, r)
	productRoutes(app, r)
	return r
}

func usersRoutes(app *app.Application, r *chi.Mux) {

	r.Post("/admin/create", app.UserHandler.HandleCreateAdmin)
	r.Post("/user/create", app.UserHandler.HandleCreateUser)
	r.Post("/admin/login", app.TokenHandler.HandleGetAdminTokenByEmailPassword)
	r.Post("/user/login", app.TokenHandler.HandleGetUserTokenByEmailPassword)

	r.Route("/admin", func(r chi.Router) {
		r.Use(middlewares.AuthorizationMiddleware)
		r.Get("/get/admin/{id}", app.UserHandler.HandleGetAdminDetailsByID)
		r.Put("/update/image/{id}", app.BlobHandler.HandleUpdateAdminProfileImage)
		r.Put("/update/details/{id}", app.UserHandler.HandleUpdateAdminDetails)
		r.Put("/update/password/{id}", app.UserHandler.HandleUpdateAdminPassword)
		r.Delete("/delete", app.UserHandler.HandleDeleteAdmin)
	})

	r.Route("/user", func(r chi.Router) {
		r.Use(middlewares.AuthorizationMiddleware)
		r.Get("/get/all", app.UserHandler.HandleGetAllUsers)
		r.Put("/update/image/{id}", app.BlobHandler.HandleUpdateUserProfileImage)
		r.Put("/update/details/{id}", app.UserHandler.HandleUpdateUserDetails)
		r.Put("/update/password/{id}", app.UserHandler.HandleUpdateUserPassword)
		r.Delete("/delete", app.UserHandler.HandleDeleteUser)
		r.Get("/get/user/{id}", app.UserHandler.HandleGetUserDetailsByID)
		r.Put("/block/{id}", app.UserHandler.HandleBlockUser)
	})
}

func businessRoutes(app *app.Application, r *chi.Mux) {

	r.Route("/business", func(r chi.Router) {
		r.Use(middlewares.AuthorizationMiddleware)
		r.Post("/create", app.BusinessHandler.HandleCreateBusiness)
		r.Get("/get/all", app.BusinessHandler.HandleGetAllBusinesses)
		r.Get("/get/{id}", app.BusinessHandler.HandleGetBusinessDetails)
		r.Get("/get/complete/{id}", app.BusinessHandler.HandleGetCompleteBusinessDetails)
		r.Get("/get/user/{id}", app.BusinessHandler.HandleGetBusinessIDByUserID)
		r.Put("/update/{id}", app.BusinessHandler.HandleUpdateBusiness)
		r.Put("/update/image/{id}", app.BlobHandler.HandleUpdateBusinessProfileImage)
		r.Delete("/delete/{id}", app.BusinessHandler.HandleDeleteBusiness)
		r.Post("/social/create", app.BusinessHandler.HandleCreateSocial)
		r.Get("/social/get/{id}", app.BusinessHandler.HandleGetSocialDetails)
		r.Put("/social/update/{id}", app.BusinessHandler.HandleUpdateSocials)
		r.Post("/legal/create", app.BusinessHandler.HandleCreateLegal)
		r.Get("/legal/get/{id}", app.BusinessHandler.HandleGetLegalDetails)
		r.Put("/legal/update/{id}", app.BusinessHandler.HandleUpdateLegals)
		r.Post("/application/create", app.BusinessHandler.HandleCreateBusinessApplication)
		r.Get("/application/get/{id}", app.BusinessHandler.HandleGetBusinessApplicationDetails)
		r.Put("/application/accept/{id}", app.BusinessHandler.HandleAcceptBusinessApplication)
		r.Put("/application/reject/{id}", app.BusinessHandler.HandleRejectBusinessApplication)
		r.Put("/status/verify/{id}", app.BusinessHandler.HandleUpdateVerifyBusinessStatus)
		r.Put("/status/trust/{id}", app.BusinessHandler.HandleUpdateTrustBusinessStatus)
		r.Put("/status/block/{id}", app.BusinessHandler.HandleUpdateBlockBusinessStatus)
		r.Get("/status/{id}", app.BusinessHandler.HandleIsBusinessApproved)
	})
}

func categoryRoutes(app *app.Application, r *chi.Mux) {
	r.Route("/category", func(r chi.Router) {
		r.Use(middlewares.AuthorizationMiddleware)
		r.Post("/create", app.CategoryHandler.HandleCreateCategory)
		r.Post("/sub/create", app.CategoryHandler.HandleCreateSubCategory)
		r.Put("/update/image/{id}", app.BlobHandler.HandleUpdateCategoryImage)
		r.Put("/sub/update/image/{id}", app.BlobHandler.HandleUpdateSubCategoryImage)
		r.Put("/update/{id}", app.CategoryHandler.HandleUpdateCategory)
		r.Put("/sub/update/{id}", app.CategoryHandler.HandleUpdateSubCategory)
		r.Delete("/delete/{id}", app.CategoryHandler.HandleDeleteCategory)
		r.Delete("/sub/delete/{id}", app.CategoryHandler.HandleDeleteSubCategory)
		r.Get("/get/{id}", app.CategoryHandler.HandleGetCategoryByID)
		r.Get("/sub/get/{id}", app.CategoryHandler.HandleGetSubCategoryByID)
		r.Get("/get/all", app.CategoryHandler.HandleGetAllCategories)
		r.Get("/sub/get/all", app.CategoryHandler.HandleGetAllSubCategories)
		r.Get("/sub/get/category/{id}", app.CategoryHandler.HandleGetSubCategoriesByCategoryID)
	})
}

func followerRoutes(app *app.Application, r *chi.Mux) {
	r.Route("/follower", func(r chi.Router) {
		r.Use(middlewares.AuthorizationMiddleware)
		r.Post("/follow", app.FollowHandler.HandleCreateFollower)
		r.Post("/unfollow", app.FollowHandler.HandleRemoveFollower)
		r.Get("/get/followers/count/{id}", app.FollowHandler.HandleGetFollowersCount)
		r.Get("/get/following/count/{id}", app.FollowHandler.HandleGetFollowingCount)
		r.Get("/get/followers/{id}", app.FollowHandler.HandleGetAllFollowers)
		r.Get("/get/followings/{id}", app.FollowHandler.HandleGetAllFollowing)
	})
}

func rfqRoutes(app *app.Application, r *chi.Mux) {
	r.Route("/rfq", func(r chi.Router) {
		r.Use(middlewares.AuthorizationMiddleware)
		r.Post("/create", app.RFQHandler.HandleCreateRFQ)
		r.Get("/get/all", app.RFQHandler.HandleGetAllRFQ)
		r.Get("/get/{id}", app.RFQHandler.HandleGetRFQByBusinessID)
		r.Put("/update/{id}", app.RFQHandler.HandleUpdateRFQ)
		r.Put("/update/status/{id}", app.RFQHandler.HandleActivateRFQ)
		r.Delete("/delete/{id}", app.RFQHandler.HandleDeleteRFQ)
	})
}

func productRoutes(app *app.Application, r *chi.Mux) {
	r.Route("/product", func(r chi.Router) {
		r.Use(middlewares.AuthorizationMiddleware)
		r.Post("/create", app.ProductHandler.HandleCreateProduct)
		r.Put("/update/image", app.BlobHandler.HandleUpdateProductImage)
		r.Delete("/delete/image", app.BlobHandler.HandleDeleteProductImage)
		r.Put("/update/{id}", app.ProductHandler.HandleUpdateProduct)
		r.Delete("/delete/{id}", app.ProductHandler.HandleDeleteProduct)
		r.Get("/get/all", app.ProductHandler.HandleGetAllProducts)
		r.Get("/get/business/{id}", app.ProductHandler.HandleGetBusinessProducts)
		r.Get("/get/category/{id}", app.ProductHandler.HandleGetCategoryBasedProducts)
		r.Get("/get/sub/category/{id}", app.ProductHandler.HandleGetSubCategoryBasedProducts)
		r.Get("/get/followers/{id}", app.ProductHandler.HandleGetFollowersProducts)
		r.Get("/get/{id}", app.ProductHandler.HandleGetProductDetailsByID)
		r.Patch("/update/status/{id}", app.ProductHandler.HandleChangeProductActivateStatus)
	})
}
