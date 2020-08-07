package logging


type Logger interface {
	Log(keyvals ...interface{}) error
}