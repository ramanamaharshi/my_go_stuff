



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
	
	aArgs := base.GetArgs(map[string]int{"?": 0, "h": 0});
	
	_, bHelp := aArgs["h"];
	_, bQuestionMark := aArgs["?"];
	if bHelp || bQuestionMark {
		fmt.Println("this tool is for grepping when there are no line breaks");
		fmt.Println("usage: streamo [regular expression] <[output per match (default='\\0')]>");
		fmt.Println("options:");
		fmt.Println("-?, -h: this help page");
	}
	
	sRegExp := aArgs[""][0];
	
	sOutputTemplate := "\\0";
	if len(aArgs[""]) > 1 {
		sOutputTemplate = aArgs[""][1];
	}
	
	oRegExA := regexp.MustCompile(sRegExp);
	
	var iChunk int64 = 1;
	aData := make([]byte, iChunk);
	var iRead int = 0;
	sAll := "";
	sPart := "";
	var err error;
	
	oIn := bufio.NewReader(os.Stdin);
	
	for {
		iRead, err = oIn.Read(aData);
		if err != nil {break}
		sAppend := fmt.Sprintf("%s", aData[:iRead]);
		sAll += sAppend;
		sPart += sAppend;
		aMatches := oRegExA.FindAllStringSubmatch(sPart, -1);
		if len(aMatches) > 0 {
			sPart = "";
			aLastMatch := aMatches[len(aMatches) - 1];
			sOutput := sOutputTemplate;
			for iNr , sValue := range aLastMatch {
				sOutput = strings.Replace(sOutput, "\\" + strconv.Itoa(iNr), sValue, -1);
			}
			fmt.Print(sOutput + "\n");
		}
	}
	
}



