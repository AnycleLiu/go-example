package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

const (
	NUMS     = 10000000
	SORT_SEG = 20 * 1024 * 1024 //20M
)

func genBigFile(fp string) {
	f, err := os.Create(fp)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Printf("生成%d个uint64整数，共%f M\n", NUMS, 8*NUMS/1024.0/1024.0)

	for i := 0; i < NUMS; i++ {
		//binary.Write(f, binary.BigEndian, rand.Uint64())
		f.WriteString(strconv.Itoa(int(rand.Int63())))
		f.WriteString("\n")
	}

	info, err := f.Stat()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s 文件大小: %f M\n", fp, float64(info.Size()/1024.0/1024.0))
}

func sort_file(fp string) string {
	segdir := "./segs"
	spit(fp, segdir)

	/*fs, err := ioutil.ReadDir(segdir)
	if err != nil {
		panic(err)
	}
	*/

	return fp
}

func spit(fp, segdir string) {
	f, err := os.Open(fp)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		panic(err)
	}
	fmt.Printf(" 待排序文件:%s, 大小: %f M\n", fp, float64(info.Size()/1024.0/1024.0))
	segnum := int(math.Ceil(float64(info.Size()) / SORT_SEG))
	fmt.Printf("文件拆分成大小为 %f M的小文件，共 %d个", float64(SORT_SEG/1024.0/1024.0), segnum)

	os.RemoveAll(segdir)
	err = os.Mkdir(segdir, 777)
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(f)
	sfidx := 0
	var sf *os.File

	var ns []uint64 = make([]uint64, 0)
	var rn int

	savesf := func() {
		sf, err = os.Create(fmt.Sprintf("./segs/f_%d", sfidx))
		if err != nil {
			panic(err)
		}
		defer sf.Close()

		sort.Slice(ns, func(i, j int) bool {
			return ns[i] < ns[j]
		})
		for _, n := range ns {
			sf.WriteString(strconv.FormatUint(n, 10))
			sf.WriteString("\n")
		}
		ns = ns[0:0]
		sfidx++
		rn = 0
	}

	for {
		if rn >= SORT_SEG {
			savesf()
		}
		buf, err := r.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		rn += len(buf) //已读字节数
		ns = append(ns, binary.LittleEndian.Uint64(buf))
	}

	if rn > 0 {
		savesf()
	}
	fmt.Println("文件拆分完毕")
}

func main() {
	newsource := flag.Bool("newsource", false, "是否需要生成待排序文件?")
	source := flag.String("source", "./source.data", "要排序的文件路径")
	help := flag.Bool("h", false, "print usages")

	flag.Parse()

	if *help {
		flag.PrintDefaults()
	}

	//fmt.Println(*newsource, *source)

	if *newsource {
		genBigFile(*source)
	}

	sort_file(*source)
}
