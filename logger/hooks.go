package logger

type Hook interface {
	Levels() []Level
	Fire(*Logger) error
}

// Internal type for storing the hooks on a logger instance.
type LevelHooks map[Level][]Hook

// Add a hook to an instance of logger. This is called with
// `log.Hooks.Add(new(MyHook))` where `MyHook` implements the `Hook` interface.
func (hooks LevelHooks) Add(hook Hook) {
	for _, level := range hook.Levels() {
		hooks[level] = append(hooks[level], hook)
	}
}

// Fire all the hooks for the passed level. Used by `logger.log` to fire
// appropriate hooks for a log entry.
func (hooks LevelHooks) Fire(level Level, l *Logger) error {
	for _, hook := range hooks[level] {
		if err := hook.Fire(l); err != nil {
			return err
		}
	}
	return nil
}


var Loggerlevels = [4]Level{LevelDebug, LevelInfo, LevelWarn, LevelError}

func (l *Logger) AddLfsHook(){
//	var lvIdx int
//	for i, lv := range Loggerlevels{
//		if l.Level == lv{
//			lvIdx = i
//		}
//	}
//	lvs := Loggerlevels[lvIdx:]
//	hookMap := make(lfshook.WriterMap)
//	for _, lv := range lvs{
//		hookMap[lv] = io.Writer()
//	}
//	lfHook := lfshook.NewHook(hookMap, l.Formatter)
//	l.Hooks.Add(lfHook)
}
