package main

import (
    "olli/base"
    "github.com/0xe2-0x9a-0x9b/Go-SDL/sdl"
    "strconv"
    "math/rand"
    "unsafe"
    "math"
    "time"
    "log"
    "fmt"
)
 
const WIDTH = 1200;
const HEIGHT = 900;

func main() {
    
    iWidth := int32(WIDTH);
    iHeight := int32(HEIGHT);
    
    var screen = sdl.SetVideoMode(int(iWidth), int(iHeight), 32, sdl.RESIZABLE);
    
    if screen == nil {
        log.Fatal(sdl.GetError());
    }
    
    iIterationNr := 0;
    iIterations := int(math.Min(float64(333333), float64((WIDTH) * (HEIGHT))));
    iStartNS := time.Now().UnixNano();
    
    //mainloop:
    for {
        
        var iX int32;
        var iY int32;
        var iColor uint32;
        for i := 0; i < iIterations; i ++ {
            iX = rand.Int31() % (iWidth);
            iY = rand.Int31() % (iHeight);
            iColor = rand.Uint32();
            iColor &= 0x00ffffff;
            iColor += 0x88000000;
            draw_point(iX, iY, iColor, screen);
            //break mainloop;
        }
        
        screen.Flip();
        
        iIterationNr ++;
        if iIterationNr % 99 == 0 {
            iIterationNS := time.Now().UnixNano();
            iDeltaNS := iIterationNS - iStartNS;
            nDeltaS := float64(iDeltaNS) / float64(1000000000);
            nFPS := float64(iIterationNr) / nDeltaS;
            fmt.Printf("fps: %v\n", nFPS);
            fmt.Printf("%v x %v\n", screen.W, screen.H);
        }
        
    }
    base.Dump("");
    
}

func print_color (iColor uint32) {
    iRed := (iColor >> 16) & 0x000000ff;
    iGreen := (iColor >> 8) & 0x000000ff;
    iBlue := (iColor >> 0) & 0x000000ff;
    sRed := strconv.Itoa((int)(iRed));
    sGreen := strconv.Itoa((int)(iGreen));
    sBlue := strconv.Itoa((int)(iBlue));
    base.Dump(sRed + "," + sGreen + "," + sBlue);
}

func merge_colors (iColorA, iColorB uint32) uint32 {
    iRedA := (iColorA >> 16) & 0x000000ff;
    iGreenA := (iColorA >> 8) & 0x000000ff;
    iBlueA := (iColorA >> 0) & 0x000000ff;
    iAlphaB := iColorB >> 24;
    iRedB := (iColorB >> 16) & 0x000000ff;
    iGreenB := (iColorB >> 8) & 0x000000ff;
    iBlueB := (iColorB >> 0) & 0x000000ff;
    iRed := ((256 - iAlphaB) * iRedA + iAlphaB * iRedB) / 256;
    iGreen := ((256 - iAlphaB) * iGreenA + iAlphaB * iGreenB) / 256;
    iBlue := ((256 - iAlphaB) * iBlueA + iAlphaB * iBlueB) / 256;
    return (iRed << 16) + (iBlue << 8) + (iGreen << 0);
}

func draw_point (iX, iY int32, iColor uint32, screen *sdl.Surface) {
    
    iPixelNr := (uintptr)(screen.Pixels) + (uintptr)(4 * (iY * screen.W + iX));
    pPixelPointer := (*uint32)(unsafe.Pointer(iPixelNr));
    iOldColor := *pPixelPointer;
    
    //iOpacity := iColor >> 24;
    //iNewColor := ((256 - iOpacity) * iOldColor + iOpacity * iColor) / 256;
    iNewColor := merge_colors(iOldColor, iColor);
    
    //print_color(iOldColor);
    //print_color(iColor);
    //base.Dump(iOpacity);
    //print_color(iNewColor);
    
    *pPixelPointer = iNewColor;
    
}

func draw_rect (iX, iY, iW, iH int32, iColor uint32, screen *sdl.Surface) {
    
    var iPixelBytes = unsafe.Sizeof(iColor);
    var iPixStart = uintptr(screen.Pixels);
    var pix = uintptr(0);
    for iPY := iY; iPY < iY + iH; iPY ++ {
        pix = iPixStart + (uintptr) ((iPY * screen.W) + iX) * iPixelBytes;
        for iPX := iX; iPX < iX + iW; iPX ++ {
            var pu = unsafe.Pointer(pix);
            var pp *uint32;
            pp = (*uint32) (pu);
            *pp = iColor;
            pix += iPixelBytes;
        }
    }
    
}
