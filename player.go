package main


import(
	ray "github.com/gen2brain/raylib-go/raylib"
	"math"
)


const REACH_DISTANCE = 5.0
const FRICTION = 0.2
const GRAVITY = 0.25
const JUMP_CUTOFF = 0.64


type Player struct {
	texture ray.Texture2D

	position ray.Vector2
	velocity ray.Vector2

	moveSpeed float32
	jumpHeight float32

	grounded float32

	jumpReleased bool
	jumpJustReleased bool

	xOffset float32
	yOffset float32
	width float32
	height float32

	xScale float32

	animationTimer float32
	animationFrameIndex int
	animations [][]float32
	currentAnimation []float32

	footstepTime float32
	footstepTimer float32

	landed bool

	jumpSfx ray.Sound
	footstepSfx ray.Sound
}


func changeAnimation(game *Game, index int) {
	if game.player.animations[index][0] != game.player.currentAnimation[0] {
		game.player.currentAnimation = game.player.animations[index]
		game.player.animationFrameIndex = 0
		game.player.animationTimer = 0
	}
}

func GetMousePositionAsTileIndex(game *Game) (int, int, bool) {
	mousePosition := ray.GetMousePosition()
	mousePosition = ray.Vector2DivideV(mousePosition, ray.Vector2 { SCREEN_SCALE, SCREEN_SCALE })
	mousePosition = ray.Vector2Add(mousePosition, game.camera.Target)

	xPosition := int(math.Floor(float64(mousePosition.X))) / 8
	yPosition := int(math.Floor(float64(mousePosition.Y))) / 8

	if xPosition < 0 || xPosition > len(game.level.values[0]) - 1 || yPosition < 0 || yPosition > len(game.level.values) - 1 {
		return 0, 0, true
	}
	
	return xPosition, yPosition, false
}

func updatePlayer(game *Game) {
	newIndex := 0

	var inputX float32 = 0
	if ray.IsKeyDown(ray.KeyA) || ray.IsKeyDown(ray.KeyLeft) {
		inputX = -1
		game.player.xScale = -1
		newIndex = 1
	} else if ray.IsKeyDown(ray.KeyD) || ray.IsKeyDown(ray.KeyRight) {
		inputX = 1
		game.player.xScale = 1
		newIndex = 1
	} else {
		newIndex = 0
	}

	game.player.velocity.Y += GRAVITY
	game.player.position.Y += game.player.velocity.Y
	
	playerRectangle := ray.Rectangle { game.player.position.X + game.player.xOffset, game.player.position.Y + game.player.yOffset, game.player.width, game.player.height }
	tileRectangle := ray.Rectangle { 0, 0, game.level.tileSize, game.level.tileSize }

	yStart := int(float32(math.Floor(float64(game.camera.Target.Y))) / game.level.tileSize)
	xStart := int(float32(math.Floor(float64(game.camera.Target.X))) / game.level.tileSize)

	game.player.grounded -= ray.GetFrameTime()
	for y := yStart; y < yStart + 26; y++ {
		if y >= len(game.level.values) { continue }
		for x := xStart; x < xStart + 41; x++ {
			if x >= len(game.level.values[0]) { continue }
			if game.level.values[y][x] > AIR_BLOCK && game.level.values[y][x] < TREE_TRUNK {
				tileRectangle.X = float32(x) * game.level.tileSize
				tileRectangle.Y = float32(y) * game.level.tileSize
				if ray.CheckCollisionRecs(playerRectangle, tileRectangle) {
					if game.player.velocity.Y < 0 {
						game.player.position.Y = tileRectangle.Y + tileRectangle.Height - game.player.yOffset;
					}
					if game.player.velocity.Y > 0 {
						game.player.grounded = 0.2
						game.player.position.Y = tileRectangle.Y - game.player.height - game.player.yOffset;
					}
					game.player.velocity.Y = 0
				}
			}
		}
	}
	
	game.player.velocity.X += inputX * game.player.moveSpeed
	game.player.velocity.X *= 1 - FRICTION
	game.player.position.X += game.player.velocity.X
	
	playerRectangle.X = game.player.position.X + game.player.xOffset;
	playerRectangle.Y = game.player.position.Y + game.player.yOffset;

	for y := yStart; y < yStart + 26; y++ {
		if y >= len(game.level.values) { continue }
		for x := xStart; x < xStart + 41; x++ {
			if x >= len(game.level.values[0]) { continue }
			if game.level.values[y][x] > AIR_BLOCK && game.level.values[y][x] < TREE_TRUNK {
				tileRectangle.X = float32(x) * game.level.tileSize
				tileRectangle.Y = float32(y) * game.level.tileSize
				if ray.CheckCollisionRecs(playerRectangle, tileRectangle) {
					if game.player.velocity.X < 0 {
						game.player.position.X = tileRectangle.X + tileRectangle.Width - game.player.xOffset;
					}
					if game.player.velocity.X > 0 {
						game.player.position.X = tileRectangle.X - game.player.width - game.player.xOffset;
					}
					game.player.velocity.X = 0
				}
			}
		}
	}

	if game.player.position.X < -game.level.tileSize { game.player.position.X = -game.level.tileSize }
	if game.player.position.X > (WORLD_WIDTH - 1) * game.level.tileSize { game.player.position.X = (WORLD_WIDTH - 1) * game.level.tileSize }


	if ray.IsKeyDown(ray.KeyW) || ray.IsKeyDown(ray.KeyUp) {
		if game.player.grounded > 0 && game.player.jumpReleased {
			game.player.grounded = 0
			game.player.velocity.Y = game.player.jumpHeight
			ray.PlaySound(game.player.jumpSfx)
		}
		game.player.jumpReleased = false
		game.player.jumpJustReleased = false
	} else {
		if !game.player.jumpJustReleased && game.player.velocity.Y < 0 {
			game.player.velocity.Y *= JUMP_CUTOFF
		}
		game.player.jumpReleased = true
		game.player.jumpJustReleased = true
	}

	if game.player.grounded <= 0 {
		newIndex = 2
		game.player.footstepTimer = game.player.footstepTime
		game.player.landed = false
	} else {
		if inputX != 0 {
			if game.player.footstepTimer <= 0 {
				game.player.footstepTimer = game.player.footstepTime
				ray.PlaySound(game.player.footstepSfx)
			}
			game.player.footstepTimer -= ray.GetFrameTime()
		} else {
			game.player.footstepTimer = game.player.footstepTime
		}

		if !game.player.landed {
			ray.PlaySound(game.player.footstepSfx)
		}
		game.player.landed = true
	}

	if ray.IsMouseButtonDown(ray.MouseLeftButton) {
		xPosition, yPosition, err := GetMousePositionAsTileIndex(game)
		if !err {
			newIndex = 3
			if ray.Vector2Distance(ray.Vector2 { float32(xPosition), float32(yPosition) }, ray.Vector2Add(ray.Vector2DivideV(game.player.position, ray.Vector2 { 8, 8 }), ray.Vector2 { 1, 1 })) < REACH_DISTANCE {	
				mineBlock(xPosition, yPosition, game)
			}
		}
	}

	changeAnimation(game, newIndex)

	game.player.animationTimer += ray.GetFrameTime()
	if game.player.animationTimer >= game.player.currentAnimation[len(game.player.currentAnimation) - 1] {
		game.player.animationFrameIndex += 1
		if game.player.animationFrameIndex >= len(game.player.currentAnimation) - 1 {
			game.player.animationFrameIndex = 0
		}
		game.player.animationTimer = 0
	}
}

func mineBlock(xPosition int, yPosition int, game *Game) {
	if game.slots.inventory[game.slots.selected]._type == "block" {
		if game.level.values[yPosition][xPosition] == AIR_BLOCK {
			if (yPosition - 1 > 0 && game.level.values[yPosition - 1][xPosition] != AIR_BLOCK) ||
				(yPosition + 1 < len(game.level.values) && game.level.values[yPosition + 1][xPosition] != AIR_BLOCK) ||
				(xPosition - 1 > 0 && game.level.values[yPosition][xPosition - 1] != AIR_BLOCK) ||
				(xPosition + 1 < len(game.level.values[0]) && game.level.values[yPosition][xPosition + 1] != AIR_BLOCK) {
				playerRectangle := ray.Rectangle { (game.player.position.X + game.player.xOffset) - 4, (game.player.position.Y + game.player.yOffset) - 4, game.player.width + 8, game.player.height + 8 }
				mouseRectangle := ray.Rectangle { float32(ray.GetMouseX()) / SCREEN_SCALE + game.camera.Target.X - 4, float32(ray.GetMouseY()) / SCREEN_SCALE + game.camera.Target.Y - 4, 8, 8 }
				if !ray.CheckCollisionRecs(playerRectangle, mouseRectangle) {
					game.level.values[yPosition][xPosition] = game.slots.inventory[game.slots.selected].textureIndex
					if game.level.values[yPosition + 1][xPosition] == GRASS_BLOCK {
						game.level.values[yPosition + 1][xPosition] = DIRT_BLOCK
					}
				}
			}
		}
	} else {
		if game.level.values[yPosition][xPosition] != BEDROCK_BLOCK {
			if (game.level.values[yPosition][xPosition] > AIR_BLOCK &&
				game.level.values[yPosition][xPosition] < TREE_TRUNK) || (
				game.level.values[yPosition][xPosition] > TREE_TOP_TOP_RIGHT) {
				
				if game.slots.inventory[game.slots.selected]._type == "pickaxe" {
					game.level.values[yPosition][xPosition] = AIR_BLOCK
				}
			} else {
				if game.slots.inventory[game.slots.selected]._type == "axe" &&
					game.level.values[yPosition][xPosition] >= TREE_TRUNK &&
					game.level.values[yPosition][xPosition] < TREE_LEAVES_LEFT {
					i := 0
					for game.level.values[yPosition - i][xPosition] >= TREE_TRUNK && game.level.values[yPosition - i][xPosition] <= TREE_TOP_TOP_RIGHT {
						mineTreeBlock(yPosition - i, xPosition, game)
						i++
					}
				}
			}
		}
	}
}

func mineTreeBlock(y int, x int, game *Game) {
	if game.level.values[y][x] < TREE_TRUNK || game.level.values[y][x] > TREE_TOP_TOP_RIGHT {
		return
	}
	game.level.values[y][x] = AIR_BLOCK
	if game.level.values[y][x - 1] >= TREE_TRUNK && game.level.values[y][x - 1] <= TREE_TOP_TOP_RIGHT {
		game.level.values[y][x - 1] = AIR_BLOCK
	}
	if game.level.values[y][x + 1] >= TREE_TRUNK && game.level.values[y][x + 1] <= TREE_TOP_TOP_RIGHT {
		game.level.values[y][x + 1] = AIR_BLOCK
	}
}

func drawPlayer(game *Game) {
	ray.DrawTextureRec(game.player.texture, ray.Rectangle { game.player.currentAnimation[game.player.animationFrameIndex] * 16.0, 0.0, 16 * game.player.xScale, 16 }, game.player.position, ray.White)
}