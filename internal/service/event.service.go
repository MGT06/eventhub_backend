package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MGT06/EventHub_Backend.git/internal/dto"
	"github.com/MGT06/EventHub_Backend.git/internal/model"
	"github.com/MGT06/EventHub_Backend.git/internal/repo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventService struct {
	er *repo.EventRepo
	db *pgxpool.Pool
}

func NewEventService(er *repo.EventRepo, db *pgxpool.Pool) *EventService {
	return &EventService{
		er: er,
		db: db,
	}
}

func (e *EventService) GetEvents(ctx context.Context, eventId int) (dto.Event, error) {
	res, err := e.er.GetEvents(ctx, e.db, eventId)
	if err != nil {
		return dto.Event{}, err
	}

	speaker, er := e.er.GetSpeakers(ctx, e.db, eventId)
	if er != nil {
		return dto.Event{}, er
	}

	discuss, err := e.er.GetDiscussions(ctx, e.db, eventId)

	byID := make(map[int]dto.DiscussionEvent, len(discuss))

	for _, v := range discuss {
		author := dto.Author{
			Id:        v.AccountId,
			Name:      v.Name,
			AvatarURL: v.Avatar_url,
		}

		byID[v.DiscussId] = dto.DiscussionEvent{
			Id:        v.DiscussId,
			EventId:   v.EventId,
			ParentId:  v.ParentId,
			Message:   v.Message,
			Author:    author,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Replies:   []dto.DiscussionEvent{},
		}
	}

	var roots []dto.DiscussionEvent
	for _, val := range byID {
		if val.ParentId != nil {
			if parent, ok := byID[*val.ParentId]; ok {
				parent.Replies = append(parent.Replies, val)
				continue
			}
		}
		roots = append(roots, val)
	}

	category := strings.Split(res.Category_name, ", ")
	organizer := dto.Organizer{
		Name:       res.Name,
		Position:   res.Position,
		Avatar_url: res.Avatar_url,
	}

	var speakers []dto.Speaker
	for _, v := range speaker {
		speakers = append(speakers, dto.Speaker{
			Name:       v.Name,
			Position:   v.Position,
			Avatar_url: v.Avatar_url,
		})
	}

	data := dto.Event{
		Id:          res.Id_event,
		Organizer:   organizer,
		Title:       res.Title,
		Description: res.Description,
		Image_url:   res.Image_event_url,
		Start_at:    res.Start_at,
		End_at:      res.End_at,
		Location:    res.City,
		Capacity:    res.Capacity,
		Attendees:   res.Attendees,
		Speakers:    speakers,
		Discuss:     roots,
		Categories:  category,
	}

	return data, nil
}

func (e *EventService) GetEventBySearchFilter(ctx context.Context, search string, category string, location string) ([]dto.Event, error) {
	res, err := e.er.GetEventBySearchFilter(ctx, e.db, search, category, location)

	if err != nil {
		return nil, err
	}

	data := make([]dto.Event, 0, len(res))
	for _, v := range res {
		category := strings.Split(v.Category_name, ", ")

		data = append(data, dto.Event{
			Id:          v.Id_event,
			Title:       v.Title,
			Description: v.Description,
			Image_url:   v.Image_event_url,
			Start_at:    v.Start_at,
			End_at:      v.End_at,
			Location:    v.City,
			Capacity:    v.Capacity,
			Attendees:   v.Attendees,
			Categories:  category,
		})
	}

	return data, nil
}

func (e *EventService) ToggleJoinEvent(ctx context.Context, idUser int, idEvent int) (bool, error) {
	res, err := e.er.IsCanJoin(ctx, e.db, idUser, idEvent)
	if err != nil {
		return res.IsJoin, err
	}

	
	if res.IsJoin {
		if err := e.er.LeaveEvent(ctx, e.db, idUser, idEvent); err != nil {
			return res.IsJoin, err
		}
		return res.IsJoin, nil
	}
	
	if (res.Cap - res.Attendess) == 0 {
		return res.IsJoin, fmt.Errorf("Event is fuul")
	}

	if err := e.er.JoinEvent(ctx, e.db, idUser, idEvent); err != nil {
		return res.IsJoin, err
	}

	return res.IsJoin, nil
}

func (e *EventService) ToggleSavedEvent(ctx context.Context, idUser int, idEvent int) (bool, error) {
	isSaved, err := e.er.IsSaved(ctx, e.db, idUser, idEvent)
	if err != nil {
		return false, err
	}

	if isSaved {
		if err := e.er.UnSavedEvent(ctx, e.db, idUser, idEvent); err != nil {
			return false, err
		}
		return isSaved, nil
	}

	if err := e.er.SavedEvent(ctx, e.db, idUser, idEvent); err != nil {
		return false, err
	}

	return isSaved, nil
}

func (e *EventService) GetUpComingEvents(ctx context.Context) ([]dto.Event, error) {
	res, err := e.er.GetUpComingEvents(ctx, e.db)
	if err != nil {
		return nil, err
	}

	data := make([]dto.Event, 0, len(res))
	for _, v := range res {
		category := strings.Split(v.Category_name, ", ")

		data = append(data, dto.Event{
			Id:          v.Id_event,
			Title:       v.Title,
			Description: v.Description,
			Image_url:   v.Image_event_url,
			Start_at:    v.Start_at,
			End_at:      v.End_at,
			Location:    v.City,
			Capacity:    v.Capacity,
			Attendees:   v.Attendees,
			Categories:  category,
		})
	}

	return data, nil
}

func (e *EventService) GetMyEvent(ctx context.Context, idUser int) ([]dto.Event, error) {
	res, err := e.er.GetMyEvent(ctx, e.db, idUser)
	if err != nil {
		return nil, err
	}

	data := make([]dto.Event, 0, len(res))
	for _, v := range res {
		category := strings.Split(v.Category_name, ", ")

		data = append(data, dto.Event{
			Id:          v.Id_event,
			Title:       v.Title,
			Description: v.Description,
			Image_url:   v.Image_event_url,
			Start_at:    v.Start_at,
			End_at:      v.End_at,
			Location:    v.City,
			Capacity:    v.Capacity,
			Attendees:   v.Attendees,
			Categories:  category,
		})
	}

	return data, nil
}

func (e *EventService) AddEvent(ctx context.Context, body dto.AddEvent, userId int, imgPath string) error {
	// speakers := "[]"
	// if body.Speakers != "" {
	// 	var list []dto.Speaker
	// 	if err := json.Unmarshal([]byte(body.Speakers), &list); err != nil {
	// 		return err
	// 	}
	// 	res, err := json.Marshal(list)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	speakers = string(res)
	// }

	tx, err := e.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Println(err.Error())
		}
	}()

	eventId, err := e.er.AddEvent(ctx, tx, model.Event{
		Organizer_id:      userId,
		Community_id:      body.CommunityId,
		Title:             body.Title,
		Description:       body.Description,
		Image_event_url:   imgPath,
		Start_at:          body.Start_at,
		End_at:            body.End_at,
		Format:            body.Format,
		Location_event_id: body.LocationId,
		Capacity:          body.Capacity,
	})
	if err != nil {
		return err
	}

	if err := e.er.AddCategory(ctx, tx, eventId, body.Categories); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (e *EventService) EditEvent(ctx context.Context, eventId int, body dto.EditEvent, userId int, imgPath string) error {
	// speakers := "[]"
	// if body.Speakers != nil {
	// 	var list []dto.Speaker
	// 	if err := json.Unmarshal([]byte(*body.Speakers), &list); err != nil {
	// 		return err
	// 	}
	// 	res, err := json.Marshal(list)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	speakers = string(res)
	// }

	tx, err := e.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Println(err.Error())
		}
	}()

	oldImg, err := e.er.GetImageForUpdate(ctx, tx, eventId, userId)
	if err != nil {
		return err
	}

	imageURL := oldImg
	if imgPath != "" {
		imageURL = imgPath
	}

	if err := e.er.UpdateEvent(ctx, tx, eventId, userId, model.Event{
		Community_id:      body.CommunityId,
		Location_event_id: *body.LocationId,
		Title:             *body.Title,
		Description:       *body.Description,
		Image_event_url:   imageURL,
		Start_at:          *body.Start_at,
		End_at:            *body.End_at,
		Format:            *body.Format,
		Capacity:          *body.Capacity,
	}); err != nil {
		return err
	}

	if err := e.er.DeleteCategories(ctx, tx, eventId); err != nil {
		return err
	}

	if err := e.er.AddCategory(ctx, tx, eventId, body.Categories); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	if imgPath != "" && oldImg != "" {
		if err := os.Remove(oldImg); err != nil {
			log.Println(err.Error())
		}
	}

	return nil
}
