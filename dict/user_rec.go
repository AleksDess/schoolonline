package dict

import (
	"fmt"
	"net/url"
	"schoolonline/config"
	"strconv"
	"strings"
)

type UserLink struct {
	School     int
	SchoolName string
	Role       string
}

func createLink(a UserLink) string {
	return fmt.Sprintf("%s://%s/input/user/%d?%s?%s", config.C.UrlType, config.C.UrlSite, a.School, a.SchoolName, a.Role)
}

func new(school int, name, role string) (rs UserLink) {
	rs.Role = role
	rs.School = school
	rs.SchoolName = name
	return
}

func GetLinkUserRegister(school int, name, role string) string {
	return createLink(new(school, name, role))
}

func ParseURL(input string) (string, int, string, string, error) {
	// Парсинг URL
	u, err := url.Parse(input)
	if err != nil {
		fmt.Println(err)
		return "", -1, "", "", err
	}

	// Разделяем путь
	pathParts := strings.Split(u.Path, "/")
	if len(pathParts) < 3 {
		return "", -1, "", "", fmt.Errorf("недостаточно элементов в пути")
	}

	// Извлекаем параметры
	inputUser := pathParts[1] // input/user
	id := pathParts[2]        // 5

	id_int, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println(err)
		return "", -1, "", "", err
	}

	// Извлечение строки запроса
	queryParams := u.RawQuery
	queryParts := strings.Split(queryParams, "?")
	if len(queryParts) < 2 {
		return "", -1, "", "", fmt.Errorf("недостаточно параметров в строке запроса")
	}

	// Первый параметр (Красота)
	param1 := queryParts[0]
	// Второй параметр (parent)
	param2 := queryParts[1]

	return inputUser, id_int, param1, param2, nil
}
