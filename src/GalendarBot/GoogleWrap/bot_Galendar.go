package GoogleWrap

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

var srv *calendar.Service

//использует Context и Config для извлечения токена
//затем генеруруем клиент, return возвращает сгенерированного клиента
func getConfig() *oauth2.Config {
	b, err := ioutil.ReadFile("./GalendarBot/GoogleWrap/client_secret.json")
	if err != nil {
		log.Fatalf("Не удается прочитать секретный файл клиента: %v", err)
		return nil
	}

	//если при именении этих областей, удалите ранее созданные учетные данные в
	// at ~/.credentials/calendar-go-quickstart.json
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Не удается проанализировать секретный файл клиента для настройки: %v", err)
		return nil
	}

	return config
}

/*
func getClient(ctx context.Context, config *oauth2.Config, clientID string) *http.Client {
	cacheFile, err := tokenCacheFile(clientID)
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
*/

func getClient(ctx context.Context, config *oauth2.Config, clientID string) *http.Client {
	cacheFile, err := tokenCacheFile(clientID)
	if err != nil {
		log.Fatalf("Не удается получить путь к кэшированному файлу учетных данных. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		return nil
	}
	return config.Client(ctx, tok)
}

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

func GetTokenURL() string {
	config := getConfig()
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return authURL
}

func SaveToken(code string, clientID string) bool {
	config := getConfig()
	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Не удается получить маркер из интернета %v", err)
		return false
	}

	cacheFile, err := tokenCacheFile(clientID)
	fmt.Printf("Сохранение файла учетных данных в: %s\n", cacheFile)
	f, err := os.OpenFile(cacheFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Не удается кэшировать маркер oauth: %v", err)
		return false
	}
	defer f.Close()
	json.NewEncoder(f).Encode(tok)

	return true
}

//генерирует  путь, имя файла учетных данных
//return возвращает сгенерированные учетные данные
func tokenCacheFile(username string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	log.Println(usr.HomeDir)
	tokenCacheDir := filepath.Join("./temp/", ".credentials")
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
func saveToken1(file string, token *oauth2.Token) {
	fmt.Printf("Сохранение файла учетных данных в: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Не удается кэшировать маркер oauth: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func ShowEvents( /*srv *calendar.Service*/ ) string {
	t := time.Now().Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Не удается получить следующие десять событий пользователя. %v", err)
	}

	buffer := ""
	if len(events.Items) > 0 {
		for _, i := range events.Items {
			var when string
			//если время пустая строка событие является на весь день, так как доступна только дата
			if i.Start.DateTime != "" {
				when = i.Start.DateTime
			} else {
				when = i.Start.Date
			}
			fmt.Printf("%s (%s) %s\n", i.Summary, when, i.Id)
			buffer = fmt.Sprint(buffer + i.Summary + " " + when + " " + i.Id + " " + "\n")
		}
		return buffer
	} else {
		fmt.Printf("Предстоящие события не найдены.\n")
		return "Дел нет"
	}
}


func AddEvent( /*Event ToDo*/ event *calendar.Event /*, srv *calendar.Service*/) bool {
	event, err := srv.Events.Insert("primary", event).Do()
	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
		return false
	}
	fmt.Printf("Event created: %s\n", event.HtmlLink)
	return true
}

func UpdateEvent(event *calendar.Event /*, srv *calendar.Service*/, eventID string) {
	event, err := srv.Events.Patch("primary", eventID, event).Do()
	if err != nil {
		log.Fatalf("Unable to create event. %v\n", err)
	}
	fmt.Printf("Event Updated: %s\n", event.HtmlLink)
}

func DeleteEvent(event *calendar.Event /*, srv *calendar.Service*/, eventID string) {
	err := srv.Events.Delete("primary", eventID).Do()
	if err != nil {
		log.Fatalf("Unable to delete event. %v", err)
	}
	fmt.Printf("Event Deleted")
}

func Auth(clientID string) bool {
	ctx := context.Background()

	config := getConfig()

	client := getClient(ctx, config, clientID)

	if client == nil {
		return false
	}

	var err error
	srv, err = calendar.New(client)
	if err != nil {
		log.Fatalf("Не удается получить клиента календаря %v", err)
		//return false
	}

	return true
	/*
		event := &calendar.Event{
			Summary:     "и3менить",
			Description: "sdasdsadasd",
			Start: &calendar.EventDateTime{
				DateTime: "2018-03-11T08:00:00+08:00",
				TimeZone: "Asia/Irkutsk",
			},
			End: &calendar.EventDateTime{
				DateTime: "2018-03-11T09:00:00+08:00",
				TimeZone: "Asia/Irkutsk",
			},
		}

		addEvent(event, srv)
	*/

}
