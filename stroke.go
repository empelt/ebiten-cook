package main

type Stroke struct {
	source StrokeSource

	// offsetX and offsetY represents a relative value from the dragItem's upper-left position to the cursor position.
	offsetX int
	offsetY int

	initX int
	initY int

	dragItem DraggableSprite
	oven     *JetOven
	plates   []*Plate
}

func NewStroke(source StrokeSource, dragItem DraggableSprite, oven *JetOven, plates []*Plate) *Stroke {
	dragItem.SetDragged(true)
	x, y := source.Position()
	return &Stroke{
		source:   source,
		offsetX:  x - dragItem.GetX(),
		offsetY:  y - dragItem.GetY(),
		initX:    dragItem.GetX(),
		initY:    dragItem.GetY(),
		dragItem: dragItem,
		oven:     oven,
		plates:   plates,
	}
}

func (s *Stroke) Update() {
	if !s.dragItem.GetDragged() {
		return
	}

	switch s.dragItem.(type) {
	case *Food:
		if s.source.IsJustReleased() {
			s.dragItem.SetDragged(false)
			for _, p := range s.plates {
				if p.In(s.source.Position()) {
					s.dragItem.(*Food).MoveTo(p.x, p.y)
					p.AddFood(s.dragItem.(*Food))
					s.dragItem.SetDraggable(false)
					p.SetDraggable(true)
				}
			}
			s.dragItem.(*Food).MoveTo(s.initX, s.initY)

		} else {
			x, y := s.source.Position()
			x -= s.offsetX
			y -= s.offsetY
			s.dragItem.(*Food).MoveTo(x, y)
		}
	case *Plate:
		if s.source.IsJustReleased() {
			s.dragItem.SetDragged(false)
			if s.oven.In(s.source.Position()) {
				s.oven.AddPlate(s.dragItem.(*Plate))
				s.dragItem.SetDraggable(false)
			}

		} else {
			x, y := s.source.Position()
			x -= s.offsetX
			y -= s.offsetY
			s.dragItem.(*Plate).MoveTo(x, y)
			for _, f := range s.dragItem.(*Plate).GetFoods() {
				f.MoveTo(x, y)
			}
		}
	}
}
