package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

// сюда писать код

// process Нужен для того чтобы, сравнивать процесс работы своей реализации с верной.
var process = `
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

// Для запуска тестов использовать: go test -v -race

// crcHash Структура, которая используется для восстановления порядка посчитанных данных
// хэшей внутри пайплайна, за счет поля Id можно все данные отсортировать, и работать с ними,
// в нужном порядке.
type crcHash struct {
	Id   int
	Hash string
}

// ExecutePipeline Функция имплементирующая работу Unix pipeline, только в языке Go.
func ExecutePipeline(j ...job) {

	in := make(chan interface{})
	out := make(chan interface{})

	for _, job := range j {
		go callJob(job, in, out)
		in = out
		out = make(chan interface{})
	}

	// Нужно для того, чтобы функция не завершила свою работу, до окончания работы вызванных ею горутин.
	for _ = range in {

	}

}

// TODO: Написать комментарии к решению, а также сделать код понятным и читаемым.

func callJob(j job, in, out chan interface{}) {
	defer close(out)
	j(in, out)
}

// SingleHash Считает значение crc32(data)+"~"+crc32(md5(data)) ( конкатенация двух строк через ~),
// где data - то что пришло на вход (по сути - числа из первой функции).
func SingleHash(in, out chan interface{}) {
	wg := new(sync.WaitGroup)

	for i := range in {
		// Преобразование значения в тип string
		data := i.(int)
		dataStr := strconv.Itoa(data)

		// Получение md5 хэша
		md := DataSignerMd5(dataStr)

		wg.Add(1)
		go concurenceSingleHash(md, dataStr, out, wg)
	}

	wg.Wait()
}

// Позволяет конкурентно вычислять участок работы функции SingleHash
func concurenceSingleHash(md, dataStr string, out chan interface{}, wg *sync.WaitGroup) {
	results := make([]crcHash, 0, 2)
	calcs := make(chan crcHash, 2)

	// Получение crc32(md5) хэша
	wg2 := new(sync.WaitGroup)
	wg2.Add(1)
	// Получение crc32 хэша
	wg2.Add(1)
	go concurenceCrc(dataStr, 0, calcs, wg2)

	go concurenceCrc(md, 1, calcs, wg2)

	wg2.Wait()

	close(calcs)
	for v := range calcs {
		results = append(results, v)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Id < results[j].Id
	})

	// Конкатенация crc и crc(md5) хэшей
	res := conStrs("~", results[0].Hash, results[1].Hash)

	// Запись результата в канал out
	out <- res

	wg.Done()
}

func concurenceCrc(data string, i int, calcs chan crcHash, wg *sync.WaitGroup) {
	crc := DataSignerCrc32(data)
	hash := crcHash{Id: i, Hash: crc}

	calcs <- hash
	wg.Done()
}

// MultiHash Считает значение crc32(th+data)) (конкатенация цифры, приведённой к строке и строки),
// где th=0..5 (т.е. 6 хешей на каждое входящее значение), потом берёт конкатенацию результатов в порядке расчета (0..5),
// где data - то что пришло на вход (и ушло на выход из SingleHash).
func MultiHash(in, out chan interface{}) {
	wg := new(sync.WaitGroup)

	for i := range in {
		wg.Add(1)
		data := i.(string)
		go calculateMultiHash(data, out, wg)
	}

	wg.Wait()
}

// calculateMultiHash Создан для конкурентного вызова DataSingerCrc32 внутри MultiHash
func calculateMultiHash(data string, out chan interface{}, wg *sync.WaitGroup) {
	wg2 := new(sync.WaitGroup)
	res := ""
	results := make([]crcHash, 0, 6)
	hashs := make(chan crcHash, 6)

	for i := 0; i < 6; i++ {
		wg2.Add(1)
		go calculateHash(i, data, hashs, wg2)
	}
	wg2.Wait()

	close(hashs)
	for v := range hashs {
		results = append(results, v)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Id < results[j].Id
	})

	for _, v := range results {
		res = conStrs("", res, v.Hash)
	}

	out <- res
	wg.Done()
}

// calculateHash Считает CRC32 Хэш, нужна для распаралеливания рассчетов хэша.
func calculateHash(i int, data string, hashs chan crcHash, wg2 *sync.WaitGroup) {
	crc := DataSignerCrc32(strconv.Itoa(i) + data)
	hashs <- crcHash{Id: i, Hash: crc}

	wg2.Done()
}

// CombineResults Получает все результаты, сортирует (https://golang.org/pkg/sort/),
// объединяет отсортированный результат через _ (символ подчеркивания) в одну строку.
func CombineResults(in, out chan interface{}) {
	preResult := make([]string, 0, MaxInputDataLen)

	for i := range in {
		data := i.(string)

		preResult = append(preResult, data)
	}

	sort.Strings(preResult)
	result := ""
	for i, s := range preResult {
		if i == 0 {
			result = s
		} else {
			result = conStrs("_", result, s)
		}
	}

	out <- result
}

// conStrs складывает строки переданную функцию и возвращает одну единственную строку как результат
func conStrs(sep string, strs ...string) string {
	var builder strings.Builder

	for i, str := range strs {
		builder.WriteString(str)
		if i != len(strs)-1 {
			builder.WriteString(sep)
		}
	}

	return builder.String()
}
