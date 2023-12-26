smoketest: FORCE
	go build .
	mkdir -p out
	tinygo build -o ./out/fric10key.uf2       --target xiao-rp2040           --size short ./targets/fric10key/
	tinygo build -o ./out/gobadge.uf2         --target gobadge               --size short ./targets/gobadge/
	tinygo build -o ./out/gopher-badge.uf2    --target gopher-badge          --size short ./targets/gopher-badge/
	tinygo build -o ./out/macropad-rp2040.uf2 --target macropad-rp2040       --size short ./targets/macropad-rp2040/
	tinygo build -o ./out/sg48key.uf2         --target xiao                  --size short ./targets/sg48key/
	tinygo build -o ./out/sgh60.uf2           --target waveshare-rp2040-zero --size short ./targets/sgh60/
	tinygo build -o ./out/sgkb-left.uf2       --target xiao-rp2040           --size short ./targets/sgkb/left/
	tinygo build -o ./out/sgkb-left-0.3.0.uf2 --target xiao-rp2040           --size short ./targets/sgkb/left-0.3.0/
	tinygo build -o ./out/sgkb-right.uf2      --target xiao-rp2040           --size short ./targets/sgkb/right/
	tinygo build -o ./out/sgkey.uf2           --target xiao-rp2040           --size short ./targets/sgkey/
	tinygo build -o ./out/wiokey.uf2          --target wioterminal           --size short ./targets/wiokey/
	tinygo build -o ./out/xiao-kb01.uf2       --target xiao-rp2040           --size short ./targets/xiao-kb01/
	tinygo build -o ./out/tut-gpio.uf2        --target xiao-rp2040           --size short ./tutorial/gpio/
	tinygo build -o ./out/tut-gpio-vial.uf2   --target xiao-rp2040           --size short ./tutorial/gpio-vial/

FORCE:

gen-def-with-find:
	find . -name vial.json | xargs -n 1 go run ./cmd/gen-def

gen-def:
	go run ./cmd/gen-def/main.go ./targets/fric10key/vial.json
	go run ./cmd/gen-def/main.go ./targets/gopher-badge/vial.json
	go run ./cmd/gen-def/main.go ./targets/gobadge/vial.json
	go run ./cmd/gen-def/main.go ./targets/macropad-rp2040/vial.json
	go run ./cmd/gen-def/main.go ./targets/sg48key/vial.json
	go run ./cmd/gen-def/main.go ./targets/sgh60/vial.json
	go run ./cmd/gen-def/main.go ./targets/sgkb/left/vial.json
	go run ./cmd/gen-def/main.go ./targets/sgkb/left-0.3.0/vial.json
	go run ./cmd/gen-def/main.go ./targets/wiokey/vial.json
	go run ./cmd/gen-def/main.go ./targets/sgkey/vial.json
	go run ./cmd/gen-def/main.go ./targets/xiao-kb01/vial.json
	go run ./cmd/gen-def/main.go ./tutorial/gpio-vial/vial.json

test-gen-def: gen-def-with-find
	git status
	test -z "$$(git status -s)"

test-gen-def-uno: gen-def-with-find
	test -z "$$(git status -s -uno)"
