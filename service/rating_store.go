package service

import "sync"

type Rating struct {
	Count uint32
	Sum   float64
}
type RatingStore interface {
	Add(laptopID string, score float64) (*Rating, error)
}

type InMemoryRatingStore struct {
	mutex  sync.Mutex
	rating map[string]*Rating
}

func NewInMemoryRatingStore() *InMemoryRatingStore {
	return &InMemoryRatingStore{
		rating: make(map[string]*Rating),
	}
}

func (store *InMemoryRatingStore) Add(laptopID string, score float64) (*Rating, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	rating := store.rating[laptopID]
	if rating == nil {
		rating = &Rating{
			Count: 1,
			Sum:   score,
		}
	} else {
		rating.Count++
		rating.Sum += score
	}
	store.rating[laptopID] = rating
	return rating, nil
}
