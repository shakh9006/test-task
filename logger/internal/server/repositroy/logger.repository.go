package repository

type LoggerRepository interface {
	WriteLog(string, string) error
}
