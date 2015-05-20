



package main;




import "os";
import "log";
import "fmt";
import "bufio";
import "olli/base";
import "regexp";
import "strings";
import "strconv";




func Dummy () {
	
	_ = strconv.Itoa(0);
	_ = strings.Index("", "");
	log.Println("");
	fmt.Printf("");
	base.Dump("");
	
}




func main () {
	
	aArgs := base.GetArgs(map[string]int{"? h help": 0, "p plain": 0, "r replace": 0, "c count": 0});
	
	_ , bHelp := aArgs["? h help"];
	//_ , bPlain := aArgs["p plain"];
	_ , bReplace := aArgs["r replace"];
	_ , bCount := aArgs["c count"];
	
	if bHelp {
		fmt.Println("this tool is for grepping when there are no line breaks");
		fmt.Println("usage: streamo [regular expression] <[output template (default='\\0')]>");
		fmt.Println("options:");
		fmt.Println("-? -h --help: this help page");
		fmt.Println("-r --replace : output is input with matches replaced according to output template");
		fmt.Println("-c --count : only count matches, print count");
		fmt.Println("planned options:");
		fmt.Println("-p --plain : search for string instead of regular expression");
	}
	
	sRegExp := aArgs[""][0];
	oRegExA := regexp.MustCompile(sRegExp);
	
	sOutputTemplate := "\\0";
	if len(aArgs[""]) > 1 {
		sOutputTemplate = aArgs[""][1];
	}
	
	sPart := "";
	iChunkSize := int64(1);
	aChunk := make([]byte, iChunkSize);
	iCount := 0;
	iRead := 0;
	var e error;
	
	oIn := bufio.NewReader(os.Stdin);
	
	for {
		iRead, e = oIn.Read(aChunk);
		if e != nil {break}
		sAppend := fmt.Sprintf("%s", aChunk[:iRead]);
		sPart += sAppend;
		aMatches := oRegExA.FindAllStringSubmatch(sPart, -1);
		if len(aMatches) > 0 {
			iCount ++;
			aLastMatch := aMatches[len(aMatches) - 1];
			sOutput := sOutputTemplate;
			for iNr , sValue := range aLastMatch {
				sOutput = strings.Replace(sOutput, "\\" + strconv.Itoa(iNr), sValue, -1);
			}
			if !bCount {
				if bReplace {
					fmt.Print(oRegExA.ReplaceAllLiteralString(sPart, sOutput));
				} else {
					fmt.Print(sOutput + "\n");
				}
			}
			sPart = "";
		}
	}
	
	if bReplace {
		fmt.Print(sPart);
	}
	
	if bCount {
		fmt.Println(iCount);
	}
	
}




/*func aGetFlag (aArgs [string][]string, aFlagGroup []string) []string {
	
	aReturn := []string{};
	
	//for _ , sFlag in 
	
}*/



