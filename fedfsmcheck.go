package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func getHtmlPage(url, userAgent string) ([]byte, int, error) {
	// функция получения ресурсов по указанному адресу url с использованием User-Agent
	// возвращает загруженный HTML контент
	client := &http.Client{}
	var Scode int
	Scode = 0

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Cannot create new request  %s, error: %v\n", url, err)
		return nil, Scode, err
	}

	// с User-agent по умолчанию контент не отдается, используем свой
	req.Header.Set("User-Agent", userAgent)

	// Отправляем запрос
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error with GET request: %v\n", err)
		return nil, Scode, err
	}
	Scode = resp.StatusCode
	defer resp.Body.Close()

	// Получаем контент и возвращаем его, как результат работы функции
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error ReadAll")
		return nil, Scode, err
	}
	return body, Scode, nil
}

func getListFsm(body []byte, tag string) []string {
	// Функция получения списка из html контента
	// Список достается из тега tag
	// Возвращает список
	tkn := html.NewTokenizer(bytes.NewReader(body))
	depth := 0
	var flist []string
	errorCode := false
	var trimedstr string

	// Проходим по всему дереву тегов (пока не встретится html.ErrorToken)
	for !errorCode {
		tt := tkn.Next()
		switch tt {
		case html.ErrorToken:
			errorCode = true
		case html.TextToken:
			if depth > 0 {
				trimedstr = strings.Trim(string(tkn.Text()), " \n\t")
				if len(trimedstr) > 0 { //Пустые строки не забираем
					flist = append(flist, trimedstr) // Если внутри нужного тега, забираем текст из блока
				}
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

func getListMinjust(body []byte) []string {
	// Функция получения списка из html контента
	// Список достается из тега tag
	// Возвращает список
	tkn := html.NewTokenizer(bytes.NewReader(body))
	depth := 0
	var flist []string
	errorCode := false
	var trimedstr string
	acctext := ""
	validNum := regexp.MustCompile(`^[ ]*[0-9]+.+`)

	// Проходим по всему дереву тегов (пока не встретится html.ErrorToken)
	for !errorCode {
		tt := tkn.Next()
		switch tt {
		case html.ErrorToken:
			errorCode = true
		case html.TextToken:
			if depth > 0 {
				trimedstr = strings.Trim(string(tkn.Text()), " \n\t")
				if len(trimedstr) > 0 { //Пустые строки не забираем
					acctext += trimedstr + " " // Если внутри нужного тега, забираем текст из блока
				}
			}
		case html.StartTagToken, html.EndTagToken:
			tn, _ := tkn.TagName()
			if string(tn) == "tr" { // выбираем нужный tag
				if tt == html.StartTagToken {
					depth++ // нужный тег открывается
				} else if tt == html.EndTagToken && depth >= 1 {
					if validNum.MatchString(acctext) { // Строка должна начинаться с числа
						flist = append(flist, acctext) // При закрытии тега добавляем в список
					}
					acctext = ""
					depth--
				}
			}
		}
	}
	return flist
}

func getListSpimex(body []byte) []string {
	// Функция получения списка из html контента
	// Список достается из тега tag
	// Возвращает список
	tkn := html.NewTokenizer(bytes.NewReader(body))
	depth := 0
	othertag := 0
	var flist []string
	errorCode := false
	var trimedstr string
	acctext := ""

	// Проходим по всему дереву тегов (пока не встретится html.ErrorToken)
	for !errorCode {
		tt := tkn.Next()
		switch tt {
		case html.ErrorToken:
			errorCode = true
		case html.TextToken:
			if depth > 0 {
				trimedstr = strings.Trim(string(tkn.Text()), " \n\t")
				if len(trimedstr) > 0 { //Пустые строки не забираем
					acctext += trimedstr + " " // Если внутри нужного тега, забираем текст из блока
				}
			}
		case html.StartTagToken, html.EndTagToken:
			tn, tAttr := tkn.TagName()
			if string(tn) == "div" { // выбираем нужный tag
				// fmt.Println("depth:", depth, "     othertag:", othertag)
				if tAttr {
					key, attr, _ := tkn.TagAttr()
					if tt == html.StartTagToken {
						// fmt.Println("key:", string(key), "     attr:", string(attr))
						if depth == 1 {
							othertag++ // считаем другие такие же теги внутри нужного
						}
						if string(key) == "class" && string(attr) == "news-item" {
							depth++ // нужный тег открывается
						}
					}
				} else if tt == html.EndTagToken && depth == 1 {
					if othertag == 0 {
						flist = append(flist, acctext) // При закрытии тега добавляем в список
						acctext = ""
						depth--
					} else {
						othertag--
					}
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

func mail(newlist []string, listName, urlList string, addressList []string) {
	var title string
	var titleLink string
	var subject string
	if strings.Contains(listName, "UL") {
		title = "Юридические лица"
	} else if strings.Contains(listName, "FL") {
		title = "Физические лица"
	}
	if strings.Contains(listName, "UL") || strings.Contains(listName, "FL") {
		subject = "Subject: New list Federal Financial Monitoring Service: " + title + "\n"
	} else if listName == "Minjust" {
		subject = "Subject: New list Minjust: нежелательные организации\n"
	} else if listName == "Spimex" {
		subject = "Subject: Биржевая информация: статистические материалы\n"
	}
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
	if strings.Contains(listName, "UL") || strings.Contains(listName, "FL") {
		htmlhead += "<head><title>New list Federal Financial Monitoring Service</title>"
	} else if listName == "Minjust" {
		htmlhead += "<head><title>New list Minjust: нежелательные организации</title>"
	} else if listName == "Spimex" {
		htmlhead += "<head><title>Биржевая информация: статистические материалы</title>"
	}
	htmlhead += "<meta charset=\"utf-8\">"
	htmlhead += "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">"
	if strings.Contains(listName, "add") {
		title += " (включенные)"
		htmlhead += "</head><body><h1>" + title + "</h1>"
	} else if strings.Contains(listName, "FL") {
		title += " (исключённые)"
		htmlhead += "</head><body><h1>" + title + "</h1>"
	} else if listName == "Minjust" {
		htmlhead += "</head><body><h1>Перечень иностранных и международных неправительственных организаций, деятельность которых признана нежелательной на территории Российской Федерации</h1>"
	} else if listName == "Spimex" {
		htmlhead += "</head><body><h1>Статистические материалы</h1>"
	}
	if strings.Contains(listName, "add") {
		titleLink = "Перечень террористов и экстремистов (включённые)"
	} else if strings.Contains(listName, "del") {
		titleLink = "Перечень террористов и экстремистов (исключённые)"
	} else if listName == "Minjust" {
		titleLink = "Перечень нежелательных организаций"
	} else if listName == "Spimex" {
		titleLink = "Биржевая информация: статистические материалы"
	}
	htmlhead += "<a href=\"" + urlList + "\">" + titleLink + "</a><br><br><br><div><ul>"
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
	var userAgent string

	// Ключи для командной строки
	flag.StringVar(&userAgent, "uagent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit", "User Agent")

	flag.Parse()

	path, _ := os.Executable()
	path = path[:strings.LastIndex(path, "/")+1]

	var byteValue_list []byte

	type Configlist struct {
		List   string
		Emails []string
		Url    string
	}

	var configlist []Configlist
	// Читаем файл с адресами
	if _, err := os.Stat(path + "/emails.json"); err == nil {
		// Open our jsonFile
		byteValue, err := os.ReadFile(path + "/emails.json")
		// if we os.ReadFile returns an error then handle it
		if err != nil {
			fmt.Println(err)
		}
		// defer the closing of our jsonFile so that we can parse it later on
		// var listHash []ArticleH
		err = json.Unmarshal(byteValue, &configlist)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Читаем файлы со списками. Файлы в порядке, указанном в конфигурационном файле
	for key, el := range configlist {
		keystr := strconv.Itoa(key)
		if _, err := os.Stat(path + "/file_" + keystr + ".txt"); err == nil {
			// Open our jsonFile
			byteValue_list, err = os.ReadFile(path + "/file_" + keystr + ".txt")
			// if we os.ReadFile returns an error then handle it
			if err != nil {
				fmt.Println(err)
			}
		}

		body, statuscode, err := getHtmlPage(el.Url, userAgent)
		if err != nil || statuscode != 200 {
			fmt.Printf("Error getHtmlPage - %v\n", err)
			fmt.Printf("Status Code - %d\n", statuscode)
		} else {
			// Получаем список
			var get_list []string
			if el.List == "ULadd" || el.List == "ULdel" {
				get_list = getListFsm(body, "russianUL")
			} else if el.List == "FLadd" || el.List == "FLdel" {
				get_list = getListFsm(body, "russianFL")
			} else if el.List == "Minjust" {
				get_list = getListMinjust(body)
			} else if el.List == "Spimex" {
				get_list = getListSpimex(body)
			}
			combined_string := []byte(strings.Join(get_list, ""))
			if !testEq(byteValue_list, combined_string) {
				err := os.WriteFile(path+"/file_"+keystr+".txt", combined_string, 0666)
				if err != nil {
					fmt.Println("Error : ", err)
				}
				// fmt.Println("FL=>", get_list)
				mail(get_list, el.List, el.Url, el.Emails)
			}
		}
	}
}
