package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

// followerRequest represents the follow/unfollow payload
type followerRequest struct {
	UserID     string    `json:"user_id"     example:"user-uuid-001"`
	BusinessID string    `json:"business_id" example:"biz-uuid-001"`
	CreatedAT  time.Time `json:"created_at"`
}

type FollowerHandler struct {
	followerStore store.FollowerStore
	logger        *log.Logger
}

func NewFollowerHandler(followerStore store.FollowerStore, logger *log.Logger) *FollowerHandler {
	return &FollowerHandler{
		followerStore: followerStore,
		logger:        logger,
	}
}

func (fh *FollowerHandler) validateCreateAndRemoveFollowerRequest(req *followerRequest) error {
	if req.UserID == "" {
		return errors.New("invalid request user id is required")
	}
	if req.BusinessID == "" {
		return errors.New("invalid request business id is required")
	}
	return nil
}

// HandleCreateFollower godoc
// @Summary      Follow a business
// @Description  Creates a follower relationship between a user and a business
// @Tags         followers
// @Accept       json
// @Produce      json
// @Param        body body followerRequest true "Follow payload"
// @Success      201 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid payload or missing fields"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /follower/follow [post]
func (fh *FollowerHandler) HandleCreateFollower(w http.ResponseWriter, r *http.Request) {
	var req followerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fh.logger.Printf("ERROR: create follower: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	err = fh.validateCreateAndRemoveFollowerRequest(&req)
	if err != nil {
		fh.logger.Printf("ERROR: create follower: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	follower := &store.Follower{
		UserID:     req.UserID,
		BusinessID: req.BusinessID,
	}

	err = fh.followerStore.CreateFollower(follower)
	if err != nil {
		fh.logger.Printf("ERROR: create follower: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "follow successfull"})
}

// HandleRemoveFollower godoc
// @Summary      Unfollow a business
// @Description  Removes the follower relationship between a user and a business
// @Tags         followers
// @Accept       json
// @Produce      json
// @Param        body body followerRequest true "Unfollow payload"
// @Success      201 {object} MessageResponse
// @Failure      400 {object} ErrorResponse "Invalid payload or missing fields"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /follower/unfollow [post]
func (fh *FollowerHandler) HandleRemoveFollower(w http.ResponseWriter, r *http.Request) {
	var req followerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fh.logger.Printf("ERROR: remove follower: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	err = fh.validateCreateAndRemoveFollowerRequest(&req)
	if err != nil {
		fh.logger.Printf("ERROR: remove follower: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	follower := &store.Follower{
		UserID:     req.UserID,
		BusinessID: req.BusinessID,
	}

	err = fh.followerStore.RemoveFollower(follower)
	if err != nil {
		fh.logger.Printf("ERROR: remove follower: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "unfollow successfull"})
}

// HandleGetFollowersCount godoc
// @Summary      Get followers count
// @Description  Returns the total number of followers for a given business
// @Tags         followers
// @Produce      json
// @Param        id path int true "Business ID"
// @Success      200 {object} map[string]interface{} "followers count"
// @Failure      400 {object} ErrorResponse "Invalid or missing business ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /follower/get/followers/count/{id} [get]
func (fh *FollowerHandler) HandleGetFollowersCount(w http.ResponseWriter, r *http.Request) {
	businessId, err := utils.ReadParamID(r)
	if err != nil {
		fh.logger.Printf("ERROR: get follower count: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	if businessId == "" {
		fh.logger.Println("ERROR: get follower count: empty business id")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request business id is required"})
		return
	}

	count, err := fh.followerStore.GetFollowersCount(businessId)
	if err != nil {
		fh.logger.Printf("ERROR: get follower count: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "follower count fetched successfully", "followers_count": count})
}

// HandleGetFollowingCount godoc
// @Summary      Get following count
// @Description  Returns the total number of businesses a user is following
// @Tags         followers
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} map[string]interface{} "following count"
// @Failure      400 {object} ErrorResponse "Invalid or missing user ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /follower/get/following/count/{id} [get]
func (fh *FollowerHandler) HandleGetFollowingCount(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadParamID(r)
	if err != nil {
		fh.logger.Printf("ERROR: get following count: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	if userId == "" {
		fh.logger.Println("ERROR: get following count: empty user id")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request user id is required"})
		return
	}

	count, err := fh.followerStore.GetFollowingCount(userId)
	if err != nil {
		fh.logger.Printf("ERROR: get following count: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "following count fetched successfully", "following_count": count})
}

// HandleGetAllFollowers godoc
// @Summary      Get all followers
// @Description  Returns a list of all users following the given business
// @Tags         followers
// @Produce      json
// @Param        id path int true "Business ID"
// @Success      200 {object} map[string]interface{} "list of followers"
// @Failure      400 {object} ErrorResponse "Invalid or missing business ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /follower/get/followers/{id} [get]
func (fh *FollowerHandler) HandleGetAllFollowers(w http.ResponseWriter, r *http.Request) {
	businessId, err := utils.ReadParamID(r)
	if err != nil {
		fh.logger.Printf("ERROR: get all followers: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	if businessId == "" {
		fh.logger.Println("ERROR: get all followers: empty business id")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request business id is required"})
		return
	}

	res, err := fh.followerStore.GetAllFollowers(businessId)
	if err != nil {
		fh.logger.Printf("ERROR: get all followers: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "followers fetched successfully", "followers": res})
}

// HandleGetAllFollowing godoc
// @Summary      Get all followings
// @Description  Returns a list of all businesses a user is following
// @Tags         followers
// @Produce      json
// @Param        id path int true "User ID"
// @Success      200 {object} map[string]interface{} "list of followings"
// @Failure      400 {object} ErrorResponse "Invalid or missing user ID"
// @Failure      500 {object} ErrorResponse "Internal server error"
// @Router       /follower/get/followings/{id} [get]
func (fh *FollowerHandler) HandleGetAllFollowing(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.ReadParamID(r)
	if err != nil {
		fh.logger.Printf("ERROR: get all following: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	if userId == "" {
		fh.logger.Println("ERROR: get all following: empty business id")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request user id is required"})
		return
	}

	res, err := fh.followerStore.GetAllFollowing(userId)
	if err != nil {
		fh.logger.Printf("ERROR: get all following: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "followings fetched successfully", "followings": res})
}