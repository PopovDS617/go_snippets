func postHandler(w http.ResponseWriter, r *http.Request) {
	var animal Animal

	if err := json.NewDecoder(r.Body).Decode(&animal); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Header().Add("Content-Type", "application/json")
		return json.NewEncoder(rw).Encode(v)
		return
	}
	defer r.Body.Close()

	// write json
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)

}
