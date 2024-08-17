package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	image      *ebiten.Image
	alphaImage *image.Alpha
	x          int
	y          int
}

// In returns true if (x, y) is in the sprite, and false otherwise.
func (s *Sprite) In(x, y int) bool {
	// Check the actual color (alpha) value at the specified position
	// so that the result of In becomes natural to users.
	//
	// Use alphaImage (*image.Alpha) instead of image (*ebiten.Image) here.
	// It is because (*ebiten.Image).At is very slow as this reads pixels from GPU,
	// and should be avoided whenever possible.
	return s.alphaImage.At(x-s.x, y-s.y).(color.Alpha).A > 0
}

// MoveTo moves the sprite to the position (x, y).
func (s *Sprite) MoveTo(x, y int) {
	w, h := s.image.Bounds().Dx(), s.image.Bounds().Dy()

	s.x = x
	s.y = y
	if s.x < 0 {
		s.x = 0
	}
	if s.x > screenWidth-w {
		s.x = screenWidth - w
	}
	if s.y < 0 {
		s.y = 0
	}
	if s.y > screenHeight-h {
		s.y = screenHeight - h
	}
}

// Draw draws the sprite.
func (s *Sprite) Draw(screen *ebiten.Image, alpha float32) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(s.x), float64(s.y))
	op.ColorScale.ScaleAlpha(alpha)
	screen.DrawImage(s.image, op)
}

type Food struct {
	*Sprite
	dragged bool
}
