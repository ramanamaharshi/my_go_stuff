

package list;


import (
    "os"
    "io/ioutil"
    "strings"
    "regexp"
);




type VideoList struct {
    File string
    Videos map[string]string
}




func MakeVideoList (sFile string) *VideoList {
    
    oList := &VideoList{File: sFile, Videos: map[string]string{}};
    
    if !bFileExists(oList.File) {
        ioutil.WriteFile(oList.File, []byte{}, 0744);
    }
    
    bListContent, _ := ioutil.ReadFile(oList.File);
    sListContent := string(bListContent);
    aListRows := strings.Split(sListContent, "\n");
    
    rRow := regexp.MustCompile("(^[_a-zA-Z0-9\\-]+)( (.*))?$");
    for _, sRow := range aListRows {
        aRowParts := rRow.FindStringSubmatch(sRow);
        if len(aRowParts) > 1 {
            sVideoID := aRowParts[1];
            sFlag := "";
            if len(aRowParts) > 3 {
                sFlag = aRowParts[3];
            }
            oList.AddVideo(sVideoID);
            oList.FlagVideo(sVideoID, sFlag);
        }
    }
    
    return oList;
    
}




func (oList *VideoList) AddVideo (sVideoID string) {
    
    if _, bSet := oList.Videos[sVideoID]; !bSet {
        oList.Videos[sVideoID] = "";
    }
    
}


func (oList *VideoList) FlagVideo (sVideoID, sFlag string) {
    
    if _, bSet := oList.Videos[sVideoID]; bSet {
        oList.Videos[sVideoID] = sFlag;
    }
    
}




func (oList *VideoList) GetFreshVideos () []string {
    
    aReturn := []string{};
    for sVideoID, sFlag := range oList.Videos {
        if sFlag == "" {
            aReturn = append(aReturn, sVideoID);
        }
    }
    return aReturn;
    
}




func (oList *VideoList) Write () {
    
    sContent := "";
    for sVideoID, sFlag := range oList.Videos {
        sContent += sVideoID + " " + sFlag + "\n";
    }
    ioutil.WriteFile(oList.File, []byte(sContent), 0744);
    
}




func bFileExists (sFile string) bool {
    
    _, err := os.Stat(sFile);
    return !(err != nil && os.IsNotExist(err));
    
}



