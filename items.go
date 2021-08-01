package main

import(
	ray "github.com/gen2brain/raylib-go/raylib"
)

const AIR_BLOCK = 0
const DIRT_BLOCK = 1
const GRASS_BLOCK = 2
const STONE_BLOCK = 3
const BEDROCK_BLOCK = 4

const WOODEN_PICKAXE = 5
const STONE_PICKAXE = 6
const IRON_PICKAXE = 7
const GOLDEN_PICKAXE = 8
const DIAMOND_PICKAXE = 9
const OBSIDIAN_PICKAXE = 10

const WOODEN_AXE = 11
const STONE_AXE = 12
const IRON_AXE = 13
const GOLDEN_AXE = 14
const DIAMOND_AXE = 15
const OBSIDIAN_AXE = 16

const WOODEN_SWORD = 17
const STONE_SWORD = 18
const IRON_SWORD = 19
const GOLDEN_SWORD = 20
const DIAMOND_SWORD = 21
const OBSIDIAN_SWORD = 22

const TREE_TRUNK = 23
const TREE_NO_BRANCH = 24
const TREE_BRANCH_LEFT = 25
const TREE_BRANCH_RIGHT = 26
const TREE_LEAVES_LEFT = 27
const TREE_LEAVES_RIGHT = 28
const TREE_TOP_BOTTOM = 29
const TREE_TOP_BOTTOM_LEFT = 30
const TREE_TOP_BOTTOM_RIGHT = 31
const TREE_TOP_MIDDLE = 32
const TREE_TOP_MIDDLE_LEFT = 33
const TREE_TOP_MIDDLE_RIGHT = 34
const TREE_TOP_TOP = 35
const TREE_TOP_TOP_LEFT = 36
const TREE_TOP_TOP_RIGHT = 37

var TEXTURES = [64]ray.Texture2D {}


func initializeTextures() {
	TEXTURES[AIR_BLOCK] = ray.LoadTexture("assets/img/missing.png")
	TEXTURES[DIRT_BLOCK] = ray.LoadTexture("assets/img/blocks/dirt.png")
	TEXTURES[GRASS_BLOCK] = ray.LoadTexture("assets/img/blocks/grass.png")
	TEXTURES[STONE_BLOCK] = ray.LoadTexture("assets/img/blocks/stone.png")
	TEXTURES[BEDROCK_BLOCK] = ray.LoadTexture("assets/img/blocks/bedrock.png")

	TEXTURES[WOODEN_PICKAXE] = ray.LoadTexture("assets/img/items/woodenPickaxe.png")
	TEXTURES[STONE_PICKAXE] = ray.LoadTexture("assets/img/items/stonePickaxe.png")
	TEXTURES[IRON_PICKAXE] = ray.LoadTexture("assets/img/items/ironPickaxe.png")
	TEXTURES[GOLDEN_PICKAXE] = ray.LoadTexture("assets/img/items/goldenPickaxe.png")
	TEXTURES[DIAMOND_PICKAXE] = ray.LoadTexture("assets/img/items/diamondPickaxe.png")
	TEXTURES[OBSIDIAN_PICKAXE] = ray.LoadTexture("assets/img/items/obsidianPickaxe.png")

	TEXTURES[WOODEN_AXE] = ray.LoadTexture("assets/img/items/woodenAxe.png")
	TEXTURES[STONE_AXE] = ray.LoadTexture("assets/img/items/stoneAxe.png")
	TEXTURES[IRON_AXE] = ray.LoadTexture("assets/img/items/ironAxe.png")
	TEXTURES[GOLDEN_AXE] = ray.LoadTexture("assets/img/items/goldenAxe.png")
	TEXTURES[DIAMOND_AXE] = ray.LoadTexture("assets/img/items/diamondAxe.png")
	TEXTURES[OBSIDIAN_AXE] = ray.LoadTexture("assets/img/items/obsidianAxe.png")

	TEXTURES[WOODEN_SWORD] = ray.LoadTexture("assets/img/items/woodenSword.png")
	TEXTURES[STONE_SWORD] = ray.LoadTexture("assets/img/items/stoneSword.png")
	TEXTURES[IRON_SWORD] = ray.LoadTexture("assets/img/items/ironSword.png")
	TEXTURES[GOLDEN_SWORD] = ray.LoadTexture("assets/img/items/goldenSword.png")
	TEXTURES[DIAMOND_SWORD] = ray.LoadTexture("assets/img/items/diamondSword.png")
	TEXTURES[OBSIDIAN_SWORD] = ray.LoadTexture("assets/img/items/obsidianSword.png")
	
	TEXTURES[TREE_TRUNK] = ray.LoadTexture("assets/img/blocks/tree/trunk.png")
	TEXTURES[TREE_NO_BRANCH] = ray.LoadTexture("assets/img/blocks/tree/noBranch.png")
	TEXTURES[TREE_BRANCH_LEFT] = ray.LoadTexture("assets/img/blocks/tree/branchRight.png")
	TEXTURES[TREE_BRANCH_RIGHT] = ray.LoadTexture("assets/img/blocks/tree/branchLeft.png")
	TEXTURES[TREE_LEAVES_LEFT] = ray.LoadTexture("assets/img/blocks/tree/leavesLeft.png")
	TEXTURES[TREE_LEAVES_RIGHT] = ray.LoadTexture("assets/img/blocks/tree/leavesRight.png")
	TEXTURES[TREE_TOP_BOTTOM] = ray.LoadTexture("assets/img/blocks/tree/topBottom.png")
	TEXTURES[TREE_TOP_BOTTOM_LEFT] = ray.LoadTexture("assets/img/blocks/tree/topBottomLeft.png")
	TEXTURES[TREE_TOP_BOTTOM_RIGHT] = ray.LoadTexture("assets/img/blocks/tree/topBottomRight.png")
	TEXTURES[TREE_TOP_MIDDLE] = ray.LoadTexture("assets/img/blocks/tree/topMiddle.png")
	TEXTURES[TREE_TOP_MIDDLE_LEFT] = ray.LoadTexture("assets/img/blocks/tree/topMiddleLeft.png")
	TEXTURES[TREE_TOP_MIDDLE_RIGHT] = ray.LoadTexture("assets/img/blocks/tree/topMiddleRight.png")
	TEXTURES[TREE_TOP_TOP] = ray.LoadTexture("assets/img/blocks/tree/topTop.png")
	TEXTURES[TREE_TOP_TOP_LEFT] = ray.LoadTexture("assets/img/blocks/tree/topTopLeft.png")
	TEXTURES[TREE_TOP_TOP_RIGHT] = ray.LoadTexture("assets/img/blocks/tree/topTopRight.png")
}