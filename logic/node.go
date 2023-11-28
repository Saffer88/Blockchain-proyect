package logic

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/signal"
	"strings"
	"syscall"

	"os"

	"github.com/syndtr/goleveldb/leveldb"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

const protocolID = "/example/1.0.0"
const discoveryNamespace = "example"

func InitNode() {
	// Add -peer-address flag
	// Create the libp2p host.
	//
	// Note that we are explicitly passing the listen address and restricting it to IPv4 over the
	// loopback interface (127.0.0.1).
	//
	// Setting the TCP port as 0 makes libp2p choose an available port for us.
	// You could, of course, specify one if you like.
	host, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"), // Especifica la dirección de escucha
	)
	if err != nil {
		panic(err)
	}
	defer host.Close()
<<<<<<< HEAD
	hostG = host

	// Print this node's addresses and ID
	//fmt.Println("Addresses:", host.Addrs())
	//fmt.Println("ID:", host.ID())
=======

	// Print this node's addresses and ID
	fmt.Println("Addresses:", host.Addrs())
	fmt.Println("ID:", host.ID())
>>>>>>> 539ddc384b98a1512728b8c9b87789e5cca69d1d

	// Setup a stream handler.
	//
	// This gets called every time a peer connects and opens a stream to this node.
	host.SetStreamHandler(protocolID, func(s network.Stream) {
		//go writeCounter(s)
		go readTransaction(s)
	})

	// Setup peer discovery.
	s := mdns.NewMdnsService(host, discoveryNamespace, &discoveryNotifee{})
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
	defer s.Close()
<<<<<<< HEAD
=======
	hostG = host
>>>>>>> 539ddc384b98a1512728b8c9b87789e5cca69d1d

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGKILL, syscall.SIGINT)
	<-sigCh
}

var globalDB *leveldb.DB

func SetGlobalDB(db *leveldb.DB) {
	globalDB = db
}

func writeTransaction(s network.Stream, txJSON string) {
	writer := bufio.NewWriter(s)

	_, err := writer.WriteString(txJSON)
	if err != nil {
		fmt.Println("Error al escribir en el stream:", err)
		return
	}

	err = writer.Flush()
	if err != nil {
		fmt.Println("Error al hacer flush en el stream:", err)
		return
	}
}

func readTransaction(s network.Stream) error {
	buffer := make([]byte, 1024)
	var data []byte

	for {
		bytesRead, err := s.Read(buffer)
		if err != nil {
			fmt.Println("Error al leer desde el stream:", err)
			return err
		}

		if bytesRead > 0 {
			data = append(data, buffer[:bytesRead]...)

			if !containsSeparator(data) {

				msg := string(data)
				msg = strings.TrimSpace(msg)

				var jsonData map[string]interface{}
				if err := json.Unmarshal([]byte(msg), &jsonData); err != nil {
					fmt.Println("Error al decodificar el mensaje JSON:", err)
					return err
				}

				modeValue, ok := jsonData["Mode"]
				if !ok {
					fmt.Println("Campo 'mode' no encontrado en el mensaje JSON.")
					data = nil
					continue
				}

				mode, ok := modeValue.(float64)
				if !ok {
					fmt.Println("Valor no válido para 'mode' en el mensaje JSON.")
					data = nil
					continue
				}

				switch int(mode) {
				case 1:
<<<<<<< HEAD
					fmt.Println("Modo 1 - Mensaje JSON recibido:", msg)
=======
					//fmt.Println("Modo 1 - Mensaje JSON recibido:", msg)
>>>>>>> 539ddc384b98a1512728b8c9b87789e5cca69d1d

					var jsonData map[string]interface{}
					if err := json.Unmarshal([]byte(msg), &jsonData); err != nil {
						fmt.Println("Error al decodificar el mensaje JSON:", err)
						break
					}

					publicKey, publicKeyOK := jsonData["public key"].(string)
					addressHex, addressHexOK := jsonData["Address"].(string)

					if !publicKeyOK || !addressHexOK {
						fmt.Println("Valores 'public key' y/o 'Address' no válidos en el mensaje JSON.")

						break
					}

					err := SaveAccountToDB(addressHex, publicKey, 1000)
					if err != nil {
						fmt.Println("Error al guardar la cuenta en la base de datos:", err)
					}

				case 2:

					fmt.Println("Modo 2 - Mensaje JSON recibido:", msg)

<<<<<<< HEAD
					lastBlock, err := GetLastBlock(globalDB)
					if err != nil {
						fmt.Println("\nError al obtener el último bloque", err)
						panic(err)
					}

					var transactionData transactionJSON
					if err := json.Unmarshal([]byte(msg), &transactionData); err != nil {
						fmt.Println("Error al decodificar el mensaje JSON:", err)
						break
					}
					tx := &Transaction{
						Sender:    transactionData.Sender,
						Receiver:  transactionData.Receiver,
						Amount:    transactionData.Amount,
						Signature: transactionData.Signature,
						Nonce:     transactionData.Nonce,
					}
					AddTransaction(&lastBlock, *tx)
					UpdateBlockHash(&lastBlock)
					err = SaveBlockToDB(lastBlock, globalDB)
					if err != nil {
						panic(err)
					}
					err = UpdateBalance(transactionData.Sender, -transactionData.Amount)
					if err != nil {
						panic(err)
					}
					err = UpdateBalance(transactionData.Receiver, transactionData.Amount)
					if err != nil {
						panic(err)
					}

=======
>>>>>>> 539ddc384b98a1512728b8c9b87789e5cca69d1d
				case 3:

					fmt.Println("Modo 3 - Mensaje JSON recibido:", msg)

<<<<<<< HEAD
					lastBlock, err := GetLastBlock(globalDB)
					if err != nil {
						fmt.Println("\nError al obtener el último bloque", err)
						panic(err)
					}

					var transactionData transactionJSON
					if err := json.Unmarshal([]byte(msg), &transactionData); err != nil {
						fmt.Println("Error al decodificar el mensaje JSON:", err)
						break
					}
					tx := &Transaction{
						Sender:    transactionData.Sender,
						Receiver:  transactionData.Receiver,
						Amount:    transactionData.Amount,
						Signature: transactionData.Signature,
						Nonce:     transactionData.Nonce,
					}
					var transactions []Transaction
					transactions = append(transactions, *tx)

					newblock := GenerateBlock(globalDB, transactions, transactionData.Limit)
					err = SaveBlockToDB(newblock, globalDB)
					if err != nil {
						panic(err)
					}

					AddTransaction(&lastBlock, *tx)
					UpdateBlockHash(&lastBlock)
					err = SaveBlockToDB(lastBlock, globalDB)
					if err != nil {
						panic(err)
					}
					err = UpdateBalance(transactionData.Sender, -transactionData.Amount)
					if err != nil {
						panic(err)
					}
					err = UpdateBalance(transactionData.Receiver, transactionData.Amount)
					if err != nil {
						panic(err)
					}

=======
>>>>>>> 539ddc384b98a1512728b8c9b87789e5cca69d1d
				default:
					fmt.Println("Modo no válido:", mode)
				}

				data = nil
			}
		}
	}
}

func containsSeparator(data []byte) bool {
	for _, b := range data {
		if b == '\n' {
			return true
		}
	}
	return false
}

type discoveryNotifee struct {
	h host.Host
}

var (
	hostG host.Host

	newPeerFound bool

	peersMap = make(map[peer.ID]peer.AddrInfo)
)

func (n *discoveryNotifee) HandlePeerFound(peerInfo peer.AddrInfo) {
	//fmt.Println("found peer", peerInfo.String())
	peersMap[peerInfo.ID] = peerInfo
}

func Broadcast(txJSON string) {
	for peerID, peerInfo := range peersMap {
		if peerID == hostG.ID() {
			continue
		}

<<<<<<< HEAD
		err := hostG.Connect(context.Background(), peerInfo)
		if err != nil {
			//fmt.Println("Error al conectar con el peer", peerID, ":", err)
			continue
		}

		//fmt.Println("Conectado exitosamente con el peer", peerID)
=======
		// Conectar al par
		err := hostG.Connect(context.Background(), peerInfo)
		if err != nil {
			fmt.Println("Error al conectar con el peer", peerID, ":", err)
			continue
		}

		fmt.Println("Conectado exitosamente con el peer", peerID)
>>>>>>> 539ddc384b98a1512728b8c9b87789e5cca69d1d

		stream, err := hostG.NewStream(context.Background(), peerID, protocolID)
		if err != nil {
			fmt.Println("Error al crear stream:", err)
			continue
		}

		if stream == nil {
			fmt.Println("Stream es nil después de intentar crearlo.")
			continue
		}
		go writeTransaction(stream, txJSON)

	}
}
<<<<<<< HEAD

type transactionJSON struct {
	Mode      int     `json:"Mode"`
	Sender    string  `json:"sender"`
	Receiver  string  `json:"receiver"`
	Amount    float64 `json:"amount"`
	Signature string  `json:"signature"`
	Nonce     int     `json:"nonce"`
	Limit     int     `json:"limit"`
}

func NewTransactionNodes(sender, receiver string, amount float64, signature string, nonce int, limit int, mode int) error {

	jsonTransaction := transactionJSON{
		Mode:      mode,
		Sender:    sender,
		Receiver:  receiver,
		Amount:    amount,
		Nonce:     nonce,
		Signature: signature,
		Limit:     limit,
	}

	jsonBytes, err := json.Marshal(jsonTransaction)
	if err != nil {
		return err
	}

	Broadcast(string(jsonBytes))

	return nil

}
=======
>>>>>>> 539ddc384b98a1512728b8c9b87789e5cca69d1d
