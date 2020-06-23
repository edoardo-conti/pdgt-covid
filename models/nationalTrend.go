package models

import "database/sql"

// NationalTrend struttura dati per la gestione di trend nazionale del Covid-19
type NationalTrend struct {
	Data                     string         `json:"data"`
	Stato                    string         `json:"stato"`
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

// NationalTrendPOST struttura dati per la gestione di richieste POST di trend nazionale del Covid-19
type NationalTrendPOST struct {
	Data                     string `json:"data" binding:"required"`
	RicoveratiConSintomi     string `json:"ricoverati_con_sintomi" binding:"required"`
	TerapiaIntensiva         string `json:"terapia_intensiva" binding:"required"`
	TotaleOspedalizzati      string `json:"totale_ospedalizzati" binding:"required"`
	IsolamentoDomiciliare    string `json:"isolamento_domiciliare" binding:"required"`
	TotalePositivi           string `json:"totale_positivi" binding:"required"`
	VariazioneTotalePositivi string `json:"variazione_totale_positivi" binding:"required"`
	NuoviPositivi            string `json:"nuovi_positivi" binding:"required"`
	DimessiGuariti           string `json:"dimessi_guariti" binding:"required"`
	Deceduti                 string `json:"deceduti" binding:"required"`
	TotaleCasi               string `json:"totale_casi" binding:"required"`
	Tamponi                  string `json:"tamponi" binding:"required"`
	CasiTestati              string `json:"casi_testati" binding:"required"`
}

// NationalTrendPATCH struttura dati per la gestione di richieste PATCH di trend nazionale del Covid-19
type NationalTrendPATCH struct {
	RicoveratiConSintomi     *int `json:"ricoverati_con_sintomi,omitempty"`
	TerapiaIntensiva         *int `json:"terapia_intensiva,omitempty"`
	TotaleOspedalizzati      *int `json:"totale_ospedalizzati,omitempty"`
	IsolamentoDomiciliare    *int `json:"isolamento_domiciliare,omitempty"`
	TotalePositivi           *int `json:"totale_positivi,omitempty"`
	VariazioneTotalePositivi *int `json:"variazione_totale_positivi,omitempty"`
	NuoviPositivi            *int `json:"nuovi_positivi,omitempty"`
	DimessiGuariti           *int `json:"dimessi_guariti,omitempty"`
	Deceduti                 *int `json:"deceduti,omitempty"`
	TotaleCasi               *int `json:"totale_casi,omitempty"`
	Tamponi                  *int `json:"tamponi,omitempty"`
	CasiTestati              *int `json:"casi_testati,omitempty"`
}
