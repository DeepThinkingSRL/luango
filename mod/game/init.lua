-- Game Module
-- Central game logic and coordination between other modules

local game = {}

-- Game state
game.state = {
    mode = "menu", -- menu, playing, paused, game_over
    level = 1,
    score = 0,
    lives = 3,
    started_time = 0,
    elapsed_time = 0
}

-- Game settings
game.settings = {
    difficulty = "normal", -- easy, normal, hard
    sound_enabled = true,
    music_enabled = true,
    auto_save = true
}

-- Initialize game module
function game.init()
    log("🎮 Game module initialized!")
    log("⚙️ Difficulty: " .. game.settings.difficulty)
    log("🔊 Sound: " .. (game.settings.sound_enabled and "enabled" or "disabled"))
    log("🎵 Music: " .. (game.settings.music_enabled and "enabled" or "disabled"))
    
    game.state.started_time = os.time()
end

-- Update game logic
function game.update()
    if game.state.mode == "playing" then
        game.state.elapsed_time = os.time() - game.state.started_time
        
        -- Example game logic
        if game.state.elapsed_time % 60 == 0 and game.state.elapsed_time > 0 then
            game.state.score = game.state.score + 10
            log("⭐ Score increased! Current score: " .. game.state.score)
        end
        
        -- Level progression
        if game.state.score > 0 and game.state.score % 100 == 0 then
            local new_level = math.floor(game.state.score / 100) + 1
            if new_level > game.state.level then
                game.advance_level(new_level)
            end
        end
    end
end

-- Start the game
function game.start()
    if game.state.mode == "menu" or game.state.mode == "paused" then
        game.state.mode = "playing"
        game.state.started_time = os.time()
        log("🎮 Game started!")
        emit("game_started", game.state.level)
    end
end

-- Pause the game
function game.pause()
    if game.state.mode == "playing" then
        game.state.mode = "paused"
        log("⏸️ Game paused")
        emit("game_paused", "user_action")
    end
end

-- Resume the game
function game.resume()
    if game.state.mode == "paused" then
        game.state.mode = "playing"
        log("▶️ Game resumed")
        emit("game_resumed", "user_action")
    end
end

-- Game over
function game.game_over(reason)
    game.state.mode = "game_over"
    log("💀 Game Over! Reason: " .. reason)
    log("📊 Final Score: " .. game.state.score)
    log("⏱️ Time Played: " .. game.state.elapsed_time .. " seconds")
    emit("game_over", {score = game.state.score, time = game.state.elapsed_time, reason = reason})
end

-- Restart the game
function game.restart()
    game.state.mode = "menu"
    game.state.level = 1
    game.state.score = 0
    game.state.lives = 3
    game.state.started_time = 0
    game.state.elapsed_time = 0
    
    log("🔄 Game restarted")
    emit("game_restarted", "user_action")
end

-- Advance to next level
function game.advance_level(new_level)
    game.state.level = new_level
    log("🎊 Level Up! Now on level " .. game.state.level)
    emit("level_advanced", game.state.level)
    
    -- Bonus score for advancing
    game.state.score = game.state.score + (game.state.level * 50)
    log("🎁 Level bonus: " .. (game.state.level * 50) .. " points")
end

-- Add score
function game.add_score(points)
    game.state.score = game.state.score + points
    log("⭐ +" .. points .. " points! Total: " .. game.state.score)
    emit("score_added", points)
end

-- Lose a life
function game.lose_life()
    game.state.lives = game.state.lives - 1
    log("💔 Lost a life! Lives remaining: " .. game.state.lives)
    
    if game.state.lives <= 0 then
        game.game_over("no_lives_remaining")
    else
        emit("life_lost", game.state.lives)
    end
end

-- Gain a life
function game.gain_life()
    game.state.lives = game.state.lives + 1
    log("💖 Gained a life! Lives: " .. game.state.lives)
    emit("life_gained", game.state.lives)
end

-- Change difficulty
function game.set_difficulty(difficulty)
    if difficulty == "easy" or difficulty == "normal" or difficulty == "hard" then
        game.settings.difficulty = difficulty
        log("⚙️ Difficulty changed to: " .. difficulty)
        emit("difficulty_changed", difficulty)
    else
        log("❌ Invalid difficulty: " .. difficulty)
    end
end

-- Toggle sound
function game.toggle_sound()
    game.settings.sound_enabled = not game.settings.sound_enabled
    log("🔊 Sound " .. (game.settings.sound_enabled and "enabled" or "disabled"))
    emit("sound_toggled", game.settings.sound_enabled)
end

-- Toggle music
function game.toggle_music()
    game.settings.music_enabled = not game.settings.music_enabled
    log("🎵 Music " .. (game.settings.music_enabled and "enabled" or "disabled"))
    emit("music_toggled", game.settings.music_enabled)
end

-- Get game statistics
function game.get_stats()
    return {
        mode = game.state.mode,
        level = game.state.level,
        score = game.state.score,
        lives = game.state.lives,
        time_played = game.state.elapsed_time,
        difficulty = game.settings.difficulty
    }
end

-- Save game state (simplified)
function game.save()
    if game.settings.auto_save then
        log("💾 Game saved automatically")
        emit("game_saved", game.get_stats())
    end
end

-- Load game state (simplified)
function game.load()
    log("📁 Game loaded")
    emit("game_loaded", "save_file")
end

-- Export the module
return game
