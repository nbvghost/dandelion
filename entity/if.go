package entity

import "github.com/nbvghost/dandelion/library/dao"

type SessionMappingData struct {
	OID dao.PrimaryKey
}

func (m *SessionMappingData) GetOID() dao.PrimaryKey {
	return m.OID
}
