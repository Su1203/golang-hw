package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Открываем файл
	src, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("ошибка открытия источника: %w", err)
	}
	defer src.Close()

	info, err := src.Stat()
	if err != nil {
		return fmt.Errorf("ошибка Stat: %w", err)
	}
	fileSize := info.Size()

	if offset > fileSize {
		return fmt.Errorf("%w: %d (размер файла %d)", ErrOffsetExceedsFileSize, offset, fileSize)
	}

	toCopy := fileSize - offset
	if limit > 0 && limit < toCopy {
		toCopy = limit
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("ошибка создания копии: %w", err)
	}
	defer dst.Close()

	if _, err := src.Seek(offset, io.SeekStart); err != nil {
		return fmt.Errorf("ошибка Seek: %w", err)
	}

	bar := pb.Full.Start64(toCopy)
	bar.Set(pb.Bytes, true)

	limitedReader := io.LimitReader(src, toCopy)
	proxyWriter := bar.NewProxyWriter(dst)

	_, err = io.Copy(proxyWriter, limitedReader)

	bar.Finish() // Завершаем прогресс-бар

	if err != nil {
		return fmt.Errorf("error copy: %w", err)
	}

	return nil
}
