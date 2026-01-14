package db

import (
	"air_router/models"
	"air_router/utils/common"
	"database/sql"
	"encoding/json"
	"fmt"
)

// ModelDB represents the database operations for models
type ModelDB struct {
	DB *sql.DB
}

// scanModels scans model rows from the database
func scanModels(rows *sql.Rows) ([]models.Model, error) {
	defer rows.Close()

	var modelsList []models.Model
	for rows.Next() {
		var model models.Model
		var assModelIDsJSON sql.NullString
		var provider string

		err := rows.Scan(&model.ID, &model.ModelID, &assModelIDsJSON, &provider, &model.Enabled, &model.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Parse associated model IDs from JSON
		if assModelIDsJSON.Valid {
			var assModelIDs []string
			if err := json.Unmarshal([]byte(assModelIDsJSON.String), &assModelIDs); err == nil {
				model.AssModelIDs = assModelIDs
			}
		}

		// Parse provider
		model.Provider = models.Provider(provider)

		modelsList = append(modelsList, model)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return modelsList, nil
}

// ModelIDExists checks if a model with the given model_id already exists
func (m *ModelDB) ModelIDExists(modelID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM models WHERE model_id = ?`
	err := m.DB.QueryRow(query, modelID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ModelIDExistsExcludingID checks if a model with the given model_id already exists, excluding the specified ID
func (m *ModelDB) ModelIDExistsExcludingID(modelID string, id int) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM models WHERE model_id = ? AND id != ?`
	err := m.DB.QueryRow(query, modelID, id).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CreateModel inserts a new model into the database
func (m *ModelDB) CreateModel(model models.Model) (int64, error) {
	// Check if a model with the same model_id already exists
	exists, err := m.ModelIDExists(model.ModelID)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, fmt.Errorf("model with model_id '%s' already exists", model.ModelID)
	}

	// Convert associated model IDs to JSON
	var assModelIDsJSON sql.NullString
	if len(model.AssModelIDs) > 0 {
		jsonData, err := json.Marshal(model.AssModelIDs)
		if err != nil {
			return 0, err
		}
		assModelIDsJSON = sql.NullString{String: string(jsonData), Valid: true}
	}

	query := `INSERT INTO models (model_id, ass_model_ids, provider, enabled, updated_at) VALUES (?, ?, ?, ?, ?)`
	result, err := m.DB.Exec(query, model.ModelID, assModelIDsJSON, string(model.Provider), model.Enabled, common.GetCurrentTimestamp())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetModels retrieves all models from the database
func (m *ModelDB) GetModels() ([]models.Model, error) {
	query := `SELECT id, model_id, ass_model_ids, provider, enabled, updated_at FROM models`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	return scanModels(rows)
}

// GetEnabledModels retrieves all enabled models from the database
func (m *ModelDB) GetEnabledModels() ([]models.Model, error) {
	query := `SELECT id, model_id, ass_model_ids, provider, enabled, updated_at FROM models WHERE enabled = 1`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	return scanModels(rows)
}

// GetEnabledModelsByProvider retrieves all enabled models for a specific provider from the database
func (m *ModelDB) GetEnabledModelsByProvider(provider models.Provider) ([]models.Model, error) {
	query := `SELECT id, model_id, ass_model_ids, provider, enabled, updated_at FROM models WHERE enabled = 1 AND provider = ?`
	rows, err := m.DB.Query(query, provider)
	if err != nil {
		return nil, err
	}
	return scanModels(rows)
}

// GetModel retrieves a specific model by ID
func (m *ModelDB) GetModel(id int) (models.Model, error) {
	return m.getModelByField("id", id)
}

// GetModelByModelID retrieves a specific model by model_id
func (m *ModelDB) GetModelByModelID(modelID string) (models.Model, error) {
	return m.getModelByField("model_id", modelID)
}

// getModelByField retrieves a specific model by field name and value
func (m *ModelDB) getModelByField(field string, value interface{}) (models.Model, error) {
	var model models.Model
	var assModelIDsJSON sql.NullString
	var provider string

	query := fmt.Sprintf(`SELECT id, model_id, ass_model_ids, provider, enabled, updated_at FROM models WHERE %s = ?`, field)
	err := m.DB.QueryRow(query, value).Scan(&model.ID, &model.ModelID, &assModelIDsJSON, &provider, &model.Enabled, &model.UpdatedAt)
	if err != nil {
		return model, err
	}

	// Parse associated model IDs from JSON
	if assModelIDsJSON.Valid {
		var assModelIDs []string
		if err := json.Unmarshal([]byte(assModelIDsJSON.String), &assModelIDs); err == nil {
			model.AssModelIDs = assModelIDs
		}
	}

	// Parse provider
	model.Provider = models.Provider(provider)

	return model, nil
}

// UpdateModel updates an existing model
func (m *ModelDB) UpdateModel(model models.Model) error {
	// Check if another model with the same model_id already exists
	exists, err := m.ModelIDExistsExcludingID(model.ModelID, model.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("model with model_id '%s' already exists", model.ModelID)
	}

	// Convert associated model IDs to JSON
	var assModelIDsJSON sql.NullString
	if len(model.AssModelIDs) > 0 {
		jsonData, err := json.Marshal(model.AssModelIDs)
		if err != nil {
			return err
		}
		assModelIDsJSON = sql.NullString{String: string(jsonData), Valid: true}
	}

	query := `UPDATE models SET model_id = ?, ass_model_ids = ?, provider = ?, enabled = ?, updated_at = ? WHERE id = ?`
	_, err = m.DB.Exec(query, model.ModelID, assModelIDsJSON, string(model.Provider), model.Enabled, common.GetCurrentTimestamp(), model.ID)
	return err
}

// DeleteModel deletes a model by ID
func (m *ModelDB) DeleteModel(id int) error {
	query := `DELETE FROM models WHERE id = ?`
	_, err := m.DB.Exec(query, id)
	return err
}

// ToggleModel toggles the enabled status of a model
func (m *ModelDB) ToggleModel(id int) error {
	// First get the current status
	var enabled bool
	query := `SELECT enabled FROM models WHERE id = ?`
	err := m.DB.QueryRow(query, id).Scan(&enabled)
	if err != nil {
		return err
	}

	// Toggle the status and update updated_at
	newEnabled := !enabled
	updateQuery := `UPDATE models SET enabled = ?, updated_at = ? WHERE id = ?`
	_, err = m.DB.Exec(updateQuery, newEnabled, common.GetCurrentTimestamp(), id)
	return err
}

// SearchModels searches for models by model_id or provider
func (m *ModelDB) SearchModels(search string) ([]models.Model, error) {
	query := `SELECT id, model_id, ass_model_ids, provider, enabled, updated_at FROM models WHERE model_id LIKE ? OR provider LIKE ?`
	searchPattern := "%" + search + "%"
	rows, err := m.DB.Query(query, searchPattern, searchPattern)
	if err != nil {
		return nil, err
	}
	return scanModels(rows)
}
