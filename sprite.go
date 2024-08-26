package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type DraggableSprite interface {
	GetDragged() bool
	SetDragged(bool)
	GetDraggable() bool
	SetDraggable(bool)
	MoveTo(int, int)
	GetX() int
	GetY() int
}

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

func (s *Sprite) GetX() int {
	return s.x
}

func (s *Sprite) GetY() int {
	return s.y
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
	dragged   bool
	draggable bool
}

func NewFood(image *ebiten.Image, alphaImage *image.Alpha, x, y int) *Food {
	return &Food{
		Sprite:    &Sprite{image: image, alphaImage: alphaImage, x: x, y: y},
		dragged:   false,
		draggable: true,
	}
}

func (f *Food) GetDragged() bool {
	return f.dragged
}

func (f *Food) SetDragged(dragged bool) {
	f.dragged = dragged
}

func (f *Food) GetDraggable() bool {
	return f.draggable
}

func (f *Food) SetDraggable(draggable bool) {
	f.draggable = draggable
}

type Plate struct {
	*Sprite
	foods     []*Food
	dragged   bool
	draggable bool
}

func NewPlate(image *ebiten.Image, alphaImage *image.Alpha, x, y int) *Plate {
	return &Plate{
		Sprite:    &Sprite{image: image, alphaImage: alphaImage, x: x, y: y},
		foods:     []*Food{},
		dragged:   false,
		draggable: false,
	}
}

func (p *Plate) GetDragged() bool {
	return p.dragged
}

func (p *Plate) SetDragged(dragged bool) {
	p.dragged = dragged
}

func (p *Plate) GetDraggable() bool {
	return p.draggable
}

func (p *Plate) SetDraggable(draggable bool) {
	p.draggable = draggable
}

// FIXME: 食べ物を取り除く可能性がないなら、imageの更新でもいいかも。そうするとこのコードはいらない
func (p *Plate) GetFoods() []*Food {
	return p.foods
}

func (p *Plate) AddFood(f *Food) {
	p.foods = append(p.foods, f)
}

func (p *Plate) RemoveFood(f *Food) {
	for i, v := range p.foods {
		if v == f {
			p.foods = append(p.foods[:i], p.foods[i+1:]...)
			return
		}
	}
}

type JetOven struct {
	*Sprite
	plates   []*Plate
	velocity int
}

func NewJetOven(image *ebiten.Image, alphaImage *image.Alpha, x, y, velocity int) *JetOven {
	return &JetOven{
		Sprite:   &Sprite{image: image, alphaImage: alphaImage, x: x, y: y},
		plates:   []*Plate{},
		velocity: velocity,
	}
}

func (j *JetOven) Update() {
	for _, p := range j.plates {
		for _, f := range p.foods {
			f.MoveTo(f.x+j.velocity, f.y)
		}
		p.MoveTo(p.x+j.velocity, p.y)
		if p.x > j.x+200 {
			j.RemovePlate(p)
			p.SetDraggable(true)
		}
	}
}

func (j *JetOven) AddPlate(p *Plate) {
	j.plates = append(j.plates, p)
}

func (j *JetOven) RemovePlate(p *Plate) {
	for i, v := range j.plates {
		if v == p {
			j.plates = append(j.plates[:i], j.plates[i+1:]...)
			return
		}
	}
}
