package comet

import (
	"github.com/33cn/chat33/proto"
	"github.com/33cn/chat33/result"
	"github.com/33cn/chat33/types"
	"github.com/33cn/chat33/utility"
)

func HttpPush(appId, userId, device string, msg []byte) error {
	var msgParse proto.Proto

	err := msgParse.FromBytes(msg)
	if err != nil {
		return result.NewError(result.MsgFormatError)
	}

	msgTime := utility.NowMillionSecond()
	dispatcher := createDispatcher(msgParse.GetChannelType(), &msgParse, userId, msgTime)
	if dispatcher == nil {
		return result.NewError(result.MsgFormatError)
	}

	err = dispatcher.intercept(appId)
	if err != nil {
		if err == types.ERR_REPEAT_MSG {
			err = nil
		}
		return err
	}

	err = dispatcher.appendLog()
	if err != nil {
		return err
	}

	err = dispatcher.pushMsg()
	if err != nil {
		return err
	}
	logParser.Info("rev msg", "raw", msgParse.GetMsg(), "userId", userId, "device", device)
	return nil
}
