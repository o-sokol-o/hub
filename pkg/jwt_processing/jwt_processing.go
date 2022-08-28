package jwt_processing

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	signingKey = "c3ert#5vqzpm34&s6hhie8ngt[is" // Случайные символы для
	tokenTTL   = 12 * time.Hour
)

// Структура со стандартным Claims и с добавленным полем id пользователя
// Где будет сохранятся всё о токене
type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

// Запросить токен пользователя
func GenerateToken(user_id int) (string, error) {

	// Генерируем токен из Стандартной подписи и Claims
	// Claims - JSON объект с набором полей
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			// ExpiresAt = на 12 часов болше текущего времени
			// т.е. токен перестанет быть валидным через 12 часов
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			// Время генерации токена
			IssuedAt: time.Now().Unix(),
		},
		user_id,
	})

	// Подпишем и вернём токен с ключём подписи signingKey. Для расшифровки он же.
	return token.SignedString([]byte(signingKey))
}

// Метод ParseToken принимает token в качестве аргумента
// и возвращает id пользователя при успешном парcинге
func ParseToken(accessToken string) (int, error) {

	// Вызываем функцию ParseWithClaims из библиотеки jwt, которая принимает:
	//   - token,
	//   - структуру Claims,
	//   - функцию которая возвращает ключ подписи или ошибку
	// В этой функции нам нужно проверить метод подписи токена.
	// Если это не HMAC то мы возвращаем ошибку
	// А если всё О'кей то ключ подписи
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	// Функция ParseWithClaims возвращает объект token в котором есть поле Claims типа интерфейс
	// Приведём его к нашей структуре и проверим всё ли хорошо
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	// Возвращаем id пользователя при успешном парcинге token
	return claims.UserId, nil
}
