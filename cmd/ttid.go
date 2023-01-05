package cmd

import (
	"log"

	"github.com/sony/sonyflake"
)

func fakeMachineID(uint16) bool {
	return true
}

func NextID() uint64 {
	// Sonyflake Id
	var st sonyflake.Settings
	st.CheckMachineID = fakeMachineID
	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		log.Fatal("New Sonyflake failed!")
	}

	id, err := sf.NextID()
	if err != nil {
		log.Fatal("NextID failed!")
	}
	return id
}
