

package main;


import(
	"os"
	"time"
	"math"
	"strconv"
	"math/rand"
	"olli/audio"
	"olli/base"
);




func main () {
	
	aArgs := base.GetArgs(map[string]int{"v": 1, "mode": 1});
	sCurrentDir, _ := os.Getwd();
	
	base.Dump(aArgs);
	base.Dump(sCurrentDir);
	
	var nVolume float32 = 1.0;
	if len(aArgs["v"]) > 0 {
		var nVolume64, _ = strconv.ParseFloat(aArgs["v"][0], 32);
		nVolume = float32(nVolume64);
	}
	base.Dump("nVolume = " + strconv.FormatFloat(float64(nVolume), 'f', 4, 32));
	
	oStream := audio.NewStream();
	oStream.SetVolume(nVolume);
	
	var iSampleRate = 44100;
	
	//if _, bSet := aArgs["mode"]; !bSet {
	//	aArgs["mode"] = []string{"a"};
	//}
	//var sMode = aArgs["mode"];
	
	
	var cSamplesWritten = oStream.Start(iSampleRate, time.Duration(100) * time.Millisecond);
	defer oStream.Stop();
	
	var sDuration = "1h";
	if len(aArgs[""]) > 0 {
		sDuration = aArgs[""][0];
	}
	var oDuration, errParse = time.ParseDuration(sDuration);
	if errParse != nil {
		oDuration = time.Duration(60) * time.Minute;
	}
	base.Dump(oDuration);
	
	
	fSound := func (iSampleRate int) []float32 {
		return SoundA(iSampleRate, time.Duration(2) * time.Second, 16, 32);
	};
	if (len(aArgs["mode"]) > 0 && aArgs["mode"][0] == "b") {
		fSound = func (iSampleRate int) []float32 {
			return SoundB(iSampleRate/*, time.Duration(2) * time.Second*/, 16, 32);
		};
	}
	
	
	var iNextPlayStart = 0;
	var oStartTime = time.Now();
	var oLastLoggedTime = time.Time{};
	var iTotalSamples = int(float64(iSampleRate) * oDuration.Seconds());
	var iSamplesPushed = 0;
	
	for {
		var iSamplesWritten = <- cSamplesWritten;
		if iSamplesWritten > iNextPlayStart - 5 * iSampleRate {
			var oNow = time.Now();
			if oNow.Sub(oLastLoggedTime) > (time.Duration(1) * time.Minute) {
				base.Dump(time.Now());
				oLastLoggedTime = oNow;
			}
			var aSound = fSound(iSampleRate);
			for iKey := range(aSound) {
				var nWavePart = float64(iKey) / float64(len(aSound));
				var nVolume = 0.75 + 0.25 * float32(math.Sin(math.Pi * nWavePart));
				aSound[iKey] = nVolume * aSound[iKey];
			}
			var iSamplesLeft = iTotalSamples - iSamplesPushed;
			if iSamplesLeft < 0 {
				break;
			}
			if iSamplesLeft < len(aSound) {
				for iS := iSamplesLeft; iS < len(aSound); iS ++ {
					aSound[iS] = 0;
				}
			}
			oStream.Play(iNextPlayStart, &aSound);
			iSamplesPushed += len(aSound);
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




func SoundB (iSampleRate int, iMinFreq, iMaxFreq int) []float32 {
	
	var oRand = rand.New(rand.NewSource(time.Now().Unix()));
	
	var iDelaySamples = int(float32(4) * float32(iSampleRate)) + oRand.Intn(1 * iSampleRate);
	var iSoundSamples = iSampleRate / 5;
	
	var oDuration = time.Duration(1000 * iSoundSamples / iSampleRate) * time.Millisecond;
	var aMainSound = SoundA(iSampleRate, oDuration, iMinFreq, iMaxFreq);
	
	var aSound = make([]float32, iDelaySamples + iSoundSamples);
	
	for iS := 0; iS < iDelaySamples; iS ++ {
		aSound[iS] = 0;
	}
	
	for iS := 0; iS < iSoundSamples; iS ++ {
		if iS < len(aMainSound) {
			aSound[iDelaySamples + iS] = aMainSound[iS];
		}
	}
	
	return aSound;
	
}



