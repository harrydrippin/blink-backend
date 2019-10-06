package main

import (
	pb "blink-backend/blink"
	"blink-backend/storage"
	"blink-backend/utils"
	"context"
	"github.com/pkg/errors"
	"io"
	"log"
)

type BlinkServer struct {
	pb.UnimplementedBlinkServer
	queue *TaskQueue
}

func (s *BlinkServer) CheckNickname(ctx context.Context, in *pb.Nickname) (*pb.NicknameResp, error) {
	log.Printf("CALL: CheckNickname")
	log.Printf("Nickname: %v", in.GetNickname())
	result, err := utils.CheckNickname(in.GetNickname())
	if err != nil {
		return &pb.NicknameResp{Result: false}, err
	}
	return &pb.NicknameResp{Result: result}, nil
}

func (s *BlinkServer) SubmitNickname(ctx context.Context, in *pb.Nickname) (*pb.NicknameResp, error) {
	log.Printf("CALL: SubmitNickname")
	log.Printf("Nickname: %v", in.GetNickname())
	result, err := utils.SubmitNickname(in.GetNickname())
	if err != nil {
		return &pb.NicknameResp{Result: false}, err
	}
	return &pb.NicknameResp{Result: result}, nil
}

func (s *BlinkServer) SetReceiverStream(ctx context.Context, in *pb.ReceiverInfo) (*pb.ReceiveRequest, error) {
	log.Printf("CALL: SetReceiverStream")
	log.Printf("Nickname: %v", in.GetNickname());
	log.Printf("Location: %v, %v", in.GetLocation().GetLatitude(), in.GetLocation().GetLongitude())

	location := in.GetLocation()
	name := in.GetNickname()

	err := utils.SetReceiverInfo(location, name)
	if err != nil {
		log.Printf("%v", err)
		return &pb.ReceiveRequest{
			Nickname:         "Nil",
			ReceiverNickname: "Nil",
			Filename:         "Nil",
			Uuid:             "Nil",
		}, err
	}

	length := len(s.queue.receiveRequestQueue[name])
	log.Printf("There are %d items on %v", length, name)
	if length != 0 {
		request := s.queue.receiveRequestQueue[name][0]
		s.queue.receiveRequestQueue[name] = s.queue.receiveRequestQueue[name][1:length]
		return &request, nil;
	}

	return &pb.ReceiveRequest{
		Nickname:         "Nil",
		ReceiverNickname: "Nil",
		Filename:         "Nil",
		Uuid:             "Nil",
	}, nil
}

func (s *BlinkServer) GetClientsByLocation(in *pb.Location, stream pb.Blink_GetClientsByLocationServer) error {
	log.Printf("CALL: GetClientsByLocation")
	log.Printf("Lat: %v", in.GetLatitude())
	log.Printf("Lng: %v", in.GetLongitude())

	lat := in.GetLatitude()
	lng := in.GetLongitude()

	result := utils.GetClientsByLocation(lat, lng)

	for _, item := range result {
		if err := stream.Send(&pb.Client{
			Nickname: item.Nickname,
			Location: &pb.Location{
				Latitude:  item.Latitude,
				Longitude: item.Longitude,
			},
			AwayDistance: item.AwayDistance,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *BlinkServer) GetClientsByName(in *pb.Nickname, stream pb.Blink_GetClientsByNameServer) error {
	log.Printf("CALL: GetClientsByName")
	log.Printf("Nickname: %v", in.GetNickname())

	nickname := in.GetNickname()

	result := utils.GetClientsByName(nickname)

	for _, item := range result {
		if err := stream.Send(&pb.Client{
			Nickname: item.Nickname,
			Location: &pb.Location{
				Latitude:  item.Latitude,
				Longitude: item.Longitude,
			},
			AwayDistance: item.AwayDistance,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *BlinkServer) UploadFileRequest(ctx context.Context, in *pb.UploadFileRequestReq) (*pb.UploadFileRequestResp, error) {
	log.Printf("CALL: UploadFileRequest")
	log.Printf("Nickname: %v", in.GetNickname())
	log.Printf("Filename: %v", in.GetFilename())

	result, err := storage.MakeFile(in.GetNickname(), in.GetFilename())
	if err != nil {
		return &pb.UploadFileRequestResp{
			Result: false,
			Uuid:   "",
		}, err
	}
	return &pb.UploadFileRequestResp{
		Result: true,
		Uuid:   result,
	}, nil
}

func (s *BlinkServer) UploadFile(stream pb.Blink_UploadFileServer) error {
	log.Printf("CALL: UploadFile")
	var uuid = ""
	for {
		buf, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				log.Printf("Got EOF")
				break
			}

			err = errors.Wrapf(err,
				"failed unexpectedly while reading chunks from stream")
			return err
		}

		uuid = buf.GetUuid()
		chunk := buf.GetChunk()
		err = storage.WriteToFile(uuid, chunk)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	err := stream.SendAndClose(&pb.UploadFileResp{
		Code: pb.UploadStatusCode_OK,
		Uuid: uuid,
	})

	if err != nil {
		return err
	}
	return nil
}

func (s *BlinkServer) SendRequest(ctx context.Context, in *pb.SendRequestReq) (*pb.SendRequestResp, error) {
	log.Printf("CALL: SendRequest")
	log.Printf("Nickname: %v", in.GetNickname())
	log.Printf("ReceiverNickname: %v", in.GetReceiverNickname())
	log.Printf("UUID: %v", in.GetUuid())

	nickname := in.GetNickname()
	receiver := in.GetReceiverNickname()
	uuid := in.GetUuid()

	reqResult := s.queue.SendRequest(nickname, receiver, uuid)

	var result bool
	if reqResult == Grant {
		result = true
	} else {
		result = false
	}

	return &pb.SendRequestResp{Result: result}, nil
}

func (s *BlinkServer) RespondGrant(ctx context.Context, in *pb.ReceiveRequest) (*pb.FileLink, error) {
	log.Printf("CALL: RespondGrant")
	log.Printf("Nickname: %v", in.GetNickname())
	log.Printf("Receiver: %v", in.GetReceiverNickname())
	log.Printf("Filename: %v", in.GetFilename())
	log.Printf("UUID: %v", in.GetUuid())

	link := s.queue.ProcessGrant(*in)

	return &pb.FileLink{Link: link}, nil
}

func (s *BlinkServer) RespondDenial(ctx context.Context, in *pb.ReceiveRequest) (*pb.Empty, error) {
	log.Printf("CALL: RespondDenial")
	log.Printf("Nickname: %v", in.GetNickname())
	log.Printf("Receiver: %v", in.GetReceiverNickname())
	log.Printf("Filename: %v", in.GetFilename())
	log.Printf("Uuid: %v", in.GetUuid())

	s.queue.ProcessDenial(*in)

	return &pb.Empty{}, nil
}

func (s *BlinkServer) MakeSpot(ctx context.Context, in *pb.MakeSpotReq) (*pb.MakeSpotResp, error) {
	log.Printf("CALL: MakeSpot")
	log.Printf("Nickname: %v", in.GetNickname())
	log.Printf("Location: %v, %v", in.GetLocation().GetLatitude(), in.GetLocation().GetLongitude())
	log.Printf("UUID: %v", in.GetUuid())

	id := utils.MakeSpot(in.GetNickname(), in.GetLocation().GetLatitude(), in.GetLocation().GetLongitude(), in.GetUuid())

	return &pb.MakeSpotResp{
		Id: id,
	}, nil
}

func (s *BlinkServer) GetSpotById(ctx context.Context, in *pb.GetSpotByIdReq) (*pb.Spot, error) {
	log.Printf("CALL: GetSpotById")
	log.Printf("ID: %v", in.GetId())

	spot := utils.GetSpotById(in.GetId())

	return &pb.Spot{
		Id:       uint32(spot.ID),
		Nickname: spot.Nickname,
		Location: &pb.Location{
			Latitude:  spot.Latitude,
			Longitude: spot.Longitude,
		},
		Uuid: spot.Uuid,
	}, nil

}

func (s *BlinkServer) GetSpotsByNickname(in *pb.Nickname, stream pb.Blink_GetSpotsByNicknameServer) error {
	log.Printf("CALL: GetSpotsByNickname")
	log.Printf("Nickname: %v", in.GetNickname())

	spots := utils.GetSpotsByNickname(in.GetNickname())

	for _, spot := range spots {
		if err := stream.Send(&pb.Spot{
			Id:       uint32(spot.ID),
			Nickname: spot.Nickname,
			Location: &pb.Location{
				Latitude: spot.Latitude,
				Longitude: spot.Longitude,
			},
			Uuid:     spot.Uuid,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *BlinkServer) GetSpotsByLocation(in *pb.Location, stream pb.Blink_GetSpotsByLocationServer) error {
	log.Printf("CALL: GetSpotsByLocation")
	log.Printf("Location: %v, %v", in.GetLatitude(), in.GetLongitude())

	spots := utils.GetSpotsByLocation(in.GetLatitude(), in.GetLongitude())

	for _, spot := range spots {
		if err := stream.Send(&pb.Spot{
			Id:       uint32(spot.ID),
			Nickname: spot.Nickname,
			Location: &pb.Location{
				Latitude: spot.Latitude,
				Longitude: spot.Longitude,
			},
			Uuid:     spot.Uuid,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *BlinkServer) GetFileFromSpot(ctx context.Context, in *pb.Spot) (*pb.FileLink, error) {
	log.Printf("CALL: GetFileFromSpot")
	log.Printf("ID: %v", in.GetId())
	log.Printf("Nickname: %v", in.GetNickname())

	link := utils.GetFileFromSpot(in.GetId())

	return &pb.FileLink{Link: link}, nil
}
