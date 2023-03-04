package main

import (
    "github.com/curzodo/perlin"
    "github.com/df-mc/dragonfly/server"
    "github.com/df-mc/dragonfly/server/player"
    "github.com/df-mc/dragonfly/server/block"
    "github.com/df-mc/dragonfly/server/world"
    "github.com/df-mc/dragonfly/server/world/chunk"
    "gonum.org/v1/gonum/interp"
    "math/rand"
)

func main() {
    conf, _ := server.DefaultConfig().Config(nil)

    gen := NewGenerator(12345)

    conf.Generator = func(world.Dimension) world.Generator {
        return gen
    }

	srv := conf.New()
	srv.CloseOnProgramEnd()

	srv.Listen()
	for srv.Accept(func(p *player.Player) {
        // p.Teleport([3]float64{10000, 64, 10000})
    }) {
	}
}

type Generator struct {
    perlin.Generator
}

func NewGenerator(seed int64) world.Generator {
    return Generator { perlin.NewGenerator(seed) }
}

var grass = world.BlockRuntimeID(block.Grass{})
var log = world.BlockRuntimeID(block.Log{})
var leaves = world.BlockRuntimeID(block.Leaves{})
var dirt = world.BlockRuntimeID(block.Dirt{})
var stone = world.BlockRuntimeID(block.Stone{})
var water = world.BlockRuntimeID(block.Water{
    Still: true,
    Depth: 1,
    Falling: true,
})

var xs = []float64{-2.0, 0.0, 2.0}
var ys = []float64{0, 3, 6}
var pc interp.ClampedCubic
var _ = pc.Fit(xs, ys)

func (g Generator) GenerateChunk(pos world.ChunkPos, chunk *chunk.Chunk) {
    for x := 0; x < 16; x++ {
            for z := 0; z < 16; z++ {
                worldx := float64(x + 15 * int(pos[0])) + 10000.01
                worldz := float64(z + 15 * int(pos[1])) + 10000.01

                spread := 40.0
                octaves := 3
                amplitude := 2.00
                persistence := 0.5
                lacunarity := 2.0

                noise := 0.0

                for i := 0; i < octaves; i++ {
                    noise += amplitude * g.Noise2D(worldx/spread, worldz/spread)
                    amplitude *= persistence
                    spread *= 1/lacunarity
                }

                y := pc.Predict(noise)

                chunk.SetBlock(uint8(x), int16(y), uint8(z), 0, grass)

                treeNoise := g.Noise2D(worldx/25, worldz/25)

                // Spawn tree
                if treeNoise > 0.15 && rand.Intn(100) == 0 {
                }
            }
    }
}
