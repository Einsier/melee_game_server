package codec

import (
	"log"
	pb "melee_game_server/api/proto"

	"google.golang.org/protobuf/proto"
)

func Decode(sourceData []byte) *pb.TopMessage {
	result := &pb.TopMessage{}
	if err := proto.Unmarshal(sourceData, result); err == nil {
		return result
	} else {
		log.Println(err)
		return nil
	}

}
