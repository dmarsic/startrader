package starsystem

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/abhinav-TB/dantdb"
)

const StoreDir = "./data"

type Systems map[string]*System

type System struct {
	Name     string   `json:"name"`
	Position Position `json:"position"`
	Market   Market   `json:"market"`
	Traders  []string `json:"traders"`
}

type Position struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
	Z int64 `json:"z"`
}

type Market map[string]Item

type Item struct {
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

func ReadAllSystems() (*Systems, error) {
	db, err := dantdb.New(StoreDir)
	if err != nil {
		return nil, err
	}

	records, err := db.ReadAll("systems")
	if err != nil {
		return nil, err
	}

	systems := Systems{}
	for _, r := range records {
		s := System{}
		if err := json.Unmarshal([]byte(r), s); err != nil {
			return nil, err
		}
		systems[s.Name] = &s
	}
	return &systems, nil
}

func ReadSystem(systemID string) (*System, error) {
	db, err := dantdb.New(StoreDir)
	if err != nil {
		return nil, err
	}

	system := &System{}
	err = db.Read("systems", systemID, system)
	if err != nil {
		return &System{}, err
	}

	return system, nil
}

func WriteSystemState(s *System) error {
	db, err := dantdb.New(StoreDir)
	if err != nil {
		return err
	}

	if err := db.Write("systems", s.Name, s); err != nil {
		return err
	}

	return nil
}

func (m *Market) Reduce(item string, quantity float64) error {
	i, ok := (*m)[item]
	if !ok {
		return errors.New(fmt.Sprintf("No item: %s in market", item))
	}

	if i.Quantity < quantity {
		return errors.New(fmt.Sprintf("Not enough %ss in market: required: %.1f, available: %.1f", item, quantity, i.Quantity))
	}

	i.Quantity -= quantity
	(*m)[item] = i

	return nil
}

func RegisterTrader(traderName, systemName string) error {
	s, err := ReadSystem(systemName)
	if err != nil {
		return err
	}

	s.Traders = append(s.Traders, traderName)
	if err := WriteSystemState(s); err != nil {
		return err
	}

	return nil
}

func DeregisterTrader(traderName, systemName string) error {
	s, err := ReadSystem(systemName)
	if err != nil {
		return err
	}

	var pos int
	for i, t := range s.Traders {
		if t == traderName {
			pos = i
			break
		}
	}
	if len(s.Traders) > 0 {
		if pos == len(s.Traders)-1 {
			s.Traders = s.Traders[:pos]
		} else {
			s.Traders = append(s.Traders[:pos], s.Traders[pos+1:]...)
		}
	}
	if err := WriteSystemState(s); err != nil {
		return err
	}
	return nil
}
