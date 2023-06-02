package libraries

import (
	"github.com/d5/tengo/v2"
)

func AddSceneEngineLibraries(m *tengo.ModuleMap) {
	m.AddBuiltinModule("vec3", Vec3Module)
	m.AddBuiltinModule("color", ColorModule)
	m.AddBuiltinModule("shape", ShapeModule)
	m.AddBuiltinModule("material", MaterialModule)
}
