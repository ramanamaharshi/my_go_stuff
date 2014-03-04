

package main;


import(
    "os"
    "olli/base"
    "olli/audio"
    "time"
    //"math"
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
    
    var iCooldown = iSampleRate * 44 / 100;
    var iNextPlayStart = 0;
    for {
        var iSamplesWritten = <- cSamplesWritten;
        if iSamplesWritten > iNextPlayStart - 2 * iSampleRate {
            var aSound = SoundA(iSampleRate, time.Duration(500) * time.Millisecond, 16, 32);
            if oRand.Intn(32) > 0 {
                oStream.Play(iNextPlayStart, &aSound);
            }
            iNextPlayStart += iCooldown;
        }
        if iSamplesWritten > 60 * 60 * iSampleRate {
            break;
        }
    }
    
    /*oStream.Play(0, aSound);
    oStream.Play(0, aSound);*/
    
    /*for {
        //aSound = []float32{-1, -1, -1, -1, 1, 1, 1, 1, 1, 1, 1, 1, -1, -1, -1, -1};
//base.Dump("play");
        //oStream.PlayNow(aSound);
        time.Sleep(time.Duration(500) * time.Millisecond);
    }*/
    
    time.Sleep(time.Duration(15) * time.Second);
    
}



func SoundB (iSampleRate int, oDuration time.Duration, iRepetitions int) {
    
    var oRand = rand.New(rand.NewSource(time.Now().Unix()));
    var iSoundSamples = int(float64(iSampleRate) * oDuration.Seconds());
    var aSound = make([]float32, iSoundSamples);
    for iR := 0; iR < iRepetitions; iR ++ {
        
    }
    
}



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
        if iWaveSamples > iSoundSamples - iSoundAt {
            iWaveSamples = iSoundSamples - iSoundAt;
        }
        for iWS := 0; iWS < iWaveSamples; iWS ++ {
            nWavePart = float64(iWS) / float64(iWaveSamples);
            //aSound[iSoundAt + iWS] = float32(0.8) * float32(math.Sin(math.Pi * nWavePart));
            if nWavePart < 0.5 {
                aSound[iSoundAt + iWS] = nVolume;
            } else {
                aSound[iSoundAt + iWS] = -nVolume;
            }
        }
        nVolume /= 4;
        if nVolume < 0.2 {
            nVolume = 0.2;
        }
        iSoundAt += iWaveSamples;
    }
    
    return aSound;
    
}



