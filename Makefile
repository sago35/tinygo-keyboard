smoketest: FORCE
	go build .
	tinygo build -o /tmp/out.uf2 --target xiao-rp2040     --size short ./targets/fric10key/
	tinygo build -o /tmp/out.uf2 --target gopher-badge    --size short ./targets/gopher-badge/
	tinygo build -o /tmp/out.uf2 --target macropad-rp2040 --size short ./targets/macropad-rp2040/
	tinygo build -o /tmp/out.uf2 --target xiao-rp2040     --size short ./targets/sgkb/left/
	tinygo build -o /tmp/out.uf2 --target xiao-rp2040     --size short ./targets/sgkb/left-0.3.0/
	tinygo build -o /tmp/out.uf2 --target xiao-rp2040     --size short ./targets/sgkb/right/
	tinygo build -o /tmp/out.uf2 --target xiao-rp2040     --size short ./targets/sgkey/
	tinygo build -o /tmp/out.uf2 --target wioterminal     --size short ./targets/wiokey/
	tinygo build -o /tmp/out.uf2 --target xiao-rp2040     --size short ./targets/xiao-kb01/

FORCE:

gen-def-with-find:
	find . -name vial.json | xargs -n 1 go run ./cmd/gen-def

gen-def:
	go run ./cmd/gen-def/main.go ./targets/fric10key/vial.json
	go run ./cmd/gen-def/main.go ./targets/gopher-badge/vial.json
	go run ./cmd/gen-def/main.go ./targets/macropad-rp2040/vial.json
	go run ./cmd/gen-def/main.go ./targets/sgkb/left/vial.json
	go run ./cmd/gen-def/main.go ./targets/sgkb/left-0.3.0/vial.json
	go run ./cmd/gen-def/main.go ./targets/sgkey/vial.json
	go run ./cmd/gen-def/main.go ./targets/xiao-kb01/vial.json
