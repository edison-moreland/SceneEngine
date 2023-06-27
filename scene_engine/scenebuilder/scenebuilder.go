package scenebuilder

import (
	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
)

func emptyScene() messages.Scene {
	return messages.Scene{
		Camera: messages.Camera{
			Aperture: 0.1,
			Fov:      90,
			LookAt: messages.Position{
				X: 0,
				Y: 0,
				Z: 0,
			},
			LookFrom: messages.Position{
				X: 4,
				Y: 0,
				Z: 0,
			},
		},
	}
}

// SceneBuilder is used to progressively build a scene
type SceneBuilder struct {
	currentScene messages.Scene
	sceneCache   *SceneCache
}

func New(sc *SceneCache) *SceneBuilder {
	return &SceneBuilder{
		currentScene: messages.Scene{},
		sceneCache:   sc,
	}
}

func (sb *SceneBuilder) Reset(config messages.Config) {
	sb.sceneCache.Reset(config)
	sb.currentScene = emptyScene()
}

func (sb *SceneBuilder) Commit(frame uint64) {
	sb.sceneCache.CacheScene(frame, sb.currentScene)
	sb.currentScene = emptyScene()
}

func (sb *SceneBuilder) AddObject(o messages.Object) {
	sb.currentScene.Objects = append(sb.currentScene.Objects, o)
}

func (sb *SceneBuilder) SetCamera(c messages.Camera) {
	sb.currentScene.Camera = c
}
