package main

import (
	"fmt"

	"github.com/BurntSushi/toml"

	"openpose_dance/conf"
	"openpose_dance/model"
	"openpose_dance/util"
)

func scanFrames(cfg *conf.VideoInfo) {
	var (
		frames  = make(map[int]*model.PersonKeypoints)
		actions []*model.Action
		action  *model.Action
	)

	for i := 0; i <= cfg.FrameNo; i++ {
		file := fmt.Sprintf(cfg.JsonPathFormat, i)
		f := model.JsonToFrame(file)
		if len(f.People) == 0 { // 空白帧
			continue
		}

		p := new(model.PersonKeypoints)
		p.InitFromFrame(f)
		frames[i] = p

		if i < cfg.ActionLen { // 前置部分忽略
			continue
		}

		beforeFrame, ok := frames[i-cfg.ActionLen]
		if !ok {
			continue
		}
		fgap := p.FrameGap(beforeFrame)
		if fgap <= cfg.GapMin { // 动作幅度不够
			continue
		}
		if action == nil { // 初始化
			action = &model.Action{
				Begin: i - cfg.ActionLen,
				End:   i,
				Gap:   fgap,
			}
			continue
		}
		if i-action.End <= cfg.MergeInterval { // 合并连续区间
			action.End = i
			continue
		}
		// 结算上一区间
		action.Gap = frames[action.End].FrameGap(frames[action.Begin])
		action.SecondBegin = float64(action.Begin) / float64(25)
		action.SecondEnd = float64(action.End) / float64(25)
		actions = append(actions, action)
		action = nil // 重置
		i--          // 重新计算本帧
	}

	fmt.Println(fmt.Sprintf("Treat %d Frames, Actions Got %d", cfg.FrameNo, len(actions)))

	util.WriteWithIoutil(cfg.KeyFramesPath, model.ActionsToKeyFrames(frames, actions, cfg))
	util.WriteWithIoutil(cfg.ActionsPath, model.ActionsOutput(frames, actions, cfg))

}

var Conf *conf.Config

func main() {
	if _, err := toml.DecodeFile("config.toml", &Conf); err != nil {
		panic(err)
	}
	scanFrames(Conf.V1)
}
