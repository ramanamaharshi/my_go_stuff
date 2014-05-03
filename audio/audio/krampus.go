

package main;


import(
    "os"
    "time"
    "math"
    "math/rand"
    "olli/audio"
    "olli/base"
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
    
    var sDuration = "1h";
    if len(aArgs[""]) > 0 {
        sDuration = aArgs[""][0];
    }
    var oDuration, errParse = time.ParseDuration(sDuration);
    if errParse != nil {
        oDuration = time.Duration(60) * time.Minute;
    }
    base.Dump(oDuration);
    var iNextPlayStart = 0;
    var oLastLoggedTime = time.Time{};
    for {
        var iSamplesWritten = <- cSamplesWritten;
        if iSamplesWritten > iNextPlayStart - 5 * iSampleRate {
            var oNow = time.Now();
            if oNow.Sub(oLastLoggedTime) > (time.Duration(1) * time.Minute) {
                base.Dump(time.Now());
                oLastLoggedTime = oNow;
            }
            var aSound = SoundA(iSampleRate, time.Duration(2) * time.Second / time.Duration(8), 16, 32);
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
        if oTimePassed > oDuration {
             break;
        }
        time.Sleep(time.Duration(50) * time.Millisecond);
    }
    
    time.Sleep(time.Second);
    
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
            if nWavePart < 0.5 {
                aSound[iSoundAt + iWS] = nVolume;
            } else {
                aSound[iSoundAt + iWS] = -nVolume;
            }
        }
        iSoundAt += iWaveSamples;
    }
    
    return aSound;
    
}



