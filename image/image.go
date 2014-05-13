

package image;


import(
    "os"
    "log"
    "strconv"
    "olli/base"
    goImage "image"
    goImagePNG "image/png"
    goImageColor "image/color"
);


type Image struct {
    Image *goImage.RGBA
}

func NewImage (iW, iH int) *Image {
    oRect := goImage.Rect(0, 0, iW, iH);
    oGoImage := goImage.NewRGBA(oRect);
    return &Image{oGoImage};
}

func NewColor (iR, iG, iB, iA int) *goImageColor.RGBA {
    return &goImageColor.RGBA{uint8(iR), uint8(iG), uint8(iB), uint8(iA)};
}

func (oImage *Image) Width () int {
    return oImage.Image.Bounds().Dx();
}

func (oImage *Image) Height () int {
    return oImage.Image.Bounds().Dy();
}

func (oImage *Image) DrawPoint (iX, iY int, oColor goImageColor.Color) {
    oImage.Image.Set(iX, iY, oColor);
}

func (oImage *Image) FillRect (iX, iY, iW, iH int, oColor goImageColor.Color) {
    for iPX := 0; iPX < iW; iPX ++ {
        for iPY := 0; iPY < iH; iPY ++ {
            oImage.DrawPoint(iX + iPX, iY + iPY, oColor);
        }
    }
}

func (oImage *Image) SaveAsPNG (sFile string) {
    f, err := os.OpenFile(sFile, os.O_CREATE | os.O_WRONLY, 0666);
    if err != nil {log.Fatal(err);}
    err = goImagePNG.Encode(f, oImage.Image);
    if err != nil {log.Fatal(err);}
}

func Pretext () {
    base.Dump(0);
    log.Fatal(strconv.Itoa(1));
}

