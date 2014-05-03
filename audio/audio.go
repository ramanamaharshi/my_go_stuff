

package audio;


import(
    "olli/base"
    "olli/audio/pulsego"
    "log"
    "time"
    //"math"
    "strconv"
);


type Stream struct {
    nVolume float32
    oPendingSounds *pendingSounds
    cNewSounds chan *pendingSound
    cSamplesWritten chan int
    iSamplesWritten int
    bKill bool
}


type pendingSounds struct {
    oFirstPendingSound *pendingSound
}


type pendingSound struct {
    aSamples *[]float32
    iStartSampleNr int
    oPrev *pendingSound
    oNext *pendingSound
}


func NewStream () *Stream {
    var oStream = &Stream{};
    oStream.nVolume = 1;
    oStream.oPendingSounds = &pendingSounds{};
    oStream.cNewSounds = make(chan *pendingSound, 999999);
    oStream.cSamplesWritten = make(chan int, 999999);
    return oStream;
}


func (oStream *Stream) SetVolume (nVolume float32) {
    oStream.nVolume = nVolume;
}


func (oStream *Stream) PlayNow (aSound *[]float32) {
    
    oStream.Play(-1, aSound);
    
}


func (oStream *Stream) Play (iStartSampleNr int, aSound *[]float32) {
    
    var oNewPendingSound = &pendingSound{
        aSamples: aSound,
        iStartSampleNr: iStartSampleNr,
    };
    oStream.oPendingSounds.vAdd(oNewPendingSound);
    //oStream.cNewSounds <- oNewPendingSound;
    
}


func (oStream *Stream) vDrainNewSounds () {
    
    var bSoundsLeft = true;
    for bSoundsLeft {
        select {
            case oNewSound := <- oStream.cNewSounds:
                if (oNewSound.iStartSampleNr == -1) {
                    oNewSound.iStartSampleNr = oStream.iSamplesWritten;
                }
                oStream.oPendingSounds.vAdd(oNewSound);
            default:
                bSoundsLeft = false;
        }
    }
    
}


func (oPendingSounds *pendingSounds) iCount () int {
    
    var iSounds = 0;
    var oPendingSound = oPendingSounds.oFirstPendingSound;
    for oPendingSound != nil {
        oPendingSound = oPendingSound.oNext;
        iSounds ++;
    }
    return iSounds;
    
}


func (oPendingSounds *pendingSounds) vAdd (oPendingSound *pendingSound) {
    
    var oPrev *pendingSound;
    var oNext = oPendingSounds.oFirstPendingSound;
    for oNext != nil && oNext.iStartSampleNr < oPendingSound.iStartSampleNr {
        oPrev = oNext;
        oNext = oNext.oNext;
    }
    oPendingSound.oPrev = oPrev;
    oPendingSound.oNext = oNext;
    if oPrev != nil {
        oPrev.oNext = oPendingSound;
    }
    if oNext != nil {
        oNext.oPrev = oPendingSound;
    }
    if oPendingSound.oPrev == nil {
        oPendingSounds.oFirstPendingSound = oPendingSound;
    }
    
}


func (oPendingSounds *pendingSounds) vRemove (oPendingSound *pendingSound) {
    
    if oPendingSound.oNext != nil {
        oPendingSound.oNext.oPrev = oPendingSound.oPrev;
    }
    if oPendingSound.oPrev != nil {
        oPendingSound.oPrev.oNext = oPendingSound.oNext;
    }
    if oPendingSound == oPendingSounds.oFirstPendingSound {
        oPendingSounds.oFirstPendingSound = oPendingSound.oNext;
    }
    
}


func (oPendingSounds *pendingSounds) vWriteOnBuffer (iStartSample int, aBuffer []float32) {
    
    var iBufferSize = len(aBuffer);
    for iB := range aBuffer {
        aBuffer[iB] = float32(0);
    }
    var oPendingSound = oPendingSounds.oFirstPendingSound;
    for oPendingSound != nil {
        if oPendingSound.iStartSampleNr <= iStartSample + iBufferSize {
            var iStartSoundAt = iStartSample - oPendingSound.iStartSampleNr;
            var iVon = oPendingSound.iStartSampleNr - iStartSample;
            if iVon < -iStartSoundAt {
                iVon = -iStartSoundAt;
            }
            if iVon < 0 {
                iVon = 0;
            }
            var iSamples = len(*oPendingSound.aSamples);
            var iBis = oPendingSound.iStartSampleNr - iStartSample + iSamples;
            if iBis > iStartSoundAt + iSamples {
                iBis = iStartSoundAt + iSamples;
            }
            if iBis > iBufferSize {
                iBis = iBufferSize;
            }
            for iB := iVon; iB < iBis; iB ++ {
                aBuffer[iB] += (*oPendingSound.aSamples)[iStartSoundAt + iB];
            }
            if oPendingSound.iStartSampleNr + len(*oPendingSound.aSamples) < iStartSample {
                oPendingSounds.vRemove(oPendingSound);
            }
        }
        oPendingSound = oPendingSound.oNext;
    }
    
}


func (oStream *Stream) Start (iSampleRate int, oBufferDuration time.Duration) <-chan int {
    
    go oStream.vLoop(iSampleRate, oBufferDuration);
    return oStream.cSamplesWritten;
    
}

func (oStream *Stream) Stop () {
    
    oStream.bKill = true;
    
}


func (oStream *Stream) vLoop (iSampleRate int, oBufferDuration time.Duration) {
    
    pa := pulsego.NewPulseMainLoop();
    defer pa.Dispose();
    pa.Start();
    
    oPulseContext := pa.NewContext("default", 0);
    if oPulseContext == nil {
        log.Fatal("Failed to create a new context");
    }
    defer oPulseContext.Dispose();
    oPulseStream := oPulseContext.NewStream("default", &pulsego.PulseSampleSpec {
        Format:pulsego.SAMPLE_FLOAT32LE,
        Rate: iSampleRate,
        Channels: 1,
    });
    if oPulseStream == nil {
        log.Fatal("Failed to create a new stream");
    }
    defer oPulseStream.Dispose();
    oPulseStream.ConnectToSink();
    
    var oStartTime = time.Now();
    oStream.iSamplesWritten = 0;
    
    var iBufferSize = int(float64(iSampleRate) * oBufferDuration.Seconds());
    aBuffer := make([]float32, iBufferSize);
    
    var iLoop = 0;
    for !oStream.bKill {
        iLoop ++;
//base.Dump("for loop #" + strconv.Itoa(iLoop));
        oStream.vDrainNewSounds();
        oStream.oPendingSounds.vWriteOnBuffer(oStream.iSamplesWritten, aBuffer);
        for iB := range(aBuffer) {
            aBuffer[iB] *= oStream.nVolume;
        }
        oPulseStream.Write(aBuffer, pulsego.SEEK_RELATIVE);
        oStream.iSamplesWritten += iBufferSize;
        oStream.cSamplesWritten <- oStream.iSamplesWritten;
        var oTimePassed = time.Now().Sub(oStartTime);
        var oTimeWritten = time.Duration(oStream.iSamplesWritten) * time.Second / time.Duration(iSampleRate);
        var oTimeDelta = oTimeWritten - oTimePassed;
//base.Dump("----");
//base.Dump(oTimeWritten);
//base.Dump(oTimePassed);
//base.Dump("----");
        if (oTimeDelta > oBufferDuration / time.Duration(2)) {
            time.Sleep(oTimeDelta - oBufferDuration / time.Duration(2));
        }
    }
    
}


func Pretext () {
    base.Dump(strconv.Itoa(1));
}

