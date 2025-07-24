-- Slime Enemy Module
-- Handles slime enemy behavior and AI

local slime = {}

-- Slime properties
slime.stats = {
    health = 50,
    damage = 10,
    speed = 2,
    aggro_range = 100
}

slime.instances = {}
slime.next_id = 1

-- Slime instance structure
function slime.create_instance(x, y)
    local instance = {
        id = slime.next_id,
        position = {x = x or 0, y = y or 0},
        health = slime.stats.health,
        state = "idle", -- idle, patrol, chase, attack
        direction = {x = 1, y = 0},
        target = nil,
        patrol_timer = 0,
        attack_cooldown = 0
    }
    
    slime.instances[slime.next_id] = instance
    slime.next_id = slime.next_id + 1
    
    log("👾 Slime #" .. instance.id .. " spawned at (" .. x .. ", " .. y .. ")")
    return instance
end

-- Initialize slime module
function slime.init()
    log("👾 Slime enemy module loaded!")
    log("📊 Slime stats - HP: " .. slime.stats.health .. ", DMG: " .. slime.stats.damage .. ", SPD: " .. slime.stats.speed)
    
    -- Spawn some test slimes
    slime.create_instance(300, 200)
    slime.create_instance(400, 300)
    slime.create_instance(150, 350)
end

-- Update all slime instances
function slime.update()
    for id, instance in pairs(slime.instances) do
        slime.update_instance(instance)
    end
end

-- Update a single slime instance
function slime.update_instance(instance)
    -- Reduce attack cooldown
    if instance.attack_cooldown > 0 then
        instance.attack_cooldown = instance.attack_cooldown - 1
    end
    
    -- AI behavior based on state
    if instance.state == "idle" then
        slime.idle_behavior(instance)
    elseif instance.state == "patrol" then
        slime.patrol_behavior(instance)
    elseif instance.state == "chase" then
        slime.chase_behavior(instance)
    elseif instance.state == "attack" then
        slime.attack_behavior(instance)
    end
    
    -- Update patrol timer
    instance.patrol_timer = instance.patrol_timer + 1
end

-- Idle behavior
function slime.idle_behavior(instance)
    -- Switch to patrol after some time
    if instance.patrol_timer > 120 then -- 2 seconds at 60fps
        instance.state = "patrol"
        instance.patrol_timer = 0
        -- Random direction
        local angle = math.random() * 2 * math.pi
        instance.direction.x = math.cos(angle)
        instance.direction.y = math.sin(angle)
    end
end

-- Patrol behavior
function slime.patrol_behavior(instance)
    -- Move in current direction
    instance.position.x = instance.position.x + instance.direction.x * slime.stats.speed
    instance.position.y = instance.position.y + instance.direction.y * slime.stats.speed
    
    -- Change direction occasionally
    if instance.patrol_timer > 180 then -- 3 seconds
        instance.state = "idle"
        instance.patrol_timer = 0
    end
end

-- Chase behavior (simplified - would need player position)
function slime.chase_behavior(instance)
    -- For now, just move randomly toward a target
    if instance.target then
        local dx = instance.target.x - instance.position.x
        local dy = instance.target.y - instance.position.y
        local distance = math.sqrt(dx*dx + dy*dy)
        
        if distance > 0 then
            instance.direction.x = dx / distance
            instance.direction.y = dy / distance
            
            instance.position.x = instance.position.x + instance.direction.x * slime.stats.speed * 1.5
            instance.position.y = instance.position.y + instance.direction.y * slime.stats.speed * 1.5
        end
        
        -- Switch to attack if close enough
        if distance < 30 then
            instance.state = "attack"
        end
    else
        -- Lost target, go back to patrol
        instance.state = "patrol"
    end
end

-- Attack behavior
function slime.attack_behavior(instance)
    if instance.attack_cooldown <= 0 then
        log("⚔️ Slime #" .. instance.id .. " attacks for " .. slime.stats.damage .. " damage!")
        emit("slime_attack", instance.id)
        instance.attack_cooldown = 60 -- 1 second cooldown
    end
    
    -- Return to chase after attack
    instance.state = "chase"
end

-- Take damage
function slime.take_damage(instance_id, damage)
    local instance = slime.instances[instance_id]
    if instance then
        instance.health = instance.health - damage
        log("💔 Slime #" .. instance_id .. " took " .. damage .. " damage. Health: " .. instance.health)
        
        if instance.health <= 0 then
            slime.kill_slime(instance_id)
        end
    end
end

-- Kill slime
function slime.kill_slime(instance_id)
    local instance = slime.instances[instance_id]
    if instance then
        log("💀 Slime #" .. instance_id .. " died!")
        emit("slime_died", instance_id)
        slime.instances[instance_id] = nil
    end
end

-- Get slime count
function slime.get_count()
    local count = 0
    for _ in pairs(slime.instances) do
        count = count + 1
    end
    return count
end

-- Export the module
return slime
