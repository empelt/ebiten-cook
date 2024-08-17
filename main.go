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
	sprite   *Sprite
}

var (
	ebitenImage      *ebiten.Image
	ebitenAlphaImage *image.Alpha
)

func resizeImage(img image.Image, width, height int) image.Image {
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.BiLinear.Scale(newImage, newImage.Bounds(), img, img.Bounds(), draw.Over, nil)
	return newImage
}

func init() {
	file, _ := os.Open("assets/images/cook.png")
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	img = resizeImage(img, 50, 50)

	ebitenImage = ebiten.NewImageFromImage(img)

	// Clone an image but only with alpha values.
	// This is used to detect a user cursor touches the image.
	b := img.Bounds()
	ebitenAlphaImage = image.NewAlpha(b)
	for j := b.Min.Y; j < b.Max.Y; j++ {
		for i := b.Min.X; i < b.Max.X; i++ {
			ebitenAlphaImage.Set(i, j, img.At(i, j))
		}
	}
}

func NewGame() *Game {
	w, h := ebitenImage.Bounds().Dx(), ebitenImage.Bounds().Dy()
	s := &Sprite{
		image:      ebitenImage,
		alphaImage: ebitenAlphaImage,
		x:          rand.Intn(screenWidth - w),
		y:          rand.Intn(screenHeight - h),
	}

	return &Game{
		strokes: map[*Stroke]struct{}{},
		sprite:  s,
	}
}

func (g *Game) spriteAt(x, y int) *Sprite {
	if g.sprite.In(x, y) {
		return g.sprite
	}
	return nil
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if sp := g.spriteAt(ebiten.CursorPosition()); sp != nil {
			s := NewStroke(&MouseStrokeSource{}, sp)
			g.strokes[s] = struct{}{}
		}
	}
	g.touchIDs = inpututil.AppendJustPressedTouchIDs(g.touchIDs[:0])
	for _, id := range g.touchIDs {
		if sp := g.spriteAt(ebiten.TouchPosition(id)); sp != nil {
			s := NewStroke(&TouchStrokeSource{id}, sp)
			g.strokes[s] = struct{}{}
		}
	}

	for s := range g.strokes {
		s.Update()
		if !s.sprite.dragged {
			delete(g.strokes, s)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.sprite.dragged {
		g.sprite.Draw(screen, 0.5)
	} else {
		g.sprite.Draw(screen, 1)
	}
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
