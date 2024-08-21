package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Очень помогло то что в тесты добавил функцию отображающую самую первую разницу строк и индекс места в котором разница.
func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}

}

// dirTree Вызывает recursiveCall и добавляет к результату переход строки
func dirTree(output io.Writer, path string, print bool) error {
	var f = true
	defer fmt.Fprintf(output, `
`)
	return recursiveCall(output, path, "", print, f)
}

// recursiveCall Записывает в io.Writer Все пути под директорий указанного пути, а также выводит файлы
// с расширением и размером файла, если флаг print = true.
func recursiveCall(output io.Writer, path, offset string, print, firstTime bool) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	path = conStrs(string(filepath.Separator), path)
	dir = conStrs(dir, path)

	dir2, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	dir3 := make([]os.DirEntry, 0, len(dir2))
	for i := 0; i < len(dir2); i++ {
		if print {
			dir3 = append(dir3, dir2[i])
		} else {
			if dir2[i].IsDir() {
				dir3 = append(dir3, dir2[i])
			}
		}
	}

	for i, v := range dir3 {
		// Нужно для того, чтобы на первой строке не было лишнего отступа.
		if !firstTime {
			fmt.Fprintf(output, "\n")
		}
		firstTime = false
		if v.IsDir() && i != len(dir3)-1 {
			fmt.Fprintf(output, "%s├───%s", offset, v.Name())
			offset2 := offset + "│\t"
			pathC := path + string(filepath.Separator) + v.Name()
			recursiveCall(output, pathC, offset2, print, firstTime)
		} else if v.IsDir() && i == len(dir3)-1 {
			pathC := path + string(filepath.Separator) + v.Name()
			fmt.Fprintf(output, "%v└───%s", offset, v.Name())
			offset2 := offset + "\t"
			recursiveCall(output, pathC, offset2, print, firstTime)
		} else if i == len(dir3)-1 {
			size, err := getSize(v)
			if err != nil {
				return err
			}
			fmt.Fprintf(output, "%s└───%s %v", offset, v.Name(), size)
		} else {
			size, err := getSize(v)
			if err != nil {
				return err
			}
			fmt.Fprintf(output, "%s├───%s %v", offset, v.Name(), size)
		}
	}

	return nil
}

// conStrs складывает строки переданную функцию и возвращает одну единственную строку как результат
func conStrs(strs ...string) string {
	var builder strings.Builder

	for _, str := range strs {
		builder.WriteString(str)
	}

	return builder.String()
}

// getSize Получение строки в которой будет указан размер в байтах, или empty если файл пустой.
func getSize(v os.DirEntry) (string, error) {
	info, err := v.Info()
	if err != nil {
		return "", err
	}

	size := "(empty)"

	if info.Size() == 0 {
		return size, nil
	}

	size = fmt.Sprintf("(%db)", info.Size())

	return size, nil
}

// Символы ранжирования - ├───  │  └───
