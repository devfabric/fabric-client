package main

type PutInfo struct {
	SeqNum string `json:"seqnum"`
	Data   string `json:"data"`
}

type GetInfo struct {
	SeqNum string `json:"seqnum"`
}

const (
	EvPutInfo string = "EvPutInfo"
)
