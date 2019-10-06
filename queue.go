package main

import (
	pb "blink-backend/blink"
	"blink-backend/database"
	"blink-backend/database/model"
)

type TaskQueue struct {
	receiveRequestQueue map[string][]pb.ReceiveRequest
	channelMap          map[string]chan int
}

const (
	Unknown = iota
	Grant
	Denial
)

func serialize(request pb.ReceiveRequest) string {
	return "req;" + request.GetNickname() + ";" + request.GetReceiverNickname() + ";" + request.GetUuid()
}

func (q *TaskQueue) SendRequest(sender string, receiver string, uuid string) int {
	db := database.GetInstance().DB
	var file model.File
	db.Where("uuid = ?", uuid).First(&file)
	request := pb.ReceiveRequest{
		Nickname:         sender,
		ReceiverNickname: receiver,
		Filename:         file.Filename,
		Uuid:             uuid,
	}

	c := make(chan int)

	q.receiveRequestQueue[receiver] = append(q.receiveRequestQueue[receiver], request)
	q.channelMap[serialize(request)] = c

	result := <-c
	return result
}

func (q *TaskQueue) ProcessGrant(request pb.ReceiveRequest) string {
	db := database.GetInstance().DB
	var file model.File
	db.Where("uuid = ?", request.Uuid).First(&file)
	if &file == nil {
		return ""
	}

	q.channelMap[serialize(request)] <- Grant
	return file.CraftFileLink()
}

func (q *TaskQueue) ProcessDenial(request pb.ReceiveRequest) {
	q.channelMap[serialize(request)] <- Denial
}
