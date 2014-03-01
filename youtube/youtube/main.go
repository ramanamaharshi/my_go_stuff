

package main;


import(
    "os"
    "log"
    "fmt"
    "strconv"
    "olli/base"
    "olli/youtube"
    "olli/youtube/list"
);




const VideoSuccessFlag = "✔";
const VideoErrorFlag = "✖";




func main () {
    
    aArgs := base.GetArgs(map[string]int{"tel": 0, "mp3": 0, "gain": 0, "user_uploads": 1, "format": 1, "q": 1});
    sCurrentDir, _ := os.Getwd();
    
    sTargetDir := sCurrentDir;
    aVideoIDs := aArgs[""];
    
    _, bTel := aArgs["tel"];
    _, bMP3 := aArgs["mp3"];
    _, bGain := aArgs["gain"];
    if bTel {
        bMP3 = true;
        bGain = true;
        sTargetDir = "/aaa/downloads/telefon"
    }
    
    bUseList := false;
    oList := &list.VideoList{};
    if _, bSet := aArgs["user_uploads"]; bSet {
        bUseList = true;
        sUserID := aArgs["user_uploads"][0];
        sTargetDir = sCurrentDir + "/" + sUserID;
        os.Mkdir(sTargetDir, 0744);
        sVideoListFile := sTargetDir + "/download_list";
        oList = list.MakeVideoList(sVideoListFile);
        if len(oList.Videos) == 0 {
            aVideoIDs := youtube.GetUserVideoIDs(sUserID);
            for _, sVideoID := range aVideoIDs {
                oList.AddVideo(sVideoID);
            }
            oList.Write();
        }
        aVideoIDs = oList.GetFreshVideos();
    }
    
    var iMaxQuality = 999999;
    if _, bSet := aArgs["q"]; bSet {
        if len(aArgs["q"]) > 0 {
            iMaxQuality64, _ := strconv.ParseInt(aArgs["q"][0], 10, 64);
            iMaxQuality = int(iMaxQuality64);
        }
    } else {
        if bMP3 {
            iMaxQuality = 360;
        }
    }
    
    for _, sVideoID := range aVideoIDs {
        err := youtube.YoutubeDownload(sVideoID, iMaxQuality, bMP3, bGain, sTargetDir);
        fmt.Println(sVideoID + " fertig");
        if bUseList {
            if err == nil {
                oList.FlagVideo(sVideoID, VideoSuccessFlag);
            } else {
                oList.FlagVideo(sVideoID, VideoErrorFlag + err.Error());
            }
            oList.Write();
        } else {
            if err != nil {
                log.Fatal(err);
            }
        }
    }
    
}
