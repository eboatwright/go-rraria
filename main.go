package main


import(
	ray "github.com/gen2brain/raylib-go/raylib"
	perlinNoise "github.com/aquilax/go-perlin"
	"math"
	"math/rand"
)


const WINDOW_WIDTH = 960
const WINDOW_HEIGHT = 600
const WINDOW_TITLE = "GO-rraria - eboatwright"
const FPS = 60
const SCREEN_SCALE = 3.0

const WORLD_WIDTH = 128
const WORLD_HEIGHT = 256

const WORLD_GEN_ALPHA = 4
const WORLD_GEN_BETA = 2
const WORLD_GEN_N = 7

const SURFACE_OFFSET = 96
const SURFACE_ROUGHNESS = 0.018
const SURFACE_HILLINESS = 20

const STONE_OFFSET = 8
const STONE_ROUGHNESS = 0.011
const STONE_HILLINESS = 24

const TREE_PROBABILITY = 16
const MIN_TREE_HEIGHT = 8
const MAX_TREE_HEIGHT = 12


type Level struct {
	values [WORLD_HEIGHT][WORLD_WIDTH]int
	tileSize float32
	lastTreeXPosition int
}

type Item struct {
	_type string
	textureIndex int
}

type Slots struct {
	texture ray.Texture2D
	amount int
	selected int
	inventory [8]Item
	itemHolding int
}

type Game struct {
	camera ray.Camera2D

	level Level
	player Player
	slots Slots

	seed int64

	debug bool
	debugReleased bool

	skyTexture ray.Texture2D
}


func generateLevel(game *Game) {
	perlin := perlinNoise.NewPerlin(WORLD_GEN_ALPHA, WORLD_GEN_BETA, WORLD_GEN_N, game.seed)
	for x := 0; x < len(game.level.values[0]); x++ {
		noise := int(perlin.Noise1D(float64(x) * SURFACE_ROUGHNESS) * SURFACE_HILLINESS) + SURFACE_OFFSET

		if noise < 0 { noise = 0 }
		if noise > len(game.level.values) - 1 { noise = len(game.level.values) - 1 }

		if x == int(WORLD_WIDTH / 2) {
			game.player.position = ray.Vector2 { WORLD_WIDTH * game.level.tileSize / 2 - 8, float32(noise) * game.level.tileSize }
		}

		for y := 0; y < len(game.level.values); y++ {
			if y == noise {
				game.level.values[y][x] = GRASS_BLOCK
				if rand.Intn(TREE_PROBABILITY) == 1 && x - game.level.lastTreeXPosition > 3 {
					game.level.lastTreeXPosition = x
					if x - 1 < 0 || x + 1 >= len(game.level.values[0]) { continue }
					game.level.values[y - 1][x] = TREE_TRUNK
					treeHeight := rand.Intn(MAX_TREE_HEIGHT - MIN_TREE_HEIGHT) + MIN_TREE_HEIGHT
					for i := 0; i < treeHeight; i++ {
						_type := rand.Intn(3)
						if _type == 0 {
							game.level.values[y - 2 - i][x] = TREE_NO_BRANCH
						} else if _type == 1 {
							game.level.values[y - 2 - i][x] = TREE_BRANCH_RIGHT
							game.level.values[y - 2 - i][x - 1] = TREE_LEAVES_LEFT
						} else if _type == 2 {
							game.level.values[y - 2 - i][x] = TREE_BRANCH_LEFT
							game.level.values[y - 2 - i][x + 1] = TREE_LEAVES_RIGHT
						}
					}
					game.level.values[y - 2 - treeHeight][x - 1] = TREE_TOP_BOTTOM_LEFT
					game.level.values[y - 2 - treeHeight][x] = TREE_TOP_BOTTOM
					game.level.values[y - 2 - treeHeight][x + 1] = TREE_TOP_BOTTOM_RIGHT
					game.level.values[y - 3 - treeHeight][x - 1] = TREE_TOP_MIDDLE_LEFT
					game.level.values[y - 3 - treeHeight][x] = TREE_TOP_MIDDLE
					game.level.values[y - 3 - treeHeight][x + 1] = TREE_TOP_MIDDLE_RIGHT
					game.level.values[y - 4 - treeHeight][x - 1] = TREE_TOP_TOP_LEFT
					game.level.values[y - 4 - treeHeight][x] = TREE_TOP_TOP
					game.level.values[y - 4 - treeHeight][x + 1] = TREE_TOP_TOP_RIGHT
				}
			}
			if y > noise {
				game.level.values[y][x] = DIRT_BLOCK
			}

			stoneNoise := int(perlin.Noise1D(float64(x) * STONE_ROUGHNESS) * STONE_HILLINESS) + SURFACE_OFFSET + STONE_OFFSET

			if stoneNoise < 0 { stoneNoise = 0 }
			if stoneNoise > len(game.level.values) - 1 { stoneNoise = len(game.level.values) - 1 }

			if y > stoneNoise {
				game.level.values[y][x] = STONE_BLOCK
			}
		}

		game.level.values[WORLD_HEIGHT - 1][x] = BEDROCK_BLOCK
	}
}

func item(_type string, textureIndex int) Item {
	i := Item {
		_type: _type,
		textureIndex: textureIndex,
	}
	return i
}


func updateCamera(game *Game) {
	game.camera.Target = ray.Vector2Subtract(game.player.position, ray.Vector2 { (WINDOW_WIDTH / SCREEN_SCALE / 2) - 8, (WINDOW_HEIGHT / SCREEN_SCALE / 2) - 8 })

	if game.camera.Target.X < 0 { game.camera.Target.X = 0 }
	if game.camera.Target.Y < 0 { game.camera.Target.Y = 0 }

	if game.camera.Target.X > WORLD_WIDTH * 8 - WINDOW_WIDTH / SCREEN_SCALE { game.camera.Target.X = WORLD_WIDTH * 8 - WINDOW_WIDTH / SCREEN_SCALE }
	if game.camera.Target.Y > WORLD_HEIGHT * 8 - WINDOW_HEIGHT / SCREEN_SCALE { game.camera.Target.Y = WORLD_HEIGHT * 8 - WINDOW_HEIGHT / SCREEN_SCALE }
}

func update(game *Game) {
	updatePlayer(game)
	updateCamera(game)

	if ray.IsKeyDown(ray.KeyF3) {
		if game.debugReleased {
			game.debug = !game.debug
		}
		game.debugReleased = false
	} else {
		game.debugReleased = true
	}

	for i := 0; i < 8; i++ {
		if ray.IsKeyDown(int32(49 + i)) {
			game.slots.selected = i
		}
	}
}


func drawLevel(game *Game) {
	tileSize := game.level.tileSize
	yStart := int(float32(math.Floor(float64(game.camera.Target.Y))) / game.level.tileSize)
	xStart := int(float32(math.Floor(float64(game.camera.Target.X))) / game.level.tileSize)
	for y := yStart; y < yStart + 26; y++ {
		if y >= len(game.level.values) { continue }
		for x := xStart; x < xStart + 41; x++ {
			if x >= len(game.level.values[0]) { continue }
			if game.level.values[y][x] != AIR_BLOCK {
				ray.DrawTextureV(TEXTURES[game.level.values[y][x]], ray.Vector2 { float32(x) * tileSize, float32(y) * tileSize }, ray.White)
			}
		}
	}
}

func drawUI(game *Game) {
	for i := 0; i < game.slots.amount; i++ {
		x := 0

		if i == game.slots.selected { x = 16 }

		position := ray.Vector2Add(game.camera.Target, ray.Vector2 { 1 + float32(i) * 18, 1 })
		ray.DrawTextureRec(game.slots.texture, ray.Rectangle { float32(x), 0, 16, 16 }, position, ray.White)
		if game.slots.inventory[i]._type != "empty" {
			ray.DrawTextureV(TEXTURES[game.slots.inventory[i].textureIndex], ray.Vector2Add(position, ray.Vector2 { 4, 4 }), ray.White)
		}
	}
}

func draw(game *Game) {
	ray.DrawTextureV(game.skyTexture, game.camera.Target, ray.White)
	drawLevel(game)
	drawPlayer(game)
	drawUI(game)
}


func main() {
	ray.InitWindow(WINDOW_WIDTH, WINDOW_HEIGHT, WINDOW_TITLE)
	ray.InitAudioDevice()
	ray.SetTargetFPS(FPS)

	camera := ray.Camera2D {  }
	camera.Zoom = SCREEN_SCALE

	initializeTextures()

	game := Game {
		camera: camera,
		level: Level {
			tileSize: 8,
		},
		player: Player {
			texture: ray.LoadTexture("assets/img/player.png"),

			moveSpeed: 0.6,
			jumpHeight: -4.4,

			xOffset: 5,
			yOffset: 2,
			width: 7,
			height: 14,

			xScale: 1,

			currentAnimation: []float32{0, 0},
			animations: [][]float32 {
				{ 0, 0 },
				{ 1, 2, 3, 4, 5, 6, 0.1 },
				{ 7, 0 },
				{ 8, 9, 10, 0.1 },
			},

			footstepTime: 0.3,

			landed: true,

			jumpSfx: ray.LoadSound("assets/sfx/jump4.wav"),
			footstepSfx: ray.LoadSound("assets/sfx/footstep4.wav"),
		},
		slots: Slots {
			texture: ray.LoadTexture("assets/img/ui/slot.png"),
			amount: 8,
			inventory: [8]Item{ item("pickaxe", OBSIDIAN_PICKAXE), item("axe", OBSIDIAN_AXE), item("sword", OBSIDIAN_SWORD), item("block", DIRT_BLOCK), item("block", GRASS_BLOCK), item("block", STONE_BLOCK), item("empty", 0), item("empty", 0) },
		},

		seed: 176,

		skyTexture: ray.LoadTexture("assets/img/sky.png"),
	}
	rand.Seed(game.seed)

	generateLevel(&game)

	for !ray.WindowShouldClose() {
		update(&game)

		ray.BeginDrawing()
		ray.ClearBackground(ray.RayWhite)

		ray.BeginMode2D(game.camera)

		draw(&game)

		ray.EndMode2D()

		if game.debug {
			ray.DrawFPS(1, 1)
		}

		ray.EndDrawing()
	}

	ray.CloseWindow()
}