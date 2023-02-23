package main

import (
    "github.com/curzodo/perlin"
    "github.com/df-mc/dragonfly/server"
    "github.com/df-mc/dragonfly/server/block"
    "github.com/df-mc/dragonfly/server/world"
    "github.com/df-mc/dragonfly/server/world/chunk"
    // "math"
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
	for srv.Accept(nil) {
	}
}

type Generator struct {
    perlin.Generator
}

func NewGenerator(seed int64) world.Generator {
    return Generator { perlin.NewGenerator(seed) }
}

var grass = world.BlockRuntimeID(block.Grass{})
var stone = world.BlockRuntimeID(block.Stone{})
var water = world.BlockRuntimeID(block.Water{
    Still: true,
    Depth: 1,
    Falling: true,
})

func (g Generator) GenerateChunk(pos world.ChunkPos, chunk *chunk.Chunk) {
    for x := 0; x < 16; x++ {
        for z := 0; z < 16; z++ {
            worldx := (float64(x + 15 * int(pos[0])) + 0.01) + 10000
            worldz := (float64(z + 15 * int(pos[1])) + 0.01) + 10000
            
            spread := 500.0
            amplitude := 20.0
            octaves := 8
            persistence := 1.0
            lacunarity := 2.0

            var noise float64
            for i := 0; i < octaves; i++ {
                noise += amplitude * (g.Noise2D(worldx/spread, worldz/spread) + 1)
                amplitude *= persistence
                spread *= 1/lacunarity
            }

            chunk.SetBlock(uint8(x), int16(noise)-1, uint8(z), 0, stone)
            chunk.SetBlock(uint8(x), int16(noise), uint8(z), 0, stone)
            chunk.SetBlock(uint8(x), int16(noise)+1, uint8(z), 0, grass)
        }
    }
}
