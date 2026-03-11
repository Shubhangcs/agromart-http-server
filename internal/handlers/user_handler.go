package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

var (
	emailRegx    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	passwordRegx = regexp.MustCompile(`^[A-Za-z\d@$!%*?&]{8,}$`)
	phoneRegx    = regexp.MustCompile(`^\d{7,15}$`)
)

// registerUserRequest represents the payload for registering a user or admin
// @Description Registration request body
type registerUserRequest struct {
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name"  example:"Doe"`
	Email     string `json:"email"      example:"john@example.com"`
	Phone     string `json:"phone"      example:"9876543210"`
	Password  string `json:"password"   example:"StrongPass@1"`
}

// updateUserProfileDetailsRequest represents updatable profile fields
type updateUserProfileDetailsRequest struct {
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name"  example:"Doe"`
	Email     string `json:"email"      example:"john@example.com"`
	Phone     string `json:"phone"      example:"9876543210"`
}

type updatePasswordRequest struct {
	Password string `json:"password" example:"NewPass@123"`
}

type blockUserRequest struct {
	IsUserBlocked bool `json:"is_user_blocked" example:"true"`
}

// ErrorResponse is a generic error response
type ErrorResponse struct {
	Error string `json:"error" example:"internal server error"`
}

// MessageResponse is a generic success message response
type MessageResponse struct {
	Message string `json:"message" example:"operation successful"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

// --- Validation helpers (unchanged) ---

func (uh *UserHandler) validateRegisterRequest(req *registerUserRequest) error {
	if req.FirstName == "" {
		return errors.New("invalid request first name is required")
	}
	if !emailRegx.Match([]byte(req.Email)) {
		return errors.New("invalid request incorrect email format")
	}
	if !passwordRegx.Match([]byte(req.Password)) {
		return errors.New("invalid request password is not strong")
	}
	if !phoneRegx.Match([]byte(req.Phone)) {
		return errors.New("invalid request incorrect phone number")
	}
	return nil
}

func (uh *UserHandler) validateUpdateRequest(req *updateUserProfileDetailsRequest) error {
	if req.Email != "" && !emailRegx.Match([]byte(req.Email)) {
		return errors.New("invalid request incorrect email format")
	}
	if req.Phone != "" && !phoneRegx.Match([]byte(req.Phone)) {
		return errors.New("invalid request incorrect phone number")
	}
	return nil
}

func (uh *UserHandler) validateUpdatePasswordRequest(req *updatePasswordRequest) error {
	if !passwordRegx.Match([]byte(req.Password)) {
		return errors.New("invalid request password is not strong")
	}
	return nil
}

// HandleCreateAdmin godoc
// @Summary      Create a new admin
// @Description  Registers a new admin account
// @Tags         admins
// @Accept       json
// @Produce      json
// @Param        body body registerUserRequest true "Admin registration payload"
// @Success      201 {object} map[string]interface{} "admin created successfully"
// @Failure      500 {object} ErrorResponse
// @Router       /admin/create [post]
func (uh *UserHandler) HandleCreateAdmin(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("ERROR: decoding admin register request: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "invalid request payload"})
		return
	}
	err = uh.validateRegisterRequest(&req)
	if err != nil {
		uh.logger.Printf("ERROR: validating admin register request: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": err.Error()})
		return
	}
	admin := &store.Admin{
		FirstName: req.FirstName,
		LastName:  &req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	}
	err = admin.Password.Set(req.Password)
	if err != nil {
		uh.logger.Printf("ERROR: setting password in admin register request: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	err = uh.userStore.CreateAdmin(admin)
	if err != nil {
		uh.logger.Printf("ERROR: creating admin in admin register request: %v\n", err)
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
// @Param        body body registerUserRequest true "User registration payload"
// @Success      201 {object} map[string]interface{} "user created successfully"
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /user/create [post]
func (uh *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("ERROR: decoding user register request: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	err = uh.validateRegisterRequest(&req)
	if err != nil {
		uh.logger.Printf("ERROR: validating user register request: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	user := &store.User{
		FirstName: req.FirstName,
		LastName:  &req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	}
	err = user.Password.Set(req.Password)
	if err != nil {
		uh.logger.Printf("ERROR: setting password in user register request: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	err = uh.userStore.CreateUser(user)
	if err != nil {
		uh.logger.Printf("ERROR: creating user in user register request: %v\n", err)
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
// @Param        id   path int                            true "Admin ID"
// @Param        body body updateUserProfileDetailsRequest true "Admin update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /admin/update/details/{id} [put]
func (uh *UserHandler) HandleUpdateAdminDetails(w http.ResponseWriter, r *http.Request) {
	adminId, err := utils.ReadParamID(r)
	if err != nil {
		uh.logger.Printf("ERROR: update admin profile details: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req updateUserProfileDetailsRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("ERROR: update admin profile details: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	err = uh.validateUpdateRequest(&req)
	if err != nil {
		uh.logger.Printf("ERROR: update admin profile details: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	admin := &store.Admin{
		ID:        adminId,
		FirstName: req.FirstName,
		LastName:  &req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	}
	err = uh.userStore.UpdateAdminDetails(admin)
	if err != nil {
		uh.logger.Printf("ERROR: update admin profile details: %v\n", err)
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
// @Param        id   path int                            true "User ID"
// @Param        body body updateUserProfileDetailsRequest true "User update payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /user/update/details/{id} [put]
func (uh *UserHandler) HandleUpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadParamID(r)
	if err != nil {
		uh.logger.Printf("ERROR: update user profile details: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req updateUserProfileDetailsRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("ERROR: update user profile details: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	err = uh.validateUpdateRequest(&req)
	if err != nil {
		uh.logger.Printf("ERROR: update user profile details: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	user := &store.User{
		ID:        userId,
		FirstName: req.FirstName,
		LastName:  &req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
	}
	err = uh.userStore.UpdateUserDetails(user)
	if err != nil {
		uh.logger.Printf("ERROR: update user profile details: %v\n", err)
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
// @Param        id   path int                   true "Admin ID"
// @Param        body body updatePasswordRequest true "New password payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /admin/update/password/{id} [put]
func (uh *UserHandler) HandleUpdateAdminPassword(w http.ResponseWriter, r *http.Request) {
	adminId, err := utils.ReadParamID(r)
	if err != nil {
		uh.logger.Printf("ERROR: update admin password: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req updatePasswordRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("ERROR: update admin password: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	err = uh.validateUpdatePasswordRequest(&req)
	if err != nil {
		uh.logger.Printf("ERROR: update admin password: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	admin := &store.Admin{ID: adminId}
	err = admin.Password.Set(req.Password)
	if err != nil {
		uh.logger.Printf("ERROR: update admin password: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	err = uh.userStore.UpdateAdminPassword(admin)
	if err != nil {
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
// @Param        id   path int                   true "User ID"
// @Param        body body updatePasswordRequest true "New password payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /user/update/password/{id} [put]
func (uh *UserHandler) HandleUpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadParamID(r)
	if err != nil {
		uh.logger.Printf("ERROR: update user password: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req updatePasswordRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("ERROR: update user password: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	err = uh.validateUpdatePasswordRequest(&req)
	if err != nil {
		uh.logger.Printf("ERROR: update user password: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	user := &store.User{ID: userId}
	err = user.Password.Set(req.Password)
	if err != nil {
		uh.logger.Printf("ERROR: update user password: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	err = uh.userStore.UpdateUserPassword(user)
	if err != nil {
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
// @Param        id path int true "Admin ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /admin/delete/{id} [delete]
func (uh *UserHandler) HandleDeleteAdmin(w http.ResponseWriter, r *http.Request) {
	adminId, err := utils.ReadParamID(r)
	if err != nil {
		uh.logger.Printf("ERROR: delete admin: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	err = uh.userStore.DeleteAdmin(adminId)
	if err != nil {
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
// @Param        id path int true "User ID"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /user/delete/{id} [delete]
func (uh *UserHandler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadParamID(r)
	if err != nil {
		uh.logger.Printf("ERROR: delete user: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	err = uh.userStore.DeleteUser(userId)
	if err != nil {
		uh.logger.Printf("ERROR: delete user: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "user deleted successfully"})
}

// HandleGetAllUsers godoc
// @Summary      Get all users
// @Description  Returns a list of all registered users
// @Tags         users
// @Produce      json
// @Success      200 {object} map[string]interface{} "list of users"
// @Failure      500 {object} ErrorResponse
// @Router       /user/get/all [get]
func (uh *UserHandler) HandleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	res, err := uh.userStore.GetAllUsers()
	if err != nil {
		uh.logger.Printf("ERROR: get all users: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"users": res})
}

// HandleBlockUser godoc
// @Summary      Block or unblock a user
// @Description  Updates the block status of the user with the given ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path int              true "User ID"
// @Param        body body blockUserRequest true "Block status payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /user/block/{id} [put]
func (uh *UserHandler) HandleBlockUser(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadParamID(r)
	if err != nil {
		uh.logger.Printf("ERROR: block users: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	var req blockUserRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("ERROR: block user: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request body"})
		return
	}
	user := &store.User{
		ID:            userId,
		IsUserBlocked: req.IsUserBlocked,
	}
	err = uh.userStore.BlockUser(user)
	if err != nil {
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
// @Param        id path int true "User ID"
// @Success      200 {object} map[string]interface{} "user details"
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /user/get/user/{id} [get]
func (uh *UserHandler) HandleGetUserDetailsByID(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadParamID(r)
	if err != nil {
		uh.logger.Printf("ERROR: get user details by id: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	res, err := uh.userStore.GetUserDetailsByID(userId)
	if err != nil {
		uh.logger.Printf("ERROR: get user details by id: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "user details fetched successfully", "user": res})
}
