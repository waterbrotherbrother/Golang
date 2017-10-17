package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Args struct {
	pragrom_name      string
	start_page        int
	end_page          int
	page_length       int
	page_type         bool //true :'-f'
	input_file        string
	print_destination string
}

func main() {
	var args Args

	argsProcess(&args)
	inputProcess(&args)

}

func argsProcess(args *Args) {

	args.pragrom_name = os.Args[0]

	flag.IntVar(&args.start_page, "s", -1, "specify start page.(s >= 0)(mandatory)")
	flag.IntVar(&args.end_page, "e", -1, "specify end page.(e >= s)(mandatory)")
	flag.IntVar(&args.page_length, "l", 72, "specify page length(number of lines).")
	flag.BoolVar(&args.page_type, "f", false, "specify type of input file.")
	flag.StringVar(&args.print_destination, "d", "", "specify the destination printer.")

	flag.Usage = usage
	flag.Parse()

	if args.start_page < 1 || args.end_page < 1 || args.start_page > args.end_page {
		fmt.Fprintln(os.Stderr, "Error: -s or -e invaild\n")
		flag.Usage()
	}

	if args.page_type == true && args.page_length != 72 {
		fmt.Fprintln(os.Stderr, "Error: -l and -f conflict\n")
		flag.Usage()
	}

	if len(flag.Args()) > 1 {
		fmt.Fprintln(os.Stderr, "Error:only allow one non flag argument\n")
		flag.Usage()
	}

	if len(flag.Args()) == 1 {
		args.input_file = flag.Args()[0]
	}

}

func inputProcess(args *Args) {

	if args.input_file == "" { //标准输入
		reader := bufio.NewReader(os.Stdin)

		if args.page_type == true {
			readByPage(reader, args)
		} else {
			readByLine(reader, args)
		}
	} else { //文件输入
		inputFile, err := os.Open(args.input_file)
		check(err)
		reader := bufio.NewReader(inputFile)
		defer inputFile.Close()

		if args.page_type == true {
			readByPage(reader, args)
		} else {
			readByLine(reader, args)
		}
	}
}

func readByPage(reader *bufio.Reader, args *Args) {

	currentPage := 1 // 纪录当前页数

	for {
		page, err := reader.ReadString('\f')
		check(err)
		// 当页数在所要选取的范围时
		if currentPage >= args.start_page && currentPage <= args.end_page {
			if args.print_destination == "" {
				fmt.Printf(page)
			} else {
				cmd := exec.Command("./tmp")     // 创建命令"./tmp"
				stdin, err := cmd.StdinPipe()    // 打开./out的标准输入管道
				check(err)                       // 错误检测
				stdin.Write([]byte(page + "\n")) // 向管道中写入文本
				stdin.Close()                    // 关闭管道
				cmd.Stdout = os.Stdout           // ./tmp将会输出到屏幕
				cmd.Run()                        // 运行./tmp命令
			}
		}
		if err == io.EOF {
			break
		}
		currentPage++
	}

	if args.start_page > currentPage {
		fmt.Printf("Warning : start_page:%d 超过总页数:%d\n ", args.start_page, currentPage)
		os.Exit(1)
	}

	if args.end_page > currentPage {
		fmt.Printf("Warning : end_page:%d 超过总页数:%d\n ", args.end_page, currentPage)
		os.Exit(1)
	}
}

func readByLine(reader *bufio.Reader, args *Args) {
	currentLine := 1 //纪录当前行数

	for {
		line, err := reader.ReadString('\n')
		check(err)
		virtualPage_length := (currentLine-1)/args.page_length + 1 //表示以当前page_length的条件下，页数
		if (virtualPage_length >= args.start_page) && (virtualPage_length <= args.end_page) {
			if args.print_destination == "" {
				fmt.Printf(line)
			} else {
				cmd := exec.Command("./tmp")
				stdin, err := cmd.StdinPipe()
				check(err)
				stdin.Write([]byte(line))
				stdin.Close()
				cmd.Stdout = os.Stdout
				cmd.Run()
			}
		}
		if err == io.EOF {
			break
		}
		currentLine++
	}
	if args.start_page > currentLine/args.page_length+1 {
		fmt.Printf("\nWarning : start_page:%d 超过总页数:%d\n ", args.start_page, currentLine/args.page_length+1)
		os.Exit(1)
	}
	if args.end_page > currentLine/args.page_length+1 {
		fmt.Printf("\nWarning : end_page:%d 超过总页数:%d\n ", args.end_page, currentLine/args.page_length+1)
		os.Exit(1)
	}
}

func check(err error) {
	if err != nil && err != io.EOF {
		panic(err)
	}
}

func usage() {
	fmt.Printf("usage: selpg [flags] [filename]\n")
	fmt.Printf("Arguments are:\n")
	fmt.Printf("\t-s=Number\tStart from Page <number>.\n")
	fmt.Printf("\t-e=Number\tEnd to Page <number>.\n")
	fmt.Printf("\t-l=Number\t[options]Specify the number of line per page.Default is 72.\n")
	fmt.Printf("\t-d=Command\t[options]Execute command.\n")
	fmt.Printf("\t-f\t\t[options]Specify that the pages are sperated by \\f.\n")
	fmt.Printf("\t[filename]\t[options]Read input from the file.\n\n")
	fmt.Printf("If no file specified, selpg will read input from stdin. Control-D to end.\n\n")
	os.Exit(1)
}
