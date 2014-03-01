

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
    
    oStream := audio.NewStream()
    
    var iSampleRate = 44100;
    
    var iMinWaveSamples = iSampleRate / 32;
    var iMaxWaveSamples = iSampleRate / 16;
    
    var nWavePart = float64(0);
    aSound := make([]float32, 64 * iSampleRate);
    rand.Seed(time.Now().Unix());
    var iSoundAt = 0;
    for iSoundAt < len(aSound) {
        var iWaveSamples = iMinWaveSamples + rand.Intn(iMaxWaveSamples - iMinWaveSamples);
        if iWaveSamples > len(aSound) - iSoundAt {
            iWaveSamples = len(aSound) - iSoundAt;
        }
        for iWS := 0; iWS < iWaveSamples; iWS ++ {
            nWavePart = float64(iWS) / float64(iWaveSamples);
            aSound[iSoundAt + iWS] = float32(0.8) * float32(math.Sin(math.Pi * nWavePart));
            if nWavePart < 0.5 {
                aSound[iSoundAt + iWS] = 1;
            } else {
                aSound[iSoundAt + iWS] = -1;
            }
        }
        iSoundAt += iWaveSamples;
    }
    //base.Dump(aSound);
    oStream.Play(0, aSound);
    
    oStream.Start(iSampleRate);
    
    time.Sleep(time.Duration(100000) * time.Second);
    
}



