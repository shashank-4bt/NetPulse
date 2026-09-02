package memory

import (
	"context"
	"strings"

	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/contract"
	"github.com/shashank-4bt/NetPulse/netpulse-api/internal/storage"
)

func (s *Store) CreateUser(_ context.Context, rec storage.UserRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	email := strings.ToLower(rec.User.Email)
	if _, ok := s.emailIndex[email]; ok {
		return storage.ErrEmailTaken
	}
	copy := rec
	if copy.Alerts.Summary == "" {
		copy.Alerts = contract.EmptyAlerts()
	}
	s.users[rec.User.ID] = copy
	s.emailIndex[email] = rec.User.ID
	return nil
}

func (s *Store) GetUserByID(_ context.Context, id string) (*storage.UserRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.users[id]
	if !ok {
		return nil, nil
	}
	copy := rec
	return &copy, nil
}

func (s *Store) GetUserByEmail(_ context.Context, email string) (*storage.UserRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.emailIndex[strings.ToLower(email)]
	if !ok {
		return nil, nil
	}
	rec := s.users[id]
	copy := rec
	return &copy, nil
}

func (s *Store) UpdateUser(_ context.Context, rec storage.UserRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[rec.User.ID] = rec
	return nil
}

func (s *Store) DeleteUser(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if rec, ok := s.users[id]; ok {
		delete(s.emailIndex, strings.ToLower(rec.User.Email))
	}
	delete(s.users, id)
	delete(s.saved, id)
	for sid, session := range s.sessions {
		if session.UserID == id {
			delete(s.sessionByHash, session.TokenHash)
			delete(s.sessions, sid)
		}
	}
	kept := s.events[:0]
	for _, event := range s.events {
		if event.UserID != id {
			kept = append(kept, event)
		}
	}
	s.events = kept
	for hash, share := range s.shares {
		if share.UserID == id {
			delete(s.shares, hash)
		}
	}
	for did, rec := range s.diagnoses {
		if rec.UserID == id {
			delete(s.diagnoses, did)
			delete(s.measurements, did)
		}
	}
	return nil
}

func (s *Store) CreateSession(_ context.Context, session contract.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.ID] = session
	if session.TokenHash != "" {
		s.sessionByHash[session.TokenHash] = session.ID
	}
	return nil
}

func (s *Store) GetSession(_ context.Context, id string) (*contract.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[id]
	if !ok {
		return nil, nil
	}
	copy := session
	return &copy, nil
}

func (s *Store) GetSessionByTokenHash(_ context.Context, hash string) (*contract.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.sessionByHash[hash]
	if !ok {
		return nil, nil
	}
	session, ok := s.sessions[id]
	if !ok {
		return nil, nil
	}
	copy := session
	return &copy, nil
}

func (s *Store) ListSessions(_ context.Context, userID string) ([]contract.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.Session{}
	for _, session := range s.sessions {
		if session.UserID == userID {
			out = append(out, session)
		}
	}
	return out, nil
}

func (s *Store) UpdateSession(_ context.Context, session contract.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.ID] = session
	if session.TokenHash != "" {
		s.sessionByHash[session.TokenHash] = session.ID
	}
	return nil
}

func (s *Store) RevokeSession(_ context.Context, userID, sessionID string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[sessionID]
	if !ok || session.UserID != userID {
		return false, nil
	}
	session.Revoked = true
	s.sessions[sessionID] = session
	return true, nil
}

func (s *Store) RevokeAllSessions(_ context.Context, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, session := range s.sessions {
		if session.UserID == userID {
			session.Revoked = true
			s.sessions[id] = session
		}
	}
	return nil
}

func (s *Store) CreateToken(_ context.Context, rec storage.TokenRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tokens[rec.Kind+":"+rec.Hash] = rec
	return nil
}

func (s *Store) GetToken(_ context.Context, kind, hash string) (*storage.TokenRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.tokens[kind+":"+hash]
	if !ok {
		return nil, nil
	}
	copy := rec
	return &copy, nil
}

func (s *Store) MarkTokenUsed(_ context.Context, kind, hash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.tokens[kind+":"+hash]
	if !ok {
		return nil
	}
	rec.Used = true
	s.tokens[kind+":"+hash] = rec
	return nil
}

func (s *Store) AddEvent(_ context.Context, event contract.SecurityEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, event)
	return nil
}

func (s *Store) ListEvents(_ context.Context, userID string) ([]contract.SecurityEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []contract.SecurityEvent{}
	for _, event := range s.events {
		if event.UserID == userID {
			out = append(out, event)
		}
	}
	return out, nil
}

func (s *Store) SaveService(_ context.Context, userID string, item contract.SavedService) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	list := s.saved[userID]
	for i, existing := range list {
		if existing.Slug == item.Slug {
			list[i] = item
			s.saved[userID] = list
			return nil
		}
	}
	s.saved[userID] = append(list, item)
	return nil
}

func (s *Store) ListSavedServices(_ context.Context, userID string) ([]contract.SavedService, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]contract.SavedService{}, s.saved[userID]...), nil
}

func (s *Store) DeleteSavedService(_ context.Context, userID, slug string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	list := s.saved[userID]
	kept := list[:0]
	for _, item := range list {
		if item.Slug != slug {
			kept = append(kept, item)
		}
	}
	s.saved[userID] = kept
	return nil
}

func (s *Store) CreateShare(_ context.Context, rec storage.ShareRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.shares[rec.Hash] = rec
	return nil
}

func (s *Store) GetShare(_ context.Context, hash string) (*storage.ShareRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.shares[hash]
	if !ok {
		return nil, nil
	}
	copy := rec
	return &copy, nil
}

func (s *Store) ListSharesByUser(_ context.Context, userID string) ([]storage.ShareRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []storage.ShareRecord{}
	for _, rec := range s.shares {
		if rec.UserID == userID {
			out = append(out, rec)
		}
	}
	return out, nil
}

func (s *Store) DeleteSharesForDiagnosis(_ context.Context, diagnosisID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for hash, rec := range s.shares {
		if rec.DiagnosisID == diagnosisID {
			delete(s.shares, hash)
		}
	}
	return nil
}
