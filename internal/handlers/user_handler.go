package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/shubhangcs/agromart-server/internal/models"
	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

var (
	emailRegx    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	passwordRegx = regexp.MustCompile(`^[A-Za-z\d@$!%*?&]{8,}$`)
	phoneRegx    = regexp.MustCompile(`^\d{7,15}$`)
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
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{userStore: userStore, logger: logger}
}

func (uh *UserHandler) validateRegisterRequest(req *models.CreateAdminRequest) error {
	if req.FirstName == "" {
		return errors.New("first name is required")
	}
	if !emailRegx.MatchString(req.Email) {
		return errors.New("invalid email address")
	}
	if !passwordRegx.MatchString(req.Password) {
		return errors.New("password must be at least 8 characters and contain letters, digits, or @$!%*?&")
	}
	if !phoneRegx.MatchString(req.Phone) {
		return errors.New("phone number must be 7-15 digits")
	}
	return nil
}

func (uh *UserHandler) validateUpdateRequest(req *models.UpdateUserDetailsRequest) error {
	if req.Email != "" && !emailRegx.MatchString(req.Email) {
		return errors.New("invalid email address")
	}
	if req.Phone != "" && !phoneRegx.MatchString(req.Phone) {
		return errors.New("phone number must be 7-15 digits")
	}
	return nil
}

func (uh *UserHandler) validateUpdatePasswordRequest(req *models.UpdatePasswordRequest) error {
	if !passwordRegx.MatchString(req.NewPassword) {
		return errors.New("password must be at least 8 characters and contain letters, digits, or @$!%*?&")
	}
	return nil
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
		uh.logger.Printf("ERROR: create admin: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err := uh.validateRegisterRequest(&req); err != nil {
		uh.logger.Printf("ERROR: create admin: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	admin := &models.Admin{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	}
	if err := admin.Password.Set(req.Password); err != nil {
		uh.logger.Printf("ERROR: create admin: hash password: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if err := uh.userStore.CreateAdmin(admin); err != nil {
		uh.logger.Printf("ERROR: create admin: %v\n", err)
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
		uh.logger.Printf("ERROR: create user: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err := uh.validateRegisterRequest(&req); err != nil {
		uh.logger.Printf("ERROR: create user: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	}
	if err := user.Password.Set(req.Password); err != nil {
		uh.logger.Printf("ERROR: create user: hash password: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if err := uh.userStore.CreateUser(user); err != nil {
		uh.logger.Printf("ERROR: create user: %v\n", err)
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
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req models.UpdateUserDetailsRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err = uh.validateUpdateRequest(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
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
		uh.logger.Printf("ERROR: update admin details: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
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
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req models.UpdateUserDetailsRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err = uh.validateUpdateRequest(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
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
		uh.logger.Printf("ERROR: update user details: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
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
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req models.UpdatePasswordRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err = uh.validateUpdatePasswordRequest(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	admin := &models.Admin{ID: adminID}
	if err = admin.Password.Set(req.NewPassword); err != nil {
		uh.logger.Printf("ERROR: update admin password: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if err = uh.userStore.UpdateAdminPassword(admin); err != nil {
		uh.logger.Printf("ERROR: update admin password: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
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
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req models.UpdatePasswordRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err = uh.validateUpdatePasswordRequest(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	user := &models.User{ID: userID}
	if err = user.Password.Set(req.NewPassword); err != nil {
		uh.logger.Printf("ERROR: update user password: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	if err = uh.userStore.UpdateUserPassword(user); err != nil {
		uh.logger.Printf("ERROR: update user password: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
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
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if err = uh.userStore.DeleteAdmin(adminID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "admin not found"})
			return
		}
		uh.logger.Printf("ERROR: delete admin: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
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
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if err = uh.userStore.DeleteUser(userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "user not found"})
			return
		}
		uh.logger.Printf("ERROR: delete user: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
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
		uh.logger.Printf("ERROR: get all users: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
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
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req models.BlockUserRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err = uh.userStore.BlockUser(&models.User{ID: userID, IsUserBlocked: req.IsUserBlocked}); err != nil {
		uh.logger.Printf("ERROR: block user: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
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
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	user, err := uh.userStore.GetUserDetailsByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "user not found"})
			return
		}
		uh.logger.Printf("ERROR: get user details: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
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
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	admin, err := uh.userStore.GetAdminDetailsByID(adminID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "admin not found"})
			return
		}
		uh.logger.Printf("ERROR: get admin details: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "admin details fetched successfully", "admin": admin})
}
