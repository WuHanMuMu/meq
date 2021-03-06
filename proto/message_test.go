package proto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var mockMsgs = []*PubMsg{
	&PubMsg{[]byte("001"), []byte("test"), []byte("hello world1"), false, 0, 1, 100},
	&PubMsg{[]byte("002"), []byte("test"), []byte("hello world2"), false, 0, 1, 101},
	&PubMsg{[]byte("003"), []byte("test"), []byte("hello world3"), false, 0, 1, 102},
	&PubMsg{[]byte("004"), []byte("test"), []byte("hello world4"), true, 0, 1, 100},
	&PubMsg{[]byte("005"), []byte("test"), []byte("hello world5"), false, 0, 1, 100},
	&PubMsg{[]byte("006"), []byte("test"), []byte("hello world6"), false, 0, 1, 100},
	&PubMsg{[]byte("007"), []byte("test"), []byte("hello world7"), false, 0, 1, 100},
	&PubMsg{[]byte("008"), []byte("test"), []byte("hello world8"), false, 0, 1, 100},
	&PubMsg{[]byte("009"), []byte("test"), []byte("hello world9"), true, 0, 1, 100},
	&PubMsg{[]byte("010"), []byte("test"), []byte("hello world10"), false, 0, 1, 100},
	&PubMsg{[]byte("011"), []byte("test"), []byte("hello world11"), false, 0, 1, 100},
}

func TestPubMsgOnePackUnpack(t *testing.T) {
	packed := PackMsg(mockMsgs[0])
	unpacked, err := UnpackMsg(packed[1:])

	assert.NoError(t, err)
	assert.Equal(t, mockMsgs[0], unpacked)
}

func TestPubMsgsPackUnpack(t *testing.T) {
	packed := PackPubBatch(mockMsgs, MSG_PUB_BATCH)
	unpacked, err := UnpackPubBatch(packed[1:])

	assert.NoError(t, err)
	assert.Equal(t, mockMsgs, unpacked)
}

func TestSubPackUnpack(t *testing.T) {
	topic := []byte("test")

	packed := PackSub(topic)
	ntopic := UnpackSub(packed[5:])

	assert.Equal(t, topic, ntopic)

}

func TestMarkReadPackUnpack(t *testing.T) {
	topic := []byte("/1234567890/1/test/a")
	var msgids [][]byte
	for _, m := range mockMsgs {
		msgids = append(msgids, m.ID)
	}

	packed := PackMarkRead(topic, msgids)
	utopic, umsgids := UnpackMarkRead(packed[1:])

	assert.Equal(t, msgids, umsgids)
	assert.Equal(t, topic, utopic)
}

func TestMsgCountPackUnpack(t *testing.T) {
	count := 10923

	packed := PackMsgCount(count)
	ncount := UnpackMsgCount(packed[1:])

	assert.Equal(t, count, ncount)
}

func TestPullPackUnpack(t *testing.T) {
	msgid := []byte("00001")
	count := 11123

	packed := PackPullMsg(count, msgid)
	ncount, nmsgid := UnPackPullMsg(packed[1:])

	assert.Equal(t, count, ncount)
	assert.Equal(t, msgid, nmsgid)
}

func TestTimerMsgPackUnpack(t *testing.T) {
	tmsg := &TimerMsg{[]byte("0001"), []byte("test"), []byte("timer msg emit!"), time.Now().Unix(), 10}
	packed := PackTimerMsg(tmsg, MSG_PUB_TIMER)
	unpacked := UnpackTimerMsg(packed[5:])

	assert.Equal(t, tmsg, unpacked)
}

func TestSubAckPackUnpack(t *testing.T) {
	tp := []byte("test")
	packed := PackSubAck(tp)
	unpacked := UnpackSubAck(packed[5:])

	assert.Equal(t, tp, unpacked)
}

func TestPackAckCount(t *testing.T) {
	count := MAX_PULL_COUNT

	packed := PackReduceCount(count)
	ucount := UnpackReduceCount(packed[1:])

	assert.Equal(t, count, ucount)

	count = REDUCE_ALL_COUNT

	packed = PackReduceCount(count)
	ucount = UnpackReduceCount(packed[1:])

	assert.Equal(t, count, ucount)
}

func TestPackPresenceUsers(t *testing.T) {
	users := [][]byte{[]byte("a1"), []byte("a2"), []byte("a3")}
	packed := PackPresenceUsers(users)
	ousers := UnpackPresenceUsers(packed[1:])

	assert.EqualValues(t, users, ousers)
}

func TestPackJoinChat(t *testing.T) {
	topic := []byte("/1234567890/12/test/a")
	packet := PackJoinChat(topic)
	utopic := UnpackJoinChat(packet[1:])

	assert.Equal(t, topic, utopic)
}

func TestPackLeaveChat(t *testing.T) {
	topic := []byte("/1234567890/12/test/a")
	packet := PackLeaveChat(topic)
	utopic := UnpackLeaveChat(packet[1:])

	assert.Equal(t, topic, utopic)
}

func TestPackJoinChatNotify(t *testing.T) {
	topic := []byte("/1234567890/12/test/a")
	user := []byte("sunface")

	packet := PackJoinChatNotify(topic, user)
	utopic, uuser := UnpackJoinChatNotify(packet[1:])

	assert.Equal(t, topic, utopic)
	assert.Equal(t, user, uuser)
}

func BenchmarkPubMsgPack(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		PackPubBatch(mockMsgs, MSG_PUB_BATCH)
	}
}

func BenchmarkPubMsgUnpack(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	packed := PackPubBatch(mockMsgs, MSG_PUB_BATCH)
	for i := 0; i < b.N; i++ {
		UnpackPubBatch(packed[5:])
	}
}
