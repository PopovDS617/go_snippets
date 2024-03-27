
var AuthAdminID = "middleware.auth.AdminID"

type Middleware func(http.Handler) http.Handler

func CreateMiddlewareStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}

}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader((statusCode))
	w.statusCode = statusCode
}

// middleware
func LoggingMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)
		log.Println(r.Method, wrapped.statusCode, r.URL.Path, time.Since(start))
	})
}

func CheckPermissionsMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking Permissions . . .")
		next.ServeHTTP(w, r)
	})
}

func isAdminMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		receivedToken := r.Header.Get("Token")

		// Check that the header begins with a prefix of Bearer
		if !strings.HasPrefix(receivedToken, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}

		// Pull out the token
		encodedToken := strings.TrimPrefix(receivedToken, "Bearer ")

		// Decode the token from base 64
		token, err := base64.StdEncoding.DecodeString(encodedToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return
		}

		adminID := string(token)

		ctx := context.WithValue(r.Context(), AuthAdminID, adminID)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}

func main() {

	router := http.NewServeMux()

	// routes
	router.HandleFunc("GET /items/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		w.Write([]byte("received request for item: " + id))
	})

	router.HandleFunc("POST /monsters/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Monster created!"))
	})

	adminRouter := http.NewServeMux()
	adminRouter.HandleFunc("POST /invoice/", func(w http.ResponseWriter, r *http.Request) {
		adminID, ok := r.Context().Value(AuthAdminID).(string)

		if !ok {
			log.Println("invalid admin ID")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Invoice created by " + adminID))
	})

	router.Handle("/", isAdminMW(adminRouter))

	// middleware
	mwStack := CreateMiddlewareStack(
		LoggingMW,
		CheckPermissionsMW,
	)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router))

	// server
	server := http.Server{
		Addr:    ":9000",
		Handler: mwStack(v1),
	}

	fmt.Println("Starting server on port :9000")

	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
