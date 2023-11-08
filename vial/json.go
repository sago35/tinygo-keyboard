package vial

import (
	"encoding/json"
	"strings"
)

func NewJSON() *VialJSON {
	return &VialJSON{}
}

func (v *VialJSON) SetName(name string) {
	v.Name = name
}

func (v *VialJSON) SetVendorID(vid string) {
	v.VendorID = vid
}

func (v *VialJSON) SetProductID(pid string) {
	v.ProductID = pid
}

func (v *VialJSON) SetMatrix(rows, cols int) {
	v.Matrix.Rows = rows
	v.Matrix.Cols = cols
}

func (v *VialJSON) SetKeymap(keymap [][]string) {
	v.Layouts.Keymap = keymap
}

func (v *VialJSON) Draw() (string, error) {
	ret := Canvas{}

	y := 0
	for _, row := range v.Layouts.Keymap {
		x := 0
		tmpy := y
		for _, col := range row {
			if len(col) == 0 {
				panic("len(col) == 0")
			} else if col[0] == '{' {
				opt := KeymapOption{}
				err := json.Unmarshal([]byte(col), &opt)
				if err != nil {
					return "", err
				}
				tx, ty, _ := ret.AddSpace(x, tmpy, int(opt.X*4), int(opt.Y*4))
				x += tx
				if ty < 0 {
					tmpy += ty
				}
			} else {
				tx, _, _ := ret.Add(x, tmpy, col[0])
				x += tx
			}
		}
		y += 4
	}

	return ret.String(), nil
}

type Canvas struct {
	buf  [256][256]byte
	xmax int
	ymax int
}

func (c *Canvas) Add(x, y int, v byte) (int, int, error) {
	for c.buf[y][x] != 0 && c.buf[y][x] != ' ' {
		x++
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			c.buf[y+i][x+j] = v
		}
	}
	if c.xmax < x+4 {
		c.xmax = x + 4
	}
	if c.ymax < y+4 {
		c.ymax = y + 4
	}
	return 4, 4, nil
}

func (c *Canvas) AddSpace(x, y, w, h int) (int, int, error) {
	if w == 0 {
		w = 4
	}
	if h == 0 {
		h = 4
	}
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			c.buf[y+i][x+j] = ' '
		}
	}
	if c.xmax < x+w {
		c.xmax = x + w
	}
	if c.ymax < y+h {
		c.ymax = y + h
	}
	return w, h, nil
}

func (c *Canvas) String() string {
	ret := ""
	for _, row := range c.buf[:c.ymax] {
		ret += string(row[:c.xmax])
		ret += "\n"
	}
	return strings.ReplaceAll(ret, "\x00", " ")
}

type KeymapOption struct {
	X float64 `json:"x",omitempty`
	Y float64 `json:"y",omitempty`
}

type VialJSON struct {
	Name      string `json:"name"`
	VendorID  string `json:"vendorId"`
	ProductID string `json:"productId"`
	Matrix    struct {
		Rows int `json:"rows"`
		Cols int `json:"cols"`
	} `json:"matrix"`
	Layouts struct {
		Keymap [][]string `json:"keymap"`
	} `json:"layouts"`
}
