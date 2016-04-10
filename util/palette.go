package util

import (
	"sync"
)

var testPalette = []string{
	`#111111`,
	`#303030`,
	`#141414`,
	`#262626`,
	`#181818`,
	`#222222`,
}
var testPalletePoint int

type Palette struct {
	sync.Mutex
	Colors []string
	Ptr int
}


func NewPalette (colors []string) *Palette {
	var p Palette
	p.Colors = colors
	p.Colors = testPalette
	return &p
}


// shift and return new color
func (p *Palette) GetNext() string {
	p.Lock()
	defer p.Unlock()
	p.Ptr++
	if p.Ptr >= len(p.Colors) {
		p.Ptr = 0
	}
	return p.Colors[p.Ptr]
}
func (p *Palette) Get() string {
	p.Lock()
	defer p.Unlock()
	return p.Colors[p.Ptr]
}
