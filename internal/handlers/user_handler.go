package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/shubhangcs/agromart-server/internal/models"
	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
	"github.com/shubhangcs/agromart-server/internal/validator"
)

// ErrorResponse is the generic error response body.
type ErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}

// MessageResponse is the generic success response body.
type MessageResponse struct {
	Message string `json:"message" example:"operation successful"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *slog.Logger
}

func NewUserHandler(userStore store.UserStore, logger *slog.Logger) *UserHandler {
	return &UserHandler{userStore: userStore, logger: logger}
}

// HandleCreateAdmin godoc
// @Summary      Create a new admin
// @Description  Registers a new admin account
// @Tags         admins
// @Accept       json
// @Produce      json
// @Param        body body models.CreateAdminRequest true "Admin registration payload"
// @Success      201 {object} map[string]string "admin created successfully"
// @Failure      400 {object} ErrorResponse "Invalid payload or validation error"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /admin/create [post]
func (uh *UserHandler) HandleCreateAdmin(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, uh.logger, "invalid request payload", err)
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.BadRequest(w, uh.logger, err.Error(), err)
		return
	}
	admin := &models.Admin{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	}
	if err := admin.Password.Set(req.Password); err != nil {
		utils.ServerError(w, uh.logger, "create admin: hash password", err)
		return
	}
	if err := uh.userStore.CreateAdmin(admin); err != nil {
		uh.logger.Error("create admin", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create admin"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "admin created successfully", "admin_id": admin.ID})
}

// HandleCreateUser godoc
// @Summary      Create a new user
// @Description  Registers a new user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body body models.CreateAdminRequest true "User registration payload"
// @Success      201 {object} map[string]string "user created successfully"
// @Failure      400 {object} ErrorResponse "Invalid payload or validation error"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /user/create [post]
func (uh *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, uh.logger, "invalid request payload", err)
		return
	}
	if err := validator.Validate(&req); err != nil {
		utils.BadRequest(w, uh.logger, err.Error(), err)
		return
	}
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	}
	if err := user.Password.Set(req.Password); err != nil {
		utils.ServerError(w, uh.logger, "create user: hash password", err)
		return
	}
	if err := uh.userStore.CreateUser(user); err != nil {
		uh.logger.Error("create user", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create user"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "user created successfully", "user_id": user.ID})
}

// HandleUpdateAdminDetails godoc
// @Summary      Update admin profile details
// @Description  Updates the profile details of the admin with the given ID
// @Tags         admins
// @Accept       json
// @Produce      json
// @Param        id   path      string                              true "Admin ID"
// @Param        body body      models.UpdateUserDetailsRequest     true "Admin update payload"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /admin/update/details/{id} [put]
func (uh *UserHandler) HandleUpdateAdminDetails(w http.ResponseWriter, r *http.Request) {
	adminID, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, uh.logger, err.Error(), err)
		return
	}
	var req models.UpdateUserDetailsRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, uh.logger, "invalid request payload", err)
		return
	}
	if err = validator.Validate(&req); err != nil {
		utils.BadRequest(w, uh.logger, err.Error(), err)
		return
	}
	admin := &models.Admin{
		ID:        adminID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	}
	if err = uh.userStore.UpdateAdminDetails(admin); err != nil {
		utils.ServerError(w, uh.logger, "update admin details", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "admin details updated successfully"})
}

// HandleUpdateUserDetails godoc
// @Summary      Update user profile details
// @Description  Updates the profile details of the user with the given ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string                              true "User ID"
// @Param        body body      models.UpdateUserDetailsRequest     true "User update payload"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /user/update/details/{id} [put]
func (uh *UserHandler) HandleUpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, uh.logger, err.Error(), err)
		return
	}
	var req models.UpdateUserDetailsRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, uh.logger, "invalid request payload", err)
		return
	}
	if err = validator.Validate(&req); err != nil {
		utils.BadRequest(w, uh.logger, err.Error(), err)
		return
	}
	user := &models.User{
		ID:        userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	}
	if err = uh.userStore.UpdateUserDetails(user); err != nil {
		utils.ServerError(w, uh.logger, "update user details", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "user details updated successfully"})
}

// HandleUpdateAdminPassword godoc
// @Summary      Update admin password
// @Description  Updates the password of the admin with the given ID
// @Tags         admins
// @Accept       json
// @Produce      json
// @Param        id   path      string                       true "Admin ID"
// @Param        body body      models.UpdatePasswordRequest true "New password payload"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /admin/update/password/{id} [put]
func (uh *UserHandler) HandleUpdateAdminPassword(w http.ResponseWriter, r *http.Request) {
	adminID, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, uh.logger, err.Error(), err)
		return
	}
	var req models.UpdatePasswordRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, uh.logger, "invalid request payload", err)
		return
	}
	if err = validator.Validate(&req); err != nil {
		utils.BadRequest(w, uh.logger, err.Error(), err)
		return
	}
	admin := &models.Admin{ID: adminID}
	if err = admin.Password.Set(req.NewPassword); err != nil {
		utils.ServerError(w, uh.logger, "update admin password", err)
		return
	}
	if err = uh.userStore.UpdateAdminPassword(admin); err != nil {
		utils.ServerError(w, uh.logger, "update admin password", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "admin password updated successfully"})
}

// HandleUpdateUserPassword godoc
// @Summary      Update user password
// @Description  Updates the password of the user with the given ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string                       true "User ID"
// @Param        body body      models.UpdatePasswordRequest true "New password payload"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /user/update/password/{id} [put]
func (uh *UserHandler) HandleUpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, uh.logger, err.Error(), err)
		return
	}
	var req models.UpdatePasswordRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, uh.logger, "invalid request payload", err)
		return
	}
	if err = validator.Validate(&req); err != nil {
		utils.BadRequest(w, uh.logger, err.Error(), err)
		return
	}
	user := &models.User{ID: userID}
	if err = user.Password.Set(req.NewPassword); err != nil {
		utils.ServerError(w, uh.logger, "update user password", err)
		return
	}
	if err = uh.userStore.UpdateUserPassword(user); err != nil {
		utils.ServerError(w, uh.logger, "update user password", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "user password updated successfully"})
}

// HandleDeleteAdmin godoc
// @Summary      Delete admin
// @Description  Deletes the admin with the given ID
// @Tags         admins
// @Produce      json
// @Param        id path string true "Admin ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse "Admin not found"
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /admin/delete/{id} [delete]
func (uh *UserHandler) HandleDeleteAdmin(w http.ResponseWriter, r *http.Request) {
	adminID, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, uh.logger, err.Error(), err)
		return
	}
	if err = uh.userStore.DeleteAdmin(adminID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "admin not found"})
			return
		}
		utils.ServerError(w, uh.logger, "delete admin", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "admin deleted successfully"})
}

// HandleDeleteUser godoc
// @Summary      Delete user
// @Description  Deletes the user with the given ID
// @Tags         users
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse "User not found"
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /user/delete/{id} [delete]
func (uh *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, uh.logger, err.Error(), err)
		return
	}
	if err = uh.userStore.DeleteUser(userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "user not found"})
			return
		}
		utils.ServerError(w, uh.logger, "delete user", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "user deleted successfully"})
}

// HandleGetAllUsers godoc
// @Summary      Get all users
// @Description  Returns a paginated list of all registered users
// @Tags         users
// @Produce      json
// @Param        page  query int false "Page number (default: 1)"
// @Param        limit query int false "Items per page, max 100 (default: 20)"
// @Success      200 {object} map[string]interface{} "list of users"
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /user/get/all [get]
func (uh *UserHandler) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	pg := utils.ReadPaginationParams(r)
	users, err := uh.userStore.GetAllUsers(pg.Limit, pg.Offset())
	if err != nil {
		utils.ServerError(w, uh.logger, "get all users", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message": "users fetched successfully",
		"users":   users,
		"pagination": map[string]int{
			"page":  pg.Page,
			"limit": pg.Limit,
		},
	})
}

// HandleBlockUser godoc
// @Summary      Block or unblock a user
// @Description  Updates the block status of the user with the given ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string                        true "User ID"
// @Param        body body      models.BlockUserRequest       true "Block status payload"
// @Success      200  {object}  MessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Security     BearerAuth
// @Router       /user/block/{id} [put]
func (uh *UserHandler) HandleBlockUser(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, uh.logger, err.Error(), err)
		return
	}
	var req models.BlockUserRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, uh.logger, "invalid request payload", err)
		return
	}
	if err = uh.userStore.BlockUser(&models.User{ID: userID, IsUserBlocked: req.IsUserBlocked}); err != nil {
		utils.ServerError(w, uh.logger, "block user", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "user block status updated successfully"})
}

// HandleGetUserDetailsByID godoc
// @Summary      Get user details by ID
// @Description  Returns the details of the user with the given ID
// @Tags         users
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200 {object} map[string]interface{} "user details"
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse "User not found"
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /user/get/user/{id} [get]
func (uh *UserHandler) HandleGetUserDetailsByID(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, uh.logger, err.Error(), err)
		return
	}
	user, err := uh.userStore.GetUserDetailsByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "user not found"})
			return
		}
		utils.ServerError(w, uh.logger, "get user details", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "user details fetched successfully", "user": user})
}

// HandleGetAdminDetailsByID godoc
// @Summary      Get admin details by ID
// @Description  Returns the details of the admin with the given ID
// @Tags         admins
// @Produce      json
// @Param        id path string true "Admin ID"
// @Success      200 {object} map[string]interface{} "admin details"
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse "Admin not found"
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /admin/get/admin/{id} [get]
func (uh *UserHandler) HandleGetAdminDetailsByID(w http.ResponseWriter, r *http.Request) {
	adminID, err := utils.ReadParamID(r)
	if err != nil {
		utils.BadRequest(w, uh.logger, err.Error(), err)
		return
	}
	admin, err := uh.userStore.GetAdminDetailsByID(adminID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "admin not found"})
			return
		}
		utils.ServerError(w, uh.logger, "get admin details", err)
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "admin details fetched successfully", "admin": admin})
}
