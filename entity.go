package main

import "github.com/gopxl/pixel"

// FixedPhysicsTimestep specifies the deltatime to be used.
// Slow framerate will not modify this.
const FixedPhysicsTimestep = 1.0 / 60.0

// Entity is an interface describing all the behaviour an entity should have.
// All entities have physics, but you can choose to disable it using IsKinematic.
// All entities should also be able to be rendered.
type Entity interface {
	Renderable

	Position() pixel.Vec
	Velocity() pixel.Vec
	Acceleration() pixel.Vec
	Mass() float64
	Radius() float64
	ApplyForce(pixel.Vec)
	ApplyImpulse(pixel.Vec)
	SlideToPosition(pixel.Vec)
	SetVelocity(pixel.Vec)
	StepPhysics()
	StepLogic()
	IsKinematic() bool
	Tags() []string // Tags should stay the same after initialisation
}

// EntityBase is a useful implementation of the physics for an entity, for use with composition.
// It does not fully implement Entity, so you have to implement some behaviours yourself.
type EntityBase struct {
	currentPosition    pixel.Vec
	previousPosition   pixel.Vec
	mass               float64
	currentForce       pixel.Vec
	currentImpulse     pixel.Vec
	recentVelocity     pixel.Vec
	recentAcceleration pixel.Vec
	radius             float64
}

// NewEntityBase creates a new base for an entity, should only be used when constructing other entities
func NewEntityBase(position pixel.Vec, mass, radius float64) *EntityBase {
	return &EntityBase{
		currentPosition:    position,
		previousPosition:   position,
		mass:               mass,
		currentForce:       pixel.ZV,
		recentVelocity:     pixel.ZV,
		recentAcceleration: pixel.ZV,
		radius:             radius,
	}
}

func (p *EntityBase) Position() pixel.Vec {
	return p.currentPosition
}

func (p *EntityBase) Velocity() pixel.Vec {
	return p.recentVelocity
}

func (p *EntityBase) Acceleration() pixel.Vec {
	return p.recentAcceleration
}

func (p *EntityBase) Mass() float64 {
	return p.mass
}

func (p *EntityBase) Radius() float64 {
	return p.radius
}

// Will add a force to be applied on the next update step. This is a FORCE not an IMPULSE.
func (p *EntityBase) ApplyForce(force pixel.Vec) {
	p.currentForce = p.currentForce.Add(force)
}

func (p *EntityBase) ApplyImpulse(impulse pixel.Vec) {
	p.currentImpulse = p.currentImpulse.Add(impulse)
}

func (p *EntityBase) SlideToPosition(newPos pixel.Vec) {
	p.currentPosition = newPos
}

func (p *EntityBase) SetVelocity(vel pixel.Vec) {
	p.previousPosition = p.currentPosition.Sub(vel.Scaled(FixedPhysicsTimestep))
}

// Verlet :)
func (p *EntityBase) StepPhysics() {
	totalForce := p.currentForce.Add(p.currentImpulse.Scaled(1 / FixedPhysicsTimestep))
	acceleration := totalForce.Scaled(1 / p.mass)
	p.recentAcceleration = acceleration
	nextPosition := p.currentPosition.Scaled(2).Sub(p.previousPosition).Add(acceleration.Scaled(FixedPhysicsTimestep * FixedPhysicsTimestep))
	p.recentVelocity = nextPosition.Sub(p.previousPosition).Scaled(0.5 / FixedPhysicsTimestep) // Sub one infront from one behind
	p.previousPosition = p.currentPosition
	p.currentPosition = nextPosition
	p.currentForce = pixel.ZV
	p.currentImpulse = pixel.ZV
}

// Default all entities to be physics-enabled
func (p *EntityBase) IsKinematic() bool { return false }

type EntitiesContainer struct {
	allEntities    []Entity
	taggedEntities map[string][]Entity
}

func NewEntitiesContainer() *EntitiesContainer {
	return &EntitiesContainer{
		allEntities:    make([]Entity, 0),
		taggedEntities: make(map[string][]Entity),
	}
}

func (ec *EntitiesContainer) Add(e Entity) {
	ec.allEntities = append(ec.allEntities, e)
	for _, t := range e.Tags() {
		if _, ok := ec.taggedEntities[t]; ok {
			ec.taggedEntities[t] = append(ec.taggedEntities[t], e)
		} else {
			ec.taggedEntities[t] = []Entity{e}
		}
	}
}

func (ec *EntitiesContainer) All() []Entity {
	return ec.allEntities
}

func (ec *EntitiesContainer) WithTag(tag string) []Entity {
	return ec.taggedEntities[tag]
}
