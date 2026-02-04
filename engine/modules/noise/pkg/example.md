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
		AddPerlin(
			noise.NewLayer(50, func(f float64) float64 { return f * .5 }),
			noise.NewLayer(15, func(f float64) float64 { return f * .3 }),
			noise.NewLayer(5, func(f float64) float64 { return f * .2 }),
		).
		Build()

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
			value := noise.Read(mgl64.Vec2{normalizedX, normalizedY}.Mul(100))
			minVal, maxVal = min(minVal, value), max(maxVal, value)
			ChangeColor(message, mgl64.Vec3{value, value, value})
			fmt.Fprintf(message, "█")
		}
		fmt.Fprintf(message, "\n")
	}
	ChangeColor(message, mgl64.Vec3{.5, .5, 1})
	fmt.Fprintf(message, "min: %v, max %v\n", minVal, maxVal)

	print(message.String())
}
```
