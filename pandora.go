package main

import (
    "github.com/curzodo/perlin"
    "github.com/df-mc/dragonfly/server"
    "github.com/df-mc/dragonfly/server/block"
    "github.com/df-mc/dragonfly/server/world"
    "github.com/df-mc/dragonfly/server/world/chunk"
    "math"
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
	for srv.Accept(nil) {
	}
}

type Generator struct {
    perlin.Generator
}

func NewGenerator(seed int64) world.Generator {
    return Generator { perlin.NewGenerator(seed) }
}

func (g Generator) GenerateChunk(pos world.ChunkPos, chunk *chunk.Chunk) {
    // Generate top layer
    for x := 0; x < 16; x++ {
        for z := 0; z < 16; z++ {
            xsign := int(math.Copysign(float64(x), float64(pos[0])))
            zsign := int(math.Copysign(float64(z), float64(pos[1])))
            worldx := (float64(xsign + 15 * int(pos[0])) + 0.01)/20
            worldz := (float64(zsign + 15 * int(pos[1])) + 0.01)/20

            y := (g.Noise2D(worldx, worldz) * 3) + 2
            b := world.BlockRuntimeID(block.Grass{})

            if y < 0 { 
                b2 := world.BlockRuntimeID(block.Water{
                    Still: true,
                    Depth: 1,
                    Falling: true,
                }) 

                if rand.Intn(5) == 0 {
                    chunk.SetBlock(uint8(x), int16(y+1), uint8(z), 0, b2)
                    continue
                }

                if chunk.Block(uint8(x+1), int16(y+1), uint8(z), 0) == b2 { 
                    chunk.SetBlock(uint8(x), int16(y+1), uint8(z), 0, b2)
                }
                if chunk.Block(uint8(x-1), int16(y+1), uint8(z), 0) == b2 { 
                    chunk.SetBlock(uint8(x), int16(y+1), uint8(z), 0, b2)
                }
                if chunk.Block(uint8(x), int16(y+1), uint8(z+1), 0) == b2 { 
                    chunk.SetBlock(uint8(x), int16(y+1), uint8(z), 0, b2)
                }
                if chunk.Block(uint8(x), int16(y+1), uint8(z-1), 0) == b2 { 
                    chunk.SetBlock(uint8(x), int16(y+1), uint8(z), 0, b2)
                }
            }

            chunk.SetBlock(uint8(x), int16(y), uint8(z), 0, b)
        }
    }
}
