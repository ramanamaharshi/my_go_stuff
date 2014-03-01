

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
    
    oStream := audio.NewStream()
    
    var iSampleRate = 44100;
    
    var iMinWaveSamples = iSampleRate / 32;
    var iMaxWaveSamples = iSampleRate / 16;
    
    var nWavePart = float64(0);
    aSound := make([]float32, 1 * iSampleRate / 10);
    rand.Seed(time.Now().Unix());
    var iSoundAt = 0;
    for iSoundAt < len(aSound) {
        var iWaveSamples = iMinWaveSamples + rand.Intn(iMaxWaveSamples - iMinWaveSamples);
        if iWaveSamples > len(aSound) - iSoundAt {
            iWaveSamples = len(aSound) - iSoundAt;
        }
        for iWS := 0; iWS < iWaveSamples; iWS ++ {
            nWavePart = float64(iWS) / float64(iWaveSamples);
            //aSound[iSoundAt + iWS] = float32(0.8) * float32(math.Sin(math.Pi * nWavePart));
            if nWavePart < 0.5 {
                aSound[iSoundAt + iWS] = 1;
            } else {
                aSound[iSoundAt + iWS] = -1;
            }
        }
        iSoundAt += iWaveSamples;
    }
    
    defer oStream.Stop();
    oStream.Start(iSampleRate, time.Duration(10) * time.Millisecond);
    
    /*oStream.Play(0, aSound);
    oStream.Play(0, aSound);*/
    
    for {
        time.Sleep(time.Duration(1) * time.Second);
        //aSound = []float32{-1, -1, -1, -1, 1, 1, 1, 1, 1, 1, 1, 1, -1, -1, -1, -1};
base.Dump("play");
        oStream.PlayNow(aSound);
    }
    
    time.Sleep(time.Duration(999999) * time.Second);
    
}



