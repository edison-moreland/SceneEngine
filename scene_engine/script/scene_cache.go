package script

import (
	"github.com/edison-moreland/SceneEngine/scene_engine/core/messages"
)

type SceneCache struct {
	config     messages.Config
	scenes     []messages.Scene
	sceneCount uint64
}

func (s *SceneCache) Reset(config messages.Config) {
	// Clear the current cache, and allocate enough space for the next scenes
	s.scenes = make([]messages.Scene, 0, config.FrameCount)
	s.scenes = s.scenes[:config.FrameCount]
	s.config = config
	s.sceneCount = 0
}

func (s *SceneCache) CacheScene(frame uint64, scene messages.Scene) {
	s.scenes[frame-1] = scene
	s.sceneCount += 1
}

func (s *SceneCache) Scene(frame uint64) messages.Scene {
	return s.scenes[frame-1]
}

func (s *SceneCache) Config() messages.Config {
	return s.config
}

func (s *SceneCache) Full() bool {
	return s.sceneCount == s.config.FrameCount
}
