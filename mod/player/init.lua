log("👤 Player module loaded!")

player = {
  name = "Hero",
  hp = 100,
  speed = 2
}

function player.update()
  if is_key_pressed("ArrowRight") then move_player(player.speed, 0) end
  if is_key_pressed("ArrowLeft")  then move_player(-player.speed, 0) end
  if is_key_pressed("ArrowUp")    then move_player(0, -player.speed) end
  if is_key_pressed("ArrowDown")  then move_player(0, player.speed) end
end
