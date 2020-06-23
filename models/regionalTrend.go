package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// RegionalTrend struttura dati per la gestione di trend regionale del Covid-19
type RegionalTrend struct {
	Data                     string         `json:"data"`
	Stato                    string         `json:"stato"`
	CodiceRegione            string         `json:"codice_regione"`
	DenominazioneRegione     string         `json:"denominazione_regione"`
	Lat                      float32        `json:"lat"`
	Long                     float32        `json:"long"`
	RicoveratiConSintomi     int            `json:"ricoverati_con_sintomi"`
	TerapiaIntensiva         int            `json:"terapia_intensiva"`
	TotaleOspedalizzati      int            `json:"totale_ospedalizzati"`
	IsolamentoDomiciliare    int            `json:"isolamento_domiciliare"`
	TotalePositivi           int            `json:"totale_positivi"`
	VariazioneTotalePositivi int            `json:"variazione_totale_positivi"`
	NuoviPositivi            int            `json:"nuovi_positivi"`
	DimessiGuariti           int            `json:"dimessi_guariti"`
	Deceduti                 int            `json:"deceduti"`
	TotaleCasi               int            `json:"totale_casi"`
	Tamponi                  int            `json:"tamponi"`
	CasiTestati              sql.NullInt64  `json:"casi_testati"`
	NoteIT                   sql.NullString `json:"note_it"`
	NoteEN                   sql.NullString `json:"note_en"`
}

// RegionalTrendCollect struttura dati per la gestione di trend regionale filtrato del Covid-19
type RegionalTrendCollect struct {
	Data time.Time       `json:"data"`
	Info json.RawMessage `json:"info"`
}
