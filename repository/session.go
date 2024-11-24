package repository

import (
	"a21hc3NpZ25tZW50/db/filebased"
	"a21hc3NpZ25tZW50/model"
	"time"
)

type SessionRepository interface {
	AddSessions(session model.Session) error
	DeleteSession(token string) error
	UpdateSessions(session model.Session) error
	SessionAvailEmail(email string) (model.Session, error)
	SessionAvailToken(token string) (model.Session, error)
	TokenExpired(session model.Session) bool
}

type sessionsRepo struct {
	filebasedDb *filebased.Data
}

func NewSessionsRepo(filebasedDb *filebased.Data) *sessionsRepo {
	return &sessionsRepo{filebasedDb}
}

// AddSessions: Menyimpan sesi baru ke database
func (u *sessionsRepo) AddSessions(session model.Session) error {
	// Gunakan fungsi AddSession dari filebased
	return u.filebasedDb.AddSession(session)
}

// DeleteSession: Menghapus sesi berdasarkan token
func (u *sessionsRepo) DeleteSession(token string) error {
	// Gunakan fungsi DeleteSession dari filebased
	return u.filebasedDb.DeleteSession(token)
}

func (u *sessionsRepo) UpdateSessions(session model.Session) error {
	// Cek apakah sesi dengan email tersebut ada
	existingSession, err := u.filebasedDb.SessionAvailEmail(session.Email)
	if err != nil {
		return err // Jika tidak ditemukan, kembalikan error
	}

	// Hapus sesi lama berdasarkan token
	err = u.filebasedDb.DeleteSession(existingSession.Token)
	if err != nil {
		return err
	}

	// Tambahkan sesi baru
	return u.filebasedDb.AddSession(session)
}

func (u *sessionsRepo) SessionAvailEmail(email string) (model.Session, error) {
	// Gunakan fungsi SessionAvailEmail dari filebased
	return u.filebasedDb.SessionAvailEmail(email)
}

func (u *sessionsRepo) SessionAvailToken(token string) (model.Session, error) {
	// Gunakan fungsi SessionAvailToken dari filebased
	return u.filebasedDb.SessionAvailToken(token)
}

func (u *sessionsRepo) TokenValidity(token string) (model.Session, error) {
	session, err := u.SessionAvailToken(token)
	if err != nil {
		return model.Session{}, err
	}

	if u.TokenExpired(session) {
		err := u.DeleteSession(token)
		if err != nil {
			return model.Session{}, err
		}
		return model.Session{}, err
	}

	return session, nil
}

func (u *sessionsRepo) TokenExpired(session model.Session) bool {
	return session.Expiry.Before(time.Now())
}
