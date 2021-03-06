

package main;


import(
    "os"
    "olli/base"
    "olli/audio"
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
    var iNextPlayStart = 0;
base.Dump(time.Now());
    for {
        var iSamplesWritten = <- cSamplesWritten;
        if iSamplesWritten > iNextPlayStart - 5 * iSampleRate {
            var aSound = SoundA(iSampleRate, time.Duration(2) * time.Second / time.Duration(4), 16, 32);
            for iKey := range(aSound) {
                var nWavePart = float64(iKey) / float64(len(aSound));
                var nVolume = 0.75 + 0.25 * float32(math.Sin(math.Pi * nWavePart));
                aSound[iKey] = nVolume * aSound[iKey];
            }
            for iS := 0; iS < len(aSound); iS ++ {
            }
            if true || oRand.Intn(32) > 0 {
                oStream.Play(iNextPlayStart, &aSound);
            }
            iNextPlayStart += len(aSound);
        }
        var oTimePassed = time.Now().Sub(oStartTime);
        if oTimePassed > time.Duration(1) * time.Hour {
             break;
        }
        //time.Sleep(time.Duration(50) * time.Millisecond);
    }
base.Dump(time.Now());
    
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
        if iWaveSamples > iSoundSamples - iSoundAt {
            iWaveSamples = iSoundSamples - iSoundAt;
        }
        for iWS := 0; iWS < iWaveSamples; iWS ++ {
            nWavePart = float64(iWS) / float64(iWaveSamples);
            if nWavePart < 0.5 {
                aSound[iSoundAt + iWS] = nVolume;
            } else {
                aSound[iSoundAt + iWS] = -nVolume;
            }
            //aSound[iSoundAt + iWS] = nVolume * float32(math.Sin(math.Pi * nWavePart));
        }
        //nVolume /= 4;
        /*if nVolume < 0.05 {
            nVolume = 0.05;
        }*/
        iSoundAt += iWaveSamples;
    }
    
    return aSound;
    
}



