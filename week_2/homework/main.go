package main

import (
	"fmt"
	"sort"
	"strconv"
)

func main() {
	results := make([]string, 0, 100)

	for i := 0; i < 2; i++ {
		// имитация работы singleHash
		i := strconv.Itoa(i)
		fmt.Println("i is: ", i)
		data := i
		md := DataSignerMd5(data)
		fmt.Println(md)
		crcmd := DataSignerCrc32(md)
		fmt.Println(crcmd)
		crc := DataSignerCrc32(data)
		fmt.Println(crc)
		res1 := conStrs("~", crc, crcmd)
		fmt.Println(res1)

		// имитация работы Multihash
		res2 := ""
		for i := 0; i < 6; i++ {
			sData := strconv.Itoa(i) + res1
			scrc := DataSignerCrc32(sData)
			fmt.Println(scrc)
			res2 = conStrs("", res2, scrc)
		}
		fmt.Println(res2)
		results = append(results, res2)
	}

	sort.Strings(results)
	globalResult := ""
	for i, v := range results {
		if i == 0 {
			globalResult = v
		} else {
			globalResult = conStrs("_", globalResult, v)
		}
	}
	fmt.Println("===========================================")
	fmt.Println(globalResult)
	fmt.Println(globalResult == "29568666068035183841425683795340791879727309630931025356555_4958044192186797981418233587017209679042592862002427381542")
}

// process_ Нужен для того чтобы, сравнивать процесс работы своей реализации с верной.
var process_ = `
0 SingleHash data 0
0 SingleHash md5(data) cfcd208495d565ef66e7dff9f98764da
0 SingleHash crc32(md5(data)) 502633748
0 SingleHash crc32(data) 4108050209
0 SingleHash result 4108050209~502633748
4108050209~502633748 MultiHash: crc32(th+step1)) 0 2956866606
4108050209~502633748 MultiHash: crc32(th+step1)) 1 803518384
4108050209~502633748 MultiHash: crc32(th+step1)) 2 1425683795
4108050209~502633748 MultiHash: crc32(th+step1)) 3 3407918797
4108050209~502633748 MultiHash: crc32(th+step1)) 4 2730963093
4108050209~502633748 MultiHash: crc32(th+step1)) 5 1025356555
4108050209~502633748 MultiHash result: 29568666068035183841425683795340791879727309630931025356555

1 SingleHash data 1
1 SingleHash md5(data) c4ca4238a0b923820dcc509a6f75849b
1 SingleHash crc32(md5(data)) 709660146
1 SingleHash crc32(data) 2212294583
1 SingleHash result 2212294583~709660146
2212294583~709660146 MultiHash: crc32(th+step1)) 0 495804419
2212294583~709660146 MultiHash: crc32(th+step1)) 1 2186797981
2212294583~709660146 MultiHash: crc32(th+step1)) 2 4182335870
2212294583~709660146 MultiHash: crc32(th+step1)) 3 1720967904
2212294583~709660146 MultiHash: crc32(th+step1)) 4 259286200
2212294583~709660146 MultiHash: crc32(th+step1)) 5 2427381542
2212294583~709660146 MultiHash result: 4958044192186797981418233587017209679042592862002427381542

CombineResults 29568666068035183841425683795340791879727309630931025356555_4958044192186797981418233587017209679042592862002427381542
`
