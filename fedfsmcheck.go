package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/net/html"
)

func getHtmlPage(url, userAgent string) ([]byte, error) {
	// функция получения ресурсов по указанному адресу url с использованием User-Agent
	// возвращает загруженный HTML контент
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Cannot create new request  %s, error: %v\n", url, err)
		return nil, err
	}

	// с User-agent по умолчанию контент не отдается, используем свой
	req.Header.Set("User-Agent", userAgent)

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error with GET request: %v\n", err)
		return nil, err
	}

	defer resp.Body.Close()

	// Получаем контент и возвращаем его, как результат работы функции
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error ReadAll")
		return nil, err
	}
	return body, nil
}

func getList(body []byte, tag string) []string {
	// Функция получения списка из html контента
	// Список достается из тега tag
	// Возвращает список
	tkn := html.NewTokenizer(bytes.NewReader(body))
	depth := 0
	var flist []string
	errorCode := false

	// Проходим по всему дереву тегов (пока не встретится html.ErrorToken)
	for !errorCode {
		tt := tkn.Next()
		switch tt {
		case html.ErrorToken:
			errorCode = true
		case html.TextToken:
			if depth > 0 {
				flist = append(flist, strings.Trim(string(tkn.Text()), " \n\t")) // Если внутри нужного тега, забираем текст из блока
			}
		case html.StartTagToken, html.EndTagToken:
			tn, tAttr := tkn.TagName()
			if string(tn) == "div" { // выбираем нужный tag
				if tAttr {
					key, attr, _ := tkn.TagAttr()
					if tt == html.StartTagToken && string(key) == "id" && string(attr) == tag {
						depth++ // нужный тег открывается
					}
				} else if tt == html.EndTagToken && depth >= 1 {
					depth--
				}
			}
		}
	}
	return flist
}

func testEq(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func mail(newlist []string, listName, urlList string) {
	addressList := []string{""}
	subject := "Subject: New list " + listName + " Federal Financial Monitoring Service\n"
	address := "To: "
	n_address := 0
	for _, a := range addressList {
		if n_address > 0 {
			address += ", "
		}
		address += a
		n_address += 1
	}
	address += "\n"
	htmlhead := "<html>"
	htmlhead += "<head><title>New list " + listName + " Federal Financial Monitoring Service</title>"
	htmlhead += "<meta charset=\"utf-8\">"
	htmlhead += "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">"
	htmlhead += "</head><body><h1>New list Federal Financial Monitoring Service</h1>"
	if listName == "UL" {
		htmlhead += "<h2>Организации</h2><br>"
	}
	if listName == "FL" {
		htmlhead += "<h2>Физические лица</h2><br>"
	}
	htmlhead += "<a href=\"" + urlList + "\">Перечень террористов и экстремистов (включённые)</a><br><div><ul>"
	headers := []byte(subject + address + "Content-Type: text/html\nMIME-Version: 1.0\n\n" + htmlhead)
	htmlfooter := []byte("</ul></div></body></html>")
	combined_string := []byte(strings.Join(newlist, "<br>"))
	headers = append(headers, combined_string...)
	msg := append(headers, htmlfooter...)
	sendmail := exec.Command("/usr/sbin/sendmail", "-t")
	stdin, err := sendmail.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}
	_, err = sendmail.StdoutPipe()
	if err != nil {
		fmt.Println(err)
	}

	sendmail.Start()
	stdin.Write([]byte(msg))
	stdin.Close()
	// sentBytes, _ := ioutil.ReadAll(stdout)
	sendmail.Wait()

}

func main() {
	var urlList string
	var userAgent string

	// Ключи для командной строки
	flag.StringVar(&urlList, "url", "https://fedsfm.ru/documents/terrorists-catalog-portal-add", "The URL to lists")
	flag.StringVar(&userAgent, "uagent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit", "User Agent")

	flag.Parse()

	// path, _ := os.Executable()
	// path = path[:strings.LastIndex(path, "/")+1]
	// fmt.Println("Path: ", path)

	var byteValue_ul []byte
	var byteValue_fl []byte

	// Читаем файлы со списками
	if _, err := os.Stat("ul_file.txt"); err == nil {
		// Open our jsonFile
		byteValue_ul, err = os.ReadFile("ul_file.txt")
		// if we os.ReadFile returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}
	}

	if _, err := os.Stat("fl_file.txt"); err == nil {
		// Open our jsonFile
		byteValue_fl, err = os.ReadFile("fl_file.txt")
		// if we os.ReadFile returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}
	}

	body, err := getHtmlPage(urlList, userAgent)
	if err != nil {
		fmt.Printf("Error getHtmlPage - %v\n", err)
	} else {
		// Получаем список
		fl_list := getList(body, "russianFL")
		ul_list := getList(body, "russianUL")
		combined_string_fl := []byte(strings.Join(fl_list, ""))
		combined_string_ul := []byte(strings.Join(ul_list, ""))
		if !testEq(byteValue_fl, combined_string_fl) {
			err := os.WriteFile("fl_file.txt", combined_string_fl, 0666)
			if err != nil {
				fmt.Println("Error : ", err)
			}
			// fmt.Println("FL=>", fl_list)
			mail(fl_list, "FL", urlList)
		}
		if !testEq(byteValue_ul, combined_string_ul) {
			err := os.WriteFile("ul_file.txt", combined_string_ul, 0666)
			if err != nil {
				fmt.Println("Error : ", err)
			}
			// fmt.Println("UL=>", ul_list)
			mail(ul_list, "UL", urlList)
		}
	}
}
