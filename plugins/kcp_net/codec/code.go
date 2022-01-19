package codec

import (
	"log"
	pb "melee_game_server/api/proto"

	"google.golang.org/protobuf/proto"
)

func Code(msg *pb.TopMessage) []byte {

	if result, err := proto.Marshal(msg); err == nil {
		return result
	} else {
		log.Println(err)
		return nil
	}
}
