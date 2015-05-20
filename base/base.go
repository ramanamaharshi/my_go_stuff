



package base;




import(
	"os"
	"log"
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
);




func GetArgs (aFlagGroups map[string]int) map[string][]string {
	
	aReturn := make(map[string][]string);
	
	sCurrentFlagGroup := "";
	
	for iNr, sArg := range os.Args {
		if iNr == 0 {
			continue;
		}
		bFlag := false;
		for sFlagGroup , _ := range aFlagGroups {
			aFlags := strings.Split(sFlagGroup, " ");
			for _ , sFlag := range aFlags {
				if sArg == "-" + sFlag || sArg == "--" + sFlag {
					sCurrentFlagGroup = sFlagGroup;
					if _ , bIsset := aReturn[sCurrentFlagGroup]; !bIsset {
						aReturn[sCurrentFlagGroup] = []string{};
					}
					bFlag = true;
				}
			}
		}
		if !bFlag {
			if len(aReturn[sCurrentFlagGroup]) == aFlagGroups[sCurrentFlagGroup] {
				sCurrentFlagGroup = "";
			}
			aReturn[sCurrentFlagGroup] = append(aReturn[sCurrentFlagGroup], sArg);
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



