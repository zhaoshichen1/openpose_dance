package model

import (
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
)

// 13个点对应openpose中点坐标的idx
var (
	_headIdx = 0
	_neckIdx = 1
	_lsIdx   = 2
	_lhIdx   = 4
	_rsIdx   = 5
	_rhIdx   = 7
	_hlpIdx  = 8
	_llIdx   = 9
	_lkIdx   = 10
	_lfIdx   = 23
	_rlIdx   = 12
	_rkIdx   = 13
	_rfIdx   = 20
)

// 帧信息，可含有多人信息
type Frame struct {
	People []*Person
}

// 人员信息，包含关键点信息和PersonID
type Person struct {
	PersonID []int     `json:"person_id"`
	Points   []float64 `json:"pose_keypoints_2d"`
}

func (p *PersonKeypoints) InitFromFrame(f *Frame) {
	// 每3个数字代表一个关键点，假设frame中只有一个人来处理
	for i := 0; i < len(f.People[0].Points); i += 3 {
		idx := i / 3
		x := f.People[0].Points[i]
		y := f.People[0].Points[i+1]
		kp := &Keypoint{X: x, Y: y}
		switch idx {
		case _headIdx:
			kp.Type = "head"
			p.Head = kp
		case _neckIdx:
			kp.Type = "neck"
			p.Neck = kp
		case _lsIdx:
			kp.Type = "lshoulder"
			p.LeftShoulder = kp
		case _lhIdx:
			kp.Type = "lhand"
			p.LeftHand = kp
		case _rsIdx:
			kp.Type = "rshoulder"
			p.RightShoulder = kp
		case _rhIdx:
			kp.Type = "rhand"
			p.RightHand = kp
		case _hlpIdx:
			kp.Type = "hp"
			p.Hlp = kp
		case _llIdx:
			kp.Type = "lleg"
			p.LeftLeg = kp
		case _lkIdx:
			kp.Type = "lknee"
			p.LeftKnee = kp
		case _lfIdx:
			kp.Type = "lfoot"
			p.LeftFoot = kp
		case _rlIdx:
			kp.Type = "rleg"
			p.RightLeg = kp
		case _rkIdx:
			kp.Type = "rknee"
			p.RightKnee = kp
		case _rfIdx:
			kp.Type = "rfoot"
			p.RightFoot = kp
		}
	}
}

// 读取json文件，转为frame结构
func JsonToFrame(path string) *Frame {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(path)
		panic(err)
	}
	f := new(Frame)
	if err := json.Unmarshal(content, f); err != nil {
		panic(err)
	}
	return f
}

// Draw 将帧数据上的关键点连成线，保存为图片数据，用于校验数据识别准确性
func Draw(frames map[int]*PersonKeypoints, pathFormat string, width, height int) {
	for k, v := range frames {
		func() {
			imgfile, _ := os.Create(fmt.Sprintf(pathFormat, k))
			defer imgfile.Close()

			img := image.NewNRGBA(image.Rect(0, 0, width, height))
			v.Draw(img)

			// 以PNG格式保存文件
			err := png.Encode(imgfile, img)
			if err != nil {
				panic(err)
			}
		}()
	}
}
