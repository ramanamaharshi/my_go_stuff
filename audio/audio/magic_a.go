

package main;


import(
    "os"
    "olli/base"
    "olli/audio"
    "olli/image"
    "time"
    "math"
    "math/rand"
);




func main () {
    
    aArgs := base.GetArgs(map[string]int{});
    sCurrentDir, _ := os.Getwd();
    
    base.Dump(aArgs);
    base.Dump(sCurrentDir);
    
    oStream := audio.NewStream();
    
    var iSampleRate = 44100;
    
    var oRand = rand.New(rand.NewSource(time.Now().Unix()));
    
    var cSamplesWritten = oStream.Start(iSampleRate, time.Duration(100) * time.Millisecond);
    defer oStream.Stop();
    
    var oStartTime = time.Now();
    
    //var iCooldownSamples = iSampleRate * 44 / 100;
    //var oCooldownDuration = (1000 * time.Duration(iCooldownSamples) / time.Duration(iSampleRate)) * time.Millisecond;
    var sDuration = "1h";
    if len(aArgs[""]) > 0 {
        sDuration = aArgs[""][0];
    }
    var oDuration, errParse = time.ParseDuration(sDuration);
    if errParse != nil {
        oDuration = time.Duration(60) * time.Minute;
    }
base.Dump(oDuration);
    
    var oImage = image.NewImage(20000, 200);
    oImage.FillRect(0, 0, oImage.Width(), oImage.Height(), image.NewColor(255, 255, 255, 255));
    oImage.DrawPoint(10, 10, image.NewColor(255, 0, 0, 255));
    
    var iNextPlayStart = 0;
    for {
        var iSamplesWritten = <- cSamplesWritten;
        if iSamplesWritten > iNextPlayStart - 5 * iSampleRate {
base.Dump(time.Now());
            var aSound = SoundA(iSampleRate, time.Duration(2) * time.Second / time.Duration(8), 16, 32);
            for iKey := range(aSound) {
                var nWavePart = float64(iKey) / float64(len(aSound));
                var nVolume = 0.75 + 0.25 * float32(math.Sin(math.Pi * nWavePart));
                //nVolume = float32(1.0);
                aSound[iKey] = nVolume * aSound[iKey];
            }
            // for iS := 0; iS < len(aSound); iS ++ {}
            if true || oRand.Intn(32) > 0 {
                oStream.Play(iNextPlayStart, &aSound);
                drawSound(oImage, &aSound, iNextPlayStart, 8);
            }
            iNextPlayStart += len(aSound);
        }
        var oTimePassed = time.Now().Sub(oStartTime);
        if oTimePassed > oDuration {
             break;
        }
        time.Sleep(time.Duration(50) * time.Millisecond);
    }
    
    /*oStream.Play(0, aSound);
    oStream.Play(0, aSound);*/
    
    /*for {
        //aSound = []float32{-1, -1, -1, -1, 1, 1, 1, 1, 1, 1, 1, 1, -1, -1, -1, -1};
//base.Dump("play");
        //oStream.PlayNow(aSound);
        time.Sleep(time.Duration(500) * time.Millisecond);
    }*/
    
    time.Sleep(time.Second);
    
}


func drawSound (oImage *image.Image, aSound *[]float32, iSoundStart int, nCompression float32) {
    if int(float32(iSoundStart) / nCompression) > oImage.Width() {return;}
    for iSampleNr, nSampleValue := range(*aSound) {
        iAt := iSoundStart + iSampleNr;
        iImageAt := int(float32(iAt) / nCompression);
        if (iImageAt >= oImage.Width()) {
            oImage.SaveAsPNG("sound.png");
            break;
        }
        oImage.DrawPoint(iImageAt, 100 + int(100.0 * nSampleValue), image.NewColor(0, 0, 0, 255));
    }
}


/*func SoundB (iSampleRate int, oDuration time.Duration, iRepetitions int) {
    
    var oRand = rand.New(rand.NewSource(time.Now().Unix()));
    var iSoundSamples = int(float64(iSampleRate) * oDuration.Seconds());
    var aSound = make([]float32, iSoundSamples);
    for iR := 0; iR < iRepetitions; iR ++ {
        
    }
    
}*/



func SoundA (iSampleRate int, oDuration time.Duration, iMinFreq, iMaxFreq int) []float32 {
    
    var oRand = rand.New(rand.NewSource(time.Now().Unix()));
    var iSoundSamples = int(float64(iSampleRate) * oDuration.Seconds());
    var aSound = make([]float32, iSoundSamples);
    var iMinWaveSamples = iSampleRate / iMaxFreq;
    var iMaxWaveSamples = iSampleRate / iMinFreq;
    var nWavePart = float64(0);
    var nVolume = float32(1);
    var iSoundAt = 0;
    for iSoundAt < iSoundSamples {
        var iWaveSamples = iMinWaveSamples + oRand.Intn(iMaxWaveSamples - iMinWaveSamples);
        iWaveSamples = 888;
        if iWaveSamples > iSoundSamples - iSoundAt {
            iWaveSamples = iSoundSamples - iSoundAt;
        }
        if (iSoundSamples - iSoundAt - iWaveSamples < iMinWaveSamples) {
            iWaveSamples = iSoundSamples - iSoundAt;
        }
        for iWS := 0; iWS < iWaveSamples; iWS ++ {
            nWavePart = float64(iWS) / float64(iWaveSamples);
            if nWavePart < 0.5 {
                aSound[iSoundAt + iWS] = nVolume;
            } else {
                aSound[iSoundAt + iWS] = -nVolume;
            }
            aSound[iSoundAt + iWS] = nVolume * float32(math.Sin(2 * math.Pi * nWavePart));
        }
        //nVolume /= 4;
        /*if nVolume < 0.05 {
            nVolume = 0.05;
        }*/
        iSoundAt += iWaveSamples;
    }
    
    return aSound;
    
}



