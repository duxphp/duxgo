package chart

import (
	"encoding/json"
	"github.com/duxphp/duxgo/v2/util/function"
	"github.com/golang-module/carbon/v2"
	"github.com/jianfengye/collection"
	"github.com/spf13/cast"
	"time"
)

type Chart struct {
	way      string
	mode     string
	option   func(*map[string]any)
	height   string
	width    string
	zoom     bool
	toolbar  bool
	dataTime bool
	legend   *ChatLegend
	date     *ChatDate
	series   []map[string]any
	labels   *[]string
	data     []ChatData
	title    *ChatTitle
	subTitle *ChatTitle
}

type ChatLegend struct {
	status bool
	x      string
	y      string
}

type ChatDate struct {
	start  string
	stop   string
	split  string
	format string
}

type ChatTitle struct {
	text  string
	align string
}

type ChatData struct {
	name   string
	data   []map[string]any
	format string
}

// New 新建图表
func New(opt ...string) *Chart {
	way := ""
	if len(opt) > 0 {
		way = "date"
	}
	switch opt[0] {
	case "custom":
		way = "custom"
	default:
		way = "date"
	}

	return &Chart{
		way:    way,
		height: "200",
		width:  "100%",
		legend: &ChatLegend{
			status: false,
			x:      "right",
			y:      "top",
		},
		date: &ChatDate{
			split:  "day",
			format: "2006-01-02",
		},
	}
}

// Width 宽度
func (t *Chart) Width(size string) *Chart {
	t.width = size
	return t
}

// Height 高度
func (t *Chart) Height(size string) *Chart {
	t.height = size
	return t
}

// Title 标题
func (t *Chart) Title(title string, align ...string) *Chart {
	opt := "left"
	if len(align) > 0 {
		opt = align[0]
	}
	t.title = &ChatTitle{
		text:  title,
		align: opt,
	}
	return t
}

// SubTitle 子标题
func (t *Chart) SubTitle(title string, align ...string) *Chart {
	opt := "left"
	if len(align) > 0 {
		opt = align[0]
	}
	t.subTitle = &ChatTitle{
		text:  title,
		align: opt,
	}
	return t
}

// Zoom 放大缩小
func (t *Chart) Zoom(status bool) *Chart {
	t.zoom = status
	return t
}

// Toolbar 工具状态
func (t *Chart) Toolbar(status bool) *Chart {
	t.toolbar = status
	return t
}

// Legend 类型状态
func (t *Chart) Legend(status bool, x string, y string) *Chart {
	t.legend = &ChatLegend{
		status: status,
		x:      x,
		y:      y,
	}
	return t
}

// DataTime 时间轴功能
func (t *Chart) DataTime(status bool) *Chart {
	t.dataTime = status
	return t
}

// Date 时间设置
func (t *Chart) Date(start string, stop string, split string, format string) *Chart {
	t.date = &ChatDate{
		start:  start,
		stop:   stop,
		split:  split,
		format: format,
	}
	return t
}

// Data 数据设置
func (t *Chart) Data(name string, data []map[string]any) *Chart {

	t.data = append(t.data, ChatData{
		name: name,
		data: data,
	})
	return t
}

// Column 柱状图
func (t *Chart) Column(stacked ...bool) *Chart {
	t.mode = "bar"
	xType := "category"
	if t.dataTime {
		xType = "datetime"
	}
	t.option = func(options *map[string]any) {
		data := *options
		data["plotOptions"] = map[string]any{
			"bar": map[string]any{
				"columnWidth": "50%",
			},
		}
		data["dataLabels"] = map[string]any{
			"enabled": false,
		}
		data["fill"] = map[string]any{
			"opacity": 1,
		}
		data["xaxis"] = map[string]any{
			"type":       xType,
			"categories": t.labels,
		}
		if len(stacked) > 0 {
			data["chart"].(map[string]any)["stacked"] = stacked[0]
		}
		options = &data
	}
	return t
}

// Line 曲线图
func (t *Chart) Line() *Chart {
	t.mode = "line"
	xType := "category"
	if t.dataTime {
		xType = "datetime"
	}
	t.option = func(options *map[string]any) {
		data := *options
		data["dataLabels"] = map[string]any{
			"enabled": false,
		}
		data["fill"] = map[string]any{
			"opacity": 1,
		}
		data["xaxis"] = map[string]any{
			"type":       xType,
			"categories": t.labels,
		}
		data["stroke"] = map[string]any{
			"curve": "straight",
		}
		options = &data
	}
	return t
}

func (t *Chart) Render() string {
	t.processData()

	options := map[string]any{}
	chart := map[string]any{
		"id": "vuechart-" + function.RandString(5),
	}
	options["grid"] = map[string]any{
		"strokeDashArray": 4,
	}
	if t.title != nil {
		options["title"] = map[string]any{
			"text":  t.title.text,
			"align": t.title.align,
			"style": map[string]any{
				"fontSize":   "16px",
				"fontWeight": "normal",
			},
		}
	}

	if t.subTitle != nil {
		options["title"] = map[string]any{
			"text":  t.title.text,
			"align": t.title.align,
			"style": map[string]any{
				"fontSize":   "16px",
				"fontWeight": "normal",
			},
		}
	}
	if t.toolbar {
		chart["toolbar"] = map[string]any{
			"show":         true,
			"autoSelected": true,
		}
	} else {
		chart["toolbar"] = map[string]any{
			"show": false,
		}
	}
	if t.zoom {
		chart["zoom"] = map[string]any{
			"enabled":        true,
			"type":           "x",
			"autoScaleYaxis": false,
		}
	} else {
		chart["zoom"] = map[string]any{
			"enabled": false,
		}
	}

	if t.legend.status {
		options["legend"] = map[string]any{
			"show":            true,
			"position":        t.legend.y,
			"horizontalAlign": t.legend.x,
			"floating":        true,
			"offsetY":         0,
			"offsetX":         -5,
		}
	} else {
		options["legend"] = map[string]any{
			"show": false,
		}
	}

	options["chart"] = chart

	t.option(&options)

	jsonOptions, _ := json.Marshal(&options)
	jsonSeries, _ := json.Marshal(&t.series)
	return `<apexchart
              ref="chart"
              width="` + t.width + `"
              height="` + t.height + `"
              type="` + t.mode + `"
              :options='` + string(jsonOptions) + `'
              :series='` + string(jsonSeries) + `'
            ></apexchart>`
}

func (t *Chart) processData() {

	var labels []string

	// 数据类型为日期
	if t.way == "date" {
		currentTime := time.Now()
		start := currentTime.AddDate(0, 0, -7).Format("2006-01-02")
		stop := currentTime.Format("2006-01-02")
		if t.date.start != "" {
			start = t.date.start
		}
		if t.date.stop != "" {
			stop = t.date.stop
		}
		labels = t.splitDate(start, stop, t.date.format)
	}

	if t.way == "custom" {
		for _, datum := range t.data {
			for _, item := range datum.data {
				labels = append(labels, cast.ToString(item["label"]))
			}
		}

		collect := collection.NewStrCollection(labels)
		labels, _ = collect.Unique().ToStrings()
	}

	t.labels = &labels

	for _, datum := range t.data {
		group := map[string]float64{}

		for _, item := range datum.data {
			if t.way == "date" {
				tmpLabel := carbon.Parse(cast.ToString(item["label"])).Carbon2Time()
				group[tmpLabel.Format(t.date.format)] += cast.ToFloat64(item["value"])
			}
			if t.way == "custom" {
				group[cast.ToString(item["label"])] += cast.ToFloat64(item["value"])
			}
		}

		var tmpArr []float64
		for _, label := range labels {
			tmpArr = append(tmpArr, group[label])
		}

		t.series = append(t.series, map[string]any{
			"name": datum.name,
			"data": &tmpArr,
		})
	}
}

func (t *Chart) splitDate(beginDate, endDate, format string) []string {
	bDate := carbon.Parse(beginDate).Carbon2Time()
	eDate := carbon.Parse(endDate).Carbon2Time()
	day := int(eDate.Sub(bDate).Hours() / 24)
	list := make([]string, 0)
	list = append(list, bDate.Format(format))
	for i := 1; i < day; i++ {
		result := bDate.AddDate(0, 0, i)
		list = append(list, result.Format(format))
	}
	list = append(list, eDate.Format(format))
	return list
}
