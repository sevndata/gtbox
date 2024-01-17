/*
Package gtbox_log Log工具
*/
package gtbox_log

import (
	"fmt"
	"github.com/george012/gtbox/gtbox_color"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

// GTLogStyle 日志样式
type GTLogStyle int

const (
	GTLogStyleDebug   GTLogStyle = iota // Debug
	GTLogStyleError                     // Error
	GTLogStyleWarning                   // Warning
	GTLogStyleInfo                      // Info
	GTLogStyleTrace                     // Trace
	GTLogStyleFatal                     // Fatal
)

func (aStyle GTLogStyle) String() string {
	switch aStyle {
	case GTLogStyleFatal:
		return "fatal"
	case GTLogStyleTrace:
		return "trace"
	case GTLogStyleInfo:
		return "info"
	case GTLogStyleWarning:
		return "warning"
	case GTLogStyleError:
		return "error"
	case GTLogStyleDebug:
		return "debug"
	default:
		return "debug"
	}
}

// GTLogSaveType 日志分片类型
type GTLogSaveType int

const (
	GTLogSaveTypeDays GTLogSaveType = iota //按日分片
	GTLogSaveHours                         //按小时分片
)

func (aFlag GTLogSaveType) String() string {
	switch aFlag {
	case GTLogSaveTypeDays:
		return "Days"
	case GTLogSaveHours:
		return "Hours"
	default:
		return "Unknown"
	}
}

type GTLog struct {
	mux               sync.RWMutex
	EnableSaveLogFile bool
	ProjectName       string
	LogLevel          GTLogStyle
	LogSaveMaxDays    int64
	LogSaveFlag       GTLogSaveType
	logDir            string
}

var (
	currentLog *GTLog
	gtLogOnce  sync.Once
)

func Instance() *GTLog {
	gtLogOnce.Do(func() {
		currentLog = &GTLog{}
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})

		logrus.SetLevel(logrus.TraceLevel)
		// 设置默认日志输出为控制台
		logrus.SetOutput(os.Stdout)
	})
	return currentLog
}

func GetProjectName() string {
	return Instance().ProjectName
}

func GetLogLevel() GTLogStyle {
	return Instance().LogLevel
}

func GetLogFilePath() string {
	return Instance().logDir
}

func (aLog *GTLog) infof(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Infof(format, args...)
}

func (aLog *GTLog) warnf(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Warnf(format, args...)
}

func (aLog *GTLog) errorf(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Errorf(format, args...)
}
func (aLog *GTLog) debugf(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Debugf(format, args...)
}
func (aLog *GTLog) tracef(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Tracef(format, args...)
}
func (aLog *GTLog) fatalf(format string, args ...interface{}) {
	aLog.mux.Lock()
	defer aLog.mux.Unlock()

	logrus.Fatalf(format, args...)
}

// logF 快捷日志Function，含模块字段封装
// Params [style] log类型  fatal、trace、info、warning、error、debug
// Params [format] 模块名称：自定义字符串
// Params [args...] 模块名称：自定义字符串
func logF(style GTLogStyle, format string, args ...interface{}) {
	colorFormat := format
	if Instance().EnableSaveLogFile != true {
		// 对每个占位符、非占位符片段和'['、']'进行迭代，为它们添加相应的颜色
		re := regexp.MustCompile(`(%[vTsdfqTbcdoxXUeEgGp]+)|(\[|\])|([^%\[\]]+)`)
		colorFormat = re.ReplaceAllStringFunc(format, func(s string) string {
			switch {
			case strings.HasPrefix(s, "%"):
				return gtbox_color.ANSIColorForegroundBrightYellow + s + gtbox_color.ANSIColorReset
			case s == "[" || s == "]":
				return s // 保持 `[` 和 `]` 的原始颜色
			default:
				if style == GTLogStyleError {
					return gtbox_color.ANSIColorForegroundBrightRed + s + gtbox_color.ANSIColorReset
				} else if style == GTLogStyleInfo {
					return gtbox_color.ANSIColorForegroundBrightGreen + s + gtbox_color.ANSIColorReset
				} else {
					return gtbox_color.ANSIColorForegroundBrightCyan + s + gtbox_color.ANSIColorReset
				}
			}
		})
	}

	if style != GTLogStyleInfo {
		pc, _, _, _ := runtime.Caller(2)
		fullName := runtime.FuncForPC(pc).Name()

		lastDot := strings.LastIndex(fullName, ".")
		if lastDot == -1 || lastDot == 0 || lastDot == len(fullName)-1 {
			return
		}
		callerClass := fullName[:lastDot]
		method := fullName[lastDot+1:]

		prefixFormat := fmt.Sprintf("[pkg--%s--][method--%s--] ", callerClass, method)
		colorFormat = prefixFormat + colorFormat
	}

	switch style {
	case GTLogStyleFatal:
		Instance().fatalf(colorFormat, args...)
	case GTLogStyleTrace:
		Instance().tracef(colorFormat, args...)
	case GTLogStyleInfo:
		Instance().infof(colorFormat, args...)
	case GTLogStyleWarning:
		Instance().warnf(colorFormat, args...)
	case GTLogStyleError:
		Instance().errorf(colorFormat, args...)
	case GTLogStyleDebug:
		Instance().debugf(colorFormat, args...)
	}
}

// LogInfof format格式化log--info信息
func LogInfof(format string, args ...interface{}) {
	logF(GTLogStyleInfo, format, args...)
}

// LogErrorf format格式化log--error信息
func LogErrorf(format string, args ...interface{}) {
	logF(GTLogStyleError, format, args...)
}

// LogDebugf format格式化log--debug信息
func LogDebugf(format string, args ...interface{}) {
	logF(GTLogStyleDebug, format, args...)
}

// LogTracef format格式化log--Trace信息
func LogTracef(format string, args ...interface{}) {
	logF(GTLogStyleTrace, format, args...)
}

// LogFatalf format格式化log--Fatal信息
func LogFatalf(format string, args ...interface{}) {
	logF(GTLogStyleFatal, format, args...)
}

// LogWarnf format格式化log--Warning信息
func LogWarnf(format string, args ...interface{}) {
	logF(GTLogStyleWarning, format, args...)
}

// SetupLogTools 初始化日志
func SetupLogTools(productName string, enableSaveLogFile bool, log_dir string, settingLogLeve GTLogStyle, logMaxSaveDays int64, logSaveType GTLogSaveType) {

	Instance().ProjectName = productName
	Instance().EnableSaveLogFile = enableSaveLogFile

	Instance().LogLevel = settingLogLeve
	switch settingLogLeve {
	case GTLogStyleFatal:
		logrus.SetLevel(logrus.FatalLevel)
	case GTLogStyleTrace:
		logrus.SetLevel(logrus.TraceLevel)
	case GTLogStyleInfo:
		logrus.SetLevel(logrus.InfoLevel)
	case GTLogStyleWarning:
		logrus.SetLevel(logrus.WarnLevel)
	case GTLogStyleError:
		logrus.SetLevel(logrus.ErrorLevel)
	case GTLogStyleDebug:
		logrus.SetLevel(logrus.DebugLevel)
	}

	Instance().LogSaveMaxDays = logMaxSaveDays
	Instance().LogSaveFlag = logSaveType
	//	设置Log
	Instance().logDir = log_dir
	if Instance().EnableSaveLogFile == true {
		if log_dir == "" {
			if runtime.GOOS == "linux" {
				Instance().logDir = "/var/log"
			} else {
				Instance().logDir = "./logs"
			}
		}

		log_file_path := Instance().logDir + "/" + strings.ToLower(Instance().ProjectName) + "/run" + "_" + Instance().ProjectName
		/* 日志轮转相关函数
		   `WithLinkName` 为最新的日志建立软连接
		   `WithRotationTime` 设置日志分割的时间，隔多久分割一次
		   WithMaxAge 和 WithRotationCount二者只能设置一个
		    `WithMaxAge` 设置文件清理前的最长保存时间
		    `WithRotationCount` 设置文件清理前最多保存的个数
		*/
		// 下面配置日志每隔 1 分钟轮转一个新文件，保留最近 3 分钟的日志文件，多余的自动清理掉。
		logRotaionFlag := time.Hour * 24

		if Instance().LogSaveFlag == GTLogSaveHours {
			logRotaionFlag = time.Hour
		} else if Instance().LogSaveFlag == GTLogSaveTypeDays {
			logRotaionFlag = time.Hour * 24
		}

		writer, _ := rotatelogs.New(
			log_file_path+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(log_file_path),
			rotatelogs.WithMaxAge(time.Duration(Instance().LogSaveMaxDays)*24*time.Hour),
			rotatelogs.WithRotationTime(logRotaionFlag),
		)
		logrus.SetOutput(writer)
	}
}
