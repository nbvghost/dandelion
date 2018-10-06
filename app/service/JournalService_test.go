package service

import (
	"testing"
)

func TestJournalService_ListUserJournalLeveBrokerage(t *testing.T) {
	service := JournalService{}

	userservice := UserService{}

	UserID := uint64(1000)

	leveaIDs := userservice.Leve1(UserID)

	service.ListUserJournalLeveBrokerage(UserID, leveaIDs)

}
