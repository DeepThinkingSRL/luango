-- World Module
-- Handles world generation, environment, and global game state

local world = {}

-- World configuration
world.config = {
    width = 2000,
    height = 2000,
    tile_size = 32,
    spawn_point = {x = 100, y = 100}
}

-- World state
world.state = {
    time_of_day = 0, -- 0-1440 (minutes in a day)
    weather = "clear", -- clear, rain, storm, fog
    temperature = 20, -- celsius
    entities_spawned = 0,
    max_entities = 50
}

-- Environmental zones
world.zones = {
    {
        name = "Starting Area",
        bounds = {x = 0, y = 0, width = 500, height = 500},
        biome = "grassland",
        spawn_rate = 0.1
    },
    {
        name = "Dark Forest",
        bounds = {x = 500, y = 0, width = 500, height = 800},
        biome = "forest",
        spawn_rate = 0.3
    },
    {
        name = "Mountain Pass",
        bounds = {x = 0, y = 500, width = 1000, height = 500},
        biome = "mountain",
        spawn_rate = 0.2
    }
}

-- Initialize world module
function world.init()
    log("🌍 World module initialized!")
    log("🗺️ World size: " .. world.config.width .. "x" .. world.config.height)
    log("🏠 Spawn point: (" .. world.config.spawn_point.x .. ", " .. world.config.spawn_point.y .. ")")
    log("🌤️ Weather: " .. world.state.weather)
    log("🌡️ Temperature: " .. world.state.temperature .. "°C")
    
    -- Initialize zones
    for i, zone in ipairs(world.zones) do
        log("🗺️ Zone " .. i .. ": " .. zone.name .. " (" .. zone.biome .. ")")
    end
end

-- Update world state
function world.update()
    -- Update time of day (1 minute = 1 frame, so 1 day = 1440 frames = 24 minutes real time at 60fps)
    world.state.time_of_day = world.state.time_of_day + 1
    if world.state.time_of_day >= 1440 then
        world.state.time_of_day = 0
        log("🌅 New day has begun!")
        emit("new_day", "world")
    end
    
    -- Update weather occasionally
    if world.state.time_of_day % 300 == 0 then -- Every 5 minutes
        world.update_weather()
    end
    
    -- Spawn entities based on zones
    if world.state.entities_spawned < world.state.max_entities then
        if math.random() < 0.01 then -- 1% chance per frame
            world.try_spawn_entity()
        end
    end
end

-- Update weather system
function world.update_weather()
    local weather_options = {"clear", "rain", "storm", "fog"}
    local old_weather = world.state.weather
    
    -- Simple weather change logic
    local weather_chance = math.random()
    if weather_chance < 0.6 then
        world.state.weather = "clear"
    elseif weather_chance < 0.8 then
        world.state.weather = "rain"
    elseif weather_chance < 0.95 then
        world.state.weather = "fog"
    else
        world.state.weather = "storm"
    end
    
    if world.state.weather ~= old_weather then
        log("🌤️ Weather changed from " .. old_weather .. " to " .. world.state.weather)
        emit("weather_changed", world.state.weather)
        
        -- Adjust temperature based on weather
        if world.state.weather == "rain" then
            world.state.temperature = world.state.temperature - math.random(2, 5)
        elseif world.state.weather == "storm" then
            world.state.temperature = world.state.temperature - math.random(5, 10)
        elseif world.state.weather == "clear" then
            world.state.temperature = world.state.temperature + math.random(1, 3)
        end
        
        -- Keep temperature in reasonable range
        world.state.temperature = math.max(-10, math.min(35, world.state.temperature))
    end
end

-- Try to spawn an entity in an appropriate zone
function world.try_spawn_entity()
    -- Pick a random zone
    local zone = world.zones[math.random(#world.zones)]
    
    if math.random() < zone.spawn_rate then
        -- Generate spawn position within zone bounds
        local spawn_x = zone.bounds.x + math.random(0, zone.bounds.width)
        local spawn_y = zone.bounds.y + math.random(0, zone.bounds.height)
        
        log("🐾 Spawning entity in " .. zone.name .. " at (" .. spawn_x .. ", " .. spawn_y .. ")")
        emit("entity_spawned", {zone = zone.name, x = spawn_x, y = spawn_y})
        
        world.state.entities_spawned = world.state.entities_spawned + 1
    end
end

-- Get current zone for a position
function world.get_zone_at(x, y)
    for _, zone in ipairs(world.zones) do
        if x >= zone.bounds.x and x <= zone.bounds.x + zone.bounds.width and
           y >= zone.bounds.y and y <= zone.bounds.y + zone.bounds.height then
            return zone
        end
    end
    return nil
end

-- Get time as string
function world.get_time_string()
    local hours = math.floor(world.state.time_of_day / 60)
    local minutes = world.state.time_of_day % 60
    return string.format("%02d:%02d", hours, minutes)
end

-- Get current time of day phase
function world.get_time_phase()
    local hour = math.floor(world.state.time_of_day / 60)
    
    if hour >= 6 and hour < 12 then
        return "morning"
    elseif hour >= 12 and hour < 18 then
        return "afternoon"
    elseif hour >= 18 and hour < 22 then
        return "evening"
    else
        return "night"
    end
end

-- Check if it's safe to travel
function world.is_safe_to_travel()
    local phase = world.get_time_phase()
    local safe_weather = world.state.weather == "clear" or world.state.weather == "fog"
    local safe_time = phase ~= "night"
    
    return safe_weather and safe_time
end

-- Get world information
function world.get_info()
    return {
        time = world.get_time_string(),
        phase = world.get_time_phase(),
        weather = world.state.weather,
        temperature = world.state.temperature,
        entities_count = world.state.entities_spawned,
        safe_travel = world.is_safe_to_travel()
    }
end

-- Export the module
return world
