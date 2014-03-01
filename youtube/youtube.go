

package youtube;


import(
    "os"
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "strings"
    "net/url"
    "regexp"
    "os/exec"
    "strconv"
    "errors"
    "time"
    "math"
    "encoding/xml"
    "olli/base"
);




type Video struct {
    sID string
    sAuthorID string
    sAuthorName string
    sTitle string
    sDescription string
    oUploaded time.Time
    oPublished time.Time
    oUpdated time.Time
    oDuration time.Duration
    aSources map[string]*Source
}

type Source struct {
    sSourceTypeID string
    sFileType string
    sQuality string
    iQuality int
    sUrl string
}

type ProcessStatus struct {
    sTimeLeft string
    nPercent float64
}

type XmlData struct {
    XMLName xml.Name `xml:"feed"`
    Entries []XmlVideo `xml:"entry"`
    Links []XmlLink `xml:"link"`
    ID string `xml:"id"`
    Logo string `xml:"logo"`
}

type XmlLink struct {
    Rel string `xml:"rel,attr"`
    Href string `xml:"href,attr"`
}

type XmlVideo struct {
    XMLName xml.Name `xml:"entry"`
    ID string `xml:"media group>videoid,yt"`
    AuthorID string `xml:"author>userId,yt"`
    AuthorName string `xml:"author>name"`
    Title string `xml:"title"`
    Description string `xml:"media group>description,yt"`
    Uploaded string `xml:"media group>uploaded,yt"`
    Published string `xml:"published"`
    Updated string `xml:"updated"`
    Duration XmlDuration `xml:"media group>duration,yt"`
}

type XmlDuration struct {
    Seconds string `xml:"seconds,attr"`
}




func YoutubeDownload (sVideoID string, iMaxQuality int, bMP3 bool, bGain bool, sTargetDir string) error {
    
    if bMP3 {
        _, errAvconv := exec.LookPath("avconv");
        if errAvconv != nil {
            return errors.New("the program 'avconv' is not installed. please install it and try again.");
        }
        _, errLame := exec.LookPath("lame");
        if errLame != nil {
            return errors.New("the program 'lame' is not installed. please install it and try again.");
        }
    }
    
    oVideo, errVideo := oGetVideoData(sVideoID);
    if errVideo != nil {
        return errVideo;
    }
    sUrl, sFileName := oVideo.aMakeDownloadData(iMaxQuality);
    sTargetFileName := sFileName;
    if bMP3 {
        sTargetFileName = sReplaceExtension(sTargetFileName, "mp3");
    }
    fmt.Println(sTargetFileName);
    
    oNow := time.Now();
    sNonce := fmt.Sprintf("%02d%02d%02d_%02d%02d%02d", oNow.Year() % 100, oNow.Month(), oNow.Day(), oNow.Hour(), oNow.Minute(), oNow.Second());
    
    sTempDir := "/tmp";
    sTempFile := sTempDir + "/youtube_download" + "_" + sNonce + "_" + oVideo.sID;
    
    fmt.Println("downloading");
    cDownloadStatus := cDownload(sUrl, sTempFile);
    for {
        oDownloadStatus, bDownloadStatusOpen := <-cDownloadStatus;
        if !bDownloadStatusOpen {break}
        fmt.Printf("  %3v%%  %s       \r", oDownloadStatus.nPercent, oDownloadStatus.sTimeLeft);
    }
    fmt.Printf("\n");
    
    if bMP3 {
        fmt.Println("converting to mp3");
        sTempFileMP3 := sTempFile + ".mp3";
        cConvertStatus, cReplayGain := cConvert(sTempFile, sTempFileMP3, 0);
        for {
            oConvertStatus, bConvertStatusOpen := <-cConvertStatus;
            if !bConvertStatusOpen {break}
            fmt.Printf("  %3v%%  %s       \r", oConvertStatus.nPercent, oConvertStatus.sTimeLeft);
        }
        fmt.Printf("\n");
        nReplayGain := <- cReplayGain;
        fmt.Printf("ReplayGain: " + strconv.FormatFloat(nReplayGain, 'f', 4, 64) + "\n");
        if (bGain) {
            fmt.Println("applying gain");
            os.Remove(sTempFileMP3);
            cConvertStatus, cReplayGain := cConvert(sTempFile, sTempFileMP3, nReplayGain);
            for {
                oConvertStatus, bConvertStatusOpen := <-cConvertStatus;
                if !bConvertStatusOpen {break}
                fmt.Printf("  %3v%%  %s       \r", oConvertStatus.nPercent, oConvertStatus.sTimeLeft);
            }
            fmt.Printf("\n");
            nReplayGain := <- cReplayGain;
            fmt.Printf("ReplayGain: " + strconv.FormatFloat(nReplayGain, 'f', 4, 64) + "\n");
        }
        os.Remove(sTempFile);
        sTempFile = sTempFileMP3;
    }
    
    sTargetFile := sTargetDir + "/" + sTargetFileName;
    os.Rename(sTempFile, sTargetFile);
    
    return nil;
    
}




func GetUserVideoIDs (sUserID string) []string {
    
    aVideoIDs := []string{};
    
    sNextApiCall := "http://gdata.youtube.com/feeds/api/users/" + sUserID + "/uploads?v=2";
    for sNextApiCall != "" {
        oData := XmlData{};
        sPart, _ := sHttpGet(sNextApiCall);
        xml.Unmarshal([]byte(sPart), &oData);
        for _, oVideo := range oData.Entries {
            aVideoID := strings.Split(oVideo.ID, "/");
            sVideoID := aVideoID[len(aVideoID) - 1];
            aVideoIDs = append(aVideoIDs, sVideoID);
        }
        sNextApiCall = "";
        for _, oLink := range oData.Links {
            if oLink.Rel == "next" {
                sNextApiCall = oLink.Href;
                break;
            }
        }
    }
    
    return aVideoIDs;
    
}




func oGetVideoData (sVideoID string) (Video, error) {
    
    sApiCall := "http://gdata.youtube.com/feeds/api/videos/" + sVideoID + "?v=2";
    sData, err := sHttpGet(sApiCall);
    if err != nil {
        return Video{}, err;
    }
    
    oData := XmlVideo{};
    xml.Unmarshal([]byte(sData), &oData);
    
    if oData.ID == "" {
        return Video{}, errors.New("ERROR for " + sVideoID + ":\n" + sData);
    }
    
    oVideo := Video{
        sID: oData.ID,
        sAuthorID: oData.AuthorID,
        sAuthorName: oData.AuthorName,
        sTitle: oData.Title,
        sDescription: oData.Description,
        oUploaded: oParseYoutubeTime(oData.Uploaded),
        oPublished: oParseYoutubeTime(oData.Published),
        oUpdated: oParseYoutubeTime(oData.Updated),
    };
    
    iDurationSeconds, _ := strconv.Atoi(oData.Duration.Seconds);
    oVideo.oDuration = time.Duration(iDurationSeconds) * time.Second;
    
    oVideo.aSources, err = aGetVideoSources(oVideo.sID);
    
    return oVideo, err;
    
}




func aGetVideoSources (sVideoID string) (map[string]*Source, error) {
    
    sApiCall := "http://www.youtube.com/get_video_info?&video_id=" + sVideoID;
    sData, errHttp := sHttpGet(sApiCall);
    
    aData, _ := url.ParseQuery(sData);
    
    aSources := map[string]*Source{};
    
    if errHttp != nil {
        return aSources, errHttp;
    }
    
    if aData["status"][0] != "ok" {
        errYoutube := errors.New("#" + aData["errorcode"][0] + ": " + aData["reason"][0]);
        return aSources, errYoutube;
    }
    
    sStreams := aData["url_encoded_fmt_stream_map"][0];
    aStreams := strings.Split(sStreams, ",");
    
    for iS := 0; iS < len(aStreams); iS ++ {
        
        aStream, _ := url.ParseQuery(aStreams[iS]);
        oSource := &Source{};
        oRegEx := regexp.MustCompile("video\\/([^\\;]*)(;[\\S\\s]*)*$");
        oSource.sFileType = oRegEx.FindStringSubmatch(aStream["type"][0])[1];
        if oSource.sFileType == "x-flv" {
            oSource.sFileType = "flv";
        }
        oSource.sQuality = aStream["quality"][0];
        oSource.iQuality = iTranslateQuality(oSource.sQuality);
        oSource.sSourceTypeID = aStream["itag"][0];
        sDecodedUrl, _ := url.QueryUnescape(aStream["url"][0]);
        oSource.sUrl = sDecodedUrl;// + "&signature=" + aStream["sig"][0];
        aSources[oSource.sSourceTypeID] = oSource;
        
    }
    
    return aSources, nil;
    
}




func (oVideo Video) aMakeDownloadData (iMaxQuality int) (sUrl, sFileName string) {
    
    sVideoID := oVideo.sID;
    oSource := oDetermineBestSource(oVideo.aSources, iMaxQuality);
    
    sUrl = oSource.sUrl;
    
    sAuthorTitle := "_[" + oVideo.sAuthorName + " - " + oVideo.sTitle + "]";
    for _, sReplaceByUnderline := range []string{" ", "/", "|", ":", "?"} {
        sAuthorTitle = strings.Replace(sAuthorTitle, sReplaceByUnderline, "_", -1);
    }
    sSourceTypeID := "_" + oSource.sSourceTypeID;
    sQuality := "_" + fmt.Sprintf("%v", oSource.iQuality) + "p";
    sFileType := "." + oSource.sFileType;
    sFileName = sVideoID + sAuthorTitle + sSourceTypeID + sQuality + sFileType;
    
    return sUrl, sFileName;
    
}




func oDetermineBestSource (aSources map[string]*Source, iMaxQuality int) *Source {
    
    aSrcMap := make(map[int]map[string]*Source);
    for sKey := range aSources {
        iQuality := aSources[sKey].iQuality;
        sFileType := aSources[sKey].sFileType;
        if _, bIsset := aSrcMap[iQuality]; !bIsset {
            aSrcMap[iQuality] = make(map[string]*Source);
        }
        aSrcMap[iQuality][sFileType] = aSources[sKey];
    }
    
    iBestQuality := -1;
    for iQuality := range aSrcMap {
        if iBestQuality < iQuality && iQuality <= iMaxQuality {
            iBestQuality = iQuality;
        }
    }
    if iBestQuality == -1 {
        log.Fatal("-q is too small");
    }
    
    var sBestFileType string;
    for sFileType := range aSrcMap[iBestQuality] {
        sBestFileType = sFileType;
        break;
    }
    for _, sFileType := range []string{"mp4", "webm", "flv"} {
        if _, bIsset := aSrcMap[iBestQuality][sFileType]; bIsset {
            sBestFileType = sFileType;
            break;
        }
    }
    
    oBestSource := aSrcMap[iBestQuality][sBestFileType];
    return oBestSource;
    
}




func oParseYoutubeTime (sYoutubeTime string) time.Time {
    
    aMonths := []time.Month{
        time.January, time.February, time.March, time.April, time.May, time.June, 
        time.July, time.August, time.September, time.October, time.November, time.December,
    };
    
    oRegEx := regexp.MustCompile("(....)-(..)-(..)T(..)\\:(..)\\:(..)\\.(\\d*)Z");
    aTimeData := oRegEx.FindStringSubmatch(sYoutubeTime);
    iYear, _ := strconv.Atoi(aTimeData[1]);
    iMonth, _ := strconv.Atoi(aTimeData[2]);
    iDay, _ := strconv.Atoi(aTimeData[3]);
    iHour, _ := strconv.Atoi(aTimeData[4]);
    iMinute, _ := strconv.Atoi(aTimeData[5]);
    iSecond, _ := strconv.Atoi(aTimeData[6]);
    iOffset, _ := strconv.Atoi(aTimeData[7]);
    oLocation := time.FixedZone("UTC", iOffset);
    oTime := time.Date(iYear, aMonths[iMonth - 1], iDay, iHour, iMinute, iSecond, 0, oLocation);
    
    //fmt.Println(sYoutubeTime);
    //fmt.Printf("%#v %#v %#v %#v %#v %#v \n", oTime.Year(), oTime.Month(), oTime.Day(), oTime.Hour(), oTime.Minute(), oTime.Second());
    
    return oTime;
    
}




func iTranslateQuality (sQuality string) int {
    
    aMap := make(map[string]int);
    aMap["small"] = 240;
    aMap["medium"] = 360;
    aMap["large"] = 480;
    aMap["hd720"] = 720;
    aMap["hd1080"] = 1080;
    iQuality := aMap[sQuality];
    return iQuality;
    
}




func sHttpGet (sUrl string) (string, error) {
    
    if !strings.Contains(sUrl, "http://") {
        sUrl = "http://" + sUrl;
    }
    resp, err := http.Get(sUrl);
    if err != nil {
        return "", err;
    }
    defer resp.Body.Close();
    body, err := ioutil.ReadAll(resp.Body);
    return string(body), err;
    
}




func cDownload (sUrl, sTargetFile string) (<-chan *ProcessStatus) {
    
    cStatus := make(chan *ProcessStatus);
    
    go func(){
        
        _, err := exec.LookPath("wget");
        if err != nil {
            log.Fatal("the program 'wget' is not installed. please install it and try again.");
        }
        
        oCommand := exec.Command("wget", sUrl, "-O", sTargetFile);
        oOutStream, err := oCommand.StderrPipe();
        if err != nil {log.Fatal(err);}
        defer oOutStream.Close();
        err = oCommand.Start();
        if err != nil {log.Fatal(err);}
        var iChunk int64 = 64 * 1024;
        aData := make([]byte, iChunk);
        var iRead int = 0;
        sPart := "";
        o1S := time.Second;//time.ParseDuration("1s");
        for {
            iRead, err = oOutStream.Read(aData);
            if err != nil {break}
            sAppend := fmt.Sprintf("%s", aData[:iRead]);
            sPart += sAppend;
            oRegEx := regexp.MustCompile("(\\d\\d?)\\% [^ ]+ ([\\S\\d]*)\\n");
            aMatches := oRegEx.FindAllStringSubmatch(sPart, -1);
            if len(aMatches) > 0 {
                aLastMatch := aMatches[len(aMatches) - 1];
                sPercent := aLastMatch[1];
                sTimeLeft := aLastMatch[2];
                nPercent, _ := strconv.ParseFloat(sPercent, 64);
                cStatus <- & ProcessStatus {
                    nPercent: nPercent,
                    sTimeLeft: sTimeLeft,
                };
                sPart = "";
            }
            time.Sleep(o1S);
        }
        cStatus <- & ProcessStatus {
            nPercent: 100,
            sTimeLeft: "0s",
        };
        close(cStatus);
        
    }();
    
    return cStatus;
    
}




func cConvert (sSourceFile, sTargetFile string, nGain float64) (<-chan *ProcessStatus, <-chan float64) {
    
    cStatus := make(chan *ProcessStatus);
    cReplayGain := make(chan float64);
    
    go func(){
        
        _, errAvconv := exec.LookPath("avconv");
        if errAvconv != nil {
            log.Fatal("the program 'avconv' is not installed. please install it and try again.");
        }
        _, errLame := exec.LookPath("lame");
        if errLame != nil {
            log.Fatal("the program 'lame' is not installed. please install it and try again.");
        }
        
        sWavFile := sSourceFile + ".wav";
        oCommandWAV := exec.Command("avconv", "-i", sSourceFile, "-vn", sWavFile);
        errWAV := oCommandWAV.Run();
        if errWAV != nil {
            base.Dump("wav error:");
            log.Fatal(errWAV);
        }
        
        nScale := float64(1);
        if nGain != 0 {
            nScale = math.Pow(10, 0.05 * nGain);
        }
        sScale := strconv.FormatFloat(nScale, 'f', 4, 64);
        
        oCommandMP3 := exec.Command("lame", "-b", "320", "--cbr", "--scale", sScale, sWavFile, sTargetFile);
        oOutStream, errOut := oCommandMP3.StderrPipe();
        if errOut != nil {log.Fatal(errOut);}
        defer oOutStream.Close();
        errStart := oCommandMP3.Start();
        if errStart != nil {log.Fatal(errStart);}
        
        var iChunk int64 = 64 * 1024;
        aData := make([]byte, iChunk);
        var iRead int = 0;
        sAll := "";
        sPart := "";
        o1S := time.Second;
        var err error;
        for {
            iRead, err = oOutStream.Read(aData);
            if err != nil {break}
            sAppend := fmt.Sprintf("%s", aData[:iRead]);
            sAll += sAppend;
            sPart += sAppend;
            oRegExA := regexp.MustCompile(" +(\\d+\\/\\d+) +\\(([^\\)]+)\\%\\)[^\\n]+(\\d+\\:\\d+)[^\\n]+\\n\\-*(\\d+\\:\\d+)\\-+");
            aMatches := oRegExA.FindAllStringSubmatch(sPart, -1);
            if len(aMatches) > 0 {
                aLastMatch := aMatches[len(aMatches) - 1];
                sPercent := aLastMatch[2];
                sTimeLeft := aLastMatch[3];
                nPercent, _ := strconv.ParseFloat(sPercent, 64);
                cStatus <- & ProcessStatus {
                    nPercent: nPercent,
                    sTimeLeft: sTimeLeft,
                };
                sPart = "";
            }
            time.Sleep(o1S);
        }
        cStatus <- & ProcessStatus {
            nPercent: 100,
            sTimeLeft: "0s",
        };
        close(cStatus);
        oRegExB := regexp.MustCompile("\nReplayGain: (.+)dB\n");
        aMatches := oRegExB.FindStringSubmatch(sAll);
        sReplayGain := "+0dB";
        if (len(aMatches) > 1) {
            sReplayGain = aMatches[1];
        }
        nReplayGain, _ := strconv.ParseFloat(sReplayGain, 64);
        cReplayGain <- nReplayGain;
        close(cReplayGain)
        
        /*if bGain {
            nScale := math.Pow(10, 0.05 * nReplayGain);
            sScale := strconv.FormatFloat(nScale, 'f', 4, 64);
base.Dump(nReplayGain);
base.Dump("scale:");
base.Dump(nScale);
base.Dump(sScale);
            sTargetFileBeforeScale := sTargetFile + "_before_gain.mp3";
            os.Rename(sTargetFile, sTargetFileBeforeScale);
base.Dump([]string{"lame", "--scale", sScale, sTargetFileBeforeScale, sTargetFile});
            oCommandGain := exec.Command("lame", "--scale", sScale, sTargetFileBeforeScale, sTargetFile);
            errGain := oCommandGain.Run();
            if errGain != nil {
                base.Dump("gain error:");
                log.Fatal(errGain);
            }
        }*/
        
        os.Remove(sWavFile);
        
    }();
    
    return cStatus, cReplayGain;
    
}




func sReplaceExtension (sFile, sNewExtension string) string {
    
    sDir, sFileName, _ := aFileSplit(sFile);
    return sDir + sFileName + "." + sNewExtension;
    
}




func aFileSplit (sFile string) (sDir, sFileNameWithoutExtension, sExtension string) {
    
    sDir = "";
    sFileName := sFile;
    oRegExA := regexp.MustCompile("^(.*\\/)([^\\/]+)$");
    aMatchesA := oRegExA.FindStringSubmatch(sFile);
    if len(aMatchesA) > 0 {
        sDir = aMatchesA[1];
        sFileName = aMatchesA[2];
    }
    sFileNameWithoutExtension = sFileName;
    sExtension = "";
    oRegExB := regexp.MustCompile("^(.*)\\.([^\\.]*)$");
    aMatchesB := oRegExB.FindStringSubmatch(sFileName);
    if len(aMatchesB) > 0 {
        sFileNameWithoutExtension = aMatchesB[1];
        sExtension = aMatchesB[2];
    }
    
    return;
    
}



func Dummy () {
    
    base.Dump("");
    
}



