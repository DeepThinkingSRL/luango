-- 🤖 Agent Demo Script
-- Ejemplo de uso del agente generativo desde Lua

log("🤖 Agent Demo loaded!")

-- Función para demostrar generación desde Lua
function demo_generate_enemy()
    log("🔄 Generating enemy from Lua...")
    local resultID = generate("script", "create a fire dragon enemy that breathes fire and has 200 health")
    
    if resultID then
        log("✅ Enemy generated with ID: " .. resultID)
        log("💡 Use F4 to open agent console and apply with: apply " .. resultID)
    else
        log("❌ Failed to generate enemy")
    end
end

-- Función para cambiar modo del agente
function demo_set_automatic()
    log("⚡ Setting agent to automatic mode")
    set_agent_mode("automatic")
end

function demo_set_interactive()
    log("🔄 Setting agent to interactive mode")
    set_agent_mode("interactive")
end

-- Auto-demostración al iniciar
function on_start()
    log("🎯 Agent Demo starting...")
    log("💡 Press 'G' to generate an enemy")
    log("💡 Press 'M' to toggle agent mode")
    log("💡 Press 'F4' to open agent console")
end

-- Manejar input de demo
function on_update()
    if is_key_pressed("G") then
        demo_generate_enemy()
    end
    
    if is_key_pressed("M") then
        demo_set_automatic()
    end
end
