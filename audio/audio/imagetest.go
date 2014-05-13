package main;

import (
        olliImage "olli/image"
);

func main() {
        oOlliImage := olliImage.NewImage(256, 256);
        for iX := 0; iX < oOlliImage.Width(); iX ++ {
            for iY := 0; iY < oOlliImage.Height(); iY ++ {
                oOlliImage.DrawPoint(iX, iY, olliImage.NewColor(iX, iY, (iX + iY / 2), 255));
            }
        }
        oOlliImage.SaveAsPNG("test.png");
}