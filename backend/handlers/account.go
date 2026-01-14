package handlers

import (
	"database/sql"
	"net/http"

	"air_router/db"
	"air_router/models"
	"air_router/utils"
	"air_router/utils/common"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	AccountDB *db.AccountDB
}

func NewAccountHandler(accountDB *db.AccountDB) *AccountHandler {
	return &AccountHandler{
		AccountDB: accountDB,
	}
}

// GetAccounts handles GET /api/accounts with pagination and search support
func (h *AccountHandler) GetAccounts(c *gin.Context) {
	// Parse pagination parameters
	params := utils.ParsePaginationParams(c)

	accounts, total, err := h.AccountDB.GetPaginatedAccounts(params.Page, params.PageSize, params.Search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, utils.BuildPaginatedResponse(accounts, total, params.Page, params.PageSize, params.Search))
}

// CreateAccount handles POST /api/accounts
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var account models.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters: " + err.Error()})
		return
	}

	if !account.Enabled && account.ID == 0 {
		account.Enabled = true
	}

	id, err := h.AccountDB.CreateAccount(account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account.ID = int(id)
	c.JSON(http.StatusCreated, account)
}

// GetAccount handles GET /api/accounts/:id
func (h *AccountHandler) GetAccount(c *gin.Context) {
	id, err := common.ParseIDParam(c, "id")
	if err != nil {
		common.SendAPIError(c, http.StatusBadRequest, common.ErrMsgInvalidID, common.ErrTypeInvalidRequest)
		return
	}

	account, err := h.AccountDB.GetAccount(id)
	if err != nil {
		if err == sql.ErrNoRows {
			common.SendAPIError(c, http.StatusNotFound, common.ErrMsgAccountNotFound, common.ErrTypeNotFound)
		} else {
			common.SendAPIError(c, http.StatusInternalServerError, err.Error(), common.ErrTypeInternalServer)
		}
		return
	}

	common.SendJSONResponse(c, http.StatusOK, account)
}

// UpdateAccount handles PUT /api/accounts/:id
func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	id, err := common.ParseIDParam(c, "id")
	if err != nil {
		common.SendAPIError(c, http.StatusBadRequest, common.ErrMsgInvalidID, common.ErrTypeInvalidRequest)
		return
	}

	var account models.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		common.SendAPIError(c, http.StatusBadRequest, "Invalid parameters: "+err.Error(), common.ErrTypeInvalidRequest)
		return
	}

	account.ID = id
	if err := h.AccountDB.UpdateAccount(account); err != nil {
		common.SendAPIError(c, http.StatusBadRequest, err.Error(), common.ErrTypeBadRequest)
		return
	}

	common.SendJSONResponse(c, http.StatusOK, account)
}

// DeleteAccount handles DELETE /api/accounts/:id
func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	id, err := common.ParseIDParam(c, "id")
	if err != nil {
		common.SendAPIError(c, http.StatusBadRequest, common.ErrMsgInvalidID, common.ErrTypeInvalidRequest)
		return
	}

	if err := h.AccountDB.DeleteAccount(id); err != nil {
		common.SendAPIError(c, http.StatusInternalServerError, err.Error(), common.ErrTypeInternalServer)
		return
	}

	c.Status(http.StatusNoContent)
}

// ToggleAccount handles PATCH /api/accounts/:id
func (h *AccountHandler) ToggleAccount(c *gin.Context) {
	id, err := common.ParseIDParam(c, "id")
	if err != nil {
		common.SendAPIError(c, http.StatusBadRequest, common.ErrMsgInvalidID, common.ErrTypeInvalidRequest)
		return
	}

	if err := h.AccountDB.ToggleAccount(id); err != nil {
		if err == sql.ErrNoRows {
			common.SendAPIError(c, http.StatusNotFound, common.ErrMsgAccountNotFound, common.ErrTypeNotFound)
		} else {
			common.SendAPIError(c, http.StatusInternalServerError, err.Error(), common.ErrTypeInternalServer)
		}
		return
	}

	account, err := h.AccountDB.GetAccount(id)
	if err != nil {
		common.SendAPIError(c, http.StatusInternalServerError, err.Error(), common.ErrTypeInternalServer)
		return
	}

	common.SendJSONResponse(c, http.StatusOK, account)
}
