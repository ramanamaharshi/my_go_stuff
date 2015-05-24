



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
	
	iRead := 0;
	iCount := 0;
	sPart := "";
	aNewLetter := make([]byte, 1);
	var e error;
	
	oIn := bufio.NewReader(os.Stdin);
	
	for {
		
		iRead, e = oIn.Read(aNewLetter);
		
		sNewLetter := fmt.Sprintf("%s", aNewLetter[:iRead])
		
		aMatches := oRegExA.FindAllStringSubmatch(sPart, -1);
		aMatchIndices := oRegExA.FindAllStringSubmatchIndex(sPart, -1);
		
		if len(aMatches) != 0 {
			
			aLastMatch := aMatches[len(aMatches) - 1];
			aLastMatchIndices := aMatchIndices[len(aMatchIndices) - 1];
			
			aMatchIndicesNext := oRegExA.FindAllStringSubmatchIndex(sPart + sNewLetter, -1);
			
			bOutput := false;
			
			if e != nil {
				
				bOutput = true;
				
			} else {
				
				if len(aMatchIndicesNext) > 0 {
					
					aLastMatchIndicesNext := aMatchIndicesNext[len(aMatchIndicesNext) - 1];
					
					if aLastMatchIndices[0] != aLastMatchIndicesNext[0] || aLastMatchIndices[1] == aLastMatchIndices[1]{
						
						bOutput = true;
						
					}
					
				}
				
			}
			
			if (bOutput) {
				sOutput := sOutputTemplate;
				for iNr , sValue := range aLastMatch {
					sOutput = strings.Replace(sOutput, "\\" + strconv.Itoa(iNr), sValue, -1);
				}
				if bCount {
					iCount ++;
				} else if bReplace {
					fmt.Print(oRegExA.ReplaceAllLiteralString(sPart, sOutput));
				} else {
					fmt.Print(sOutput + "\n");
				}
				sPart = sPart[(aLastMatchIndices[1]):];
			}
			
		}
		
		if e != nil { break; }
		
		sPart += sNewLetter;
		//base.Dump(sPart);
		
		/*
		aLastMatch := []string{};
		aLastMatchIndices := []int{};
		bOutputPrevious := false;
		
		if e == nil {
			
			sPart += sChunk;
			
			aMatches := oRegExA.FindAllStringSubmatch(sPart, -1);
			aMatchIndices := oRegExA.FindAllStringSubmatchIndex(sPart, -1);
			
			if len(aMatchIndices) == 0 {
				bOutputPrevious = true;
			} else {
				aLastMatch = aMatches[len(aMatches) - 1];
				aLastMatchIndices = aMatchIndices[len(aMatchIndices) - 1];
				if -1 != aPreviousLastMatchIndices[0] && aLastMatchIndices[0] != aPreviousLastMatchIndices[0] {
					bOutputPrevious = true;
				}
			}
			
		} else {
			
			bOutputPrevious = true;
			
		}
		
		//base.Dump(aLastMatch);
		//base.Dump(aLastMatchIndices);
		//base.Dump(bOutputPrevious);
		
		if (-1 != aPreviousLastMatchIndices[0] && bOutputPrevious) {
			aMatch := aPreviousLastMatch;
			sOutput := sOutputTemplate;
			for iNr , sValue := range aMatch {
				sOutput = strings.Replace(sOutput, "\\" + strconv.Itoa(iNr), sValue, -1);
			}
			if !bCount {
				if bReplace {
					fmt.Print(oRegExA.ReplaceAllLiteralString(sPart, sOutput));
				} else {
					fmt.Print(sOutput + "\n");
				}
			}
base.Dump(aPreviousLastMatchIndices[1]);
base.Dump(sPart);
			sPart = sPart[aPreviousLastMatchIndices[1]:];
			aPreviousLastMatchIndices[0] = -1;
			aPreviousLastMatch = []string{};
			iCount ++;
		}
		
		if (len(aLastMatch) > 0) {
			aPreviousLastMatchIndices = aLastMatchIndices;
			aPreviousLastMatch = aLastMatch;
		}
		
		if e != nil { break; }
		*/
		
	}
	
	if bReplace {
		fmt.Print(sPart);
	}
	
	if bCount {
		fmt.Println(iCount);
	}
	
}



