
package main;


import(
	//"os"
	"log"
	"regexp"
	"strings"
	"io/ioutil"
	"olli/base"
);


func main () {
	
	aArgs := base.GetArgs(map[string]int{"?": 0, "r": 0});
	
	_,bHelp := aArgs["?"];
	bHelp = bHelp || len(aArgs[""]) < 1;
	if bHelp {
		base.Dump("usage: \nautoheader [files]");
		return;
	}
	
	_,bRecursive := aArgs["r"];
	
	aFiles := aArgs[""];
	
	rxC := regexp.MustCompile("\\.c$");
	aDirs := []string{};
	aFilesTemp := aFiles;
	aFiles = []string{};
	for _,sFile := range(aFilesTemp) {
		if !base.FileExists(sFile) {
			log.Fatal("error: file " + sFile + " don't exist");
		}
		if base.IsDirectory(sFile) {
			aDirs = append(aDirs, sFile);
		} else {
			if rxC.MatchString(sFile) {
				aFiles = append(aFiles, sFile);
			}
		}
	}
	for 0 < len(aDirs) {
		aDirsTemp := aDirs;
		aDirs = []string{};
		for _,sDir := range(aDirsTemp) {
			aDirFiles,_ := ioutil.ReadDir(sDir);
			for _,oFile := range(aDirFiles) {
				sFile := oFile.Name();
				if oFile.IsDir() {
					aDirs = append(aDirs, sDir + "/" + sFile);
				} else {
					if rxC.MatchString(sFile) {
						aFiles = append(aFiles, sDir + "/" + sFile);
					}
				}
			}
		}
		if !bRecursive {break;}
	}
	base.Dump(aFiles);
	
	rxF := regexp.MustCompile("\n[^\n]*\\S+ +(\\S+) +\\([^\\)\\;]*\\) {");
	rxAuto := regexp.MustCompile("\n///#///\n[\\s\\S]+\n///#///\n");
	aRXSFH := []*regexp.Regexp{/*regexp.MustCompile("^\n *static ")*/};
	aRXSFF := []*regexp.Regexp{regexp.MustCompile("^_")};
	for _,sFile := range(aFiles) {
		sFC := sFile;
		//sFB := rxC.ReplaceAllString(sFC, ".b");
		//sContentB := "";
		sFH := rxC.ReplaceAllString(sFC, ".h");
		aContentC, _ := ioutil.ReadFile(sFC);
		sContentC := string(aContentC);
		sDeclarations := "";
		aFunctionHeads := rxF.FindAllStringSubmatch(sContentC, -1);
		for _,aFunctionHead := range(aFunctionHeads) {
			sFunctionHead := aFunctionHead[0];
			sFunctionName := aFunctionHead[1];
			bDeclare := true;
			for _,rxSFH := range(aRXSFH) {
				if rxSFH.MatchString(sFunctionHead) {bDeclare = false;}
			}
			for _,rxSFF := range(aRXSFF) {
				if rxSFF.MatchString(sFunctionName) {bDeclare = false;}
			}
			if bDeclare {sDeclarations += strings.Replace(sFunctionHead, "{", ";", 1);}
		}
		sDeclarations += "\n";
		sWrite := "\n///#///\n" + sDeclarations + "\n///#///\n";
		sContentH := "";
		if base.FileExists(sFH) {
			aContentH, _ := ioutil.ReadFile(sFH);
			sContentH = string(aContentH);
		}
		if (rxAuto.MatchString(sContentH)) {
			sWrite = rxAuto.ReplaceAllString(sContentH, sWrite);
		} else {
			sWrite = sContentH + sWrite;
		}
		ioutil.WriteFile(sFH, []byte(sWrite), 0644);
	}
	
}
