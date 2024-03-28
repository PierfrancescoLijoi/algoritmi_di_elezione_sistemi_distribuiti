package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sort"

	"google.golang.org/grpc"
	"progetto_sdcc/interfaccia"
)

var (
	jsonList    []*interfaccia.JSONClientRegistration
	PORT_EXPOSE      = os.Getenv("PORT")
	ELEZIONE         = os.Getenv("ELEZIONE")
	Update      bool = true
)

type serverRingElection struct {
	interfaccia.UnimplementedRingElectionServer
}
type serverRaftElection struct {
	interfaccia.UnimplementedRaftElectionServer
}

// /ring election registry
func (s *serverRingElection) SendRequest(ctx context.Context, item *interfaccia.JSONClientRegistration) (*interfaccia.JSONListPeer, error) {
	for _, existingItem := range jsonList {
		if existingItem.Hostname == item.Hostname && existingItem.Port == item.Port || existingItem.IdPeer == item.IdPeer {
			Update = false
			return &interfaccia.JSONListPeer{Items: jsonList}, nil
		}
	}

	jsonList = append(jsonList, item)
	sort.Slice(jsonList, func(i, j int) bool {
		return jsonList[i].Hostname+jsonList[i].Port < jsonList[j].Hostname+jsonList[j].Port
	})
	if Update == true && len(jsonList) > 1 {

		update_peer_informazion(item.Hostname, item.Port)
	}
	return &interfaccia.JSONListPeer{Items: jsonList}, nil
}
func update_peer_informazion(hostnameOld string, portOld string) {
	for _, obj := range jsonList {
		log.Printf("sto su  %s , %s", obj.Hostname, obj.Port)
		log.Printf("messaggio vecchio %s , %s", hostnameOld, portOld)

		if obj.Port != portOld {
			log.Printf("dentro if")

			if err := sendUpdate_allPeer(obj.Hostname, obj.Port); err != nil {
				log.Printf("Errore durante l'invio dell'aggiornamento a %s:%s: %v", obj.Hostname, obj.Port, err)
			} else {
				log.Printf("Aggiornamento inviato con successo a %s:%s", obj.Hostname, obj.Port)
			}

		}

	}
}
func sendUpdate_allPeer(hostname, port string) error {
	// Connessione al serverRingElection gRPC
	conn, err := grpc.Dial(hostname+":"+port, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("connessione al serverRingElection gRPC fallita: %v", err)
	}
	defer conn.Close()

	// Creazione del client gRPC
	c := interfaccia.NewRingElectionClient(conn)

	item := &interfaccia.JSONListPeer{
		Items: jsonList,
	}

	_, err = c.SendUpdate(context.Background(), item)
	if err != nil {
		log.Fatalf("Could not call MyMethod: %v", err)
	}
	log.Println("Messaggio inviato con successo al peer")

	return nil
}

// ////
func main() {
	fmt.Println("algoritmo selezione nella variabile di ambiente: ", ELEZIONE)
	if ELEZIONE == "anello" {
		lis, err := net.Listen("tcp", ":"+PORT_EXPOSE)
		if err != nil {
			fmt.Printf("Failed to listen: %v\n", err)
			return
		}

		s := grpc.NewServer()
		interfaccia.RegisterRingElectionServer(s, &serverRingElection{})
		if err := s.Serve(lis); err != nil {
			fmt.Printf("Failed to serve: %v\n", err)
			return
		}
	} else if ELEZIONE == "raft" {
		lis, err := net.Listen("tcp", ":"+PORT_EXPOSE)
		if err != nil {
			fmt.Printf("Failed to listen: %v\n", err)
			return
		}

		s := grpc.NewServer()
		interfaccia.RegisterRaftElectionServer(s, &serverRaftElection{})
		if err := s.Serve(lis); err != nil {
			fmt.Printf("Failed to serve: %v\n", err)
			return
		}
	} else {
		log.Println("credo che hai sbagliato a digitiare il nome dell'algoritmo di elezione,\n" +
			" la mia traccia B2 aveva come algortimi l'algortimo ad anello e raft")
	}

}

/////

// /raft election registry
func update_peer_informazion_raft(hostnameOld string, portOld string) {
	for _, obj := range jsonList {
		log.Printf("sto su  %s , %s", obj.Hostname, obj.Port)
		log.Printf("messaggio vecchio %s , %s", hostnameOld, portOld)
		////modificare qua !!!!!!!
		if obj.Port != portOld {
			log.Printf("dentro if")

			if err := sendUpdate_allPeer_raft(obj.Hostname, obj.Port); err != nil {
				log.Printf("Errore durante l'invio dell'aggiornamento a %s:%s: %v", obj.Hostname, obj.Port, err)
			} else {
				log.Printf("Aggiornamento inviato con successo a %s:%s", obj.Hostname, obj.Port)
			}

		}

	}
}
func sendUpdate_allPeer_raft(hostname, port string) error {
	// Connessione al serverRingElection gRPC
	conn, err := grpc.Dial(hostname+":"+port, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("connessione al serverRingElection gRPC fallita: %v", err)
	}
	defer conn.Close()

	// Creazione del client gRPC
	c := interfaccia.NewRaftElectionClient(conn)

	item := &interfaccia.JSONListPeer{
		Items: jsonList,
	}

	_, err = c.SendUpdate(context.Background(), item)
	if err != nil {
		log.Printf("Could not call MyMethod: %v", err)
	}
	log.Println("Messaggio inviato con successo al peer")

	return nil
}
func (s *serverRaftElection) SendRequest(ctx context.Context, item *interfaccia.JSONClientRegistration) (*interfaccia.JSONListPeer, error) {
	for _, existingItem := range jsonList {
		if existingItem.Hostname == item.Hostname && existingItem.Port == item.Port || existingItem.IdPeer == item.IdPeer {
			Update = false
			return &interfaccia.JSONListPeer{Items: jsonList}, nil
		}
	}

	jsonList = append(jsonList, item)
	sort.Slice(jsonList, func(i, j int) bool {
		return jsonList[i].Hostname+jsonList[i].Port < jsonList[j].Hostname+jsonList[j].Port
	})
	if Update == true && len(jsonList) > 1 {

		update_peer_informazion_raft(item.Hostname, item.Port)
	}
	return &interfaccia.JSONListPeer{Items: jsonList}, nil
}
