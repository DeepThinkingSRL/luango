package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AgentMode representa el modo del agente generativo
type AgentMode string

const (
	ModeInteractive AgentMode = "interactive" // Modo interactivo con confirmación
	ModeAutomatic   AgentMode = "automatic"   // Modo automático sin confirmación
	ModeManual      AgentMode = "manual"      // Modo manual sin generación
)

// ResourceType define los tipos de recursos que se pueden generar
type ResourceType string

const (
	ResourceSprite     ResourceType = "sprite"
	ResourceSound      ResourceType = "sound"
	ResourceScript     ResourceType = "script"
	ResourceLevel      ResourceType = "level"
	ResourceEntity     ResourceType = "entity"
	ResourceBehavior   ResourceType = "behavior"
	ResourceDialogue   ResourceType = "dialogue"
	ResourceTerrain    ResourceType = "terrain"
	ResourceAnimation  ResourceType = "animation"
)

// GenerationRequest representa una solicitud de generación
type GenerationRequest struct {
	Type        ResourceType `json:"type"`
	Prompt      string       `json:"prompt"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	OutputPath  string       `json:"output_path,omitempty"`
	AutoApply   bool         `json:"auto_apply"`
	Timestamp   time.Time    `json:"timestamp"`
}

// GenerationResult representa el resultado de una generación
type GenerationResult struct {
	ID          string                 `json:"id"`
	Request     GenerationRequest      `json:"request"`
	Status      string                 `json:"status"` // pending, completed, failed, applied
	Content     interface{}            `json:"content"`
	FilePath    string                 `json:"file_path,omitempty"`
	Preview     string                 `json:"preview,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	AppliedAt   *time.Time             `json:"applied_at,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// GenerativeAgent es el núcleo del sistema de generación
type GenerativeAgent struct {
	mode           AgentMode
	projectPath    string
	generators     map[ResourceType]ResourceGenerator
	history        []GenerationResult
	pendingResults map[string]*GenerationResult
	callbacks      map[string]func(*GenerationResult)
}

// ResourceGenerator interfaz para diferentes tipos de generadores
type ResourceGenerator interface {
	Generate(ctx context.Context, request GenerationRequest) (*GenerationResult, error)
	Preview(content interface{}) (string, error)
	Apply(result *GenerationResult) error
	Validate(content interface{}) error
}

// NewGenerativeAgent crea una nueva instancia del agente generativo
func NewGenerativeAgent(projectPath string, mode AgentMode) *GenerativeAgent {
	return &GenerativeAgent{
		mode:           mode,
		projectPath:    projectPath,
		generators:     make(map[ResourceType]ResourceGenerator),
		history:        make([]GenerationResult, 0),
		pendingResults: make(map[string]*GenerationResult),
		callbacks:      make(map[string]func(*GenerationResult)),
	}
}

// RegisterGenerator registra un generador para un tipo de recurso específico
func (ga *GenerativeAgent) RegisterGenerator(resourceType ResourceType, generator ResourceGenerator) {
	ga.generators[resourceType] = generator
}

// SetMode cambia el modo del agente
func (ga *GenerativeAgent) SetMode(mode AgentMode) {
	ga.mode = mode
	fmt.Printf("[Agent] Mode changed to: %s\n", mode)
}

// GenerateResource solicita la generación de un recurso
func (ga *GenerativeAgent) GenerateResource(ctx context.Context, request GenerationRequest) (*GenerationResult, error) {
	// Validar que existe un generador para este tipo
	generator, exists := ga.generators[request.Type]
	if !exists {
		return nil, fmt.Errorf("no generator registered for resource type: %s", request.Type)
	}

	// Crear ID único para el resultado
	resultID := fmt.Sprintf("%s_%d", request.Type, time.Now().Unix())
	
	fmt.Printf("[Agent] Generating %s: %s\n", request.Type, request.Prompt)

	// Generar el recurso
	result, err := generator.Generate(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("generation failed: %w", err)
	}

	result.ID = resultID
	result.Request = request
	result.CreatedAt = time.Now()
	result.Status = "completed"

	// Manejar según el modo
	switch ga.mode {
	case ModeAutomatic:
		if request.AutoApply {
			err = ga.ApplyResult(result)
			if err != nil {
				result.Status = "failed"
				result.Error = err.Error()
			}
		}
	case ModeInteractive:
		ga.pendingResults[resultID] = result
		result.Status = "pending"
		fmt.Printf("[Agent] Resource generated. Use ApplyResult(\"%s\") to apply or RejectResult(\"%s\") to discard\n", resultID, resultID)
	case ModeManual:
		// En modo manual, solo guardamos el resultado sin aplicar
		result.Status = "completed"
	}

	// Agregar al historial
	ga.history = append(ga.history, *result)
	
	// Ejecutar callbacks si existen
	if callback, exists := ga.callbacks[string(request.Type)]; exists {
		callback(result)
	}

	return result, nil
}

// ApplyResult aplica un resultado generado al proyecto
func (ga *GenerativeAgent) ApplyResult(result *GenerationResult) error {
	generator, exists := ga.generators[result.Request.Type]
	if !exists {
		return fmt.Errorf("no generator found for type: %s", result.Request.Type)
	}

	err := generator.Apply(result)
	if err != nil {
		return fmt.Errorf("failed to apply result: %w", err)
	}

	result.Status = "applied"
	now := time.Now()
	result.AppliedAt = &now

	// Remover de pendientes si existe
	delete(ga.pendingResults, result.ID)

	fmt.Printf("[Agent] Applied %s: %s\n", result.Request.Type, result.Request.Prompt)
	return nil
}

// RejectResult rechaza un resultado pendiente
func (ga *GenerativeAgent) RejectResult(resultID string) error {
	result, exists := ga.pendingResults[resultID]
	if !exists {
		return fmt.Errorf("no pending result with ID: %s", resultID)
	}

	result.Status = "rejected"
	delete(ga.pendingResults, resultID)
	
	fmt.Printf("[Agent] Rejected result: %s\n", resultID)
	return nil
}

// GetPendingResults devuelve todos los resultados pendientes
func (ga *GenerativeAgent) GetPendingResults() map[string]*GenerationResult {
	return ga.pendingResults
}

// GetHistory devuelve el historial de generaciones
func (ga *GenerativeAgent) GetHistory() []GenerationResult {
	return ga.history
}

// PreviewResult muestra una previsualización de un resultado
func (ga *GenerativeAgent) PreviewResult(resultID string) (string, error) {
	result, exists := ga.pendingResults[resultID]
	if !exists {
		return "", fmt.Errorf("no pending result with ID: %s", resultID)
	}

	generator, exists := ga.generators[result.Request.Type]
	if !exists {
		return "", fmt.Errorf("no generator found for type: %s", result.Request.Type)
	}

	return generator.Preview(result.Content)
}

// RegisterCallback registra un callback para cuando se genera un tipo de recurso
func (ga *GenerativeAgent) RegisterCallback(resourceType ResourceType, callback func(*GenerationResult)) {
	ga.callbacks[string(resourceType)] = callback
}

// SaveSession guarda el estado actual de la sesión
func (ga *GenerativeAgent) SaveSession(path string) error {
	sessionData := map[string]interface{}{
		"mode":     ga.mode,
		"history":  ga.history,
		"pending":  ga.pendingResults,
		"saved_at": time.Now(),
	}

	data, err := json.MarshalIndent(sessionData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// LoadSession carga una sesión guardada
func (ga *GenerativeAgent) LoadSession(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("session file does not exist: %s", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var sessionData map[string]interface{}
	err = json.Unmarshal(data, &sessionData)
	if err != nil {
		return err
	}

	// Restaurar modo
	if mode, ok := sessionData["mode"].(string); ok {
		ga.mode = AgentMode(mode)
	}

	// TODO: Restaurar historial y pendientes
	// (requiere deserialización más compleja)

	return nil
}

// GetProjectStructure analiza y devuelve la estructura actual del proyecto
func (ga *GenerativeAgent) GetProjectStructure() (map[string]interface{}, error) {
	structure := make(map[string]interface{})
	
	err := filepath.Walk(ga.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		relPath, _ := filepath.Rel(ga.projectPath, path)
		
		if info.IsDir() {
			structure[relPath] = "directory"
		} else {
			ext := filepath.Ext(path)
			structure[relPath] = map[string]interface{}{
				"type": "file",
				"extension": ext,
				"size": info.Size(),
				"modified": info.ModTime(),
			}
		}
		
		return nil
	})
	
	return structure, err
}

// GenerateFromPrompt es un método de conveniencia para generar desde un prompt simple
func (ga *GenerativeAgent) GenerateFromPrompt(prompt string, resourceType ResourceType) (*GenerationResult, error) {
	request := GenerationRequest{
		Type:      resourceType,
		Prompt:    prompt,
		AutoApply: ga.mode == ModeAutomatic,
		Timestamp: time.Now(),
	}
	
	return ga.GenerateResource(context.Background(), request)
}

// GetMode devuelve el modo actual del agente
func (ga *GenerativeAgent) GetMode() AgentMode {
	return ga.mode
}
