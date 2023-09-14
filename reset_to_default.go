//go:build reset_to_default

package keyboard

import (
	"machine"
)

func init() {
	machine.Flash.EraseBlocks(0, 1)
}
