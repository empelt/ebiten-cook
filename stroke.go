package main

type Stroke struct {
	source StrokeSource

	// offsetX and offsetY represents a relative value from the sprite's upper-left position to the cursor position.
	offsetX int
	offsetY int

	initX int
	initY int

	food *Food

	targets []*Sprite
}

func NewStroke(source StrokeSource, food *Food, targets []*Sprite) *Stroke {
	food.dragged = true
	x, y := source.Position()
	return &Stroke{
		source:  source,
		offsetX: x - food.x,
		offsetY: y - food.y,
		initX:   food.x,
		initY:   food.y,
		food:    food,
		targets: targets,
	}
}

func (s *Stroke) Update() {
	if !s.food.dragged {
		return
	}
	if s.source.IsJustReleased() {
		s.food.dragged = false

		for _, t := range s.targets {
			if t.In(s.source.Position()) {
				s.food.MoveTo(t.x, t.y)
				return
			}
		}

		s.food.MoveTo(s.initX, s.initY)
		return
	}

	x, y := s.source.Position()
	x -= s.offsetX
	y -= s.offsetY
	s.food.MoveTo(x, y)
}

func (s *Stroke) Food() *Food {
	return s.food
}
