package models

import "database/sql"

//NationalTrend ...
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

//NationalTrendInput ...
type NationalTrendPOST struct {
	Data                     string `json:"data"`
	RicoveratiConSintomi     int    `json:"ricoverati_con_sintomi"`
	TerapiaIntensiva         int    `json:"terapia_intensiva"`
	TotaleOspedalizzati      int    `json:"totale_ospedalizzati"`
	IsolamentoDomiciliare    int    `json:"isolamento_domiciliare"`
	TotalePositivi           int    `json:"totale_positivi"`
	VariazioneTotalePositivi int    `json:"variazione_totale_positivi"`
	NuoviPositivi            int    `json:"nuovi_positivi"`
	DimessiGuariti           int    `json:"dimessi_guariti"`
	Deceduti                 int    `json:"deceduti"`
	TotaleCasi               int    `json:"totale_casi"`
	Tamponi                  int    `json:"tamponi"`
	CasiTestati              int    `json:"casi_testati"`
}

//NationalTrendPATCH ...
type NationalTrendPATCH struct {
	RicoveratiConSintomi     int `json:"ricoverati_con_sintomi"`
	TerapiaIntensiva         int `json:"terapia_intensiva"`
	TotaleOspedalizzati      int `json:"totale_ospedalizzati"`
	IsolamentoDomiciliare    int `json:"isolamento_domiciliare"`
	TotalePositivi           int `json:"totale_positivi"`
	VariazioneTotalePositivi int `json:"variazione_totale_positivi"`
	NuoviPositivi            int `json:"nuovi_positivi"`
	DimessiGuariti           int `json:"dimessi_guariti"`
	Deceduti                 int `json:"deceduti"`
	TotaleCasi               int `json:"totale_casi"`
	Tamponi                  int `json:"tamponi"`
	CasiTestati              int `json:"casi_testati"`
}
