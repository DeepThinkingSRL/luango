-- Potions Module
-- Handles potion items and their effects

local potions = {}

-- Potion types and their properties
potions.types = {
    health = {
        name = "Health Potion",
        heal_amount = 50,
        icon = "❤️",
        description = "Restores 50 health points"
    },
    energy = {
        name = "Energy Potion", 
        energy_amount = 30,
        icon = "⚡",
        description = "Restores 30 energy points"
    },
    strength = {
        name = "Strength Potion",
        damage_boost = 10,
        duration = 300, -- 5 seconds at 60fps
        icon = "💪",
        description = "Increases damage by 10 for 5 seconds"
    },
    speed = {
        name = "Speed Potion",
        speed_boost = 3,
        duration = 600, -- 10 seconds at 60fps
        icon = "🏃",
        description = "Increases movement speed by 3 for 10 seconds"
    }
}

-- Active potion effects
potions.active_effects = {}

-- Potion inventory (simplified)
potions.inventory = {
    health = 3,
    energy = 2,
    strength = 1,
    speed = 1
}

-- Initialize potions module
function potions.init()
    log("🧪 Potions module ready!")
    log("💼 Starting inventory:")
    for potion_type, count in pairs(potions.inventory) do
        local potion = potions.types[potion_type]
        if potion then
            log("  " .. potion.icon .. " " .. potion.name .. " x" .. count)
        end
    end
end

-- Update active potion effects
function potions.update()
    -- Update effect durations
    for effect_id, effect in pairs(potions.active_effects) do
        effect.remaining_time = effect.remaining_time - 1
        
        if effect.remaining_time <= 0 then
            potions.remove_effect(effect_id)
        end
    end
end

-- Use a potion
function potions.use_potion(potion_type, target)
    -- Check if we have the potion
    if not potions.inventory[potion_type] or potions.inventory[potion_type] <= 0 then
        log("❌ No " .. potion_type .. " potions available!")
        return false
    end
    
    local potion = potions.types[potion_type]
    if not potion then
        log("❌ Unknown potion type: " .. potion_type)
        return false
    end
    
    -- Consume the potion
    potions.inventory[potion_type] = potions.inventory[potion_type] - 1
    log("🧪 Used " .. potion.icon .. " " .. potion.name)
    
    -- Apply potion effects
    if potion_type == "health" then
        potions.apply_health_effect(potion.heal_amount, target)
    elseif potion_type == "energy" then
        potions.apply_energy_effect(potion.energy_amount, target)
    elseif potion_type == "strength" then
        potions.apply_strength_effect(potion.damage_boost, potion.duration, target)
    elseif potion_type == "speed" then
        potions.apply_speed_effect(potion.speed_boost, potion.duration, target)
    end
    
    emit("potion_used", potion_type)
    return true
end

-- Apply health effect
function potions.apply_health_effect(amount, target)
    log("💚 Healing for " .. amount .. " HP")
    -- In a real game, this would call target.heal(amount)
    -- For now, just log the effect
    emit("player_healed", amount)
end

-- Apply energy effect  
function potions.apply_energy_effect(amount, target)
    log("🔋 Restoring " .. amount .. " energy")
    -- In a real game, this would call target.restore_energy(amount)
    emit("player_energy_restored", amount)
end

-- Apply strength effect (temporary)
function potions.apply_strength_effect(boost, duration, target)
    local effect_id = "strength_" .. os.time()
    
    potions.active_effects[effect_id] = {
        type = "strength",
        boost = boost,
        remaining_time = duration,
        target = target
    }
    
    log("💪 Strength increased by " .. boost .. " for " .. (duration/60) .. " seconds")
    emit("strength_boost_applied", boost)
end

-- Apply speed effect (temporary)
function potions.apply_speed_effect(boost, duration, target)
    local effect_id = "speed_" .. os.time()
    
    potions.active_effects[effect_id] = {
        type = "speed",
        boost = boost,
        remaining_time = duration,
        target = target
    }
    
    log("🏃 Speed increased by " .. boost .. " for " .. (duration/60) .. " seconds")
    emit("speed_boost_applied", boost)
end

-- Remove an effect
function potions.remove_effect(effect_id)
    local effect = potions.active_effects[effect_id]
    if effect then
        log("⏰ " .. effect.type .. " effect expired")
        emit("effect_expired", effect.type)
        potions.active_effects[effect_id] = nil
    end
end

-- Add potions to inventory
function potions.add_potion(potion_type, count)
    count = count or 1
    
    if potions.types[potion_type] then
        potions.inventory[potion_type] = (potions.inventory[potion_type] or 0) + count
        local potion = potions.types[potion_type]
        log("📦 Found " .. potion.icon .. " " .. potion.name .. " x" .. count)
        emit("potion_found", potion_type)
    else
        log("❌ Unknown potion type: " .. potion_type)
    end
end

-- Get inventory count
function potions.get_count(potion_type)
    return potions.inventory[potion_type] or 0
end

-- Get all active effects
function potions.get_active_effects()
    return potions.active_effects
end

-- Export the module
return potions
