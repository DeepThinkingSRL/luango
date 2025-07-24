-- Player Module
-- Handles player-specific logic and behaviors

local player = {}

-- Player stats and state
player.stats = {
    health = 100,
    energy = 100,
    speed = 5
}

player.state = {
    is_moving = false,
    direction = "none",
    last_position = {x = 0, y = 0}
}

-- Initialize player module
function player.init()
    log("👤 Player module loaded!")
    log("🔧 Player stats initialized")
    log("❤️ Health: " .. player.stats.health)
    log("⚡ Energy: " .. player.stats.energy)
    log("🏃 Speed: " .. player.stats.speed)
end

-- Update player logic
function player.update()
    -- Handle player input and movement
    local moving = false
    local direction = "none"
    
    if is_key_pressed("W") or is_key_pressed("ArrowUp") then
        move_player(0, -player.stats.speed)
        direction = "up"
        moving = true
    end
    if is_key_pressed("S") or is_key_pressed("ArrowDown") then
        move_player(0, player.stats.speed)
        direction = "down"
        moving = true
    end
    if is_key_pressed("A") or is_key_pressed("ArrowLeft") then
        move_player(-player.stats.speed, 0)
        direction = "left"
        moving = true
    end
    if is_key_pressed("D") or is_key_pressed("ArrowRight") then
        move_player(player.stats.speed, 0)
        direction = "right"
        moving = true
    end
    
    -- Update player state
    if moving ~= player.state.is_moving then
        player.state.is_moving = moving
        if moving then
            log("🏃 Player started moving " .. direction)
        else
            log("🛑 Player stopped moving")
        end
    end
    
    player.state.direction = direction
end

-- Player actions
function player.take_damage(amount)
    player.stats.health = math.max(0, player.stats.health - amount)
    log("💔 Player took " .. amount .. " damage. Health: " .. player.stats.health)
    
    if player.stats.health <= 0 then
        log("💀 Player died!")
        emit("player_died", "game_over")
    end
end

function player.heal(amount)
    player.stats.health = math.min(100, player.stats.health + amount)
    log("💚 Player healed " .. amount .. " HP. Health: " .. player.stats.health)
end

function player.use_energy(amount)
    player.stats.energy = math.max(0, player.stats.energy - amount)
    log("⚡ Player used " .. amount .. " energy. Energy: " .. player.stats.energy)
end

function player.restore_energy(amount)
    player.stats.energy = math.min(100, player.stats.energy + amount)
    log("🔋 Player restored " .. amount .. " energy. Energy: " .. player.stats.energy)
end

-- Export the module
return player
