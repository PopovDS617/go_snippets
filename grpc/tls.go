// client
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

creds, err := credentials.NewClientTLSFromFile("service.pem", "")
if err != nil {
	log.Fatalf("could not process the credentials: %v", err)
}

conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
if err != nil {
	log.Fatalf("failed to connect to server: %v", err)
}
defer conn.Close()

c := desc.NewNoteV1Client(conn)

// server
/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
if err != nil {
	log.Fatalf("failed to listen: %v", err)
}

creds, err := credentials.NewServerTLSFromFile("service.pem", "service.key")
if err != nil {
	// log.Fatalf("failed to load TLS keys: %v", err)
}

s := grpc.NewServer(grpc.Creds(creds))
reflection.Register(s)
desc.RegisterNoteV1Server(s, &server{})

// log.Printf("server listening at %v", lis.Addr())

if err = s.Serve(lis); err != nil {
	// log.Fatalf("failed to serve: %v", err)
}