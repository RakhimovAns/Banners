package banners

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Service struct {
	mu    sync.RWMutex
	items []*Banner
}

func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

var IDCounter int64 = 1

type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
	Image   string
}

func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}
	return nil, errors.New("item not found")
}

func (s *Service) Save(ctx context.Context, banner *Banner, id int64, request *http.Request) error {
	if id == 0 {
		banner.ID = IDCounter
		file, _, err := request.FormFile("image")
		if err != nil {
			log.Println(err)
			return err
		}
		defer file.Close()
		FileName := strconv.Itoa(int(IDCounter))
		banner.Image = FileName
		f, err := os.Create("./web/banners/" + FileName + "png")
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = io.Copy(f, file)
		if err != nil {
			log.Println(err)
			return err
		}
		IDCounter++
		s.items = append(s.items, banner)
		return nil
	} else if id != 0 {
		var ban *Banner
		for _, bann := range s.items {
			if bann.ID == id {
				ban = bann
			}
		}
		if ban == nil {
			return errors.New("item not found")
		}
		ban.Link = banner.Link
		ban.Title = banner.Title
		ban.Button = banner.Button
		ban.Content = banner.Content
	}
	return errors.New("can't save")
}

func (s *Service) All() []*Banner {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*Banner, len(s.items))
	copy(result, s.items)
	return result
}
func (s *Service) DeleteByID(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, banner := range s.items {
		if banner.ID == id {
			copy(s.items[i:], s.items[i+1:])
			s.items = s.items[:len(s.items)-1]
			return nil
		}
	}
	return errors.New("item not found")
}
