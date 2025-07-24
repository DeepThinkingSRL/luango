package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// LuaScriptGenerator genera scripts Lua para el motor
type LuaScriptGenerator struct {
	projectPath   string
	templateCache map[string]*template.Template
}

// LuaScriptContent representa el contenido de un script Lua generado
type LuaScriptContent struct {
	Type        string                 `json:"type"`        // entity, behavior, system, etc.
	Name        string                 `json:"name"`        // nombre del script
	Code        string                 `json:"code"`        // código Lua generado
	Functions   []string               `json:"functions"`   // funciones principales
	Dependencies []string              `json:"dependencies"` // dependencias de otros scripts
	Parameters  map[string]interface{} `json:"parameters"`  // parámetros configurables
	Category    string                 `json:"category"`    // categoría del script
}

// NewLuaScriptGenerator crea un nuevo generador de scripts Lua
func NewLuaScriptGenerator(projectPath string) *LuaScriptGenerator {
	return &LuaScriptGenerator{
		projectPath:   projectPath,
		templateCache: make(map[string]*template.Template),
	}
}

// Generate implementa ResourceGenerator para scripts Lua
func (lg *LuaScriptGenerator) Generate(ctx context.Context, request GenerationRequest) (*GenerationResult, error) {
	// Analizar el prompt para determinar qué tipo de script generar
	scriptType := lg.determineScriptType(request.Prompt)
	
	// Generar el contenido del script
	content, err := lg.generateScriptContent(request.Prompt, scriptType, request.Parameters)
	if err != nil {
		return nil, fmt.Errorf("failed to generate script content: %w", err)
	}

	// Crear el resultado
	result := &GenerationResult{
		Request:   request,
		Status:    "completed",
		Content:   content,
		CreatedAt: time.Now(),
		Metadata: map[string]interface{}{
			"script_type": scriptType,
			"functions":   content.Functions,
			"category":    content.Category,
		},
	}

	// Determinar la ruta de salida
	if request.OutputPath == "" {
		result.FilePath = filepath.Join("mod", content.Category, content.Name+".lua")
	} else {
		result.FilePath = request.OutputPath
	}

	return result, nil
}

// Preview implementa ResourceGenerator para mostrar una previsualización
func (lg *LuaScriptGenerator) Preview(content interface{}) (string, error) {
	scriptContent, ok := content.(*LuaScriptContent)
	if !ok {
		return "", fmt.Errorf("invalid content type for Lua script preview")
	}

	preview := fmt.Sprintf(`
🔧 Script Lua Generado
━━━━━━━━━━━━━━━━━━━━━━
📁 Tipo: %s
📝 Nombre: %s
📂 Categoría: %s
🔗 Dependencias: %v
⚙️  Funciones: %v

📜 Código (primeras 10 líneas):
━━━━━━━━━━━━━━━━━━━━━━━━━━━━
%s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━
`,
		scriptContent.Type,
		scriptContent.Name,
		scriptContent.Category,
		scriptContent.Dependencies,
		scriptContent.Functions,
		lg.truncateCode(scriptContent.Code, 10),
	)

	return preview, nil
}

// Apply implementa ResourceGenerator para aplicar el script al proyecto
func (lg *LuaScriptGenerator) Apply(result *GenerationResult) error {
	scriptContent, ok := result.Content.(*LuaScriptContent)
	if !ok {
		return fmt.Errorf("invalid content type for Lua script")
	}

	// Crear directorio si no existe
	fullPath := filepath.Join(lg.projectPath, result.FilePath)
	dir := filepath.Dir(fullPath)
	if err := ensureDir(dir); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Escribir el archivo
	err := os.WriteFile(fullPath, []byte(scriptContent.Code), 0644)
	if err != nil {
		return fmt.Errorf("failed to write script file: %w", err)
	}

	fmt.Printf("[LuaGen] Script aplicado: %s\n", fullPath)
	return nil
}

// Validate implementa ResourceGenerator para validar el contenido
func (lg *LuaScriptGenerator) Validate(content interface{}) error {
	scriptContent, ok := content.(*LuaScriptContent)
	if !ok {
		return fmt.Errorf("invalid content type for Lua script")
	}

	if scriptContent.Name == "" {
		return fmt.Errorf("script name cannot be empty")
	}

	if scriptContent.Code == "" {
		return fmt.Errorf("script code cannot be empty")
	}

	// Validación básica de sintaxis Lua (básica)
	if !strings.Contains(scriptContent.Code, "function") && 
	   !strings.Contains(scriptContent.Code, "=") {
		return fmt.Errorf("script appears to have no valid Lua content")
	}

	return nil
}

// determineScriptType analiza el prompt para determinar el tipo de script
func (lg *LuaScriptGenerator) determineScriptType(prompt string) string {
	prompt = strings.ToLower(prompt)
	
	if strings.Contains(prompt, "enemy") || strings.Contains(prompt, "enemigo") {
		return "enemy"
	}
	if strings.Contains(prompt, "player") || strings.Contains(prompt, "jugador") {
		return "player"
	}
	if strings.Contains(prompt, "item") || strings.Contains(prompt, "objeto") {
		return "item"
	}
	if strings.Contains(prompt, "behavior") || strings.Contains(prompt, "comportamiento") {
		return "behavior"
	}
	if strings.Contains(prompt, "system") || strings.Contains(prompt, "sistema") {
		return "system"
	}
	if strings.Contains(prompt, "world") || strings.Contains(prompt, "mundo") {
		return "world"
	}
	if strings.Contains(prompt, "ui") || strings.Contains(prompt, "interface") {
		return "ui"
	}
	
	return "general"
}

// generateScriptContent genera el contenido real del script
func (lg *LuaScriptGenerator) generateScriptContent(prompt, scriptType string, params map[string]interface{}) (*LuaScriptContent, error) {
	// Templates básicos para diferentes tipos de scripts
	templates := map[string]string{
		"enemy": `-- 👾 Enemy: {{.Name}}
-- Generated from prompt: {{.Prompt}}

{{.Name}} = {
    name = "{{.Name}}",
    health = {{.Health}},
    speed = {{.Speed}},
    damage = {{.Damage}},
    position = { x = 0, y = 0 },
    state = "idle"
}

function {{.Name}}.init()
    log("🎯 " .. {{.Name}}.name .. " initialized!")
end

function {{.Name}}.update(dt)
    -- AI behavior logic here
    if {{.Name}}.state == "idle" then
        {{.Name}}.patrol()
    elseif {{.Name}}.state == "chase" then
        {{.Name}}.chase_player()
    end
end

function {{.Name}}.patrol()
    -- Basic patrol behavior
    {{.Name}}.position.x = {{.Name}}.position.x + {{.Name}}.speed * math.sin(world.time * 0.01)
end

function {{.Name}}.chase_player()
    -- Chase logic (needs player reference)
    log("👾 " .. {{.Name}}.name .. " is chasing!")
end

function {{.Name}}.take_damage(amount)
    {{.Name}}.health = {{.Name}}.health - amount
    if {{.Name}}.health <= 0 then
        {{.Name}}.die()
    end
end

function {{.Name}}.die()
    log("💀 " .. {{.Name}}.name .. " has been defeated!")
    emit("enemy_died", {{.Name}}.name)
end

-- Initialize the enemy
{{.Name}}.init()
`,

		"item": `-- 🎒 Item: {{.Name}}
-- Generated from prompt: {{.Prompt}}

{{.Name}} = {
    name = "{{.Name}}",
    type = "{{.Type}}",
    stackable = {{.Stackable}},
    max_stack = {{.MaxStack}},
    rarity = "{{.Rarity}}",
    description = "{{.Description}}"
}

function {{.Name}}.use(player)
    log("✨ Using " .. {{.Name}}.name)
    
    {{if eq .Type "consumable"}}
    -- Consumable item logic
    player.health = player.health + {{.HealthRestore}}
    log("❤️ Restored " .. {{.HealthRestore}} .. " health")
    return true -- Item consumed
    {{else if eq .Type "weapon"}}
    -- Weapon logic
    player.weapon = {{.Name}}
    player.damage = {{.Damage}}
    log("⚔️ Equipped " .. {{.Name}}.name)
    return false -- Item not consumed
    {{else}}
    -- Generic item use
    log("🔧 " .. {{.Name}}.name .. " was used")
    return false
    {{end}}
end

function {{.Name}}.on_pickup(player)
    log("📦 " .. player.name .. " picked up " .. {{.Name}}.name)
    emit("item_picked_up", {{.Name}}.name)
end

function {{.Name}}.get_tooltip()
    return {{.Name}}.name .. "\n" .. {{.Name}}.description .. "\nRarity: " .. {{.Name}}.rarity
end
`,

		"player": `-- 👤 Player: {{.Name}}
-- Generated from prompt: {{.Prompt}}

{{.Name}} = {
    name = "{{.Name}}",
    health = {{.Health}},
    max_health = {{.MaxHealth}},
    speed = {{.Speed}},
    position = { x = 100, y = 100 },
    inventory = {},
    experience = 0,
    level = 1
}

function {{.Name}}.init()
    log("👤 Player " .. {{.Name}}.name .. " initialized!")
end

function {{.Name}}.update(dt)
    {{.Name}}.handle_input()
    {{.Name}}.check_collisions()
end

function {{.Name}}.handle_input()
    local dx, dy = 0, 0
    
    if is_key_pressed("W") or is_key_pressed("ArrowUp") then
        dy = -{{.Name}}.speed
    end
    if is_key_pressed("S") or is_key_pressed("ArrowDown") then
        dy = {{.Name}}.speed
    end
    if is_key_pressed("A") or is_key_pressed("ArrowLeft") then
        dx = -{{.Name}}.speed
    end
    if is_key_pressed("D") or is_key_pressed("ArrowRight") then
        dx = {{.Name}}.speed
    end
    
    if dx ~= 0 or dy ~= 0 then
        move_player(dx, dy)
        {{.Name}}.position.x = {{.Name}}.position.x + dx
        {{.Name}}.position.y = {{.Name}}.position.y + dy
    end
end

function {{.Name}}.check_collisions()
    -- Collision detection logic here
end

function {{.Name}}.take_damage(amount)
    {{.Name}}.health = math.max(0, {{.Name}}.health - amount)
    log("💔 " .. {{.Name}}.name .. " took " .. amount .. " damage")
    
    if {{.Name}}.health <= 0 then
        {{.Name}}.die()
    end
end

function {{.Name}}.heal(amount)
    {{.Name}}.health = math.min({{.Name}}.max_health, {{.Name}}.health + amount)
    log("❤️ " .. {{.Name}}.name .. " healed " .. amount .. " HP")
end

function {{.Name}}.die()
    log("💀 " .. {{.Name}}.name .. " has died!")
    emit("player_died", {{.Name}}.name)
end

function {{.Name}}.gain_experience(amount)
    {{.Name}}.experience = {{.Name}}.experience + amount
    log("⭐ Gained " .. amount .. " experience!")
end

-- Initialize the player
{{.Name}}.init()
`,

		"general": `-- 🔧 General Script: {{.Name}}
-- Generated from prompt: {{.Prompt}}

{{.Name}} = {}

function {{.Name}}.init()
    log("🔧 " .. "{{.Name}}" .. " module initialized!")
end

function {{.Name}}.update(dt)
    -- Update logic here
end

-- Initialize the module
{{.Name}}.init()
`,
	}

	// Obtener template apropiado
	templateStr, exists := templates[scriptType]
	if !exists {
		templateStr = templates["general"]
	}

	// Preparar datos para el template
	data := lg.prepareTemplateData(prompt, scriptType, params)
	
	// Procesar template
	tmpl, err := template.New("script").Parse(templateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// Crear contenido del script
	content := &LuaScriptContent{
		Type:        scriptType,
		Name:        data["Name"].(string),
		Code:        buf.String(),
		Functions:   lg.extractFunctions(buf.String()),
		Category:    scriptType,
		Parameters:  params,
	}

	return content, nil
}

// prepareTemplateData prepara los datos para el template
func (lg *LuaScriptGenerator) prepareTemplateData(prompt, scriptType string, params map[string]interface{}) map[string]interface{} {
	data := map[string]interface{}{
		"Prompt": prompt,
		"Name":   lg.generateName(prompt, scriptType),
	}

	// Agregar parámetros específicos según el tipo
	switch scriptType {
	case "enemy":
		data["Health"] = getParamOrDefault(params, "health", 100)
		data["Speed"] = getParamOrDefault(params, "speed", 1.0)
		data["Damage"] = getParamOrDefault(params, "damage", 10)
	case "player":
		data["Health"] = getParamOrDefault(params, "health", 100)
		data["MaxHealth"] = getParamOrDefault(params, "max_health", 100)
		data["Speed"] = getParamOrDefault(params, "speed", 2.0)
	case "item":
		data["Type"] = getParamOrDefault(params, "type", "consumable")
		data["Stackable"] = getParamOrDefault(params, "stackable", true)
		data["MaxStack"] = getParamOrDefault(params, "max_stack", 10)
		data["Rarity"] = getParamOrDefault(params, "rarity", "common")
		data["Description"] = getParamOrDefault(params, "description", "A useful item")
		data["HealthRestore"] = getParamOrDefault(params, "health_restore", 25)
		data["Damage"] = getParamOrDefault(params, "damage", 15)
	}

	// Agregar cualquier parámetro personalizado
	for k, v := range params {
		if _, exists := data[strings.Title(k)]; !exists {
			data[strings.Title(k)] = v
		}
	}

	return data
}

// generateName genera un nombre para el script basado en el prompt
func (lg *LuaScriptGenerator) generateName(prompt, scriptType string) string {
	words := strings.Fields(strings.ToLower(prompt))
	
	// Filtrar palabras comunes
	filtered := []string{}
	commonWords := map[string]bool{
		"create": true, "make": true, "generate": true, "add": true,
		"a": true, "an": true, "the": true, "with": true, "for": true,
		"crea": true, "crear": true, "haz": true, "hacer": true, "agrega": true,
		"un": true, "una": true, "el": true, "la": true, "con": true, "para": true,
	}
	
	for _, word := range words {
		if !commonWords[word] && len(word) > 2 {
			filtered = append(filtered, strings.Title(word))
		}
	}
	
	if len(filtered) == 0 {
		return strings.Title(scriptType)
	}
	
	// Tomar las primeras 2-3 palabras significativas
	if len(filtered) > 3 {
		filtered = filtered[:3]
	}
	
	return strings.Join(filtered, "")
}

// extractFunctions extrae los nombres de las funciones del código Lua
func (lg *LuaScriptGenerator) extractFunctions(code string) []string {
	lines := strings.Split(code, "\n")
	functions := []string{}
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "function ") {
			// Extraer nombre de función
			parts := strings.Split(line, " ")
			if len(parts) >= 2 {
				funcName := strings.Split(parts[1], "(")[0]
				functions = append(functions, funcName)
			}
		}
	}
	
	return functions
}

// truncateCode trunca el código a un número específico de líneas
func (lg *LuaScriptGenerator) truncateCode(code string, maxLines int) string {
	lines := strings.Split(code, "\n")
	if len(lines) <= maxLines {
		return code
	}
	
	truncated := strings.Join(lines[:maxLines], "\n")
	return truncated + "\n... (" + fmt.Sprintf("%d", len(lines)-maxLines) + " líneas más)"
}

// getParamOrDefault obtiene un parámetro o devuelve un valor por defecto
func getParamOrDefault(params map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if val, exists := params[key]; exists {
		return val
	}
	return defaultValue
}

// ensureDir crea un directorio si no existe
func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}
