package common

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)
type Transfer struct {
	Conn net.Conn
	Buf [8092]byte
}

func (tf *Transfer) ReadPkg() (message Message, err error) {
	n, err := tf.Conn.Read(tf.Buf[:4])
	if err != nil || n != 4 {
		fmt.Println("Receive Error", err)
		return
	}
	pkgLen := binary.BigEndian.Uint32(tf.Buf[:4])
	n, err = tf.Conn.Read(tf.Buf[:pkgLen])
	if uint32(n) != pkgLen || err != nil {
		fmt.Println("Receive Error", err)
		return
	}
	err = json.Unmarshal(tf.Buf[:pkgLen], &message)
	if err != nil {
		fmt.Println("Decode Error", err)
		return
	}
	return
}


func (tf *Transfer) WritePkg(message Message) (err error) {
	data, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Write Encode Message Error", err)
		return
	}
	pkgLen := len(data)
	pkgLenData := make([]byte, 4)
	binary.BigEndian.PutUint32(pkgLenData, uint32(pkgLen))
	_, err = tf.Conn.Write(pkgLenData)
	if err != nil {
		fmt.Println("Write pkgLenData Error", err)
		return
	}

	_, err = tf.Conn.Write(data)
	if err != nil {
		fmt.Println("Write Error", err)
		return
	}

	return
}
