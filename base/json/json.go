

package json;


import(
    "io"
    "log"
    "fmt"
    "regexp"
    "strings"
    "strconv"
    "crypto/md5"
    "encoding/hex"
    "olli/base"
"os"
"io/ioutil"
);


type JsonNode struct {
    Type string
    Value string
    ArrayChildren []*JsonNode
    ObjectChildren map[string]*JsonNode
}


func (oNode *JsonNode) GetB (sPath string) (bool/*, bool*/) {
    oNode = oNode.Get(sPath);
    var bValue = oNode.Value == "true";
    return bValue/* , oNode.Type == "boolean"*/;
}


func (oNode *JsonNode) GetI (sPath string) (int/*, bool*/) {
    oNode = oNode.Get(sPath);
    var iValue64 , _ = strconv.ParseInt(oNode.Value, 10, 0);
    var iValue = int(iValue64);
    return iValue/* , oNode.Type == "number"*/;
}


func (oNode *JsonNode) GetF (sPath string) (float64/*, bool*/) {
    oNode = oNode.Get(sPath);
    var nValue , _ = strconv.ParseFloat(oNode.Value, 64);
    return nValue/* , oNode.Type == "number"*/;
}


func (oNode *JsonNode) GetS (sPath string) (string/*, bool*/) {
    oNode = oNode.Get(sPath);
    var sValue = oNode.Value[1:len(oNode.Value) - 1];
    return sValue/* , oNode.Type == "string"*/;
}


func (oNode *JsonNode) GetA (sPath string) ([]*JsonNode/*, bool*/) {
    oNode = oNode.Get(sPath);
    var aValue = oNode.ArrayChildren;
    return aValue/* , oNode.Type == "array"*/;
}


func (oNode *JsonNode) GetO (sPath string) (map[string]*JsonNode/*, bool*/) {
    oNode = oNode.Get(sPath);
    var aValue = oNode.ObjectChildren;
    return aValue/* , oNode.Type == "object"*/;
}


func (oNode *JsonNode) Get (sPath string) *JsonNode {
    
    var aHashMap = aHashMap(sPath);
    var sPathHashed = sPath;
    for sHash, sString := range(aHashMap) {
        sPathHashed = strings.Replace(sPathHashed, sString, sHash, -1);
    }
    
    var aPath = strings.Split(sPathHashed, " -> ");
    var oCurrentNode = oNode;
    for _ , sPathPartHashed := range(aPath) {
        var sPathPart , bInHashMap = aHashMap[sPathPartHashed];
        if !bInHashMap {
            sPathPart = sPathPartHashed;
        }
        var bObject = oCurrentNode.Type == "object";
        sPathPart = "\"" + strings.Trim(sPathPart, "\"") + "\"";
        if bObject {
            var _ , bExists = oCurrentNode.ObjectChildren[sPathPart];
if !bExists {
    base.Dump(sPathPart + " does not exist");
    var aExistingKeys = []string{};
    for sKey := range(oCurrentNode.ObjectChildren) {
        aExistingKeys = append(aExistingKeys, sKey);
    }
    base.Dump(aExistingKeys);
    base.Dump(oCurrentNode.Value);
}
            if !bExists {
                oCurrentNode = nil;
                break;
            }
            oCurrentNode = oCurrentNode.ObjectChildren[sPathPart];
        } else {
            var iPathPart , _ = strconv.Atoi(sPathPart);
            var bExists = len(oCurrentNode.ArrayChildren) > iPathPart;
if !bExists {
    base.Dump(sPathPart + " does not exist");
    base.Dump("length is " + strconv.Itoa(len(oCurrentNode.ArrayChildren)));
}
            if !bExists {
                oCurrentNode = nil;
                break;
            }
            oCurrentNode = oCurrentNode.ArrayChildren[iPathPart];
        }
    }
    
    return oCurrentNode;
    
}


func Read (sJson string) *JsonNode {
    
    var aHashMap = aHashMap(sJson);
    var sJsonHashed = sJson;
    for sHash, sString := range(aHashMap) {
        sJsonHashed = strings.Replace(sJsonHashed, sString, sHash, -1);
    }
ioutil.WriteFile("test", []byte(sJson + "\n" + sJsonHashed), os.FileMode(0766));
    
    return oRecReadHashed(sJsonHashed, &aHashMap);
    
}


func aHashMap (sInput string) map[string]string {
    //var aStrings = regexp.MustCompile("\"(.*?[^\\\\]|)(\\\\\\\\)*\"").FindAllString(sInput, -1);
    var aStrings = []string{};
    var aInput = []rune(sInput);
    var bInString = false;
    var iBackSlashes = 0;
    var aStringTemp = []rune{};
    for _ , rRune := range(aInput) {
        if (rRune == '"') {
            if (iBackSlashes % 2 == 0) {
                if (bInString) {
                    aStringTemp = append(aStringTemp, rRune);
                    aStrings = append(aStrings, string(aStringTemp));
                    aStringTemp = []rune{};
                }
                bInString = !bInString;
            }
        }
        if bInString {
            aStringTemp = append(aStringTemp, rRune);
        }
        if (rRune == '\\') {
            iBackSlashes ++;
        } else {
            iBackSlashes = 0;
        }
    }
    var aHashMap = map[string]string{};
    for _ , sString := range(aStrings) {
        var sHash = "s" + sHash(sString) + "s";
        aHashMap[sHash] = sString;
//base.Dump(sString);
    }
    return aHashMap;
}


func oRecReadHashed (sHashedJson string, aHashMap *map[string]string) *JsonNode {
    
    var oReturn = &JsonNode{Type: "?", Value: sHashedJson};
    
    var aRunes = []rune(sHashedJson);
    var rFirst = aRunes[0];
    if rFirst == '{' {
        oReturn.Type = "object";
    } else if rFirst == '[' {
        oReturn.Type = "array";
    } else if rFirst == 's' {
        oReturn.Type = "string";
        oReturn.Value = (*aHashMap)[oReturn.Value];
    } else if bDigit(rFirst) {
        oReturn.Type = "number";
    } else if sHashedJson == "true" || sHashedJson == "false" {
        oReturn.Type = "boolean";
    } else if sHashedJson == "null" {
        oReturn.Type = "null";
    } else if sHashedJson == "undefined" {
        oReturn.Type = "undefined";
    }
    
    if oReturn.Type == "object" || oReturn.Type == "array" {
        
        var aParts = []string{};
        var aPart = []rune{};
        var iOpenObjects = 0;
        var iOpenArrays = 0;
        var iLast = len(aRunes) -1;
        for iPos , rRune := range(aRunes) {
            if iPos != 0 && iPos != iLast{
                if rRune == '[' {
                    iOpenArrays ++;
                } else if rRune == ']' {
                    iOpenArrays --;
                } else if rRune == '{' {
                    iOpenObjects ++;
                } else if rRune == '}' {
                    iOpenObjects --;
                }
                var bSeperator = iOpenArrays == 0 && iOpenObjects == 0 && rRune == ',';
                var bEnd = iPos == iLast - 1;
                if bSeperator || bEnd {
                    aParts = append(aParts, string(aPart));
                    aPart = []rune{};
                } else {
                    aPart = append(aPart, rRune);
                }
            }
        }
        
        if oReturn.Type == "array" {
            oReturn.ArrayChildren = []*JsonNode{};
            for _ , sPart := range(aParts) {
                var sValue = sPart;
                var oChild = oRecReadHashed(sValue, aHashMap);
                oReturn.ArrayChildren = append(oReturn.ArrayChildren, oChild);
            }
        }
        
        if oReturn.Type == "object" {
            oReturn.ObjectChildren = map[string]*JsonNode{};
            for _ , sPart := range(aParts) {
                var aSplitPart = strings.SplitAfterN(sPart, ":", 2);
                var aSplitKey = strings.Split(aSplitPart[0], ":");
                var sKey = aSplitKey[0];
                sKey = (*aHashMap)[sKey];
                var sValue = aSplitPart[1];
                var oChild = oRecReadHashed(sValue, aHashMap);
                oReturn.ObjectChildren[sKey] = oChild;
            }
        }
        
    }
        
    return oReturn;
    
}


func bDigit (rInput rune) bool {
    return rInput == '0' || rInput == '1' || rInput == '2' || rInput == '3' || rInput == '4' || rInput == '5' || rInput == '6' || rInput == '7' || rInput == '8' || rInput == '9';
}


func sHash (sInput string) string {
    var oHash = md5.New();
    io.WriteString(oHash, sInput);
    var sHash16 = hex.EncodeToString(oHash.Sum(nil));
    return sHash16;
}


func NeverUsed () {
    var _, _ = strconv.ParseInt("", 0, 0);
    var _ = regexp.QuoteMeta("");
    base.Dump("");
    fmt.Println("");
    log.Fatal("");
}

