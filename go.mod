module github.com/sago35/tinygo-keyboard

go 1.19

require (
	github.com/bgould/tinygo-rotary-encoder v0.0.0-20221224155058-c26fcc9a3d20
	github.com/itchio/lzma v0.0.0-20190703113020-d3e24e3e3d49
	tinygo.org/x/drivers v0.24.1-0.20230520223205-95f0ca8c3ee0
	tinygo.org/x/tinydraw v0.3.0
	tinygo.org/x/tinyfont v0.3.0
)

require github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect

replace (
	github.com/bgould/tinygo-rotary-encoder => github.com/akif999/tinygo-rotary-encoder delete_override_definitions
)
