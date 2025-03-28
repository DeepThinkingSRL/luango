function on_start()
	log("🌟 Game is starting, mi amor!")
  end
  
  function on_update()
	world.update()
	player.update()
	slime.update()
  end
  