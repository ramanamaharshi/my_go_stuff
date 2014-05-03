

package main;


import(
    "os"
    //"io"
    "log"
    "fmt"
    "time"
    "sync"
    "regexp"
    "strconv"
    //"strings"
    "net/url"
    "net/http"
    "io/ioutil"
    "olli/base"
    "bytes"
    //"olli/base/json"
);




type InterfaceTime interface {
    Year() int
    Month() time.Month
    Day() int
    Hour() int
    Minute() int
    Second() int
}




type MyJar struct {
    lk      sync.Mutex
    cookies map[string][]*http.Cookie
}
func (jar *MyJar) SetCookies (u *url.URL, cookies []*http.Cookie) {
    jar.lk.Lock();
    jar.cookies[jar.sGetHost(u.Host)] = cookies;
//base.Dump(jar.cookies[jar.sGetHost(u.Host)]);
    jar.lk.Unlock();
}
func (jar *MyJar) Cookies (u *url.URL) []*http.Cookie {
    return jar.cookies[jar.sGetHost(u.Host)];
}
func (jar *MyJar) sGetHost (sHost string) string {
    var sMatch = regexp.MustCompile("[^\\.]+\\.[^\\.]+$").FindString(sHost);
    return sMatch;
}
func NewJar () *MyJar {
    jar := new(MyJar);
    jar.cookies = map[string][]*http.Cookie{};
    return jar;
}




type Client struct {
    sToken string
    oHttpClient *http.Client
}

func (oClient *Client) vInit () {
    
    oClient.oHttpClient = &http.Client{Jar: NewJar()};
    var sDomainContent = oClient.sGet("http://bestbookread.com/");
    var aToken = regexp.MustCompile("\\<meta content=\"([^\"]+)\" name=\"csrf-token\" /\\>").FindStringSubmatch(sDomainContent);
    if aToken == nil {base.Dump("could not find csrf-token");}
    if aToken != nil {
        oClient.sToken = aToken[1];
    }
//base.Dump(oClient.sToken);
    
}

func (oClient *Client) sGet (sUrl string) string {
    
    if oClient.oHttpClient == nil {oClient.vInit();}
    
    var oRequest, errA = http.NewRequest("GET", sUrl, nil);
    if errA != nil {base.Dump(errA); return"";}
    oRequest.Header.Add("X-CSRF-Token", oClient.sToken);
    var oResponse, errB = oClient.oHttpClient.Do(oRequest);
    if errB != nil {base.Dump(errB); return "";}
    var aBody, errC = ioutil.ReadAll(oResponse.Body);
    if errC != nil {base.Dump(errC); return "";}
    defer oResponse.Body.Close();
    var sResponse = string(aBody);
    
    return sResponse;
    
}

func (oClient *Client) sPost (sUrl string, aData url.Values) string {
    
    if oClient.oHttpClient == nil {oClient.vInit();}
    
    var oRequest, errA = http.NewRequest("POST", sUrl, bytes.NewBufferString(aData.Encode()));
    if errA != nil {base.Dump(errA); return"";}
    oRequest.Header.Add("X-CSRF-Token", oClient.sToken);
    var oResponse, errB = oClient.oHttpClient.Do(oRequest);
    if errB != nil {base.Dump(errB); return"";}
    var aBody, errC = ioutil.ReadAll(oResponse.Body);
    if errC != nil {base.Dump(errC); return"";}
    defer oResponse.Body.Close();
    var sResponse = string(aBody);
    
    return sResponse;
    
}




func main () {
    
    var aArgs = base.GetArgs(map[string]int{"tel": 0, "mp3": 0, "gain": 0, "user_uploads": 1, "format": 1, "q": 1});
    var sCurrentDir, _ = os.Getwd();
    _ = sCurrentDir;
    
    var sBookID = aArgs[""][0];
    var oStartTime = time.Now();
    var oPrevTime = time.Now();
    var iVotesPerLoop = 16;
    var iVotesCast = 0;
    for {
        for iV := 0; iV < iVotesPerLoop - 1; iV ++ {
            go vVote(sBookID);
        }
        vVote(sBookID);
        iVotesCast += iVotesPerLoop;
        var oNow = time.Now();
        var nDeltaSeconds = oNow.Sub(oPrevTime).Seconds();
        var nPerHour = float64(60 * 60) * (float64(iVotesPerLoop) / nDeltaSeconds);
        var sPerHour = strconv.FormatFloat(nPerHour, 'f', 4, 64);
        var nTotalDeltaSeconds = oNow.Sub(oStartTime).Seconds();
        var nTotalPerHour = float64(60 * 60) * (float64(iVotesCast) / nTotalDeltaSeconds);
        var sTotalPerHour = strconv.FormatFloat(nTotalPerHour, 'f', 4, 64);
        base.Dump(strconv.Itoa(iVotesCast) + " votes cast");
        base.Dump("spamming at " + sTotalPerHour + " (" + sPerHour + ") per Hour");
        oPrevTime = oNow;
    }
    
}


func vVote (sBookID string) {
    var oClient = &Client{};
    var aValues = url.Values{};
    aValues.Add("book[Item][ASIN]", sBookID);
    aValues.Add("book[Item][ItemAttributes][author]", "");
    var sReturn = "";
    sReturn += oClient.sPost("http://bestbookread.com/books/insert_or_update", aValues);
    sReturn += oClient.sGet("http://bestbookread.com/books/check_history?book=" + sBookID);
    base.Dump(sReturn);
    //return sReturn;
}


//func sDateTime (oTime time.Time) string {
//    var sReturn = "";
//    sReturn += fmt.Sprintf("%#04d", oTime.Year())[2:4];
//    sReturn += fmt.Sprintf("%#02d", int(oTime.Month()));
//    sReturn += fmt.Sprintf("%#02d", oTime.Day());
//    sReturn += fmt.Sprintf("%#02d", oTime.Hour());
//    sReturn += fmt.Sprintf("%#02d", oTime.Minute());
//    sReturn += fmt.Sprintf("%#02d", oTime.Second());
//    return sReturn;
//}


//func oParseTime (sTime string) time.Time {
//    var sRegExp = "(\\d\\d\\d\\d)\\-(\\d\\d)\\-(\\d\\d)T(\\d\\d)\\:(\\d\\d)\\:(\\d\\d)\\-(\\d\\d)\\:(\\d\\d)";
//    var aTime = regexp.MustCompile(sRegExp).FindStringSubmatch(sTime);
//    var aT = make([]int, len(aTime), len(aTime));
//    for iNr , sT := range(aTime) {
//        aT[iNr] , _ = strconv.Atoi(sT);
//    }
//    var oTime = time.Time{};
//    if len(aT) >= 7 {
//        oTime = time.Date(aT[1], time.Month(aT[2] - 1), aT[3], aT[4], aT[5], aT[6], 0, time.UTC);
//    }
//    return oTime;
//}


//func iParseInt (sInput string) int {
//    var iOutput, err = strconv.ParseInt(sInput, 10, 0);
//    if err != nil {log.Fatal(err);}
//    return int(iOutput);
//}


func NeverUsed () {
    var _, _ = os.Getwd();
    var _ = regexp.QuoteMeta("");
    var _, _ = strconv.ParseInt("", 0, 0);
    fmt.Println("");
    log.Fatal("");
    base.Dump("");
}

