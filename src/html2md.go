package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"strconv"
	"regexp"
	"io"
)

var nf *os.File

func main(){

	f,err := os.Open("F://aaa.html")
	if err != nil{
		fmt.Println(err)
	}
	defer f.Close()

	nf,_ = os.Create("F://bbb.md")
	defer nf.Close()

	r:=bufio.NewReader(f)
	for{
		s,_:=r.ReadString('\n')
		if index:=strings.Index(s,`<article`);index!=-1{
			if index := strings.Index(s,`<h1`);index!=-1{
				nf.WriteString(fmt.Sprintf("**"))
			}
			index = strings.Index(s,"</a>")
			s = s[index+4:strings.Index(s,"</h1>")]
			s = fmt.Sprintf("%s**  \n\n",s)
			nf.WriteString(s)
			break
		}
	}
	s,_:=r.ReadString('\n')
	if strings.Contains(s,`toc`) {
		nf.WriteString("[TOC]\n\n")
	}
	//var i=0
	//var target []string
	//var ix []int
	//for{
	//	s,err:=r.ReadString('\n')
	//	if strings.Contains(s,`<li>`){
	//		if index := strings.Index(s,`href`); index!=-1{
	//			end:=strings.Index(s,`">`)
	//			t := s[index+6:end]
	//			if t[:2] != "#_" {
	//				var sign string
	//				for j:=i;j>0;j-- {
	//					sign += "#"
	//				}
	//				str := s[end+2:strings.Index(s,`</a>`)]
	//				str = fmt.Sprintf("%s  %s  \n",sign,str)
	//				nf.WriteString(str)
	//			}else{
	//				ix = append(ix, i)
	//				target = append(target,t + `"`)
	//			}
	//		}
	//	}
	//
	//	if strings.Contains(s,`<ul>`) {
	//		i++
	//	}
	//	if strings.Contains(s,`</ul>`) {
	//		i--
	//	}
	//	if i==0 {
	//		content(r,target,ix)
	//		break
	//	}
	//	if err != nil && err==io.EOF {
	//		break
	//	}
	//}


	for {
		s,err:= r.ReadString('\n')

		if err != nil && err==io.EOF{
			break
		}

		if strings.EqualFold(strings.TrimSpace(s),"<hr />"){
			continue
		}
		if index := strings.Index(s,`<h`);index != -1 {

			var str string
			var i int
			str = s[index+2:3]
			//defer func() {
			//	if r := recover(); r != nil {
			//		fmt.Printf("===================%v",r)
			//	}
			//}()
			if i, err = strconv.Atoi(str);err != nil {
				panic("type convert fail")
			}

			var sign string
			for ;i>0;i-- {
				sign += "#"
			}
			var hs string
			if reg := regexp.MustCompile(`>(.*?)<`);reg!=nil {
				result := reg.FindAllStringSubmatch(s,-1)
				for _,str := range result {
					hs += string(str[1])
				}
			}
			//fmt.Println(hs)
			hs = fmt.Sprintf("\n\n%s  %s  \n\n",sign,hs)
			nf.WriteString(hs)

		}

		//nf.WriteString("\n")

		if index := strings.Index(s,`<p><strong>`);index != -1 {
			s = s[11:]

			if index := strings.Index(s,"<br />");index!=-1{
				if strings.Contains(s,`</strong>`) {
					nf.WriteString("\n**" + s[:strings.Index(s,`</strong>`)] + "**  \n")
				}else{
					nf.WriteString("\n**" + s[:index] + "  \n")
				}
				//nf.WriteString("\n**" + s[:strings.Index(s,`</strong>`)] + "**  \n")
				for{
					s,_ := r.ReadString('\n')
					if index := strings.Index(s,"<br />");index!=-1{
						nf.WriteString( s[:index]+"  \n")
					}else{
						if strings.Contains(s,`</`) {
							nf.WriteString(s[:strings.Index(s,`</`)])
						}
						nf.WriteString("**  \n")

						break
					}
				}
			}else{
				nf.WriteString("**"+s[:strings.Index(s,"</strong>")]+"**  \n")
			}
		}

		if index := strings.Index(s,`<table>`);index != -1 {
			for{
				s,_ = r.ReadString('\n')
				if strings.TrimSpace(s) == `</table>` {
					nf.WriteString("\n\n")
					break
				}
				if strings.TrimSpace(s) == `<thead>` {
					nf.WriteString("\n")
					var tableline string
					for{
						s,_ = r.ReadString('\n')
						if index:=strings.Index(s,`</th>`);index!=-1 {
							s = s[4:index]
							tableline += "|----"
							nf.WriteString("|" + s)
						}
						if strings.TrimSpace(s) == `</thead>` {
							nf.WriteString("|  \n")
							tableline += "|  "
							nf.WriteString(tableline)
							break
						}
					}
				}
				if strings.TrimSpace(s) == `<tbody>` {
					nf.WriteString("\n")
					for{
						s,_ = r.ReadString('\n')
						if strings.TrimSpace(s) == `<tr>` {
							for{
								s,_ = r.ReadString('\n')
								if index:=strings.Index(s,`</td>`);index!=-1 {
									s = s[4:index]
									nf.WriteString("|" + s)
								}else{
									nf.WriteString("|  \n")
									break
								}
							}
						}

						if strings.TrimSpace(s) == `</tbody>` {
							nf.WriteString("  \n\n")
							break
						}
					}
				}
			}
		}

		if index:=strings.Index(s,`<pre>`);index!=-1{
			nf.WriteString("\n```\n"+s[11:])
			for{
				s,_ = r.ReadString('\n')
				if index:=strings.Index(s,`</pre>`);index!=-1{
					nf.WriteString("\n```\n\n")
					break
				}
				nf.WriteString(s)
			}
		}

	}

}

//func content(r *bufio.Reader,target []string,ix []int){
//	for i,v := range target {
//		s,err := r.ReadString('\n')
//		if err != nil && err==io.EOF{
//			break
//		}
//
//
//		fmt.Println(v)
//	}
//}
