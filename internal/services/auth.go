package service

import (
	"crypto/sha1"
	"fmt"

	"github.com/AquaEngineering/AquaHub/internal/domain"
	cachememory "github.com/o-sokol-o/cache-memory"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const salt = "hjqrhjqw124617ajfhajs" // Случайные символы для генерации хеша пароля

type AuthService struct {
	repo  IStoreAuthorization
	cache cachememory.Cache
}

func NewAuthService(cache cachememory.Cache, repo IStoreAuthorization) *AuthService {
	return &AuthService{
		cache: cache,
		repo:  repo}
}

func (s *AuthService) GetUser_LoginPassword(email, password string) (domain.User, error) {

	// Запросить пользователя из БД
	return s.repo.GetUser_LoginPassword(email, generatePasswordHash(password))
}

// Authenticate finds a user by their email and verifies their password.
// On success it returns a User
func (s *AuthService) Authenticate(email, password string) (*domain.User, error) { //(Token, error) {

	// Запросить пользователя из БД
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Append the salt from the user record to the supplied password.
	saltedPassword := password + user.PasswordSalt

	// Compare the provided password with the saved hash. Use the bcrypt comparison
	// function so it is cryptographically secure. Return authentication error for
	// invalid password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(saltedPassword)); err != nil {
		// err = errors.WithStack(ErrAuthenticationFailure)
		return nil, err
	}

	return user, nil
}

// --------------------------------- Old -----------------------------------

// Генерация хеша пароля
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	// Случайные символы для генерации хеша пароля
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

// Принимает структуру User в качестве аргумента
// Возвращает ID созданного нового пользователя из БД
func (s *AuthService) CreateUser(user domain.User) (int, error) {

	// Запросить пользователя из БД
	_, err := s.repo.GetUserByEmail(user.Email)
	if err == nil {
		return 0, errors.New(fmt.Sprintf("user %s already exists", user.Email))
	}

	// Генерация хеша пароля
	user.PasswordSalt = uuid.NewRandom().String()
	saltedPassword := user.PasswordHash + user.PasswordSalt
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.Wrap(err, "generating password hash")
	}
	user.PasswordHash = string(passwordHash)

	// Создание пользователя с аккаунтом
	id, err := s.repo.CreateUser(user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// // Запросить токен пользователя
// func (s *AuthService) GenerateToken(username, password string) (string, error) {

// 	// Запросить в БД id пользователя
// 	user, err := s.repo.GetUser_LoginPassword(username, generatePasswordHash(password))
// 	if err != nil {
// 		return "", err
// 	}

// 	// Генерируем токен из Стандартной подписи и Claims
// 	// Claims - JSON объект с набором полей
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
// 		jwt.StandardClaims{
// 			// ExpiresAt = на 12 часов болше текущего времени
// 			// т.е. токен перестанет быть валидным через 12 часов
// 			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
// 			// Время генерации токена
// 			IssuedAt: time.Now().Unix(),
// 		},
// 		user.Id,
// 	})

// 	// Подпишем и вернём токен с ключём подписи signingKey. Для расшифровки он же.
// 	return token.SignedString([]byte(signingKey))
// }

// // Метод ParseToken принимает token в качестве аргумента
// // и возвращает id пользователя при успешном парcинге
// func (s *AuthService) ParseToken(accessToken string) (int, error) {

// 	// Вызываем функцию ParseWithClaims из библиотеки jwt, которая принимает:
// 	//   - token,
// 	//   - структуру Claims,
// 	//   - функцию которая возвращает ключ подписи или ошибку
// 	// В этой функции нам нужно проверить метод подписи токена.
// 	// Если это не HMAC то мы возвращаем ошибку
// 	// А если всё О'кей то ключ подписи
// 	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, errors.New("invalid signing method")
// 		}

// 		return []byte(signingKey), nil
// 	})

// 	if err != nil {
// 		return 0, err
// 	}

// 	// Функция ParseWithClaims возвращает объект token в котором есть поле Claims типа интерфейс
// 	// Приведём его к нашей структуре и проверим всё ли хорошо
// 	claims, ok := token.Claims.(*tokenClaims)
// 	if !ok {
// 		return 0, errors.New("token claims are not of type *tokenClaims")
// 	}

// 	// Возвращаем id пользователя при успешном парcинге token
// 	return claims.UserId, nil
// }

// Запросить по токенам из БД id_шки аккаунта, аквахаба, устройств и сенсоров
func (s *AuthService) GetUserHWfromTokens(h_token, u_token string) ([]domain.SensorDataSet, error) { // user_id, aquahub_id, []{device_id, sensor_id}

	var lists []domain.SensorDataSet
	if s.cache != nil {
		if lsts, err := s.cache.Get("AquaHub-HW-ID" + h_token + u_token); err == nil {

			return lsts.([]domain.SensorDataSet), nil

		} else {

			// Запросить по токенам из БД id_шки аккаунта, аквахаба, устройств и сенсоров
			lists, err = s.repo.GetUserHW_fromTokens(h_token, u_token)

			s.cache.Set("AquaHub-"+h_token+u_token, lists)

			return lists, err
		}

	}

	lists, err := s.repo.GetUserHW_fromTokens(h_token, u_token)
	return lists, err
}

// Запросить по токенам из БД id_шки аккаунта, аквахаба, устройств и сенсоров
func (s *AuthService) GetAquahubIdfromTokens(h_token, u_token string) (int, error) { // user_id, aquahub_id, []{device_id, sensor_id}

	// Запросить по токенам из БД id_шки аккаунта, аквахаба, устройств и сенсоров
	aId, err := s.repo.GetAquahubId_fromTokens(h_token, u_token)
	return aId, err
}
