# SceneEngine

SceneEngine is an in progress script based animation/rendering tool.

Here is an example of a scene file that describes the final render of "Ray Tracing in One Weekend":

```
vec3 := import("vec3")
color := import("color")

s := import("shape")
m := import("material")


export {
    config: {
        aspect_ratio: (3.0 / 2.0),
        image_width: 1200,
        samples: 500,
        depth: 50
    },

    scene: func(se, frame, seconds) {
        se.Object(
            s.Sphere(vec3.New(0, 1, 0), 1),
            m.Dielectric(1.5)
        )

        se.Object(
            s.Sphere(vec3.New(-4, 1, 0), 1),
            m.Lambert(color.New(102, 51, 25))
        )

        se.Object(
            s.Sphere(vec3.New(4, 1, 0), 1),
            m.Metal(color.New(178, 153, 127), 0.0)
        )

        se.Object(
            s.Sphere(vec3.New(0, -1000, 0), 1000),
            m.Lambert(color.New(127, 127, 127))
        )

        se.Camera(
            vec3.New(13, 2, 3),
            vec3.New(0, 0, 0),
            20
        )
    }
}
```

Which outputs:

[Final render of "Ray Tracing in One Weekend"](./images/rtow.png)

# Development

Requirements:
- [Go >= 1.20](https://go.dev)
  - [raylib dependencies](https://github.com/gen2brain/raylib-go)
- [Pony >= 0.54.1](https://www.ponylang.io)

Running:
```bash
./se run
```