-- Main Game Script
-- Coordinates all game modules and handles main game loop

-- Load modules
local agent_demo = require("agent_demo")
-- Note: Other modules are loaded automatically by the engine

local game = {}
game.initialized = false
game.modules = {}

-- Game lifecycle hooks
function on_start()
    log("🚀 Game starting...")
    
    -- Register modules
    game.modules.agent_demo = agent_demo
    
    -- Initialize all modules
    for name, module in pairs(game.modules) do
        if module.init then
            log("📦 Initializing module: " .. name)
            module.init()
        end
    end
    
    game.initialized = true
    log("✅ Game initialization complete")
end

function on_update()
    if not game.initialized then
        return
    end
    
    -- Update all modules
    for name, module in pairs(game.modules) do
        if module.update then
            module.update()
        end
    end
end

-- Export game object for debugging
_G.game = game
