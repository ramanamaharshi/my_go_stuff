

package main;


import(
    //"os"
    "log"
    //"fmt"
    "strings"
    "io/ioutil"
    "olli/base"
);




type Commit struct {
    sHash string
    sAuthor string
    sDate string
    sMessage string
    aDiffs []Diff
}

type Diff struct {
    sFile string
    bFileCreated bool
    bFileDeleted bool
    sDiff string
}




func main () {
    
    aArgs := base.GetArgs(map[string]int{});
    //sCurrentDir, _ := os.Getwd();
    
    sGitLogFile := aArgs[""][0];
    
    aContent, err := ioutil.ReadFile(sGitLogFile);
    if err != nil {
        log.Fatal("could not read file " + sGitLogFile);
    }
    sContent := string(aContent);
    aLines := strings.Split(sContent, "\n");
    
    a
    
    for _, sLine := range aLines {
        base.Dump(sLine);
    }
    
}
