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

//NationalTrendPOST verificare binding:"exists" o binding:"required" ( https://github.com/gin-gonic/gin/issues/491#issuecomment-162330541 )
type NationalTrendPOST struct {
	Data                     string `json:"data" binding:"required"`
	RicoveratiConSintomi     int    `json:"ricoverati_con_sintomi" binding:"required"`
	TerapiaIntensiva         int    `json:"terapia_intensiva" binding:"required"`
	TotaleOspedalizzati      int    `json:"totale_ospedalizzati" binding:"required"`
	IsolamentoDomiciliare    int    `json:"isolamento_domiciliare" binding:"required"`
	TotalePositivi           int    `json:"totale_positivi" binding:"required"`
	VariazioneTotalePositivi int    `json:"variazione_totale_positivi" binding:"required"`
	NuoviPositivi            int    `json:"nuovi_positivi" binding:"required"`
	DimessiGuariti           int    `json:"dimessi_guariti" binding:"required"`
	Deceduti                 int    `json:"deceduti" binding:"required"`
	TotaleCasi               int    `json:"totale_casi" binding:"required"`
	Tamponi                  int    `json:"tamponi" binding:"required"`
	CasiTestati              int    `json:"casi_testati" binding:"required"`
}

//NationalTrendPATCH ...
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
