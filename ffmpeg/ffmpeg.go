package ffmpeg

type VideoDetail struct {
	FileName        string  `json:"file_name,omitempty"`        //文件名
	FileType        int     `json:"file_type,omitempty"`        //文件类型
	VideoDuration   int64   `json:"video_duration,omitempty"`   //视频录制时长
	SubPath         string  `json:"sub_path,omitempty"`         //视频子目录
	FileSize        int64   `json:"file_size,omitempty"`        //文件大小
	VideoBirate     int64   `json:"video_birate,omitempty"`     //视频码率
	VideoFps        int64   `json:"video_fps,omitempty"`        //视频帧率
	VideoHeight     int64   `json:"video_height,omitempty"`     //视频高度
	VideoWidth      int64   `json:"video_width,omitempty"`      //视频宽度
	VideoStartTime  int64   `json:"video_start_time,omitempty"` //视频开始时间
	VideoEndTime    int64   `json:"video_end_time,omitempty"`   //视频结束时间
	IsFreezeFrame   int64   `json:"is_freeze_frame,omitempty"`  //是否定格资源
	VideoDetailId   string  `json:"video_detail_id,omitempty"`  //videoDetailID
	DetectionResult string  `json:"detection_result,omitempty"` //detectionResult
	TsDuration      float32 `json:"ts_duration,omitempty"`      //单片ts录制时长
}
