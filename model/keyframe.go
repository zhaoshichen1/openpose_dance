package model

import (
	"encoding/json"
	"math"

	"openpose_dance/conf"
)

// To20FPS 强转为20FPS的帧号
func To20FPS(no, fps int) int {
	return int(math.Floor(float64(no) / float64(fps) * float64(20)))
}

type KeyOutput struct {
	No        int         `json:"no"`
	Positions []*Keypoint `json:"positions"`
}

type KeyFinalBigData struct {
	Data  []*KeyOutput `json:"data"`
	Total int          `json:"total_frames"`
}

// 入参一堆actions，出参都是关键帧数据
func ActionsToKeyFrames(frames map[int]*PersonKeypoints, input []*Action, cfg *conf.VideoInfo) (res string) {
	var k []*KeyOutput
	for _, v := range input {
		q := new(KeyOutput)
		end := frames[v.End]
		q.No = To20FPS(v.End, cfg.FPS) // 输出时才强转
		q.Positions = end.ToSlice()
		k = append(k, q)
	}
	finalBig := &KeyFinalBigData{ // 25帧强转20帧
		Total: To20FPS(cfg.FrameNo, cfg.FPS),
		Data:  k,
	}
	str, _ := json.Marshal(finalBig)
	return string(str)
}
