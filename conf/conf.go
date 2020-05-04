package conf

type VideoInfo struct {
	Name            string
	JsonPathFormat  string
	FrameNo         int     // 视频总帧数
	FPS             int     // frame per second
	CmPerCoordinate float64 // 每个坐标对应的厘米数
	ActionLen       int     // 动作最短持续时间
	GapMin          float64 // 动作最小坐标差
	MergeInterval   int     // 可合并的区间最大差距
	KeyFramesPath   string  // 关键帧保存地址
	ActionsPath     string  // 动作集保存地址
}

type Config struct {
	V1 *VideoInfo
	V2 *VideoInfo
}
