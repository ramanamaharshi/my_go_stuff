
package base;

import(
    "os"
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




func HttpGet (sUrl string) (string, error) {
    
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




func Dump (mStuff interface{}) {
    fmt.Printf("%+v\n", mStuff);
}



