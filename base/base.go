



package base;




import(
    "os"
    "log"
    "fmt"
    "strings"
    "net/http"
    "io/ioutil"
);




func GetArgs (aAcceptedFlags map[string]int) map[string][]string {
    
    aReturn := make(map[string][]string);
    
    var sCurrentFlag string = "";
    
    for iNr, sArg := range os.Args {
        if iNr == 0 {
            continue;
        }
        bFlag := false;
        for sAcceptedFlag, _ := range aAcceptedFlags {
            if sArg == "-" + sAcceptedFlag {
                sCurrentFlag = sAcceptedFlag;
                if _, bIsset := aReturn[sAcceptedFlag]; !bIsset {
                    aReturn[sCurrentFlag] = []string{};
                }
                bFlag = true;
            }
        }
        if !bFlag {
            if len(aReturn[sCurrentFlag]) == aAcceptedFlags[sCurrentFlag] {
                sCurrentFlag = "";
            }
            aReturn[sCurrentFlag] = append(aReturn[sCurrentFlag], sArg);
        }
    }
    
    return aReturn;
    
}




func FileExists (sFile string) bool {
    
    _,err := os.Stat(sFile);
    return err == nil;
    
}




func IsDirectory (sFile string) bool {
    
    oInfo,_ := os.Stat(sFile);
    return oInfo.IsDir();
    
}




func HttpGet (sUrl string) (string, error) {
    
    if !strings.Contains(sUrl, "http://") && !strings.Contains(sUrl, "https://") {
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




func SimpleHttpGet (sUrl string) (string) {
    
    if !strings.Contains(sUrl, "http://") && !strings.Contains(sUrl, "https://") {
        sUrl = "http://" + sUrl;
    }
    resp, err := http.Get(sUrl);
    if err != nil {
        log.Fatal("could not http get \"" + sUrl + "\"\n" + err.Error());
    }
    defer resp.Body.Close();
    body, err := ioutil.ReadAll(resp.Body);
    return string(body);
    
}




func Dump (mStuff interface{}) {
    
    fmt.Printf("%+v\n", mStuff);
    
}



