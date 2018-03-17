package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"               //"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"        //"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3" //"google.golang.org/api/calendar/v3"
)

var username = "COSFAR"

//использует Context и Config для извлечения токена
//затем генеруруем клиент, return возвращает сгенерированного клиента
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Не удается получить путь к кэшированному файлу учетных данных. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

//ЭТО НАДО В БОТ ТЕЛЕГРАМА
//использует Config для запроса токена
//return возвращает полученный токен
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Перейдите по следующей ссылке в браузере и введите "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Не удается прочитать код авторизации %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Не удается получить маркер из интернета %v", err)
	}
	return tok
}

//генерирует  путь, имя файла учетных данных
//return возвращает сгенерированные учетные данные
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("calendar"+username+".json")), err
}

//извлекает токен из заданного пути к файлу
// return возвращает в извлеченный токен и любую обнаруженную ошибку чтения
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

//использует путь к файлу для создания файла и хранения токена в нем
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Сохранение файла учетных данных в: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Не удается кэшировать маркер oauth: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Не удается прочитать секретный файл клиента: %v", err)
	}

	//если при именении этих областей, удалите ранее созданные учетные данные в
	// at ~/.credentials/calendar-go-quickstart.json
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Не удается проанализировать секретный файл клиента для настройки: %v", err)
	}
	client := getClient(ctx, config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Не удается получить клиента календаря %v", err)
	}

	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Не удается получить следующие десять событий пользователя. %v", err)
	}
	
	calendars, err := srv.CalendarList.List().Do();
	if len(calendars.Items) > 0 {
		for _, i := range calendars.Items {
			fmt.Printf("%s (%s) \n", i.Id, i.Summary)
		}
	}
	
	fmt.Println("Ближайшие события:")
	if len(events.Items) > 0 {
		for _, i := range events.Items {
			var when string
			//если время пустая строка событие является на весь день, так как доступна только дата
			if i.Start.DateTime != "" {
				when = i.Start.DateTime
			} else {
				when = i.Start.Date
			}
			fmt.Printf("%s (%s) (%s)\n", i.Id, i.Summary, when)
		}
	} else {
		fmt.Printf("Предстоящие события не найдены.\n")
	}
	
		
	fmt.Println("\nВведите задание:\n0 - Добавить задачу\n1 - Изменить задачу\n2 - Удалить задачу\n3 - Показать задачи\nИ нажмите Enter\n")
	var i int;
	_, err = fmt.Scanf("%d", &i)
	if (err != nil) {
	  fmt.Println("wtf: ", err)
	}
	
	switch i {
	case 0:
		//addEvt();
		
	case 1:
		//updateEvt()
	case 2:
		//deleteEvt()
		var eventID string		
		fmt.Println("Введите ID события:")
		_, err = fmt.Scanf("%s", &eventID)		
		err = srv.Events.Delete("primary", eventID).Do()
		if (err != nil) {
			fmt.Println("wtf: ", err)		
		}
	case 3:
		//showEvts()
		fmt.Println("Ближайшие события:")
		if len(events.Items) > 0 {
			for _, i := range events.Items {
				var when string
				//если время пустая строка событие является на весь день, так как доступна только дата
				if i.Start.DateTime != "" {
					when = i.Start.DateTime
				} else {
					when = i.Start.Date
				}
				fmt.Printf("%s (%s) (%s)\n", i.Id, i.Summary, when)
			}
		} else {
			fmt.Printf("Предстоящие события не найдены.\n")
		}
	default:
		fmt.Println("wtf: incorrect value")
	}


}
