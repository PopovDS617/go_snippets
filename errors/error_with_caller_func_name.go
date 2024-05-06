fmt.Errorf("%s query get: %w", utils.GetFuncName(), err), 
где GetFuncName() функция по типу: 
func FuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}