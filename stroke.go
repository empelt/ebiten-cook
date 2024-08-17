package main

// Stroke manages the current drag state by mouse.
type Stroke struct {
	source StrokeSource

	// offsetX and offsetY represents a relative value from the sprite's upper-left position to the cursor position.
	offsetX int
	offsetY int

	// sprite represents a sprite being dragged.
	sprite *Sprite
}

func NewStroke(source StrokeSource, sprite *Sprite) *Stroke {
	sprite.dragged = true
	x, y := source.Position()
	return &Stroke{
		source:  source,
		offsetX: x - sprite.x,
		offsetY: y - sprite.y,
		sprite:  sprite,
	}
}

func (s *Stroke) Update() {
	if !s.sprite.dragged {
		return
	}
	if s.source.IsJustReleased() {
		s.sprite.dragged = false
		return
	}

	x, y := s.source.Position()
	x -= s.offsetX
	y -= s.offsetY
	s.sprite.MoveTo(x, y)
}

func (s *Stroke) Sprite() *Sprite {
	return s.sprite
}
