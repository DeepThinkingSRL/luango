package generator

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// AgentCLI proporciona una interfaz de línea de comandos para el agente generativo
type AgentCLI struct {
	agent  *GenerativeAgent
	reader *bufio.Reader
	active bool
}

// NewAgentCLI crea una nueva instancia de la CLI del agente
func NewAgentCLI(agent *GenerativeAgent) *AgentCLI {
	return &AgentCLI{
		agent:  agent,
		reader: bufio.NewReader(os.Stdin),
		active: false,
	}
}

// Start inicia la interfaz de línea de comandos
func (cli *AgentCLI) Start() {
	cli.active = true
	cli.printWelcome()
	
	for cli.active {
		cli.printPrompt()
		input := cli.readInput()
		cli.processCommand(input)
	}
}

// Stop detiene la CLI
func (cli *AgentCLI) Stop() {
	cli.active = false
	fmt.Println("\n🔄 Agente generativo detenido. ¡Hasta la vista!")
}

// printWelcome muestra el mensaje de bienvenida
func (cli *AgentCLI) printWelcome() {
	fmt.Println(`
🤖 ═══════════════════════════════════════════════════════════════
   LUANGO GENERATIVE AGENT - Motor de Videojuegos Generativo  
═══════════════════════════════════════════════════════════════
🎮 Modo actual: ` + string(cli.agent.mode) + `
🔧 Comandos disponibles:
   • generate <tipo> <prompt>  - Genera un recurso
   • mode <modo>              - Cambia el modo (interactive/automatic/manual)
   • pending                  - Muestra recursos pendientes
   • apply <id>               - Aplica un recurso pendiente
   • reject <id>              - Rechaza un recurso pendiente  
   • preview <id>             - Previsualiza un recurso
   • history                  - Muestra el historial
   • status                   - Estado del proyecto
   • help                     - Muestra ayuda
   • exit                     - Salir del agente
   
💡 Ejemplo: generate enemy "crea un slime verde que salte"
═══════════════════════════════════════════════════════════════`)
}

// printPrompt muestra el prompt de comando
func (cli *AgentCLI) printPrompt() {
	mode := cli.agent.mode
	pendingCount := len(cli.agent.GetPendingResults())
	
	modeIcon := "🤖"
	switch mode {
	case ModeInteractive:
		modeIcon = "🔄"
	case ModeAutomatic:
		modeIcon = "⚡"
	case ModeManual:
		modeIcon = "✋"
	}
	
	fmt.Printf("\n%s [%s", modeIcon, mode)
	if pendingCount > 0 {
		fmt.Printf(" | %d pendientes", pendingCount)
	}
	fmt.Print("] > ")
}

// readInput lee la entrada del usuario
func (cli *AgentCLI) readInput() string {
	input, err := cli.reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error leyendo entrada: %v\n", err)
		return ""
	}
	return strings.TrimSpace(input)
}

// processCommand procesa un comando ingresado
func (cli *AgentCLI) processCommand(input string) {
	if input == "" {
		return
	}
	
	parts := strings.Fields(input)
	command := strings.ToLower(parts[0])
	
	switch command {
	case "generate", "gen", "g":
		cli.handleGenerate(parts[1:])
	case "mode", "m":
		cli.handleMode(parts[1:])
	case "pending", "p":
		cli.handlePending()
	case "apply", "a":
		cli.handleApply(parts[1:])
	case "reject", "r":
		cli.handleReject(parts[1:])
	case "preview", "pre":
		cli.handlePreview(parts[1:])
	case "history", "h":
		cli.handleHistory()
	case "status", "s":
		cli.handleStatus()
	case "help", "?":
		cli.handleHelp()
	case "exit", "quit", "q":
		cli.Stop()
	case "clear", "cls":
		cli.clearScreen()
	default:
		fmt.Printf("❌ Comando desconocido: '%s'. Usa 'help' para ver comandos disponibles.\n", command)
	}
}

// handleGenerate maneja el comando generate
func (cli *AgentCLI) handleGenerate(args []string) {
	if len(args) < 2 {
		fmt.Println("❌ Uso: generate <tipo> <prompt>")
		fmt.Println("   Tipos: script, entity, item, enemy, player, behavior, level")
		return
	}
	
	resourceType := ResourceType(strings.ToLower(args[0]))
	prompt := strings.Join(args[1:], " ")
	
	// Validar tipo de recurso
	validTypes := []ResourceType{
		ResourceScript, ResourceEntity, ResourceSprite, ResourceSound,
		ResourceLevel, ResourceBehavior, ResourceDialogue, ResourceTerrain,
	}
	
	valid := false
	for _, vt := range validTypes {
		if resourceType == vt {
			valid = true
			break
		}
	}
	
	if !valid {
		fmt.Printf("❌ Tipo de recurso inválido: %s\n", resourceType)
		return
	}
	
	fmt.Printf("🔄 Generando %s: %s\n", resourceType, prompt)
	
	// Generar recurso
	result, err := cli.agent.GenerateFromPrompt(prompt, resourceType)
	if err != nil {
		fmt.Printf("❌ Error generando recurso: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Recurso generado: %s\n", result.ID)
	
	// Si está en modo interactivo, mostrar preview automáticamente
	if cli.agent.mode == ModeInteractive {
		cli.showResultPreview(result)
		fmt.Printf("💡 Usa 'apply %s' para aplicar o 'reject %s' para rechazar\n", 
			result.ID, result.ID)
	}
}

// handleMode maneja el comando mode
func (cli *AgentCLI) handleMode(args []string) {
	if len(args) == 0 {
		fmt.Printf("📋 Modo actual: %s\n", cli.agent.mode)
		fmt.Println("   Modos disponibles: interactive, automatic, manual")
		return
	}
	
	newMode := AgentMode(strings.ToLower(args[0]))
	
	switch newMode {
	case ModeInteractive, ModeAutomatic, ModeManual:
		cli.agent.SetMode(newMode)
		fmt.Printf("✅ Modo cambiado a: %s\n", newMode)
	default:
		fmt.Printf("❌ Modo inválido: %s\n", newMode)
		fmt.Println("   Modos válidos: interactive, automatic, manual")
	}
}

// handlePending maneja el comando pending
func (cli *AgentCLI) handlePending() {
	pending := cli.agent.GetPendingResults()
	
	if len(pending) == 0 {
		fmt.Println("📭 No hay recursos pendientes")
		return
	}
	
	fmt.Printf("📋 Recursos pendientes (%d):\n", len(pending))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	for id, result := range pending {
		fmt.Printf("🔸 %s\n", id)
		fmt.Printf("   Tipo: %s\n", result.Request.Type)
		fmt.Printf("   Prompt: %s\n", result.Request.Prompt)
		fmt.Printf("   Creado: %s\n", result.CreatedAt.Format("15:04:05"))
		fmt.Println("   ────────────────────────────────────────")
	}
}

// handleApply maneja el comando apply
func (cli *AgentCLI) handleApply(args []string) {
	if len(args) == 0 {
		fmt.Println("❌ Uso: apply <id>")
		return
	}
	
	resultID := args[0]
	pending := cli.agent.GetPendingResults()
	
	result, exists := pending[resultID]
	if !exists {
		fmt.Printf("❌ No se encontró resultado pendiente: %s\n", resultID)
		return
	}
	
	fmt.Printf("🔄 Aplicando resultado: %s\n", resultID)
	
	err := cli.agent.ApplyResult(result)
	if err != nil {
		fmt.Printf("❌ Error aplicando resultado: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Resultado aplicado exitosamente: %s\n", result.FilePath)
}

// handleReject maneja el comando reject  
func (cli *AgentCLI) handleReject(args []string) {
	if len(args) == 0 {
		fmt.Println("❌ Uso: reject <id>")
		return
	}
	
	resultID := args[0]
	
	err := cli.agent.RejectResult(resultID)
	if err != nil {
		fmt.Printf("❌ Error rechazando resultado: %v\n", err)
		return
	}
	
	fmt.Printf("🗑️ Resultado rechazado: %s\n", resultID)
}

// handlePreview maneja el comando preview
func (cli *AgentCLI) handlePreview(args []string) {
	if len(args) == 0 {
		fmt.Println("❌ Uso: preview <id>")
		return
	}
	
	resultID := args[0]
	
	preview, err := cli.agent.PreviewResult(resultID)
	if err != nil {
		fmt.Printf("❌ Error generando preview: %v\n", err)
		return
	}
	
	fmt.Println(preview)
}

// handleHistory maneja el comando history
func (cli *AgentCLI) handleHistory() {
	history := cli.agent.GetHistory()
	
	if len(history) == 0 {
		fmt.Println("📝 Historial vacío")
		return
	}
	
	fmt.Printf("📚 Historial de generaciones (%d):\n", len(history))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	// Mostrar los últimos 10 elementos
	start := 0
	if len(history) > 10 {
		start = len(history) - 10
	}
	
	for i := start; i < len(history); i++ {
		result := history[i]
		statusIcon := "⏳"
		
		switch result.Status {
		case "completed":
			statusIcon = "✅"
		case "applied":
			statusIcon = "🚀"
		case "failed":
			statusIcon = "❌"
		case "rejected":
			statusIcon = "🗑️"
		}
		
		fmt.Printf("%s %s | %s | %s\n", 
			statusIcon, 
			result.Request.Type,
			result.Request.Prompt[:min(50, len(result.Request.Prompt))],
			result.CreatedAt.Format("15:04"))
	}
	
	if len(history) > 10 {
		fmt.Printf("... y %d más (usa 'history full' para ver todo)\n", len(history)-10)
	}
}

// handleStatus maneja el comando status
func (cli *AgentCLI) handleStatus() {
	pending := cli.agent.GetPendingResults()
	history := cli.agent.GetHistory()
	
	fmt.Println("📊 Estado del Agente Generativo")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("🔧 Modo: %s\n", cli.agent.mode)
	fmt.Printf("📋 Recursos pendientes: %d\n", len(pending))
	fmt.Printf("📚 Total generaciones: %d\n", len(history))
	
	// Contar por estado
	statusCount := make(map[string]int)
	for _, result := range history {
		statusCount[result.Status]++
	}
	
	fmt.Println("\n📈 Estadísticas:")
	for status, count := range statusCount {
		fmt.Printf("   %s: %d\n", status, count)
	}
	
	// Mostrar estructura del proyecto
	fmt.Println("\n📁 Estructura del proyecto:")
	structure, err := cli.agent.GetProjectStructure()
	if err != nil {
		fmt.Printf("❌ Error obteniendo estructura: %v\n", err)
		return
	}
	
	cli.printStructure(structure, 0)
}

// handleHelp muestra la ayuda
func (cli *AgentCLI) handleHelp() {
	fmt.Println(`
🆘 Ayuda del Agente Generativo
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🎯 Comandos principales:
   generate <tipo> <prompt>     - Genera un nuevo recurso
   apply <id>                   - Aplica un recurso pendiente  
   reject <id>                  - Rechaza un recurso pendiente
   preview <id>                 - Muestra preview de un recurso

📋 Gestión:
   pending                      - Lista recursos pendientes
   history                      - Muestra historial de generaciones
   status                       - Estado actual del agente
   mode <modo>                  - Cambia modo de operación

🔧 Modos disponibles:
   • interactive - Requiere confirmación para aplicar
   • automatic   - Aplica automáticamente  
   • manual      - Solo genera, no aplica

🎮 Tipos de recursos:
   • script     - Scripts Lua generales
   • enemy      - Enemigos y NPCs
   • player     - Sistemas de jugador
   • item       - Objetos e items
   • behavior   - Comportamientos AI
   • level      - Niveles y mapas
   • dialogue   - Diálogos y narrativa

💡 Ejemplos:
   generate enemy "crea un dragón que escupa fuego"
   generate item "poción de vida que cure 50 HP"  
   generate player "sistema de inventario expandido"
   mode interactive
   apply enemy_1640123456
`)
}

// clearScreen limpia la pantalla
func (cli *AgentCLI) clearScreen() {
	fmt.Print("\033[2J\033[H")
	cli.printWelcome()
}

// showResultPreview muestra un preview de un resultado
func (cli *AgentCLI) showResultPreview(result *GenerationResult) {
	// Intentar generar preview usando el generador apropiado
	generator, exists := cli.agent.generators[result.Request.Type]
	if !exists {
		fmt.Println("⚠️ No se puede generar preview para este tipo de recurso")
		return
	}
	
	preview, err := generator.Preview(result.Content)
	if err != nil {
		fmt.Printf("⚠️ Error generando preview: %v\n", err)
		return
	}
	
	fmt.Println(preview)
}

// printStructure imprime la estructura del proyecto de forma jerárquica
func (cli *AgentCLI) printStructure(structure map[string]interface{}, depth int) {
	indent := strings.Repeat("  ", depth)
	count := 0
	maxItems := 20 // Limitar la salida
	
	for path, info := range structure {
		if count >= maxItems {
			fmt.Printf("%s... y más archivos\n", indent)
			break
		}
		
		if info == "directory" {
			fmt.Printf("%s📁 %s/\n", indent, path)
		} else {
			fmt.Printf("%s📄 %s\n", indent, path)
		}
		count++
	}
}

// min función auxiliar para obtener el mínimo de dos enteros
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
