

package main;


import(
	"os"
	"log"
	"fmt"
	"strings"
	"strconv"
	"olli/base"
	"olli/youtube"
	"olli/youtube/list"
);




const VideoSuccessFlag = "✔";
const VideoErrorFlag = "✖";




func main () {
	
	aArgs := base.GetArgs(map[string]int{"?": 0, "h": 0, "q": 1, "tel": 0, "mp3": 0, "gain": 0, "user_uploads": 1, "format": 1, "convert_only": 0});
	sCurrentDir, _ := os.Getwd();
	
	sTargetDir := sCurrentDir;
	aVideoIDs := aArgs[""];
	for sKey, sID := range(aVideoIDs) {
		aParams := aGetParams(sID);
		if sParam, bSet := aParams["v"]; bSet {
			aVideoIDs[sKey] = sParam;
		}
	}
	
	_, bHelp := aArgs["h"];
	_, bQuestionMark := aArgs["?"];
	if bHelp || bQuestionMark {
		fmt.Println("this is for downloading youtube videos.");
		fmt.Println("usage: youtube [OPTIONS] [YOUTUBE_IDS]");
		fmt.Println("options:");
		fmt.Println("-h, --help: this help page");
		fmt.Println("-q [max_quality]: set maximum quality. standard youtube qualities are: 240, 320, 480, 720, 1080");
		fmt.Println("-mp3: convert to mp3");
		fmt.Println("-gain: make mp3 louder if very quiet");
		fmt.Println("-user_uploads [user_id]: download all videos from a youtube user");
		return;
	}
	
	_, bTel := aArgs["tel"];
	_, bMP3 := aArgs["mp3"];
	_, bGain := aArgs["gain"];
	_, bConvertOnly := aArgs["convert_only"];
	
	if bConvertOnly {
		return;
	}
	
	if bTel {
		bMP3 = true;
		bGain = true;
		sTargetDir = "/aaa/downloads/telefon";
		if _, err := os.Stat(sTargetDir); os.IsNotExist(err) {
			sTargetDir = "/aaa/cache/telefon";
		}
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
	
	fFileNameMaker := func (oDownloadData *youtube.DownloadData) string {
		if bTel {
			sAuthorTitle := oDownloadData.Author + " - " + oDownloadData.Title;
			for _, sReplaceWithUnderline := range []string{" ", "/", ":", "?", "|", "\"", "'"} {
				sAuthorTitle = strings.Replace(sAuthorTitle, sReplaceWithUnderline, "_", -1);
			}
			sFileName := oDownloadData.VideoID + "_[" + sAuthorTitle + "]_";
			sFileName += oDownloadData.SourceTypeID + "_" + oDownloadData.Quality + "." + oDownloadData.FileType;
			return sFileName;
		}
		return oDownloadData.FileName;
	}
	
	for _, sVideoID := range aVideoIDs {
		err := youtube.YoutubeDownload(sVideoID, iMaxQuality, bMP3, bGain, sTargetDir, fFileNameMaker);
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




func aGetParams (sUrl string) map[string]string {
	
	aParams := map[string]string{};
	aUrl := strings.Split(sUrl, "?");
	if len(aUrl) > 1 {
		sParams := aUrl[1];
		aParamStrings := strings.Split(sParams, "&");
		for _, sParamString := range(aParamStrings) {
			aParamString := strings.Split(sParamString, "=");
			sValue := "";
			if len(aParamString) > 1 {
				sValue = aParamString[1];
			}
			aParams[aParamString[0]] = sValue;
		}
	}
	return aParams;
	
}



