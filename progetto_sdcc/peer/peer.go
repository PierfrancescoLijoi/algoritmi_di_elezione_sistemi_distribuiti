package main

import (
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"math/rand"
	"net"
	"os"
	"progetto_sdcc/interfaccia"
	"strings"
	"time"
)

const (
	ADDRESS_REGISTRY = "registry:8080" // localhostname "localhostname:8080" se in locale
	ELEZIONE         = "anello"        // anello o raft

)

var (
	HOSTNAME          string = os.Getenv("HOSTNAME") // in locale è uguale per tutti
	PORT_EXPOSE              = os.Getenv("PORT")     //da 50051 fino a quel che si vuole
	ID_PEER                  = os.Getenv("HOSTNAME")
	LeaderPeerAddress string = ""
	jsonListPeer      []Peer
	FollowerChannel          = make(chan interface{})
	ElectionChannel          = make(chan interface{})
	HeartbeatChannel         = make(chan interface{})
	StopChannel              = make(chan interface{})
	TimerElection     int    = -1
	TimerHeartBeat    int    = -1
	Term              int32  = 0
	ElectionVote      int    = 0
	LeaderID          string = ""
)

type serverRingElection struct {
	interfaccia.UnimplementedRingElectionServer
}
type serverRaftElection struct {
	interfaccia.UnimplementedRaftElectionServer
}
type Peer struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	IDpeer   string `json:"iDpeer"`
}

// /////////////algoritmo ring election
func (s *serverRingElection) SendElection(ctx context.Context, req *interfaccia.MessageElection) (*emptypb.Empty, error) {
	//LeaderPeerAddress = "" // se mi arriva una richiesta di elezione vuol dire che devo dimenticare chi era prima il leader e resettare

	if len(req.ListIdPeer) == 0 {
		return &emptypb.Empty{}, nil
	} else {

		var finished bool = false

		for _, id := range req.ListIdPeer {
			if id == ID_PEER {
				finished = true
				log.Printf("ELEZIONE FINITA")

				break
			}

		}
		if finished == false {

			idPeerPosition := findMe() // ELEZIONE
			fmt.Println("HO chiamato connessione in procedura remota")
			var listPeer []string = req.ListIdPeer //inoltro lista ricevuta
			findValidConnection(idPeerPosition, listPeer)

		} //else
		if finished == true {
			resultRingElection(req.ListIdPeer)
		}

	}
	return &emptypb.Empty{}, nil

}
func (s *serverRingElection) SendResult(ctx context.Context, req *interfaccia.ResultElection) (*emptypb.Empty, error) {
	if LeaderPeerAddress != req.IdPeers {

		myAddress := HOSTNAME + ":" + ID_PEER
		if myAddress == req.IdPeers {
			LeaderPeerAddress = myAddress

		} else {
			LeaderPeerAddress = req.IdPeers
		}
		fmt.Println("RICEVUTA NOTIZIA")
		comunicationResult()
		//return &emptypb.Empty{}, nil
	}
	return &emptypb.Empty{}, nil
}
func (s *serverRingElection) SendCheck(ctx context.Context, item *interfaccia.CheckLeader) (*interfaccia.CheckLeader, error) {

	return &interfaccia.CheckLeader{Check: "OK"}, nil
}
func (s *serverRingElection) SendUpdate(ctx context.Context, item *interfaccia.JSONListPeer) (*emptypb.Empty, error) {
	// Sovrascrivi la jsonListPeer con la lista ricevuta dalla risposta del registro
	jsonListPeer = nil
	for _, jsonItem := range item.Items {
		jsonListPeer = append(jsonListPeer, Peer{
			Hostname: jsonItem.Hostname,
			Port:     jsonItem.Port,
			IDpeer:   jsonItem.IdPeer,
		})
	}
	fmt.Println("jsonListPeer aggiornata:", jsonListPeer)

	for _, item := range item.Items {

		address := fmt.Sprintf("%s:%s", item.Hostname, item.Port)
		fmt.Println("Hostname e porta:", address)
	}
	return &emptypb.Empty{}, nil
}
func findValidConnection(startPosition int, listElection []string) {

	var peerNext int

	listElection = append(listElection, ID_PEER) //ci devo appendere il valore sia che è vuota che non !

	if startPosition == len(jsonListPeer)-1 {
		peerNext = 0
	} else {
		peerNext = startPosition + 1
	}
	fmt.Println("posizione successivo a me", peerNext)

	for jsonListPeer[peerNext].Hostname != HOSTNAME && jsonListPeer[peerNext].Port != PORT_EXPOSE { //per quando sarà distribuito
		//for jsonListPeer[peerNext].Port != PORT_EXPOSE { //per il locale

		fmt.Println("porta corrente nel for :", jsonListPeer[peerNext].Port)
		///////////////////////////
		var address string = jsonListPeer[peerNext].Hostname + ":" + jsonListPeer[peerNext].Port
		fmt.Printf("address next to peer: %s\n", address)
		conn, err := connectToPeer(address)
		defer conn.Close()

		if err != nil {
			fmt.Printf("Tentativo connessione peer successivo fallita: %v\n", err)

		} else {

			fmt.Println("Lista di stringhe:", listElection)

			election := &interfaccia.MessageElection{
				ListIdPeer: listElection,
			}

			c := interfaccia.NewRingElectionClient(conn)
			_, err = c.SendElection(context.Background(), election)
			if err != nil {
				fmt.Printf("Elezione fallita: %v\n", err)
			} else {
				fmt.Println("MESSAGGIO ELEZIONE INVIATO")
				break
			}
		}
		peerNext++
		if peerNext > len(jsonListPeer)-1 {
			peerNext = 0
			fmt.Println("correzione indice nell'anello: ", peerNext)
		}
	}
	fmt.Println("Sono uscito dal loop di find valid connection")

}
func findMe() int {
	fmt.Println("sto in findMe")
	var idPeerPosition int = 0

	for idPeerPosition = 0; idPeerPosition < len(jsonListPeer)-1; idPeerPosition++ {
		if jsonListPeer[idPeerPosition].Hostname == HOSTNAME && jsonListPeer[idPeerPosition].Port == PORT_EXPOSE {
			// Quando trovi l'oggetto uguale ad A, ottieni le informazioni sull'oggetto successivo
			break // Esci dal loop dopo aver trovato l'oggetto
		}
	}
	fmt.Println("me", jsonListPeer[idPeerPosition].Hostname)
	fmt.Println("me", jsonListPeer[idPeerPosition].Port)

	return idPeerPosition

}
func resultRingElection(ListIdPeer []string) {

	maxID := ""
	fmt.Println("chi è il migliore?")

	for _, id := range ListIdPeer {
		fmt.Println("lista considerata ora: ", id)
		// Confronta gli ID dei peer in ordine lessicografico
		if strings.Compare(id, maxID) > 0 {
			maxID = id
			//log.Printf("LEADER: %s", maxID)
		}
		fmt.Println("chi è il migliore ora:  ", maxID)

	}
	for _, obj := range jsonListPeer {
		if obj.IDpeer == maxID {
			fmt.Println("TROVATO")
			LeaderPeerAddress = obj.Hostname + ":" + obj.Port
			log.Printf("ADDRESS LEADER: %s", LeaderPeerAddress)
			break // Se trovi l'oggetto corrispondente, esce dalla funzione
		}
	}
	fmt.Println("FINITA ELEZIONE ")
	comunicationResult()
}
func comunicationResult() {
	idPeerPosition := findMe() // ELEZIONE
	fmt.Println("HO chiamato connessione comunicazione risultato")

	findValidConnectionSendLeader(idPeerPosition)
	fmt.Println("AAAAAAA")
}
func findValidConnectionSendLeader(idPeerPosition int) {

	fmt.Println("SONO DENTROOOOO")
	var peerNext int = 0
	if idPeerPosition == len(jsonListPeer)-1 {
		peerNext = 0
	} else {
		peerNext = idPeerPosition + 1
	}
	fmt.Println("INVIO NOTIZIA, valore position", jsonListPeer[idPeerPosition].Port)
	fmt.Println("INVIO NOTIZIA, valore nextpeer prima loop :", jsonListPeer[peerNext].Port)

	for jsonListPeer[peerNext].Port != jsonListPeer[idPeerPosition].Port {
		fmt.Println("SONO DENTROOOOO LOOP")
		fmt.Println("INVIO NOTIZIA, porta posizione del peer nel for :", jsonListPeer[idPeerPosition].Port)

		var address string = jsonListPeer[peerNext].Hostname + ":" + jsonListPeer[peerNext].Port
		fmt.Printf("INVIO NOTIZIA address next to peer: %s\n", address)

		conn, err := connectToPeer(address)
		defer conn.Close()

		if err != nil {
			fmt.Printf("INVIO NOTIZIA Tentativo connessione peer successivo fallita: %v\n", err)

		} else {

			resultElection := &interfaccia.ResultElection{
				IdPeers: LeaderPeerAddress,
			}

			c := interfaccia.NewRingElectionClient(conn)
			_, err = c.SendResult(context.Background(), resultElection)
			if err == nil {
				fmt.Println("INVIO NOTIZIA")
				break
			} else {
				fmt.Printf("INVIO NOTIZIA fallita: %v\n", err)
			}
		}

		peerNext = peerNext + 1
		//fmt.Println("INVIO NOTIZIA posizione aggiornata  peerNext++ ", jsonListPeer[peerNext].Port)

		if peerNext >= len(jsonListPeer) { //oppure maggiore di len-1
			peerNext = 0
			fmt.Println("INVIO NOTIZIA correzione indice nell'anello: ", peerNext)
		}

	}
	fmt.Println("INVIO NOTIZIA fuori loop")
}
func checkleader() bool {
	conn, err := connectToPeer(LeaderPeerAddress)
	if err != nil {
		fmt.Println("did not connect: %v", err)
		return false
	}
	//defer conn.Close()
	c := interfaccia.NewRingElectionClient(conn)

	check := &interfaccia.CheckLeader{Check: "OK?"}

	response, err := c.SendCheck(context.Background(), check)
	if err != nil {
		fmt.Println("did not connect: %v", err)
		return false
	}
	if response != nil {
		return true
	} else {
		return false
	}

}

/////////////////////////////////////

///funzioni/procedure condivise

func connectToPeer(ipAddress string) (*grpc.ClientConn, error) {

	conn, err := grpc.Dial(ipAddress, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("errore durante la connessione a %s: %v", ipAddress, err)
	}

	return conn, nil
}

// /////////////// algorimo raft election
func (s *serverRaftElection) SendHeartBeat(ctx context.Context, req *interfaccia.HeartBeat) (*emptypb.Empty, error) {
	fmt.Println("ricevuto heartbeat")

	LeaderID = req.LeaderID
	Term = req.Term
	StopChannel <- ""
	fmt.Printf("ricevuto heartbeat leader %s e %d\n", LeaderID, Term)

	if req.Term >= Term {
		LeaderID = req.LeaderID
		Term = req.Term
		StopChannel <- ""

	} else {
		HeartbeatChannel <- ""
	}

	return &emptypb.Empty{}, nil
}
func (s *serverRaftElection) SendUpdate(ctx context.Context, item *interfaccia.JSONListPeer) (*emptypb.Empty, error) {
	// Sovrascrivi la jsonListPeer con la lista ricevuta dalla risposta del registro
	jsonListPeer = nil
	for _, jsonItem := range item.Items {
		jsonListPeer = append(jsonListPeer, Peer{
			Hostname: jsonItem.Hostname,
			Port:     jsonItem.Port,
			IDpeer:   jsonItem.IdPeer,
		})
	}
	fmt.Println("jsonListPeer aggiornata:", jsonListPeer)

	for _, item := range item.Items {

		address := fmt.Sprintf("%s:%s", item.Hostname, item.Port)
		fmt.Println("Hostname e porta:", address)
	}
	return &emptypb.Empty{}, nil
}
func (s *serverRaftElection) SendRaftElection(ctx context.Context, req *interfaccia.RequestElection) (*interfaccia.AnswerElection, error) {
	// Implementazione della logica di voto
	if req.Term > Term {
		Term = req.Term
		// Aggiorna la propria logica per rispondere alla richiesta di voto
		return &interfaccia.AnswerElection{Vote: "granted"}, nil
	}

	// Altrimenti, rifiuta il voto
	return &interfaccia.AnswerElection{Vote: "denied"}, nil
}

func randomGenerateTimerElection() {
	rand.Seed(time.Now().UnixNano())
	TimerElection = rand.Intn(5) + 4 // nell algortimo uficiale è tra 150ms e 320ms qui invece è tra 4 e 8 secondi per semplicità
	fmt.Println("stampa timer elezione: ", TimerElection)

}
func randomGenerateTimerHeartbeat() {
	rand.Seed(time.Now().UnixNano())
	TimerHeartBeat = rand.Intn(5) + 2 // tra 2 e 6, cosi che la probbiltà che scade prima il heartbeat nonè nulla perche l'intersenzione fra i due generatori non è nulla
	fmt.Println("stampa timer Heartbeat: ", TimerElection)
}
func sendRequestVoteToPeers() {
	fmt.Printf("Voglio vincere e invio richieste di voto\n")

	ElectionVote = ElectionVote + 1 //vota per se stesso
	for _, peer := range jsonListPeer {

		if peer.Hostname != HOSTNAME && peer.Port != PORT_EXPOSE { //quando sarà distribuito
			//if peer.Port != PORT_EXPOSE {// per il locale
			address := fmt.Sprintf("%s:%s", peer.Hostname, peer.Port)
			conn, err := connectToPeer(address)
			if err != nil {
				fmt.Printf("Errore durante la connessione a %s: %v\n", address, err)
				continue
			}
			defer conn.Close()

			client := interfaccia.NewRaftElectionClient(conn)
			response, err := client.SendRaftElection(context.Background(), &interfaccia.RequestElection{
				Vote: "can i be your leader?",
				Term: Term,
			})
			if err != nil {
				fmt.Printf("Errore durante l'invio della richiesta di voto a %s: %v\n", address, err)
				continue
			}

			// Aggiornare la logica per gestire la risposta dei peer
			if response.Vote == "granted" {
				ElectionVote = ElectionVote + 1 // Incrementa il conteggio dei voti ricevuti
			}
		}
	}
}
func sendHeartBeat() {
	for _, peer := range jsonListPeer {

		//if peer.Hostname!=HOSTNAME && peer.Port!=PORT_EXPOSE { //qunado sarà distribuito
		if peer.Port != PORT_EXPOSE {
			fmt.Printf("sono io il leader STO INVIANDO HEARTBEAT \n")
			address := fmt.Sprintf("%s:%s", peer.Hostname, peer.Port)
			conn, err := connectToPeer(address)
			if err != nil {
				fmt.Printf("Errore durante la connessione a %s: %v\n", address, err)
				continue
			}
			defer conn.Close()

			client := interfaccia.NewRaftElectionClient(conn)
			_, err = client.SendHeartBeat(context.Background(), &interfaccia.HeartBeat{
				Term:     Term,
				LeaderID: LeaderID,
			})
			if err != nil {
				fmt.Printf("Errore durante l'invio della richiesta di voto a %s: %v\n", address, err)
				continue
			}

		} else {
			HeartbeatChannel <- ""
		}
	}

}
func handleElectionMessages() {
	randomGenerateTimerElection()
	Term = Term + 1
	ticker := time.Tick(time.Duration(TimerElection) * time.Second)
	fmt.Printf("elezione %d\n", Term)
	sendRequestVoteToPeers()

	//se il timer scade prima interrompe tutto
	select {
	case <-ticker:
		if ElectionVote > (len(jsonListPeer) / 2) {
			LeaderID = ID_PEER
			fmt.Printf("Ho vinto le elezioni con %d , su %d possibili e azzero i voti", ElectionVote, len(jsonListPeer))
			ElectionVote = 0
			randomGenerateTimerHeartbeat()
			HeartbeatChannel <- ""
			return

		} else if LeaderID == "" {
			fmt.Println("Ci riprovo a vincere le elezioni e azzero i voti")
			ElectionVote = 0
			ElectionChannel <- ""
			return

		} else {

			fmt.Println("Ho perso le elezioni e azzero i voti")
			ElectionVote = 0
			FollowerChannel <- ""
			return

		}
	}
}
func handleHeartbeat() {
	ticker := time.Tick(time.Duration(TimerHeartBeat) * time.Second)
	//ticker := time.Tick(2 * time.Second)
	if LeaderID == ID_PEER {
		select {
		case <-ticker:
			fmt.Printf("condizioni attuali nel loop heart beat %s e %d\n", LeaderID, Term)
			if LeaderID == ID_PEER {
				sendHeartBeat()
			} else { //	LeaderID != ID_PEER || LeaderID != ""
				FollowerChannel <- ""
				return
			}
		case <-StopChannel:
			fmt.Println("mi hanno bloccato la handlHeartbeat")
			FollowerChannel <- ""
			return
		}

	}

}
func handleFollower() {

	randomGenerateTimerElection()
	ticker := time.Tick(time.Duration(TimerElection) * time.Second)
	//ticker := time.Tick(6 * time.Second)
	fmt.Println("Sono un Follower e sono in attesa di ricevere heartbeat")
	select {
	case <-ticker:
		fmt.Println("STOP non arriva nulla!!!")
		ElectionChannel <- ""
		return
	case <-StopChannel:
		fmt.Println("restart follow!!!")
		FollowerChannel <- ""
		return
	}

}
func genesiRaft() {
	for {
		select {
		case <-FollowerChannel:
			go handleFollower()
		case <-ElectionChannel:
			go handleElectionMessages()
		case <-HeartbeatChannel:
			go handleHeartbeat()
		}
	}
}

// //////////////////////////////////////
func main() {

	if ELEZIONE == "anello" {
		///duplicazione del codice
		lis, err := net.Listen("tcp", ":"+PORT_EXPOSE)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		s := grpc.NewServer()
		interfaccia.RegisterRingElectionServer(s, &serverRingElection{})
		log.Println("Server in ascolto sulla porta  :" + PORT_EXPOSE)

		go func() {
			if err := s.Serve(lis); err != nil {
				log.Fatalf("Failed to serve: %v", err)
			}
		}()

		go func() {

			conn, err := connectToPeer(ADDRESS_REGISTRY)
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			//defer conn.Close()
			c := interfaccia.NewRingElectionClient(conn)

			item := &interfaccia.JSONClientRegistration{
				Hostname: HOSTNAME,
				Port:     PORT_EXPOSE,
				IdPeer:   ID_PEER,
			}

			response, err := c.SendRequest(context.Background(), item)
			if err != nil {
				log.Fatalf("could not update JSON list: %v", err)
			}

			marshaler := jsonpb.Marshaler{Indent: "  "}
			jsonStr, err := marshaler.MarshalToString(response)
			if err != nil {
				log.Fatalf("could not marshal JSON list: %v", err)
			}

			fmt.Println("Updated JSON list received from serverKnowledgeRegistry:")
			fmt.Println(jsonStr)

			for _, jsonItem := range response.Items {
				jsonListPeer = append(jsonListPeer, Peer{
					Hostname: jsonItem.Hostname,
					Port:     jsonItem.Port,
					IDpeer:   jsonItem.IdPeer,
				})
			}
			fmt.Println("LISTA COPIATA ", jsonListPeer)

			for _, item := range response.Items {

				address := fmt.Sprintf("%s:%s", item.Hostname, item.Port)
				fmt.Println("Hostname e porta:", address)
			}
		}()
		// termina qua

		//elezione
		go func() {
			ticker := time.Tick(10 * time.Second) // Canale di time.Tick per ottenere un'indicazione ogni 5 secondi
			//restart := make(chan struct{})
			for {
				<-ticker

				fmt.Println("sto nel loop ELEZIONE")

				//leader Assente
				if LeaderPeerAddress == "" && len(jsonListPeer) > 1 {

					idPeerPosition := findMe() // ELEZIONE
					fmt.Println("HO chiamato connessione")
					var listPeer []string
					findValidConnection(idPeerPosition, listPeer) // lista vuota

				} // Attende il prossimo tick del timer

			}
		}()

		//controllo leader se è presente
		go func() {
			ticker := time.Tick(5 * time.Second) // Canale di time.Tick per ottenere un'indicazione ogni 5 secondi
			var alive bool
			for {
				<-ticker
				fmt.Println("sto nel loop CHECK")
				if LeaderPeerAddress != "" && len(jsonListPeer) > 1 {
					fmt.Println("CHECK leader: ", LeaderPeerAddress)
					alive = checkleader()

				} // Attende il prossimo tick del timer
				if alive == false {
					LeaderPeerAddress = ""
				}
			}
		}()

	} else if ELEZIONE == "raft" {
		///è presente una dublicazione del codice rispetto l'algortimo precedente
		lis, err := net.Listen("tcp", ":"+PORT_EXPOSE)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		s := grpc.NewServer()
		interfaccia.RegisterRaftElectionServer(s, &serverRaftElection{})
		log.Println("Server in ascolto sulla porta  :" + PORT_EXPOSE)

		go func() {
			if err := s.Serve(lis); err != nil {
				log.Fatalf("Failed to serve: %v", err)
			}
		}()

		go func() {

			conn, err := connectToPeer(ADDRESS_REGISTRY)
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			//defer conn.Close()
			c := interfaccia.NewRaftElectionClient(conn)

			item := &interfaccia.JSONClientRegistration{
				Hostname: HOSTNAME,
				Port:     PORT_EXPOSE,
				IdPeer:   ID_PEER,
			}

			response, err := c.SendRequest(context.Background(), item)
			if err != nil {
				log.Fatalf("could not update JSON list: %v", err)
			}

			marshaler := jsonpb.Marshaler{Indent: "  "}
			jsonStr, err := marshaler.MarshalToString(response)
			if err != nil {
				log.Fatalf("could not marshal JSON list: %v", err)
			}

			fmt.Println("Updated JSON list received from serverKnowledgeRegistry:")
			fmt.Println(jsonStr)

			for _, jsonItem := range response.Items {
				jsonListPeer = append(jsonListPeer, Peer{
					Hostname: jsonItem.Hostname,
					Port:     jsonItem.Port,
					IDpeer:   jsonItem.IdPeer,
				})
			}
			fmt.Println("LISTA COPIATA ", jsonListPeer)

			for _, item := range response.Items {

				address := fmt.Sprintf("%s:%s", item.Hostname, item.Port)
				fmt.Println("Hostname e porta:", address)
			}
		}()
		/// termina qua

		go genesiRaft()
		FollowerChannel <- ""

	} else {

		log.Println("credo che hai sbagliato a digitiare il nome dell'algoritmo di elezione,\n" +
			" la mia traccia B2 aveva come algortimi l'algortimo ad anello e raft")
	}

	select {}

}
