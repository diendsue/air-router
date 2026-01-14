package handlers

import (
	air_router_db "air_router/db"
	air_router_models "air_router/models"
	air_router_utils "air_router/utils/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ModelHandler struct {
	modelDB *air_router_db.ModelDB
}

func NewModelHandler(modelDB *air_router_db.ModelDB) *ModelHandler {
	return &ModelHandler{
		modelDB: modelDB,
	}
}

// GetModels retrieves all models
func (h *ModelHandler) GetModels(c *gin.Context) {
	models, err := h.modelDB.GetModels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve models"})
		return
	}

	// Ensure we return an empty array instead of null when no models exist
	if models == nil {
		models = []air_router_models.Model{}
	}

	c.JSON(http.StatusOK, models)
}

// GetModel retrieves a specific model by ID
func (h *ModelHandler) GetModel(c *gin.Context) {
	id, err := air_router_utils.ParseIDParam(c, "id")
	if err != nil {
		air_router_utils.SendAPIError(c, http.StatusBadRequest, air_router_utils.ErrMsgInvalidID, air_router_utils.ErrTypeInvalidRequest)
		return
	}

	model, err := h.modelDB.GetModel(id)
	if err != nil {
		air_router_utils.SendAPIError(c, http.StatusNotFound, air_router_utils.ErrMsgModelNotFound, air_router_utils.ErrTypeNotFound)
		return
	}

	air_router_utils.SendJSONResponse(c, http.StatusOK, model)
}

// CreateModel creates a new model
func (h *ModelHandler) CreateModel(c *gin.Context) {
	var model air_router_models.Model
	if err := c.ShouldBindJSON(&model); err != nil {
		air_router_utils.SendAPIError(c, http.StatusBadRequest, err.Error(), air_router_utils.ErrTypeInvalidRequest)
		return
	}

	// Validate provider using common function
	if !air_router_utils.ValidateModelProvider(model.Provider) {
		air_router_utils.SendAPIError(c, http.StatusBadRequest, air_router_utils.ErrMsgInvalidProvider, air_router_utils.ErrTypeInvalidProvider)
		return
	}

	// Set default enabled status
	if model.Enabled == false {
		model.Enabled = true
	}

	id, err := h.modelDB.CreateModel(model)
	if err != nil {
		air_router_utils.SendAPIError(c, http.StatusInternalServerError, err.Error(), air_router_utils.ErrTypeInternalServer)
		return
	}

	model.ID = int(id)
	air_router_utils.SendJSONResponse(c, http.StatusCreated, model)
}

// UpdateModel updates an existing model
func (h *ModelHandler) UpdateModel(c *gin.Context) {
	id, err := air_router_utils.ParseIDParam(c, "id")
	if err != nil {
		air_router_utils.SendAPIError(c, http.StatusBadRequest, air_router_utils.ErrMsgInvalidID, air_router_utils.ErrTypeInvalidRequest)
		return
	}

	var model air_router_models.Model
	if err := c.ShouldBindJSON(&model); err != nil {
		air_router_utils.SendAPIError(c, http.StatusBadRequest, err.Error(), air_router_utils.ErrTypeInvalidRequest)
		return
	}

	// Validate provider using common function
	if !air_router_utils.ValidateModelProvider(model.Provider) {
		air_router_utils.SendAPIError(c, http.StatusBadRequest, air_router_utils.ErrMsgInvalidProvider, air_router_utils.ErrTypeInvalidProvider)
		return
	}

	model.ID = id
	err = h.modelDB.UpdateModel(model)
	if err != nil {
		air_router_utils.SendAPIError(c, http.StatusInternalServerError, err.Error(), air_router_utils.ErrTypeInternalServer)
		return
	}

	air_router_utils.SendJSONResponse(c, http.StatusOK, model)
}

// DeleteModel deletes a model by ID
func (h *ModelHandler) DeleteModel(c *gin.Context) {
	id, err := air_router_utils.ParseIDParam(c, "id")
	if err != nil {
		air_router_utils.SendAPIError(c, http.StatusBadRequest, air_router_utils.ErrMsgInvalidID, air_router_utils.ErrTypeInvalidRequest)
		return
	}

	err = h.modelDB.DeleteModel(id)
	if err != nil {
		air_router_utils.SendAPIError(c, http.StatusInternalServerError, air_router_utils.ErrMsgFailedToDelete, air_router_utils.ErrTypeInternalServer)
		return
	}

	air_router_utils.SendJSONResponse(c, http.StatusOK, gin.H{"message": "Model deleted successfully"})
}

// ToggleModel toggles the enabled status of a model
func (h *ModelHandler) ToggleModel(c *gin.Context) {
	id, err := air_router_utils.ParseIDParam(c, "id")
	if err != nil {
		air_router_utils.SendAPIError(c, http.StatusBadRequest, air_router_utils.ErrMsgInvalidID, air_router_utils.ErrTypeInvalidRequest)
		return
	}

	err = h.modelDB.ToggleModel(id)
	if err != nil {
		air_router_utils.SendAPIError(c, http.StatusInternalServerError, air_router_utils.ErrMsgFailedToToggle, air_router_utils.ErrTypeInternalServer)
		return
	}

	// Get the updated model
	model, err := h.modelDB.GetModel(id)
	if err != nil {
		air_router_utils.SendAPIError(c, http.StatusInternalServerError, air_router_utils.ErrMsgFailedToUpdate, air_router_utils.ErrTypeInternalServer)
		return
	}

	air_router_utils.SendJSONResponse(c, http.StatusOK, model)
}

// SearchModels searches for models by model_id or provider
func (h *ModelHandler) SearchModels(c *gin.Context) {
	search := c.Query("search")
	if search == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query parameter is required"})
		return
	}

	models, err := h.modelDB.SearchModels(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search models"})
		return
	}

	c.JSON(http.StatusOK, models)
}
