package color

import (
	"math/rand"
)

type DiscordColor int

const (
	Green  DiscordColor = 65280
	Blue   DiscordColor = 255
	Red    DiscordColor = 16711680
	Yellow DiscordColor = 16776960
	Black  DiscordColor = 0
	White  DiscordColor = 16777215
	Orange DiscordColor = 16744192
	Purple DiscordColor = 8388736
)

func RandomColor() DiscordColor {
	colors := []DiscordColor{
		Blue, Green, Black, Orange, Purple, Red, Yellow, White,
	}

	index := rand.Intn(len(colors) - 1)

	return colors[index]
}
