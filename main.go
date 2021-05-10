package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"sync"
)

// Программа должна работать на локальном компьютере и получать на вход путь до директории.
// Программа должна вывести в стандартный поток вывода список дублирующихся файлов,
// которые находятся как в директории, так и в поддиректориях директории,
// переданной через аргумент командной строки. Данная функция должна работать
// эффективно при помощи распараллеливания программы
// Программа должна принимать дополнительный ключ - возможность удаления обнаруженных
// дубликатов файлов после поиска. Дополнительно нужно придумать, как обезопасить
// пользователей от случайного удаления файлов. В качестве ключей желательно придерживаться
// общепринятых практик по использованию командных опций.

var (
	dir *string
	del *bool
)

func init() {
	dir = flag.String("dir", "", "testing directory")
	del = flag.Bool("delete", false, "testing directory")
}

type File struct {
	Path string
	Name string
	Hash string
	Size int64
}

type List struct {
	mx sync.Mutex
	m  map[int]File
}

func NewList() *List {
	return &List{
		m: make(map[int]File),
	}
}

func (l *List) Load(key int) File {
	l.mx.Lock()
	defer l.mx.Unlock()

	return l.m[key]
}

func (l *List) Store(name, path, hash string, size int64) {
	l.mx.Lock()
	l.m[len(l.m)+1] = File{
		Name: name,
		Path: path,
		Hash: hash,
		Size: size,
	}
	l.mx.Unlock()
}

func main() {
	List := NewList()
	flag.Parse()
	err := Walk(*dir, List)
	if err != nil {
		fmt.Printf("%+v", err)
	} else {
		f := Find(List)
		if !*del {
			if f != nil {
				for key, slc := range f {
					fmt.Println(key+1, "набор дубликатов:")
					for _, value := range slc {
						fmt.Printf("Файл: %+v, путь: %+v\n", List.m[value].Name, List.m[value].Path)
					}
				}
			} else {
				fmt.Println("Дубликаты не найдены")
			}
		} else {
			if f != nil {
				answ := make([]int, len(f))
				for key, slc := range f {

					fmt.Println(key+1, "набор дубликатов:")
					for k, value := range slc {
						fmt.Printf("№%d Файл: %+v, путь: %+v\n", k+1, List.m[value].Name, List.m[value].Path)
					}

					fmt.Println("Какой из этих фалов оставить? Введите номер:")
					fmt.Scanln(&answ[key])

					slc[answ[key]-1] = slc[len(slc)-1]
					slc[len(slc)-1] = 0
					slc = slc[:len(slc)-1]
					delete(List, slc)
				}
			} else {
				fmt.Println("Дубликаты не найдены")
			}
		}

	}

}

func Hash(path string) (string, int64, error) {
	h := md5.New()
	f, err := os.Open(path)
	defer f.Close()

	stat, _ := f.Stat()
	size := stat.Size()
	if err != nil {
		return "", size, fmt.Errorf("%+v", err)
	}

	_, err = io.Copy(h, f)

	if err != nil {
		return "", size, fmt.Errorf("%+v", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), size, nil
}

func Walk(path string, l *List) error {
	wg := sync.WaitGroup{}
	lst, err := ioutil.ReadDir(path)

	if err != nil {
		return fmt.Errorf("%+v", err)
	}

	for _, val := range lst {
		if val.IsDir() {
			name := path + `/` + val.Name()
			wg.Add(1)
			go func() {
				defer wg.Done()
				Walk(name, l)
			}()
		} else {
			name := path + `/` + val.Name()
			hash, size, err := Hash(name)
			if err != nil {
				log.Printf("%+v", err)
			}
			wg.Add(1)
			fileName := val.Name()
			go func() {
				defer wg.Done()
				l.Store(fileName, path, hash, size)
			}()

		}
	}
	wg.Wait()
	return nil
}

func Find(l *List) [][]int {

	p := make(PairList, len(l.m))

	i := 0
	for k, v := range l.m {
		p[i] = Pair{k, v.Hash}
		i++
	}

	sort.Sort(p)
	var slice []int
	var mslice [][]int
	next := ""
	slice = append(slice, p[0].Key)
	for i := 0; i < len(p)-1; i++ {
		next = p[i+1].Value

		if p[i].Value == next {
			slice = append(slice, p[i+1].Key)
			if i+1 == len(p)-1 {
				slice = append(slice, p[i].Key)
				mslice = append(mslice, slice)
			}
		} else {
			if len(slice) >= 2 {
				mslice = append(mslice, slice)
			}
			slice = nil
		}

	}
	return mslice
}

type Pair struct {
	Key   int
	Value string
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func delete(l *List, slc []int) {
	for _, value := range slc {
		os.Remove(l.m[value].Path + "/" + l.m[value].Name)
	}
}
