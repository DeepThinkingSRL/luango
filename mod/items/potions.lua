log("🧪 Potions module ready!")

potions = {
  health = {
    name = "Health Potion",
    restore = 25
  }
}

function use_potion()
  player.hp = player.hp + potions.health.restore
  log("💖 Used a health potion! HP: " .. player.hp)
end
