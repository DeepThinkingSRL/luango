-- Agent Demo Script
-- Main demonstration script that shows off engine capabilities

local agent_demo = {}

-- Demo variables
local demo_timer = 0
local demo_mode = "idle"
local demo_phases = {"movement", "audio", "debug"}
local current_phase = 1

-- Initialize the demo
function agent_demo.init()
    log("🤖 Agent Demo initialized")
    emit("demo_started", "main")
    demo_mode = "active"
end

-- Update demo logic
function agent_demo.update()
    demo_timer = demo_timer + 1
    
    -- Switch demo phases every 300 frames (about 5 seconds at 60fps)
    if demo_timer % 300 == 0 then
        current_phase = current_phase + 1
        if current_phase > #demo_phases then
            current_phase = 1
        end
        
        local phase_name = demo_phases[current_phase]
        log("🔄 Demo phase: " .. phase_name)
        emit("phase_changed", phase_name)
    end
    
    -- Execute current phase
    if demo_phases[current_phase] == "movement" then
        agent_demo.movement_demo()
    elseif demo_phases[current_phase] == "audio" then
        agent_demo.audio_demo()
    elseif demo_phases[current_phase] == "debug" then
        agent_demo.debug_demo()
    end
end

-- Movement demonstration
function agent_demo.movement_demo()
    -- Simple player movement pattern
    if is_key_pressed("W") then
        move_player(0, -2)
        log("⬆️ Moving up")
    end
    if is_key_pressed("S") then
        move_player(0, 2)
        log("⬇️ Moving down")
    end
    if is_key_pressed("A") then
        move_player(-2, 0)
        log("⬅️ Moving left")
    end
    if is_key_pressed("D") then
        move_player(2, 0)
        log("➡️ Moving right")
    end
end

-- Audio demonstration
function agent_demo.audio_demo()
    -- Play sound every 120 frames (about 2 seconds)
    if demo_timer % 120 == 0 then
        play_sound("assets/sounds/sonidito.wav")
        log("🔊 Playing demo sound")
    end
end

-- Debug information demonstration
function agent_demo.debug_demo()
    -- Output debug info periodically
    if demo_timer % 60 == 0 then
        debug("Demo timer: " .. demo_timer)
        debug("Current phase: " .. demo_phases[current_phase])
        debug("Demo mode: " .. demo_mode)
    end
end

-- Export the module
return agent_demo
