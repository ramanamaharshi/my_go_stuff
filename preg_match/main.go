
package main;


import(
    "os"
    "log"
    "regexp"
    "strconv"
    "io/ioutil"
    "olli/base"
);


func main () {
    
    aArgs := base.GetArgs(map[string]int{"?": 0});
    _, bHelp := aArgs["?"];
    bHelp = bHelp || len(aArgs[""]) == 0;
    if bHelp {
        base.Dump("usage: \npreg_match <regular expression> (<match nr>) < <file>\nor\nnpreg_match <regular expression> (<match nr>) <text>");
        return;
    }
    sRegExp := aArgs[""][0];
    oRegExp := regexp.MustCompile(sRegExp);
    iNr := 0;
    if len(aArgs[""]) > 1 {
        iNr64, _ := strconv.ParseInt(aArgs[""][1], 0, 64);
        iNr = int(iNr64);
    }
    
    sInput := "";
    if len(aArgs[""]) > 2 {
		sInput = aArgs[""][2];
	} else {
	    aInput, err := ioutil.ReadAll(os.Stdin);
	    if err != nil {
	        log.Fatal(err);
	    }
	    sInput = string(aInput);
	}
    aMatches := oRegExp.FindStringSubmatch(sInput);
    if len(aMatches) > 0 {
        base.Dump(aMatches[iNr]);
    }
    
}
