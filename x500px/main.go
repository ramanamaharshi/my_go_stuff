

package main;


import(
    "os"
    "io"
    "log"
    "fmt"
    "time"
    "sync"
    "regexp"
    "strconv"
    "strings"
    "net/url"
    "net/http"
    "io/ioutil"
    "olli/base"
    "olli/base/json"
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
    var sDomainContent = oClient.sGet("http://500px.com/");
    var aToken = regexp.MustCompile("\\<meta content=\"([^\"]+)\" name=\"csrf-token\" /\\>").FindStringSubmatch(sDomainContent);
    if aToken == nil {log.Fatal("could not find csrf-token");}
    oClient.sToken = aToken[1];
    
}

func (oClient *Client) sGet (sUrl string) string {
    
    if oClient.oHttpClient == nil {oClient.vInit();}
    
    var oResponse, errA = oClient.oHttpClient.Get(sUrl);
    if errA != nil {log.Fatal(errA);}
    var aBody, errB = ioutil.ReadAll(oResponse.Body);
    if errB != nil {log.Fatal(errB);}
    defer oResponse.Body.Close();
    var sResponse = string(aBody);
    
    return sResponse;
    
}

func (oClient *Client) aGetUserPhotos (sUserName string) []Photo {
    
    if oClient.oHttpClient == nil {oClient.vInit();}
    
    var sUserPageContent = oClient.sGet("http://500px.com/" + sUserName)
    var aUserID = regexp.MustCompile("data-user-id=\"(\\d+)\"").FindStringSubmatch(sUserPageContent);
    if aUserID == nil {log.Fatal("could not find data-user-id");}
    var sUserID = aUserID[1];
    
    var iPageLength = 8;
    
    var sApiCall = "https://api.500px.com/v1/photos";
    sApiCall += "?feature=user&sort=created_at&include_states=false";
    sApiCall += "&user_id=" + sUserID + "&authenticity_token=" + url.QueryEscape(oClient.sToken);
    sApiCall += "&image_size=5&rpp=" + strconv.Itoa(iPageLength) + "&page=1";
    
    var sFirstPageContent = oClient.sGet(sApiCall);
//base.Dump(sFirstPageContent);
    var oFirstPageContent = json.Read(sFirstPageContent);
base.Dump(oFirstPageContent.GetI("total_pages"));
base.Dump(oFirstPageContent.GetI("total_items"));
    var iPages = oFirstPageContent.GetI("total_pages");
    var iPhotos = oFirstPageContent.GetI("total_items");
base.Dump(iPhotos);
    
    var aPhotos = []Photo{};
    for iPage := 1; iPage <= iPages; iPage ++ {
        sApiCall = "https://api.500px.com/v1/photos";
        sApiCall += "?feature=user&sort=created_at&include_states=false";
        sApiCall += "&user_id=" + sUserID + "&authenticity_token=" + url.QueryEscape(oClient.sToken);
        sApiCall += "&image_size=5&rpp=" + strconv.Itoa(iPageLength) + "&page=" + strconv.Itoa(iPage);
        var sData = oClient.sGet(sApiCall);
        var oData = json.Read(sData);
        var aDataPhotos = oData.GetA("photos");
        for _ , oDataPhoto := range(aDataPhotos) {
//base.Dump(len(aPhotos));
            var oPhoto = Photo{
                sID: strconv.Itoa(oDataPhoto.GetI("id")),
                sUrl: oDataPhoto.GetS("image_url"),
                sUserID: strconv.Itoa(oDataPhoto.GetI("user_id")),
                sUserName: oDataPhoto.GetS("user -> username"),
                oTaken: oParseTime(oDataPhoto.GetS("taken_at")),
                oCreated: oParseTime(oDataPhoto.GetS("created_at")),
            };
            aPhotos = append(aPhotos, oPhoto);
        }
    }
base.Dump(len(aPhotos));
    
    return aPhotos;
    
}




type Photo struct {
    sID string
    sUrl string
    sUserID string
    sUserName string
    oTaken time.Time
    oCreated time.Time
}




func main () {
    
    var aArgs = base.GetArgs(map[string]int{"tel": 0, "mp3": 0, "gain": 0, "user_uploads": 1, "format": 1, "q": 1});
    var sCurrentDir, _ = os.Getwd();
    
    var sUserName = aArgs[""][0];
    fmt.Println(sUserName);
    
//var sJson = "{\"current_page\":12,\"total_pages\":68,\"total_items\":135,\"photos\":[{\"id\":45898264,\"user_id\":1171535,\"name\":\"following the rainbow\",\"description\":\"\",\"camera\":\"Canon EOS 5D Mark II\",\"lens\":\"\",\"focal_length\":\"135\",\"iso\":\"100\",\"shutter_speed\":\"1/320\",\"aperture\":\"3.2\",\"times_viewed\":33776,\"rating\":50.0,\"status\":1,\"created_at\":\"2013-09-10T04:09:20-04:00\",\"category\":0,\"location\":null,\"latitude\":null,\"longitude\":null,\"taken_at\":\"2013-09-03T19:55:30-04:00\",\"hi_res_uploaded\":0,\"for_sale\":false,\"width\":1549,\"height\":1894,\"votes_count\":2440,\"favorites_count\":1656,\"comments_count\":168,\"nsfw\":false,\"sales_count\":0,\"for_sale_date\":null,\"highest_rating\":98.8,\"highest_rating_date\":\"2013-09-11T04:45:22-04:00\",\"license_type\":0,\"converted\":27,\"collections_count\":0,\"privacy\":false,\"image_url\":\"http://ppcdn.500px.org/45898264/4b3078bd16cbebc01151b7c37d19d34b02ffba11/5.jpg\",\"images\":[{\"size\":5,\"url\":\"http://ppcdn.500px.org/45898264/4b3078bd16cbebc01151b7c37d19d34b02ffba11/5.jpg\",\"https_url\":\"https://gp1.wac.edgecastcdn.net/806614/photos/photos.500px.net/45898264/4b3078bd16cbebc01151b7c37d19d34b02ffba11/5.jpg\"}],\"positive_votes_count\":2440,\"user\":{\"id\":1171535,\"username\":\"ElenaShumilova\",\"firstname\":\"Elena\",\"lastname\":\"Shumilova\",\"city\":\"Andreapol\",\"country\":\"Russia\",\"fullname\":\"Elena Shumilova\",\"userpic_url\":\"http://pacdn.500px.org/1171535/a54cf3d035dad57d6fdbb82b14f625c16f30b5dc/1.jpg?127\",\"userpic_https_url\":\"https://gp1.wac.edgecastcdn.net/806614/avatars/avatars.500px.net/1171535/a54cf3d035dad57d6fdbb82b14f625c16f30b5dc/1.jpg?127\",\"upgrade_status\":0,\"affection\":151952,\"followers_count\":30519}},{\"id\":45200950,\"user_id\":1171535,\"name\":\"Untitled\",\"description\":null,\"camera\":\"Canon EOS 5D Mark II\",\"lens\":null,\"focal_length\":\"135\",\"iso\":\"100\",\"shutter_speed\":\"1/800\",\"aperture\":\"4\",\"times_viewed\":12717,\"rating\":49.8,\"status\":1,\"created_at\":\"2013-09-04T14:15:56-04:00\",\"category\":0,\"location\":null,\"latitude\":null,\"longitude\":null,\"taken_at\":\"2013-05-08T21:30:00-04:00\",\"hi_res_uploaded\":0,\"for_sale\":false,\"width\":1479,\"height\":1026,\"votes_count\":957,\"favorites_count\":539,\"comments_count\":36,\"nsfw\":false,\"sales_count\":0,\"for_sale_date\":null,\"highest_rating\":94.7,\"highest_rating_date\":\"2013-09-05T14:03:50-04:00\",\"license_type\":0,\"converted\":27,\"collections_count\":0,\"privacy\":false,\"image_url\":\"http://ppcdn.500px.org/45200950/75b4becc573dd1864bca7398deb1b208b100a32d/5.jpg\",\"images\":[{\"size\":5,\"url\":\"http://ppcdn.500px.org/45200950/75b4becc573dd1864bca7398deb1b208b100a32d/5.jpg\",\"https_url\":\"https://gp1.wac.edgecastcdn.net/806614/photos/photos.500px.net/45200950/75b4becc573dd1864bca7398deb1b208b100a32d/5.jpg\"}],\"positive_votes_count\":957,\"user\":{\"id\":1171535,\"username\":\"ElenaShumilova\",\"firstname\":\"Elena\",\"lastname\":\"Shumilova\",\"city\":\"Andreapol\",\"country\":\"Russia\",\"fullname\":\"Elena Shumilova\",\"userpic_url\":\"http://pacdn.500px.org/1171535/a54cf3d035dad57d6fdbb82b14f625c16f30b5dc/1.jpg?127\",\"userpic_https_url\":\"https://gp1.wac.edgecastcdn.net/806614/avatars/avatars.500px.net/1171535/a54cf3d035dad57d6fdbb82b14f625c16f30b5dc/1.jpg?127\",\"upgrade_status\":0,\"affection\":151952,\"followers_count\":30519}}],\"filters\":{\"category\":false,\"exclude\":false,\"user_id\":1171535},\"feature\":\"user\"}"
//json.Read(sJson);
//log.Fatal("---");
    
    var oClient = &Client{};
    var aPhotos = oClient.aGetUserPhotos(sUserName);
    vDownloadPhotos(aPhotos, sCurrentDir);
    
}


func vDownloadPhotos (aPhotos []Photo, sTargetDir string) {
    os.Mkdir(sTargetDir, os.FileMode(0766));
    for iNr , oPhoto := range(aPhotos) {
        var aUrlSplit = strings.Split(oPhoto.sUrl, ".");
        var sExtension = aUrlSplit[len(aUrlSplit) - 1];
        var sFileName = "500px_";
        sFileName += oPhoto.sUserName + "_" + oPhoto.sUserID;
        sFileName += "_" + sDateTime(oPhoto.oCreated);
        sFileName += "_" + oPhoto.sID + "." + sExtension;
        var sPhotoTargetDir = sTargetDir + "/" + oPhoto.sUserName;
        var sTarget = sPhotoTargetDir + "/" + sFileName;
        os.Mkdir(sPhotoTargetDir, os.FileMode(0776));
        vDownload(oPhoto.sUrl, sTarget);
base.Dump("downloaded " + strconv.Itoa(iNr + 1) + " / " + strconv.Itoa(len(aPhotos)));
    }
}


func sDateTime (oTime time.Time) string {
    var sReturn = "";
    sReturn += fmt.Sprintf("%#04d", oTime.Year())[2:4];
    sReturn += fmt.Sprintf("%#02d", int(oTime.Month()));
    sReturn += fmt.Sprintf("%#02d", oTime.Day());
    sReturn += fmt.Sprintf("%#02d", oTime.Hour());
    sReturn += fmt.Sprintf("%#02d", oTime.Minute());
    sReturn += fmt.Sprintf("%#02d", oTime.Second());
    return sReturn;
}


func oParseTime (sTime string) time.Time {
    var sRegExp = "(\\d\\d\\d\\d)\\-(\\d\\d)\\-(\\d\\d)T(\\d\\d)\\:(\\d\\d)\\:(\\d\\d)\\-(\\d\\d)\\:(\\d\\d)";
    var aTime = regexp.MustCompile(sRegExp).FindStringSubmatch(sTime);
    var aT = make([]int, len(aTime), len(aTime));
    for iNr , sT := range(aTime) {
        aT[iNr] , _ = strconv.Atoi(sT);
    }
    var oTime = time.Time{};
    if len(aT) >= 7 {
        oTime = time.Date(aT[1], time.Month(aT[2] - 1), aT[3], aT[4], aT[5], aT[6], 0, time.UTC);
    }
    return oTime;
}


func vDownload (sUrl, sFile string) {
    var oFileOut, errA = os.Create(sFile);
    if errA != nil {log.Fatal(errA);}
    defer oFileOut.Close();
    var oUrlIn, errB = http.Get(sUrl);
    if errB != nil {log.Fatal(errB);}
    defer oUrlIn.Body.Close();
    var _, errC = io.Copy(oFileOut, oUrlIn.Body);
    if errC != nil {log.Fatal(errC);}
}


func iParseInt (sInput string) int {
    var iOutput, err = strconv.ParseInt(sInput, 10, 0);
    if err != nil {log.Fatal(err);}
    return int(iOutput);
}


func NeverUsed () {
    var _, _ = os.Getwd();
    var _ = regexp.QuoteMeta("");
    var _, _ = strconv.ParseInt("", 0, 0);
    fmt.Println("");
    log.Fatal("");
    base.Dump("");
}

