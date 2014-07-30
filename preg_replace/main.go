
package main;


import(
	"os"
	"log"
	"regexp"
	"io/ioutil"
	"olli/base"
);


func main () {
	
	aArgs := base.GetArgs(map[string]int{"?": 0});
	_, bHelp := aArgs["?"];
	bHelp = bHelp || len(aArgs[""]) < 2;
	if bHelp {
		base.Dump("usage: \npreg_replace [find] [replace] < [file]\nor\npreg_replace [find] [replace] [text]");
		return;
	}
	sRegExp := aArgs[""][0];
	oRegExp := regexp.MustCompile(sRegExp);
	sReplace := aArgs[""][1];
	
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
	
	sReplaced := oRegExp.ReplaceAllString(sInput, sReplace);
	base.Dump(sReplaced);
	
}
