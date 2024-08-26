package main

import (
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/draw"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	touchIDs []ebiten.TouchID
	strokes  map[*Stroke]struct{}
	food     *Food
	oven     *JetOven
}

var (
	foodImage      *ebiten.Image
	foodAlphaImage *image.Alpha
	plates         []*Plate
)

func resizeImage(img image.Image, width, height int) image.Image {
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.BiLinear.Scale(newImage, newImage.Bounds(), img, img.Bounds(), draw.Over, nil)
	return newImage
}

func newImageWithSize(path string, width, height int) (*ebiten.Image, *image.Alpha, error) {
	file, _ := os.Open(path)
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, nil, err
	}

	img = resizeImage(img, width, height)

	ebitenImage := ebiten.NewImageFromImage(img)

	// Clone an image but only with alpha values.
	// This is used to detect a user cursor touches the image.
	b := img.Bounds()
	ebitenAlphaImage := image.NewAlpha(b)
	for j := b.Min.Y; j < b.Max.Y; j++ {
		for i := b.Min.X; i < b.Max.X; i++ {
			ebitenAlphaImage.Set(i, j, img.At(i, j))
		}
	}
	return ebitenImage, ebitenAlphaImage, nil
}

func init() {
	foodImage_, foodAlphaImage_, err := newImageWithSize("assets/images/cook.png", 50, 50)
	if err != nil {
		log.Fatal(err)
	}
	foodImage = foodImage_
	foodAlphaImage = foodAlphaImage_

	plateImg, plateAlphaImage, err := newImageWithSize("assets/images/plate.png", 60, 60)
	if err != nil {
		log.Fatal(err)
	}
	plate := &Plate{
		Sprite: &Sprite{
			image:      plateImg,
			alphaImage: plateAlphaImage,
			x:          100,
			y:          100,
		},
		foods:     []*Food{},
		draggable: false,
	}
	plate2 := &Plate{
		Sprite: &Sprite{
			image:      plateImg,
			alphaImage: plateAlphaImage,
			x:          200,
			y:          200,
		},
		foods:     []*Food{},
		draggable: false,
	}
	plate3 := &Plate{
		Sprite: &Sprite{
			image:      plateImg,
			alphaImage: plateAlphaImage,
			x:          300,
			y:          300,
		},
		foods:     []*Food{},
		draggable: false,
	}

	plates = append(plates, plate, plate2, plate3)
}

func NewGame() *Game {
	w, h := foodImage.Bounds().Dx(), foodImage.Bounds().Dy()
	f := &Food{
		Sprite: &Sprite{
			image:      foodImage,
			alphaImage: foodAlphaImage,
			x:          rand.Intn(screenWidth - w),
			y:          rand.Intn(screenHeight - h),
		},
		dragged:   false,
		draggable: true,
	}

	ovenImg, ovenAlphaImage, err := newImageWithSize("assets/images/oven.png", 300, 200)
	if err != nil {
		log.Fatal(err)
	}
	oven := NewJetOven(ovenImg, ovenAlphaImage, 50, 50, 1)

	return &Game{
		strokes: map[*Stroke]struct{}{},
		food:    f,
		oven:    oven,
	}
}

func (g *Game) dragItemAt(x, y int) DraggableSprite {
	if g.food.In(x, y) && g.food.GetDraggable() {
		return g.food
	}
	for _, p := range plates {
		if p.In(x, y) && p.GetDraggable() {
			return p
		}
	}
	return nil
}

func (g *Game) Update() error {
	g.oven.Update()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if ds := g.dragItemAt(ebiten.CursorPosition()); ds != nil {
			s := NewStroke(&MouseStrokeSource{}, ds, g.oven, plates)
			g.strokes[s] = struct{}{}
		}
	}
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		if ds := g.dragItemAt(ebiten.TouchPosition(id)); ds != nil {
			s := NewStroke(&TouchStrokeSource{id}, ds, g.oven, plates)
			g.strokes[s] = struct{}{}
		}
	}

	for s := range g.strokes {
		s.Update()
		if !s.dragItem.GetDragged() {
			delete(g.strokes, s)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for t := range plates {
		plates[t].Draw(screen, 1)
	}
	if g.food.dragged {
		g.food.Draw(screen, 0.5)
	} else {
		g.food.Draw(screen, 1)
	}
	g.oven.Draw(screen, 1)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
