package model

import (
	"image"
	"image/color"

	"openpose_dance/util"
)

// 单人的13个关键点的坐标信息
type PersonKeypoints struct {
	Head          *Keypoint
	Neck          *Keypoint
	LeftShoulder  *Keypoint
	LeftHand      *Keypoint
	RightShoulder *Keypoint
	RightHand     *Keypoint
	Hlp           *Keypoint
	LeftLeg       *Keypoint
	LeftKnee      *Keypoint
	LeftFoot      *Keypoint
	RightLeg      *Keypoint
	RightKnee     *Keypoint
	RightFoot     *Keypoint
}

// 将结构体转成数组，方便关键帧输出
func (v *PersonKeypoints) ToSlice() []*Keypoint {
	return []*Keypoint{
		v.Head, v.Neck, v.LeftShoulder, v.LeftHand, v.RightShoulder, v.RightHand, v.Hlp,
		v.LeftLeg, v.LeftKnee, v.LeftFoot, v.RightLeg, v.RightKnee, v.RightFoot,
	}
}

// 13个点位的total gap
func (v *PersonKeypoints) FrameGap(input *PersonKeypoints) float64 {
	return v.Head.Gap(input.Head) + v.Neck.Gap(input.Neck) +
		v.LeftShoulder.Gap(input.LeftShoulder) + v.LeftHand.Gap(input.LeftHand) +
		v.RightShoulder.Gap(input.RightShoulder) + v.RightHand.Gap(input.RightHand) +
		v.Hlp.Gap(input.Hlp) + v.LeftLeg.Gap(input.LeftLeg) +
		v.LeftKnee.Gap(input.LeftKnee) + v.LeftFoot.Gap(input.LeftFoot) +
		v.RightLeg.Gap(input.RightLeg) + v.RightKnee.Gap(input.RightKnee) +
		v.RightFoot.Gap(input.RightFoot)
}

func (v *PersonKeypoints) DrawLine(k1, k2 *Keypoint, img *image.NRGBA) {
	x1, x2 := k1.X, k2.X
	y1, y2 := k1.Y, k2.Y

	rgba := color.RGBA{0, 255, 0, 255} // Green

	if util.Abs(int(x2-x1)) > util.Abs(int(y2-y1)) {
		if x2-x1 < 0 { // 倒过来
			x1, x2 = x2, x1
			y1, y2 = y2, y1
		}
		stepY := float64(y2-y1) / float64(x2-x1)
		for i := x1; i <= x2; i++ {
			img.Set(int(i), int(float64(y1)+float64(i-x1)*stepY), rgba)
		}
	} else {
		if y2-y1 < 0 { // 倒过来
			x1, x2 = x2, x1
			y1, y2 = y2, y1
		}
		stepX := float64(x2-x1) / float64(y2-y1)
		for i := y1; i <= y2; i++ {
			img.Set(int(float64(x1)+float64(i-y1)*stepX), int(i), rgba)
		}
	}

	red := color.RGBA{255, 0, 0, 255} // Red
	img.Set(int(x1), int(y1), red)    // 线的两头标出
	img.Set(int(x2), int(y2), red)
}

func (v *PersonKeypoints) Draw(img *image.NRGBA) {
	rgba := color.RGBA{50, 50, 50, 255}
	img.Set(int(v.Head.X), int(v.Head.Y), rgba)

	v.DrawLine(v.Head, v.Neck, img)
	v.DrawLine(v.Neck, v.LeftShoulder, img)
	v.DrawLine(v.Neck, v.RightShoulder, img)
	v.DrawLine(v.LeftShoulder, v.LeftHand, img)
	v.DrawLine(v.RightShoulder, v.RightHand, img)
	v.DrawLine(v.Neck, v.Hlp, img)

	v.DrawLine(v.Hlp, v.LeftLeg, img)
	v.DrawLine(v.LeftLeg, v.LeftKnee, img)
	v.DrawLine(v.LeftKnee, v.LeftFoot, img)

	v.DrawLine(v.Hlp, v.RightLeg, img)
	v.DrawLine(v.RightLeg, v.RightKnee, img)
	v.DrawLine(v.RightKnee, v.RightFoot, img)
}
