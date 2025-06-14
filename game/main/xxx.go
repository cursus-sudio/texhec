package main

import (
	"fmt"
	"sync"
	"time"
)

// --- Components ---
// Components are just data.

type Position struct {
	X float64
	Y float64
}

type Velocity struct {
	DX float64
	DY float64
}

type Health struct {
	Current int
	Max     int
}

type Renderable struct {
	Sprite string // Path to sprite asset or identifier
}

type InputControlled struct{} // Marker component for player input

type AIControlled struct{} // Marker component for AI control

// --- Entities ---
// An entity is just a unique ID.
// In a real ECS, you'd likely have a more sophisticated entity manager.
type EntityID int

var nextEntityID EntityID
var mu sync.Mutex // For thread-safe ID generation

func GenerateEntityID() EntityID {
	mu.Lock()
	defer mu.Unlock()
	nextEntityID++
	return nextEntityID
}

// --- Component Storage (Simplified) ---
// In a full ECS, this would be more optimized (e.g., arrays of structs, contiguous memory).
// For simplicity, we'll use maps.
var positions = make(map[EntityID]Position)
var velocities = make(map[EntityID]Velocity)
var healths = make(map[EntityID]Health)
var renderables = make(map[EntityID]Renderable)
var inputControlled = make(map[EntityID]struct{}) // Using struct{} for set-like behavior
var aiControlled = make(map[EntityID]struct{})

// AddComponent is a helper to associate components with an entity.
func AddComponent(entity EntityID, component interface{}) {
	switch c := component.(type) {
	case Position:
		positions[entity] = c
	case Velocity:
		velocities[entity] = c
	case Health:
		healths[entity] = c
	case Renderable:
		renderables[entity] = c
	case InputControlled:
		inputControlled[entity] = struct{}{}
	case AIControlled:
		aiControlled[entity] = struct{}{}
	default:
		fmt.Printf("Warning: Unknown component type %T\n", c)
	}
}

// RemoveComponent (simplified)
func RemoveComponent(entity EntityID, componentType interface{}) {
	switch componentType.(type) {
	case Position:
		delete(positions, entity)
	case Velocity:
		delete(velocities, entity)
	case Health:
		delete(healths, entity)
	case Renderable:
		delete(renderables, entity)
	case InputControlled:
		delete(inputControlled, entity)
	case AIControlled:
		delete(aiControlled, entity)
	default:
		fmt.Printf("Warning: Cannot remove unknown component type %T\n", componentType)
	}
}

// HasComponents checks if an entity has all the specified component types.
// This is a simplified check for our example. A real ECS would use bitmasks or more efficient queries.
func HasComponents(entity EntityID, componentTypes ...interface{}) bool {
	for _, ct := range componentTypes {
		switch ct.(type) {
		case Position:
			if _, ok := positions[entity]; !ok {
				return false
			}
		case Velocity:
			if _, ok := velocities[entity]; !ok {
				return false
			}
		case Health:
			if _, ok := healths[entity]; !ok {
				return false
			}
		case Renderable:
			if _, ok := renderables[entity]; !ok {
				return false
			}
		case InputControlled:
			if _, ok := inputControlled[entity]; !ok {
				return false
			}
		case AIControlled:
			if _, ok := aiControlled[entity]; !ok {
				return false
			}
		default:
			return false // Unknown component type
		}
	}
	return true
}

// GetEntitiesWithComponents returns a slice of entity IDs that have all the specified components.
func GetEntitiesWithComponents(componentTypes ...interface{}) []EntityID {
	var entities []EntityID
	// This is highly inefficient for a real game; a real ECS would have an optimized query system.
	// We iterate over a common set of entities (e.g., all entities ever created) or an index.
	// For this example, we'll iterate over entities with positions as a common starting point.
	for id := range positions { // Assuming most entities have a position
		if HasComponents(id, componentTypes...) {
			entities = append(entities, id)
		}
	}
	return entities
}

// --- Systems ---
// Systems contain logic and operate on specific component combinations.

// MovementSystem updates position based on velocity.
type MovementSystem struct{}

func (s *MovementSystem) Update(deltaTime float64) {
	entityIDs := GetEntitiesWithComponents(Position{}, Velocity{})
	for _, id := range entityIDs {
		pos := positions[id]
		vel := velocities[id]

		pos.X += vel.DX * deltaTime
		pos.Y += vel.DY * deltaTime
		positions[id] = pos
	}
}

// RenderSystem (dummy for now) prints entity positions.
type RenderSystem struct{}

func (s *RenderSystem) Update() {
	entityIDs := GetEntitiesWithComponents(Position{}, Renderable{})
	fmt.Println("--- Rendering ---")
	for _, id := range entityIDs {
		pos := positions[id]
		rend := renderables[id]
		fmt.Printf("Entity %d (%s) at (%.2f, %.2f)\n", id, rend.Sprite, pos.X, pos.Y)
	}
	fmt.Println("-----------------")
}

// PlayerInputSystem (dummy for now) simulates player input.
type PlayerInputSystem struct {
	game *Game
}

func (s *PlayerInputSystem) Update() {
	entityIDs := GetEntitiesWithComponents(InputControlled{}, Velocity{})
	for _, id := range entityIDs {
		// Simulate player input: move right
		vel := velocities[id]
		vel.DX = 1.0 // Example: move right
		velocities[id] = vel
		fmt.Printf("Player %d input processed: moving right\n", id)
	}
}

// AISystem (dummy for now) simulates simple enemy AI.
type AISystem struct {
	game *Game
}

func (s *AISystem) Update() {
	entityIDs := GetEntitiesWithComponents(AIControlled{}, Velocity{})
	for _, id := range entityIDs {
		// Simulate simple AI: move down
		vel := velocities[id]
		vel.DY = 0.5 // Example: move down
		velocities[id] = vel
		fmt.Printf("AI %d processed: moving down\n", id)
	}
}

// HealthSystem checks for entities with 0 health.
type HealthSystem struct {
	game *Game
}

func (s *HealthSystem) Update() {
	entityIDs := GetEntitiesWithComponents(Health{})
	for _, id := range entityIDs {
		h := healths[id]
		if h.Current <= 0 {
			fmt.Printf("Entity %d (Health: %d) has died!\n", id, h.Current)
			// In a real game, you'd trigger death animations, remove the entity, etc.
			// For simplicity, we'll just mark it as dead and prevent further updates
			// by removing its movement component.
			RemoveComponent(id, Velocity{})
			RemoveComponent(id, Renderable{})
			RemoveComponent(id, Health{})
			// A real ECS would have an Entity Manager to fully remove or deactivate.
		}
	}
}

// --- Game State Management ---
type GameState int

const (
	StateLoading GameState = iota
	StateMainMenu
	StatePlaying
	StatePaused
	StateGameOver
)

// --- Game Structure ---
type Game struct {
	currentState GameState
	// Add other managers/systems here
	movementSystem    *MovementSystem
	renderSystem      *RenderSystem
	playerInputSystem *PlayerInputSystem
	aiSystem          *AISystem
	healthSystem      *HealthSystem
	quit              chan struct{}
}

func NewGame() *Game {
	game := &Game{
		currentState:      StateLoading,
		movementSystem:    &MovementSystem{},
		renderSystem:      &RenderSystem{},
		playerInputSystem: &PlayerInputSystem{}, // Will be initialized with game reference
		aiSystem:          &AISystem{},          // Will be initialized with game reference
		healthSystem:      &HealthSystem{},      // Will be initialized with game reference
		quit:              make(chan struct{}),
	}
	game.playerInputSystem.game = game
	game.aiSystem.game = game
	game.healthSystem.game = game
	return game
}

// Init initializes game resources, entities, etc.
func (g *Game) Init() {
	fmt.Println("Game Initializing...")
	g.currentState = StateMainMenu

	// Create a player entity
	player := GenerateEntityID()
	AddComponent(player, Position{X: 0, Y: 0})
	AddComponent(player, Velocity{DX: 0, DY: 0})
	AddComponent(player, Health{Current: 100, Max: 100})
	AddComponent(player, Renderable{Sprite: "player.png"})
	AddComponent(player, InputControlled{})
	fmt.Printf("Created Player Entity: %d\n", player)

	// Create an enemy entity
	enemy := GenerateEntityID()
	AddComponent(enemy, Position{X: 5, Y: 5})
	AddComponent(enemy, Velocity{DX: 0, DY: 0})
	AddComponent(enemy, Health{Current: 50, Max: 50})
	AddComponent(enemy, Renderable{Sprite: "enemy.png"})
	AddComponent(enemy, AIControlled{})
	fmt.Printf("Created Enemy Entity: %d\n", enemy)

	// Create another enemy entity (that will die soon)
	dyingEnemy := GenerateEntityID()
	AddComponent(dyingEnemy, Position{X: 10, Y: 10})
	AddComponent(dyingEnemy, Velocity{DX: 0, DY: 0})
	AddComponent(dyingEnemy, Health{Current: 1, Max: 10}) // Low health
	AddComponent(dyingEnemy, Renderable{Sprite: "dying_enemy.png"})
	AddComponent(dyingEnemy, AIControlled{})
	fmt.Printf("Created Dying Enemy Entity: %d\n", dyingEnemy)

	fmt.Println("Game Initialization Complete. Current State:", g.currentState)
}

// HandleInput processes raw input.
func (g *Game) HandleInput() {
	// In a real game, this would read keyboard/mouse/gamepad input.
	// For this example, we'll let PlayerInputSystem simulate it.
	g.playerInputSystem.Update()
}

// Update game logic based on the current state.
func (g *Game) Update(deltaTime float64) {
	switch g.currentState {
	case StatePlaying:
		g.playerInputSystem.Update() // Process player specific input
		g.aiSystem.Update()          // Update AI
		g.movementSystem.Update(deltaTime)
		g.healthSystem.Update() // Check for dead entities
		// More systems would go here (e.g., collision, animation, physics)
	case StateMainMenu:
		fmt.Println("In Main Menu. Press 'P' to play (simulated).")
		// Simulate playing
		g.currentState = StatePlaying
	case StateGameOver:
		fmt.Println("Game Over!")
		// Optionally, transition back to main menu or exit
		g.Quit()
	default:
		// Do nothing or handle other states
	}
}

// Render draws the game to the screen.
func (g *Game) Render() {
	switch g.currentState {
	case StatePlaying:
		g.renderSystem.Update()
	case StateMainMenu:
		fmt.Println("Displaying Main Menu...")
	case StateGameOver:
		fmt.Println("Displaying Game Over Screen...")
	default:
		// Do nothing
	}
}

// GameLoop
func (g *Game) Run() {
	const fixedUpdateTime = 1.0 / 60.0 // 60 FPS
	lastFrameTime := time.Now()
	accumulator := 0.0

	fmt.Println("Game Loop Starting...")
	for {
		select {
		case <-g.quit:
			fmt.Println("Game Loop Exiting...")
			return
		default:
			currentTime := time.Now()
			frameTime := currentTime.Sub(lastFrameTime).Seconds()
			lastFrameTime = currentTime

			// Prevent spiraling in case of very long frame times
			if frameTime > 0.25 {
				frameTime = 0.25
			}

			accumulator += frameTime

			// Fixed timestep updates for physics/game logic
			for accumulator >= fixedUpdateTime {
				g.Update(fixedUpdateTime)
				accumulator -= fixedUpdateTime
			}

			// Render as fast as possible, but use interpolation if needed
			// (not implemented in this simple example)
			g.Render()

			// Simple delay to prevent 100% CPU usage if game logic is too fast
			time.Sleep(time.Millisecond * 10)
		}
	}
}

// Quit signals the game loop to stop.
func (g *Game) Quit() {
	close(g.quit)
}

func main() {
	game := NewGame()
	game.Init()
	game.Run()
}
