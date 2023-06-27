# SubMsg

SubMsg is a small rpc framework used between the SceneEngine and the render core.

The render core runs as a child process to SceneEngine, messages are sent back and forth over stdin/stdout.

# TODO
- A more compact language for defining types/messages
- Track down determinism bug (sometimes the order of types/messages changes in the output)