package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Открываем исходный файл
	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Получаем информацию о файле
	info, err := src.Stat()
	if err != nil {
		return err
	}

	// Проверяем, что это обычный файл
	if !info.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	// Проверяем offset
	fileSize := info.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	// Создаем целевой файл
	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Устанавливаем позицию чтения
	_, err = src.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	// Определяем, сколько байт копировать
	var reader io.Reader = src
	var bytesToCopy int64
	if limit > 0 {
		bytesToCopy = limit
		if offset+limit > fileSize {
			bytesToCopy = fileSize - offset
		}
		reader = io.LimitReader(src, bytesToCopy)
	}

	// progress bar
	bar := pb.Default.Start64(bytesToCopy)
	barReader := bar.NewProxyReader(reader)
	defer bar.Finish()

	// Копируем данные
	_, err = io.Copy(dst, barReader)
	return err
}
