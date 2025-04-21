package engine

import (
	"errors"
	"fmt"
	"sync"
)

// SceneObject represents a scene object
type SceneObject struct {
	// Object ID
	ID string
	
	// Object name
	Name string
	
	// Object position
	Position [3]float64
	
	// Object rotation
	Rotation [3]float64
	
	// Object scale
	Scale [3]float64
	
	// Object parent
	Parent *SceneObject
	
	// Object children
	Children []*SceneObject
	
	// Object components
	Components map[string]Component
	
	// Object active
	Active bool
	
	// Object visible
	Visible bool
	
	// Object layer
	Layer int
	
	// Object tag
	Tag string
	
	// Object user data
	UserData map[string]interface{}
}

// Component represents a component
type Component interface {
	// GetID gets the component ID
	GetID() string
	
	// GetName gets the component name
	GetName() string
	
	// GetObject gets the component's object
	GetObject() *SceneObject
	
	// SetObject sets the component's object
	SetObject(object *SceneObject)
	
	// OnAttach is called when the component is attached to an object
	OnAttach()
	
	// OnDetach is called when the component is detached from an object
	OnDetach()
	
	// OnEnable is called when the component is enabled
	OnEnable()
	
	// OnDisable is called when the component is disabled
	OnDisable()
	
	// OnUpdate is called every frame
	OnUpdate(deltaTime float64)
}

// BaseComponent is a base implementation of the Component interface
type BaseComponent struct {
	// Component ID
	ID string
	
	// Component name
	Name string
	
	// Component object
	Object *SceneObject
	
	// Component enabled
	Enabled bool
}

// GetID gets the component ID
func (c *BaseComponent) GetID() string {
	return c.ID
}

// GetName gets the component name
func (c *BaseComponent) GetName() string {
	return c.Name
}

// GetObject gets the component's object
func (c *BaseComponent) GetObject() *SceneObject {
	return c.Object
}

// SetObject sets the component's object
func (c *BaseComponent) SetObject(object *SceneObject) {
	c.Object = object
}

// OnAttach is called when the component is attached to an object
func (c *BaseComponent) OnAttach() {
	// Base implementation does nothing
}

// OnDetach is called when the component is detached from an object
func (c *BaseComponent) OnDetach() {
	// Base implementation does nothing
}

// OnEnable is called when the component is enabled
func (c *BaseComponent) OnEnable() {
	// Base implementation does nothing
}

// OnDisable is called when the component is disabled
func (c *BaseComponent) OnDisable() {
	// Base implementation does nothing
}

// OnUpdate is called every frame
func (c *BaseComponent) OnUpdate(deltaTime float64) {
	// Base implementation does nothing
}

// Scene represents a scene
type Scene struct {
	// Scene ID
	ID string
	
	// Scene name
	Name string
	
	// Root objects
	RootObjects []*SceneObject
	
	// All objects
	Objects map[string]*SceneObject
	
	// Active camera
	ActiveCamera *Camera
	
	// Ambient light
	AmbientLight [3]float64
	
	// Gravity
	Gravity [3]float64
	
	// Mutex for thread safety
	mutex sync.RWMutex
}

// NewScene creates a new scene
func NewScene(id, name string) *Scene {
	return &Scene{
		ID:          id,
		Name:        name,
		RootObjects: []*SceneObject{},
		Objects:     make(map[string]*SceneObject),
		AmbientLight: [3]float64{0.1, 0.1, 0.1},
		Gravity:     [3]float64{0, -9.81, 0},
	}
}

// CreateObject creates a new scene object
func (s *Scene) CreateObject(id, name string) *SceneObject {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	object := &SceneObject{
		ID:         id,
		Name:       name,
		Position:   [3]float64{0, 0, 0},
		Rotation:   [3]float64{0, 0, 0},
		Scale:      [3]float64{1, 1, 1},
		Children:   []*SceneObject{},
		Components: make(map[string]Component),
		Active:     true,
		Visible:    true,
		Layer:      0,
		UserData:   make(map[string]interface{}),
	}
	
	s.Objects[id] = object
	s.RootObjects = append(s.RootObjects, object)
	
	return object
}

// DestroyObject destroys a scene object
func (s *Scene) DestroyObject(object *SceneObject) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if object == nil {
		return
	}
	
	// Remove from parent
	if object.Parent != nil {
		for i, child := range object.Parent.Children {
			if child == object {
				object.Parent.Children = append(object.Parent.Children[:i], object.Parent.Children[i+1:]...)
				break
			}
		}
	}
	
	// Remove from root objects
	for i, root := range s.RootObjects {
		if root == object {
			s.RootObjects = append(s.RootObjects[:i], s.RootObjects[i+1:]...)
			break
		}
	}
	
	// Destroy children
	for _, child := range object.Children {
		s.DestroyObject(child)
	}
	
	// Remove from objects
	delete(s.Objects, object.ID)
}

// GetObject gets a scene object by ID
func (s *Scene) GetObject(id string) *SceneObject {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	return s.Objects[id]
}

// GetObjectByName gets a scene object by name
func (s *Scene) GetObjectByName(name string) *SceneObject {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	for _, object := range s.Objects {
		if object.Name == name {
			return object
		}
	}
	
	return nil
}

// GetObjectsByTag gets scene objects by tag
func (s *Scene) GetObjectsByTag(tag string) []*SceneObject {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	var objects []*SceneObject
	
	for _, object := range s.Objects {
		if object.Tag == tag {
			objects = append(objects, object)
		}
	}
	
	return objects
}

// GetObjectsByLayer gets scene objects by layer
func (s *Scene) GetObjectsByLayer(layer int) []*SceneObject {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	var objects []*SceneObject
	
	for _, object := range s.Objects {
		if object.Layer == layer {
			objects = append(objects, object)
		}
	}
	
	return objects
}

// Update updates the scene
func (s *Scene) Update(deltaTime float64) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	// Update all objects
	for _, object := range s.Objects {
		if !object.Active {
			continue
		}
		
		// Update components
		for _, component := range object.Components {
			if component, ok := component.(Component); ok {
				component.OnUpdate(deltaTime)
			}
		}
	}
}

// SetParent sets the parent of a scene object
func (s *Scene) SetParent(object, parent *SceneObject) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if object == nil {
		return errors.New("object is nil")
	}
	
	// Remove from current parent
	if object.Parent != nil {
		for i, child := range object.Parent.Children {
			if child == object {
				object.Parent.Children = append(object.Parent.Children[:i], object.Parent.Children[i+1:]...)
				break
			}
		}
	}
	
	// Remove from root objects
	for i, root := range s.RootObjects {
		if root == object {
			s.RootObjects = append(s.RootObjects[:i], s.RootObjects[i+1:]...)
			break
		}
	}
	
	// Set new parent
	object.Parent = parent
	
	if parent == nil {
		// Add to root objects
		s.RootObjects = append(s.RootObjects, object)
	} else {
		// Add to parent's children
		parent.Children = append(parent.Children, object)
	}
	
	return nil
}

// AddComponent adds a component to a scene object
func (s *Scene) AddComponent(object *SceneObject, component Component) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if object == nil {
		return errors.New("object is nil")
	}
	
	if component == nil {
		return errors.New("component is nil")
	}
	
	// Check if component already exists
	if _, ok := object.Components[component.GetID()]; ok {
		return fmt.Errorf("component with ID %s already exists", component.GetID())
	}
	
	// Set component's object
	component.SetObject(object)
	
	// Add component
	object.Components[component.GetID()] = component
	
	// Call OnAttach
	component.OnAttach()
	
	return nil
}

// RemoveComponent removes a component from a scene object
func (s *Scene) RemoveComponent(object *SceneObject, componentID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if object == nil {
		return errors.New("object is nil")
	}
	
	// Check if component exists
	component, ok := object.Components[componentID]
	if !ok {
		return fmt.Errorf("component with ID %s does not exist", componentID)
	}
	
	// Call OnDetach
	component.OnDetach()
	
	// Remove component
	delete(object.Components, componentID)
	
	return nil
}

// GetComponent gets a component from a scene object
func (s *Scene) GetComponent(object *SceneObject, componentID string) Component {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	if object == nil {
		return nil
	}
	
	return object.Components[componentID]
}

// Camera represents a camera component
type Camera struct {
	BaseComponent
	
	// Field of view
	FieldOfView float64
	
	// Aspect ratio
	AspectRatio float64
	
	// Near clip plane
	NearClip float64
	
	// Far clip plane
	FarClip float64
	
	// Orthographic
	Orthographic bool
	
	// Orthographic size
	OrthographicSize float64
	
	// Clear color
	ClearColor [4]float64
	
	// Clear flags
	ClearFlags int
	
	// Viewport
	Viewport [4]float64
	
	// Projection matrix
	ProjectionMatrix [16]float64
	
	// View matrix
	ViewMatrix [16]float64
}

// NewCamera creates a new camera component
func NewCamera(id, name string) *Camera {
	return &Camera{
		BaseComponent: BaseComponent{
			ID:      id,
			Name:    name,
			Enabled: true,
		},
		FieldOfView:      60,
		AspectRatio:      1.7777,
		NearClip:         0.1,
		FarClip:          1000,
		Orthographic:     false,
		OrthographicSize: 5,
		ClearColor:       [4]float64{0.2, 0.3, 0.3, 1.0},
		ClearFlags:       0x1, // Clear color and depth
		Viewport:         [4]float64{0, 0, 1, 1},
		ProjectionMatrix: [16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1},
		ViewMatrix:       [16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1},
	}
}

// OnAttach is called when the camera is attached to an object
func (c *Camera) OnAttach() {
	// Update matrices
	c.UpdateMatrices()
}

// OnUpdate is called every frame
func (c *Camera) OnUpdate(deltaTime float64) {
	// Update matrices
	c.UpdateMatrices()
}

// UpdateMatrices updates the camera matrices
func (c *Camera) UpdateMatrices() {
	// This would normally update the projection and view matrices
	// For now, we'll just use identity matrices
	c.ProjectionMatrix = [16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
	c.ViewMatrix = [16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
}

// MeshRenderer represents a mesh renderer component
type MeshRenderer struct {
	BaseComponent
	
	// Mesh
	Mesh *Mesh
	
	// Materials
	Materials []*Material
	
	// Cast shadows
	CastShadows bool
	
	// Receive shadows
	ReceiveShadows bool
	
	// Bounds
	Bounds [6]float64 // min x, min y, min z, max x, max y, max z
}

// NewMeshRenderer creates a new mesh renderer component
func NewMeshRenderer(id, name string) *MeshRenderer {
	return &MeshRenderer{
		BaseComponent: BaseComponent{
			ID:      id,
			Name:    name,
			Enabled: true,
		},
		Materials:      []*Material{},
		CastShadows:    true,
		ReceiveShadows: true,
		Bounds:         [6]float64{-0.5, -0.5, -0.5, 0.5, 0.5, 0.5},
	}
}

// Mesh represents a mesh
type Mesh struct {
	// Mesh ID
	ID string
	
	// Mesh name
	Name string
	
	// Vertices
	Vertices [][3]float64
	
	// Normals
	Normals [][3]float64
	
	// UVs
	UVs [][2]float64
	
	// Colors
	Colors [][4]float64
	
	// Indices
	Indices []int
	
	// Submeshes
	Submeshes [][]int
	
	// Bounds
	Bounds [6]float64 // min x, min y, min z, max x, max y, max z
}

// Material represents a material
type Material struct {
	// Material ID
	ID string
	
	// Material name
	Name string
	
	// Shader
	Shader *GPUShader
	
	// Textures
	Textures map[string]*GPUTexture
	
	// Properties
	Properties map[string]interface{}
	
	// Render queue
	RenderQueue int
}

// Light represents a light component
type Light struct {
	BaseComponent
	
	// Light type
	Type string
	
	// Light color
	Color [3]float64
	
	// Light intensity
	Intensity float64
	
	// Light range
	Range float64
	
	// Light spot angle
	SpotAngle float64
	
	// Light shadows
	Shadows bool
	
	// Light shadow resolution
	ShadowResolution int
	
	// Light shadow bias
	ShadowBias float64
}

// NewLight creates a new light component
func NewLight(id, name string, type_ string) *Light {
	return &Light{
		BaseComponent: BaseComponent{
			ID:      id,
			Name:    name,
			Enabled: true,
		},
		Type:            type_,
		Color:           [3]float64{1, 1, 1},
		Intensity:       1,
		Range:           10,
		SpotAngle:       30,
		Shadows:         true,
		ShadowResolution: 1024,
		ShadowBias:      0.05,
	}
}