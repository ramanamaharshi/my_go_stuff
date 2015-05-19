



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
	
	log.Println("");
	fmt.Printf("");
	base.Dump("");
	
}




func main () {
	
	aArgs := base.GetArgs(map[string]int{"?": 0, "h": 0});
	
	sRegExp := aArgs[""][0];
	
	iOnlyMatchNr := -1;
	if len(aArgs[""]) > 1 {
		iOnlyMatchNr64, _ := strconv.ParseInt(aArgs[""][1], 10, 64);
		iOnlyMatchNr = int(iOnlyMatchNr64);
	}
	
	oRegExA := regexp.MustCompile(sRegExp);
	
	var iChunk int64 = 4;
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
			if iOnlyMatchNr == -1 {
				fmt.Print(strings.Join(aLastMatch, " ") + "\n");
			} else {
				fmt.Print(aLastMatch[iOnlyMatchNr] + "\n");
			}
		}
	}
	
}



