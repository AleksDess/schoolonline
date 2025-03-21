package crypto

// Исходная карта кодирования
var encodeMap = map[rune]string{
	'a': "Qx", 'b': "Wd", 'c': "Er", 'd': "Ty", 'e': "Ui",
	'f': "Op", 'g': "As", 'h': "Df", 'i': "Gh", 'j': "Jk",
	'k': "Lz", 'l': "Xc", 'm': "Vb", 'n': "Nm", 'o': "Qa",
	'p': "Ws", 'q': "Ed", 'r': "Rf", 's': "Tg", 't': "Yh",
	'u': "Uj", 'v': "Ik", 'w': "Ol", 'x': "Pm", 'y': "Az",
	'z': "Sc",

	'0': "Hn", '1': "Rt", '2': "Vq", '3': "Xp", '4': "Yv",
	'5': "Fb", '6': "Cj", '7': "Kg", '8': "Bl", '9': "Zm",
	':': "Ji",
}

// Обратная карта для декодирования
var decodeMap = map[string]rune{
	"Qx": 'a', "Wd": 'b', "Er": 'c', "Ty": 'd', "Ui": 'e',
	"Op": 'f', "As": 'g', "Df": 'h', "Gh": 'i', "Jk": 'j',
	"Lz": 'k', "Xc": 'l', "Vb": 'm', "Nm": 'n', "Qa": 'o',
	"Ws": 'p', "Ed": 'q', "Rf": 'r', "Tg": 's', "Yh": 't',
	"Uj": 'u', "Ik": 'v', "Ol": 'w', "Pm": 'x', "Az": 'y',
	"Sc": 'z',

	"Hn": '0', "Rt": '1', "Vq": '2', "Xp": '3', "Yv": '4',
	"Fb": '5', "Cj": '6', "Kg": '7', "Bl": '8', "Zm": '9',
	"Ji": ':',
}

// Функция кодирования строки
func EncodeTgLink(input string) string {
	var result string
	for _, char := range input {
		if code, ok := encodeMap[char]; ok {
			result += code // Заменяем символ на соответствующий код
		} else {
			result += string(char) // Если символ не найден, добавляем его как есть
		}
	}
	return result
}

// Функция декодирования строки
func DecodeTgLink(input string) string {
	var result string
	// Итерируем по входной строке, беря по 2 символа
	for i := 0; i < len(input); i += 2 {
		if i+1 >= len(input) {
			return ""
		}
		code := input[i : i+2] // Берем текущие два символа
		if char, ok := decodeMap[code]; ok {
			result += string(char) // Заменяем код на соответствующий символ
		} else {
			return ""
		}
	}
	return result
}
