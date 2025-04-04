package token

import "time"

// Общий интерфейс для создания и проверки токенов
type Maker interface {
	// Этот метод создает и подписывает новый токен для конкретного имени пользователя и допустимой длительности
	CreateToken(username string, role string, duration time.Duration) (string, *Payload, error)

	// Проверить действительный ли входной токен, если действителен - метод вернет данные хранящиеся в теле токена
	VefifyToken(token string) (*Payload, error)
}
