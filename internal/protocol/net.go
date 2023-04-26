package protocol

import (
	"io"
	"net"
	"time"
)

func SendAllAndReceiveAll(conn net.Conn, toSend []byte, timeoutDuration time.Duration) (receiveBytes []byte, err error) {
	// Write email content
	remainingBytes := toSend
	for len(remainingBytes) > 0 {
		_ = conn.SetDeadline(time.Now().Add(timeoutDuration))
		n, err := conn.Write(remainingBytes)
		_ = conn.SetDeadline(time.Time{})
		if err != nil {
			return []byte(""), err
		}
		remainingBytes = remainingBytes[n:]
	}

	// Read response
	_ = conn.SetDeadline(time.Now().Add(timeoutDuration))
	receiveBytes, err = io.ReadAll(conn)
	_ = conn.SetDeadline(time.Time{})
	if err != nil {
		return []byte(""), err
	}

	return receiveBytes, nil
}
