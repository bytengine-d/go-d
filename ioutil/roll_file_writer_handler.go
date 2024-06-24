package ioutil

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/bytengine-d/go-d/event"
	"github.com/bytengine-d/go-d/lang"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	RollDayFileOriginFileName     = "_roll_day_file_origin_file_name"
	RollDayFileOriginFileBaseName = "_roll_day_file_origin_file_base_name"
	RollDayFileOriginFileAbsPath  = "_roll_day_file_origin_file_abs_path"
	RollDayFileOriginFileExt      = "_roll_day_file_origin_file_ext"
	RollDayFileOriginDir          = "_roll_day_file_origin_dir"
)

const (
	RollDayFileEventDayChange  = "_roll_day_file_event_day_change"
	RollDayFileEventFileChange = "_roll_day_file_event_file_change"
)

const (
	RollDayFileDatetimeYearMonthDayFormat = "2006-01-02"
)

type RollDayFileOption func(wrapper *WrapperWriter, ep *event.EventGroup) error

func DefaultRollDay(outFilePath string) WrapperWriterHandler {
	return RollDayFileWriterHandler(outFilePath,
		RollDayFileCheckTime,
		RollDayFileMoveFileName(RollDayFileDatetimeYearMonthDayFormat),
		RollDayFileCompress,
		RollDayFileSetupFileWriter)
}

func RollDayWithDatetimeFormat(outFilePath, datetimeFormat string) WrapperWriterHandler {
	return RollDayFileWriterHandler(outFilePath,
		RollDayFileCheckTime,
		RollDayFileMoveFileName(datetimeFormat),
		RollDayFileCompress,
		RollDayFileSetupFileWriter)
}

func RollDayFileWriterHandler(outFilePath string, options ...RollDayFileOption) WrapperWriterHandler {
	ep := event.NewEventGroup()
	return func(wrapper *WrapperWriter) error {
		createdAt, originFileName, absOutFilePath, originFileBaseName, originFileExt, originFileDir, err := obtainFileInfo(outFilePath)
		if err != nil {
			return err
		}
		if options != nil && len(options) > 0 {
			for _, option := range options {
				err = option(wrapper, ep)
				if err != nil {
					return err
				}
			}
		}
		wrapper.Set(RollDayFileOriginFileName, originFileName)
		wrapper.Set(RollDayFileOriginFileBaseName, originFileBaseName)
		wrapper.Set(RollDayFileOriginFileAbsPath, absOutFilePath)
		wrapper.Set(RollDayFileOriginFileExt, originFileExt)
		wrapper.Set(RollDayFileOriginDir, originFileDir)

		now := time.Now()
		if createdAt.Day() != now.Day() {
			ep.Publish(RollDayFileEventDayChange)
		} else if err = rollDayFileSetupFileWriter(wrapper); err != nil {
			return err
		}
		return nil
	}
}

func obtainFileInfo(outFilePath string) (*time.Time, string, string, string, string, string, error) {
	outFileInfo, err := lang.RealFileInfo(outFilePath)
	if errors.Is(err, os.ErrNotExist) {
		err = nil
		_, err = os.Create(outFilePath)
		if err != nil {
			return nil, "", "", "", "", "", err
		}
		outFileInfo, err = lang.RealFileInfo(outFilePath)
	}
	if err != nil {
		return nil, "", "", "", "", "", err
	}
	absOutFilePath, err := filepath.Abs(outFilePath)
	if err != nil {
		return nil, "", "", "", "", "", err
	}
	originFileName := outFileInfo.Name()
	originFileBaseName := path.Base(originFileName)
	originFileExt := path.Ext(originFileName)
	originFileDir := path.Dir(absOutFilePath)
	outFileCreateTime := lang.FileCreateTime(outFileInfo)
	return &outFileCreateTime, originFileName, absOutFilePath, originFileBaseName, originFileExt, originFileDir, nil
}

func RollDayFileSetupFileWriter(wrapper *WrapperWriter, ep *event.EventGroup) error {
	return ep.RegisterAsyncSubscribe(RollDayFileEventFileChange, func() {
		_ = rollDayFileSetupFileWriter(wrapper)
	})
}

func rollDayFileSetupFileWriter(wrapper *WrapperWriter) error {
	val, has := wrapper.Get(RollDayFileOriginFileAbsPath)
	if !has {
		return errors.New("out file name not found.")
	}
	filePath := val.(string)
	outFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	wrapper.ChangeDelegate(outFile)
	return nil
}

func RollDayFileMoveFileName(datetimeFormat string) RollDayFileOption {
	return func(wrapper *WrapperWriter, ep *event.EventGroup) error {
		return ep.RegisterAsyncSubscribe(RollDayFileEventDayChange, func() {
			_ = rollDayFileMoveFileName(wrapper, ep, datetimeFormat)
		})
	}
}

func rollDayFileMoveFileName(wrapper *WrapperWriter, ep *event.EventGroup, datetimeFormat string) error {
	val, has := wrapper.Get(RollDayFileOriginFileAbsPath)
	if !has {
		return errors.New("out file name not found.")
	}
	filePath := val.(string)
	val, has = wrapper.Get(RollDayFileOriginDir)
	if !has {
		return errors.New("out file name not found.")
	}
	dir := val.(string)
	val, has = wrapper.Get(RollDayFileOriginFileBaseName)
	if !has {
		return errors.New("out file name not found.")
	}
	fileBaseName := val.(string)
	val, has = wrapper.Get(RollDayFileOriginFileExt)
	if !has {
		return errors.New("out file name not found.")
	}
	ext := val.(string)
	yesterday := time.Now().AddDate(0, 0, -1).Format(datetimeFormat)
	targetFileName := fmt.Sprintf("%s/%s-%s.%s", dir, fileBaseName, yesterday, ext)
	err := os.Rename(filePath, targetFileName)
	if err == nil {
		ep.Publish(RollDayFileEventFileChange, targetFileName)
	}
	return err
}

func RollDayFileCompress(wrapper *WrapperWriter, ep *event.EventGroup) error {
	return ep.RegisterAsyncSubscribe(RollDayFileEventFileChange, func(targetFileName string) {
		_ = rollDayFileCompress(wrapper, targetFileName)
	})
}

func rollDayFileCompress(wrapper *WrapperWriter, targetFileName string) error {
	val, has := wrapper.Get(RollDayFileOriginDir)
	if !has {
		return errors.New("out file name not found.")
	}
	dir := val.(string)
	fileBaseName := path.Base(targetFileName)
	targetGzFileName := path.Join(dir, fileBaseName+".tar.gz")
	f, err := os.Open(targetFileName)
	if err != nil {
		return errors.New("out file name not found.")
	}
	fi, err := os.Lstat(targetFileName)
	if err != nil {
		return errors.New("out file name not found.")
	}
	err = fileToTarGz(f, fi, targetGzFileName)
	if err != nil {
		return errors.New("out file name not found.")
	}
	return os.Remove(targetFileName)
}

func fileToTarGz(file *os.File, info os.FileInfo, fileName string) error {
	tarFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, os.FileMode(0644))
	if err != nil {
		return err
	}
	defer tarFile.Close()
	zr := gzip.NewWriter(tarFile)
	defer zr.Close()
	tw := tar.NewWriter(zr)
	defer tw.Close()
	header, err := tar.FileInfoHeader(info, file.Name())
	if err != nil {
		return err
	}
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}
	return nil
}

func RollDayFileCheckTime(wrapper *WrapperWriter, ep *event.EventGroup) error {
	go func() {
		for {
			nowTime := time.Now()
			nowTimeStr := nowTime.Format("2006-01-02")
			//使用Parse 默认获取为UTC时区 需要获取本地时区 所以使用ParseInLocation
			t2, _ := time.ParseInLocation("2006-01-02", nowTimeStr, time.Local)
			// 第二天零点时间戳
			next := t2.AddDate(0, 0, 1)
			after := next.UnixNano() - nowTime.UnixNano() - 1
			<-time.After(time.Duration(after) * time.Nanosecond)
			ep.Publish(RollDayFileEventDayChange)
		}
	}()
	return nil
}
