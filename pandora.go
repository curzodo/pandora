package main

import (
    "github.com/curzodo/perlin"
    "github.com/df-mc/dragonfly/server"
    "github.com/df-mc/dragonfly/server/world"
    "github.com/df-mc/dragonfly/server/world/chunk"
)

func main() {
    conf := server.DefaultConfig().Config(nil)

	srv := conf.New()
	srv.CloseOnProgramEnd()

	srv.Listen()
	for srv.Accept(nil) {
	}
}

type Generator struct {
    gen perlin.Generator
}

func NewGenerator(seed int64) world.Generator {
    return Generator { perlin.NewGenerator(seed) }
}

func (g Generator) GenerateChunk(pos world.ChunkPos, chunk *chunk.Chunk) {
    // Generate top layer
    for x := 0; x < 16; x++ {
        for z := 0; z < 16; z++ {
            trueX := x + 15 * pos[0]
            trueZ := z + 15 * pos[1]

            chunk.SetBlock(x, int16(g.gen.Noise2D(x, z)), z, 1)
        }
    }
}
