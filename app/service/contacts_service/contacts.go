package contacts_service

import (
	"goskeleton/app/model"
	"strconv"
)

type ContactsStruct struct {
	UserId   string
	FriendId string
	Status   int64
	Grouping string
}

func (m *ContactsStruct) GetContactsByBothId() (*model.Contacts, error) {
	u, err := model.GetContactsByBothId(m.UserId, m.FriendId)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (m *ContactsStruct) CreateNewContact() (*model.Contacts, error) {
	userId, _ := strconv.Atoi(m.UserId)
	frdId, _ := strconv.Atoi(m.FriendId)
	u, err := model.CreateNewContact(&model.Contacts{
		UserId:   int64(userId),
		FriendId: int64(frdId),
		Status:   m.Status,
		Grouping: m.Grouping,
	})
	if err != nil {
		return nil, err
	}
	if m.Status > 0 {
		_, err = model.CreateNewContact(&model.Contacts{
			UserId:   int64(frdId),
			FriendId: int64(userId),
			Status:   m.Status,
		})
		if err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (m *ContactsStruct) AcceptContact() (*model.Contacts, error) {
	userId, _ := strconv.Atoi(m.UserId)
	frdId, _ := strconv.Atoi(m.FriendId)
	u, err := model.Updates(&model.Contacts{
		UserId:   int64(userId),
		FriendId: int64(frdId),
	}, map[string]interface{}{
		"status": m.Status,
	})
	if err != nil {
		return nil, err
	}
	_, err = model.CreateNewContact(&model.Contacts{
		UserId:   int64(frdId),
		FriendId: int64(userId),
		Status:   m.Status,
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}
