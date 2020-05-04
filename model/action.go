package model

import (
	"encoding/json"
	"math"

	"openpose_dance/conf"
	"openpose_dance/util"
)

// Action
type Action struct {
	Gap         float64 // 坐标差
	Begin       int     // 开始帧
	End         int
	SecondBegin float64 // 对应到秒
	SecondEnd   float64
}

type ActionFinalBig struct {
	Data  []*ActionOutput `json:"data"`
	Total int             `json:"total_frames"`
}

type ActionOutput struct {
	BeginFrame   int       `json:"begin_frame"`
	EndFrame     int       `json:"end_frame"`
	LeftHandAcc  *Keypoint `json:"left_hand_acc"`  // 左手加速度 m/s2
	RightHandAcc *Keypoint `json:"right_hand_acc"` // 右手加速度 m/s2
}

func coordinateGapToV(gap float64, cfg *conf.VideoInfo) (v float64) {
	gap = math.Abs(gap)
	frameLength := float64(1) / float64(cfg.FPS)          // 每帧时长
	return gap / frameLength * cfg.CmPerCoordinate * 0.01 // cm to meter
}

// 计算动作的速度，输出加速度值
func (v *Action) actionToAcceleration(frames map[int]*PersonKeypoints, cfg *conf.VideoInfo) (left, right *Keypoint) {
	kBegin, kBeginNext := frames[v.Begin], frames[v.Begin+1]
	kEnd, kEndNext := frames[v.End], frames[v.End+1]

	vBLeftX := coordinateGapToV(kBegin.LeftHand.X-kBeginNext.LeftHand.X, cfg)
	vBLeftY := coordinateGapToV(kBegin.LeftHand.Y-kBeginNext.LeftHand.Y, cfg)

	vBRightX := coordinateGapToV(kBegin.RightHand.X-kBeginNext.RightHand.X, cfg)
	vBRightY := coordinateGapToV(kBegin.RightHand.Y-kBeginNext.RightHand.Y, cfg)

	vELeftX := coordinateGapToV(kEnd.LeftHand.X-kEndNext.LeftHand.X, cfg)
	vELeftY := coordinateGapToV(kEnd.LeftHand.Y-kEndNext.LeftHand.Y, cfg)

	vERightX := coordinateGapToV(kEnd.RightHand.X-kEndNext.RightHand.X, cfg)
	vERightY := coordinateGapToV(kEnd.RightHand.Y-kEndNext.RightHand.Y, cfg)

	left = &Keypoint{Type: "acceleration"}
	right = &Keypoint{Type: "acceleration"}

	frameLength := float64(1) / float64(cfg.FPS)     // 每帧时长
	frameGap := float64(v.End-v.Begin) * frameLength // 帧差时长
	left.X = util.FloatLimit(math.Abs(vELeftX-vBLeftX) / frameGap)
	left.Y = util.FloatLimit(math.Abs(vELeftY-vBLeftY) / frameGap)
	right.X = util.FloatLimit(math.Abs(vERightX-vBRightX) / frameGap)
	right.Y = util.FloatLimit(math.Abs(vERightY-vBRightY) / frameGap)

	return
}

// 输出对应的action数据
func ActionsOutput(frames map[int]*PersonKeypoints, input []*Action, cfg *conf.VideoInfo) (res string) {

	var aos []*ActionOutput
	for _, v := range input {
		ao := &ActionOutput{}
		ao.BeginFrame = To20FPS(v.Begin, cfg.FPS)
		ao.EndFrame = To20FPS(v.End, cfg.FPS)
		ao.LeftHandAcc, ao.RightHandAcc = v.actionToAcceleration(frames, cfg)
		aos = append(aos, ao)
	}

	q := &ActionFinalBig{
		Data:  aos,
		Total: To20FPS(cfg.FrameNo, cfg.FPS),
	}
	str, _ := json.Marshal(q)
	return string(str)
}
