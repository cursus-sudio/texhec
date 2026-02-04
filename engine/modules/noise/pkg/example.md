```go
func ChangeColor(b *strings.Builder, color mgl64.Vec3) {
	fmt.Fprintf(b,
		"\033[38;2;%d;%d;%dm",
		uint8(color[0]*255),
		uint8(color[1]*255),
		uint8(color[2]*255),
	)
}

func Example(c ioc.Dic, seed seed.Seed) {
	noiseService := ioc.Get[noise.Service](c)
	noise := noiseService.NewNoise(seed).
		AddPerlin(noise.LayerConfig{
			CellSize:        3,
			ValueMultiplier: .8,
		}).
		AddPerlin(noise.LayerConfig{
			CellSize:        8,
			ValueMultiplier: .2,
		}).
		Build().
		Wrap(func(v float64) float64 {
			// v *= 2
			// if v < 1 {
			// 	return 0.5 * v * v * v
			// } else {
			// 	v -= 2
			// 	return 0.5 * (v*v*v + 2)
			// }
			switch {
			case v == 0:
				return 0
			case v == 1:
				return 1
			case v < 0.5:
				return math.Pow(2, 20*v-10) / 2
			default:
				return (2 - math.Pow(2, -20*v+10)) / 2
			}
		})

	size := mgl64.Vec2{
		100,
		44,
	}
	message := &strings.Builder{}
	message.Grow(
		int(size.Y()) * (2 + int(size.X())*len("\033[38;2;%d;%d;%dm█")),
	)
	minVal, maxVal := 1., 0.
	for y := range int(size.Y() + 1) {
		normalizedY := float64(y) / size.Y()
		for x := range int(size.X() + 1) {
			normalizedX := float64(x) / size.X()
			value := noise.Read(mgl64.Vec2{normalizedX, normalizedY})
			minVal, maxVal = min(minVal, value), max(maxVal, value)
			ChangeColor(message, mgl64.Vec3{value, value, value}) // normalizedX, 1 - normalizedY*value, value,
			// changeColor(message, mgl64.Vec3{normalizedX, 1 - normalizedY*value, value})
			fmt.Fprintf(message, "█")
		}
		fmt.Fprintf(message, "\n")
	}
	ChangeColor(message, mgl64.Vec3{.5, .5, 1})
	fmt.Fprintf(message, "min: %v, max %v\n", minVal, maxVal)

	print(message.String())
}
```
