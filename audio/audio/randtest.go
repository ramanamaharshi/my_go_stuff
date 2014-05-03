

package main;


import(
    "os"
    "time"
    "olli/base"
    "math/rand"
    //"math"
);


func main () {
    
    aArgs := base.GetArgs(map[string]int{});
    sCurrentDir, _ := os.Getwd();
    
    base.Dump(aArgs);
    base.Dump(sCurrentDir);
    
    var oRand = rand.New(rand.NewSource(time.Now().Unix()));
    
    var iRandomIntRange = 20;
    var aRandomInts = make([]int, 999);
    var aInts = map[int]int{};
    for iI := -iRandomIntRange; iI <= iRandomIntRange; iI ++ {
        aInts[iI] = 0;
    }
base.Dump(aInts);
    
    for iKey := range aRandomInts {
        var iRandomInt = int(oRand.NormFloat64() * float64(iRandomIntRange / 2));
        aRandomInts[iKey] = iRandomInt;
        if _, bSet := aInts[iRandomInt]; bSet {
            aInts[iRandomInt] ++;
        }
    }
    
    for iI := -iRandomIntRange; iI <= iRandomIntRange; iI ++ {
        var iCount = aInts[iI];
        var sBar = "";
        for iC := 0; iC < iCount; iC ++ {
            sBar += "#";
        }
        base.Dump(sBar);
    }
    
}

