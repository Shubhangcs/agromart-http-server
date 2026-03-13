package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/shubhangcs/agromart-server/internal/models"
	"github.com/shubhangcs/agromart-server/internal/store"
	"github.com/shubhangcs/agromart-server/internal/utils"
)

// FollowerHandler handles all follower/following HTTP requests.
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

func (fh *FollowerHandler) validateFollowRequest(req *models.FollowRequest) error {
	if req.UserID == "" {
		return errors.New("user_id is required")
	}
	if req.BusinessID == "" {
		return errors.New("business_id is required")
	}
	return nil
}

// HandleCreateFollower godoc
// @Summary      Follow a business
// @Description  Creates a follower relationship between a user and a business (idempotent)
// @Tags         followers
// @Accept       json
// @Produce      json
// @Param        body body models.FollowRequest true "Follow payload"
// @Success      201 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /follower/follow [post]
func (fh *FollowerHandler) HandleCreateFollower(w http.ResponseWriter, r *http.Request) {
	var req models.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fh.logger.Printf("ERROR: create follower: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err := fh.validateFollowRequest(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if err := fh.followerStore.CreateFollower(&models.Follower{
		UserID:     req.UserID,
		BusinessID: req.BusinessID,
	}); err != nil {
		fh.logger.Printf("ERROR: create follower: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"message": "followed successfully"})
}

// HandleRemoveFollower godoc
// @Summary      Unfollow a business
// @Description  Removes the follower relationship between a user and a business
// @Tags         followers
// @Accept       json
// @Produce      json
// @Param        body body models.FollowRequest true "Unfollow payload"
// @Success      200 {object} MessageResponse
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /follower/unfollow [post]
func (fh *FollowerHandler) HandleRemoveFollower(w http.ResponseWriter, r *http.Request) {
	var req models.FollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fh.logger.Printf("ERROR: remove follower: %v\n", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}
	if err := fh.validateFollowRequest(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	if err := fh.followerStore.RemoveFollower(&models.Follower{
		UserID:     req.UserID,
		BusinessID: req.BusinessID,
	}); err != nil {
		fh.logger.Printf("ERROR: remove follower: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "unfollowed successfully"})
}

// HandleGetFollowersCount godoc
// @Summary      Get followers count
// @Description  Returns the total number of followers for a given business
// @Tags         followers
// @Produce      json
// @Param        id path string true "Business ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /follower/get/followers/count/{id} [get]
func (fh *FollowerHandler) HandleGetFollowersCount(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	count, err := fh.followerStore.GetFollowersCount(id)
	if err != nil {
		fh.logger.Printf("ERROR: get followers count: %v\n", err)
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
// @Param        id path string true "User ID"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /follower/get/following/count/{id} [get]
func (fh *FollowerHandler) HandleGetFollowingCount(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	count, err := fh.followerStore.GetFollowingCount(id)
	if err != nil {
		fh.logger.Printf("ERROR: get following count: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "following count fetched successfully", "following_count": count})
}

// HandleGetAllFollowers godoc
// @Summary      Get all followers
// @Description  Returns a paginated list of all users following the given business
// @Tags         followers
// @Produce      json
// @Param        id    path  string true  "Business ID"
// @Param        page  query int    false "Page number (default 1)"
// @Param        limit query int    false "Items per page (default 20, max 100)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /follower/get/followers/{id} [get]
func (fh *FollowerHandler) HandleGetAllFollowers(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	pg := utils.ReadPaginationParams(r)
	res, err := fh.followerStore.GetAllFollowers(id, pg.Limit, pg.Offset())
	if err != nil {
		fh.logger.Printf("ERROR: get all followers: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":    "followers fetched successfully",
		"followers":  res,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}

// HandleGetAllFollowing godoc
// @Summary      Get all followings
// @Description  Returns a paginated list of all businesses a user is following
// @Tags         followers
// @Produce      json
// @Param        id    path  string true  "User ID"
// @Param        page  query int    false "Page number (default 1)"
// @Param        limit query int    false "Items per page (default 20, max 100)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /follower/get/followings/{id} [get]
func (fh *FollowerHandler) HandleGetAllFollowing(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadParamID(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}
	pg := utils.ReadPaginationParams(r)
	res, err := fh.followerStore.GetAllFollowing(id, pg.Limit, pg.Offset())
	if err != nil {
		fh.logger.Printf("ERROR: get all following: %v\n", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{
		"message":    "followings fetched successfully",
		"followings": res,
		"pagination": map[string]int{"page": pg.Page, "limit": pg.Limit},
	})
}
