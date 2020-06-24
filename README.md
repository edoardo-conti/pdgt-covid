# 🧪 pdgt-covid 📊 #
[![Build Status](https://travis-ci.org/edoardo-conti/pdgt-covid.svg?branch=master)](https://travis-ci.org/edoardo-conti/pdgt-covid)
[![Deploy](https://heroku-badge.herokuapp.com/?app=pdgt-covid)](https://pdgt-covid.herokuapp.com/)
[![Go Report Card](https://goreportcard.com/badge/github.com/edoardo-conti/pdgt-covid)](https://goreportcard.com/report/github.com/edoardo-conti/pdgt-covid)
![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/edoardo-conti/pdgt-covid/master)

# Progetto Piattaforme Digitali per la Gestione del Territorio 

* Secondo appello sessione estiva 2019/2020
* [Edoardo Conti - 278717](https://github.com/edoardo-conti)

------------------------------------------

### Introduzione ###

Progetto finalizzato alla realizzazione di un Web Server RESTfull con lo scopo di erogare API per garantire fruizione e manipolazione di dati relativi all'andamento del Covid-19 in Italia. Il sistema prevede due strati di sicurezza: autenticazione ed autorizzazione. 

Gli obiettivi principali del progetto sono di seguito riportati: 
- **Trend Nazionale**
  - Visualizzazione trend nazionale
  - Visualizzazione picco di nuovi positivi a livello nazionale 
  - Ricerca rilevazione trend giornaliero nazionale filtrato per data
  - Aggiunta, Modifica e Rimozione rilevazione trend giornaliero nazionale 
- **Trend Regionale**
  - Visualizzazione trend regionale
  - Visualizzazione picco di nuovi positivi a livello regionale 
  - Ricerca trend regionale filtrato per data o regione
  - Ricerca picco nuovi positivi filtrato per regione
- **Utenza**
  - Visualizzazione lista utenti registrati
  - Ricerca utente per username
  - Registrazione ed Accesso utente al sistema
  - Rimozione utente dal sistema

------------------------------------------

### Architettura e Scelte Implementative ###

Il progetto è stato sviluppato seguendo un approccio client-server in quanto si è scelto di sviluppare un applicativo dimostrativo per la fruizione delle API. Il software è inoltre diviso in blocchi dalle funzionalità ben distinte rispettando il pattern architetturale **Model-View-Controller** (MVC). Da specificare che in questo caso la parte grafica è affidata al client, per tanto il server implementerà le parti di Modellazione delle strutture dati e Controller per la gestione delle funzionalità. 

Per quanto concerne il lato server, si è sposata la scelta del linguaggio di programmazione open source [Go](https://golang.org). Le motivazioni di tale scelta ricadono nell'efficienza di scrittura di codice che permettono la scrittura di software semplici ed affidabili, presenza di framework moderni per l'approccio alle comunicazioni HTTP ma sopratutto: cogliere l'occasione per imparare questo nuovo linguaggio di programmazione che da tempo mi incuriosiva. 
La gestione delle richieste HTTP è stata affidata a [Gin Web Framework](https://github.com/gin-gonic/gin) il quale vanta performance *40x* superiori rispetto a *HttpRouter*, un multiplexer alternativo sempre scritto in Go. Per installare il package Gin ed impostarlo nel workspace è sufficiente: 
1. Utilizzare il seguente comando per installare Gin (necessario Go v1.11+).
    ```sh
    $ go get -u github.com/gin-gonic/gin
    ```
2. Importare Gin nel progetto.
    ```go
    import "github.com/gin-gonic/gin"
    ```

###### Database ######
La scelta di come mantenere e manipolare i dati interessati, dopo diverse valutazioni, è ricaduta su Heroku Postgres: uno dei DBMS più popolari al mondo basato su SQL. Successivamente si discuterà della scelta d'utilizzo di Heroku per effettuare il deployment dell'applicativo, al momento è sufficiente essere a conoscenza della presenza di tale servizio. Per aggiungere l'addon Postgresql con piano hobby-dev gratuito entro soglie definite nei termini dei servizi alla propria App occorre inviare il comando sottostante (richiede Heroku CLI):
```sh
$ heroku addons:create heroku-postgresql:hobby-dev
Creating heroku-postgresql:hobby-dev on ⬢ pdgt-covid... free
Created postgresql-concentric-70860 as DATABASE_URL
```
Così facendo si avrà a disposizione la variabile d'ambiente Heroku `DATABASE_URL` che verrà sfruttata per stabilire una connessione con il database per visualizzazione e manipolare i dati relativi all'andamento del Covid-19 in Italia. Di seguito è riportato il comando per collegarsi alla console psql dell'addon Heroku (previo accesso tramite `$ heroku login`) e la lista di tabelle presenti nel database.
```sh
$ heroku pg:psql --app pdgt-covid
--> Connecting to postgresql-concentric-70860
psql (12.3)
SSL connection (protocol: TLSv1.2, cipher: ECDHE-RSA-AES256-GCM-SHA384, bits: 256, compression: off)
Type "help" for help.

pdgt-covid::DATABASE=> \dt
             List of relations
 Schema |  Name   | Type  |     Owner      
--------+---------+-------+----------------
 public | nazione | table | smmnpqlyusgxdh
 public | regioni | table | smmnpqlyusgxdh
 public | utenti  | table | smmnpqlyusgxdh
(3 rows)
```
La tabella nazione è dedicata allo storage dei trend nazionali del Covid-19. In seguito si riporta lo schema della tabella SQL e più sotto l'analoga struttura dati Go `NationalTrend` per interfacciarsi con il DB. Da notare la consistenza del nome dei campi per evitare misspelling durante la stesura del codice e migliorare la leggibilità stessa.
```sh
pdgt-covid::DATABASE=> \d+ nazione
                 Table "public.nazione"                           
           Column           |          Type          
----------------------------+------------------------
 data                       | date                   
 stato                      | character varying(3)   
 ricoverati_con_sintomi     | integer                
 terapia_intensiva          | integer                
 totale_ospedalizzati       | integer                
 isolamento_domiciliare     | integer                
 totale_positivi            | integer                
 variazione_totale_positivi | integer                
 nuovi_positivi             | integer                
 dimessi_guariti            | integer               
 deceduti                   | integer              
 totale_casi                | integer                
 tamponi                    | integer                
 casi_testati               | integer               
 note_it                    | character varying(255) 
 note_en                    | character varying(255)
```
```go
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
```
Onde evitare ridondanza e mantenere il Readme il più pulito e coinciso possibile, gli schemi e strutture delle tabelle `regioni` ed `utenti` non verranno riportate.

###### Sistema di Autenticazione ed Autorizzazione ######
La sicurezza sulla modifica dei dati archiviati è gestita su due layer: autenticazione ed autorizzazione. Il sistema è stato pensato ponendo dei limiti di lettura e scrittura su utenti non autenticati. La registrazione di un nuovo utente è effettuabile solamente da un utente già registrato nel database per evitare spiacevoli inconvenienti. Per quanto riguarda le autorizzazioni, un utente registrato gode di privilegi di lettura superiori rispetto ad un visitatore e vanta la possibilità di effettuare richieste http *POST*. Un Admin possiede tutti i privilegi di lettura e modifica, per tanto oltre a quanto possibile ad un Utente può effettuare richieste http *PATCH* e *DELETE*. Di seguito è riportata una tabella che riassume i permessi delle API:

Visitatore | Utente | Admin
------------ | ------------- | -------------
`GET /api/trend/*` | `GET /api/trend/*` | `GET /api/trend/*`
\- | `GET /api/utenti/*` | `GET /api/utenti/*`
`POST /api/utenti/signin` | `POST /api/*` | `POST /api/*`
\- | - | `PATCH /api/*`
\- | - | `DELETE /api/*`

Il login di un utente è verificato tramite **JSON Web Token (JWT)**, uno standard open che definisce uno schema JSON per lo scambio di informazioni tra vari servizi. Il token generato verrà firmato con una chiave segreta impostata come variabile d'ambiente in Heroku (`JWT_ACCESS_SECRET`) tramite l'algoritmo HMAC. Durante la fase di login le credenziali vengono cryptate secondo l'algoritmo di hashing **bcrypt** evitando di esporre password in chiaro. Ergo nel database viene salvato unicamente l'hash della password. Lato server sfruttando il metodo `CompareHashAndPassword(...)` della libreria `bcrypt` si comparerà la password in chiaro proposta dall'utente e l'hash presente nel database. 

###### Client ######
Per poter sfruttare le API messe a disposizione dal web service in questione si è scelto di sviluppare una Web App con tecnologia *React*, una libreria javascript per creare interfacce utente moderne.
L'obiettivo imposto è stato quello di offrire un'interfaccia grafica per interrogare il web service fino ad ora discusso. Il design della GUI è stato realizzato con *material-ui*, componente React per web development semplice e rapido con stile Material by Google.

Per creare il workspace di una Web App React è sufficiente un semplice comando:
```sh
$ npx create-react-app pdgt-covid-webapp
```
Terminato il processo si avrà, senza aver battuto una singola riga di codice, una **Single Sage Application (SPA)** con tutte le directory e files di default predisposti pronti ad essere modificati per realizzare la propria applicazione.
Il principale componente sfruttato per mostrare e manipolare i dati a disposizione è `material-table`, una semplice ma molto potente data-table per React basata su [Material-UI Table](https://material-ui.com/components/tables/#table).
L'installazione prevede un singolo comando `npm` :
```sh
$ npm install material-table @material-ui/core --save
```
Le comunicazioni con protocollo HTTP lato client sono effettuate sfruttando `axios`, un client HTTP promise-api-based per nodejs.
Maggiori informazioni circa il funzionamento della web app sono disponibili nella sezione dedicata con tanto di screenshots del funzionamento. 

------------------------------------------

### Dati e Servizi Esterni ###

Il recupero dei dati relativi all'andamento del Covid-19 in Italia è stato ricavato dal *Dipartimento della Protezione Civile (DPC)*  sotto licenza *Creative Commons Attribution 4.0 International* per mezzo di file .csv ospitati nel repository pubblico GitHub [pcm-dpc/COVID-19](https://github.com/pcm-dpc/COVID-19). Per correttezza si riporta un estratto di quando specificato nella [licenza](https://github.com/pcm-dpc/COVID-19/blob/master/LICENSE): 
>       1. Subject to the terms and conditions of this Public License,
>          the Licensor hereby grants You a worldwide, royalty-free,
>          non-sublicensable, non-exclusive, irrevocable license to
>          exercise the Licensed Rights in the Licensed Material to:
>
>            a. reproduce and Share the Licensed Material, in whole or
>               in part; and
>
>            b. produce, reproduce, and Share Adapted Material.

La realizzazione della mappa interattiva accessibile dalla web app, che illustra la diffusione del Covid-19 in italia mediante cerchi di dimensione in scala direttamente proporzionale al totale dei casi della regione, è possibile grazie alle API di Google Maps. [Termini di servizio](https://cloud.google.com/maps-platform/terms?_ga=2.90427935.407167450.1593004949-1198461100.1591999667).

L'immagine degli avatar utente è ricavata sfruttando l'API del servizio [DiceBear Avatars](https://avatars.dicebear.com). Il servizio permette di effettuare richieste HTTP richiedendo diverse tipologie ed immagini di avatar in base alle preferenze imposte mediante il path dell'API. [Termini di servizio](https://github.com/DiceBear/avatars/blob/v4/LICENSE).

------------------------------------------

### Documentazione API ###

###### Trend Nazionale ######

* **Visualizzare tutti i trend nazionali**

`GET https://pdgt-covid.herokuapp.com/api/trend/nazionale`
```json
{
  "count": 106,
  "data": [
    {
      "data": "2020-02-23T00:00:00Z",
      "stato": "ITA",
      "ricoverati_con_sintomi": 4,
      "terapia_intensiva": 1,
      "totale_ospedalizzati": 1,
      "isolamento_domiciliare": 1,
      "totale_positivi": 1,
      "variazione_totale_positivi": -5,
      "nuovi_positivi": 1,
      "dimessi_guariti": 0,
      "deceduti": 0,
      "totale_casi": 1,
      "tamponi": 10,
      "casi_testati": {
        "Int64": 200,
        "Valid": true
      },
      "note_it": {
        "String": "",
        "Valid": false
      },
      "note_en": {
        "String": "",
        "Valid": false
      }
    },
    {
      "data": "2020-02-24T00:00:00Z",
      "stato": "ITA",
      "ricoverati_con_sintomi": 101,
      "terapia_intensiva": 26,
      "totale_ospedalizzati": 127,
      "isolamento_domiciliare": 94,
      "totale_positivi": 221,
      "variazione_totale_positivi": 0,
      "nuovi_positivi": 221,
      "dimessi_guariti": 1,
      "deceduti": 7,
      "totale_casi": 229,
      "tamponi": 4324,
      "casi_testati": {
        "Int64": 0,
        "Valid": false
      },
      "note_it": {
        "String": "",
        "Valid": false
      },
      "note_en": {
        "String": "",
        "Valid": false
      }
    },
    [...]
  ]
}
```

* **Visualizzare il trend nazionale filtrato per data**

`GET https://pdgt-covid.herokuapp.com/api/trend/nazionale/data/:bydate`
```json
// https://pdgt-covid.herokuapp.com/api/trend/nazionale/data/2020-02-24
{
  "data": {
    "data": "2020-02-24T00:00:00Z",
    "stato": "ITA",
    "ricoverati_con_sintomi": 101,
    "terapia_intensiva": 26,
    "totale_ospedalizzati": 127,
    "isolamento_domiciliare": 94,
    "totale_positivi": 221,
    "variazione_totale_positivi": 0,
    "nuovi_positivi": 221,
    "dimessi_guariti": 1,
    "deceduti": 7,
    "totale_casi": 229,
    "tamponi": 4324,
    "casi_testati": {
      "Int64": 0,
      "Valid": false
    },
    "note_it": {
      "String": "",
      "Valid": false
    },
    "note_en": {
      "String": "",
      "Valid": false
    }
  },
  "status": 200
}
```

* **Visualizzare il picco di nuovi positivi in Italia**

`GET https://pdgt-covid.herokuapp.com/api/trend/nazionale/picco`
```json
{
  "data": {
    "data": "2020-03-21T00:00:00Z",
    "stato": "ITA",
    "ricoverati_con_sintomi": 17708,
    "terapia_intensiva": 2857,
    "totale_ospedalizzati": 20565,
    "isolamento_domiciliare": 22116,
    "totale_positivi": 42681,
    "variazione_totale_positivi": 4821,
    "nuovi_positivi": 6557,
    "dimessi_guariti": 6072,
    "deceduti": 4825,
    "totale_casi": 53578,
    "tamponi": 233222,
    "casi_testati": {
      "Int64": 0,
      "Valid": false
    },
    "note_it": {
      "String": "",
      "Valid": false
    },
    "note_en": {
      "String": "",
      "Valid": false
    }
  },
  "status": 200
}
```

* **Inserimento di un nuovo trend giornaliero nazionale** `(Richiesta Autenticazione)`

`POST https://pdgt-covid.herokuapp.com/api/trend/nazionale`
```json
// richiesta
// Authorization <token>
{
        "data": "2020-02-20",
        "ricoverati_con_sintomi": "178",
        "terapia_intensiva": "57",
        "totale_ospedalizzati": "265",
        "isolamento_domiciliare": "116",
        "totale_positivi": "481",
        "variazione_totale_positivi": "21",
        "nuovi_positivi": "67",
        "dimessi_guariti": "434",
        "deceduti": "425",
        "totale_casi": "578",
        "tamponi": "222",
        "casi_testati": "213"
}
```
```json
// risposta
{
    "info": "/api/trend/nazionale/data/2020-02-20",
    "message": "Trend giornaliero nazionale registrato con successo.",
    "status": 200
}
```

* **Aggiornamento di un trend giornaliero nazionale esistente** `(Richiesta Autorizzazione)`

`PATCH https://pdgt-covid.herokuapp.com/api/trend/nazionale/data/:bydate`
```json
// richiesta
// https://pdgt-covid.herokuapp.com/api/trend/nazionale/data/2020-02-20
// Authorization <token>
{
        "terapia_intensiva": 100,
        "totale_ospedalizzati": 300,
        "variazione_totale_positivi": -10,
        "tamponi": 101
}
```
```json
// risposta
{
    "info": "/api/trend/nazionale/data/2020-02-20",
    "message": "Trend in data 2020-02-20 aggiornato con successo.",
    "status": 200
}
```

* **Eliminazione di un trend giornaliero nazionale esistente** `(Richiesta Autorizzazione)`

`DELETE https://pdgt-covid.herokuapp.com/api/trend/nazionale/data/:bydate`
```json
// richiesta
// https://pdgt-covid.herokuapp.com/api/trend/nazionale/data/2020-02-20
// Authorization <token>
```
```json
// risposta
{
    "message": "Trend in data 2020-02-20 eliminato dal database con successo.",
    "status": 200
}
```

###### Trend Regionale ######

* **Visualizzare tutti i trend regionali**

`GET https://pdgt-covid.herokuapp.com/api/trend/regionale`
```json
{
    "count": 105,
    "data": [
        {
            "data": "2020-02-24T00:00:00Z",
            "info": [
                {
                    "stato": "ITA",
                    "codice_regione": 13,
                    "denominazione_regione": "Abruzzo",
                    "lat": 42.35122196,
                    "long": 13.39843823,
                    "ricoverati_con_sintomi": 0,
                    "terapia_intensiva": 0,
                    "totale_ospedalizzati": 0,
                    "isolamento_domiciliare": 0,
                    "totale_positivi": 0,
                    "variazione_totale_positivi": 0,
                    "nuovi_positivi": 0,
                    "dimessi_guariti": 0,
                    "deceduti": 0,
                    "totale_casi": 0,
                    "tamponi": 5,
                    "casi_testati": null,
                    "note_it": null,
                    "note_en": null
                },
                {
                    "stato": "ITA",
                    "codice_regione": 17,
                    "denominazione_regione": "Basilicata",
                    "lat": 40.63947052,
                    "long": 15.80514834,
                    "ricoverati_con_sintomi": 0,
                    "terapia_intensiva": 0,
                    "totale_ospedalizzati": 0,
                    "isolamento_domiciliare": 0,
                    "totale_positivi": 0,
                    "variazione_totale_positivi": 0,
                    "nuovi_positivi": 0,
                    "dimessi_guariti": 0,
                    "deceduti": 0,
                    "totale_casi": 0,
                    "tamponi": 0,
                    "casi_testati": null,
                    "note_it": null,
                    "note_en": null
                },
                [...]
            ]
        },
        {
            "data": "2020-02-25T00:00:00Z",
            "info": [...]
        }
    ]
}
```

* **Visualizzare tutti i trend regionali filtrati per data**

`GET https://pdgt-covid.herokuapp.com/api/trend/regionale/data/:bydata`
```json
// https://pdgt-covid.herokuapp.com/api/trend/regionale/data/2020-02-24
{
    "count": 1,
    "data": [
        {
            "data": "2020-02-24T00:00:00Z",
            "info": [
                {
                    "stato": "ITA",
                    "codice_regione": 13,
                    "denominazione_regione": "Abruzzo",
                    "lat": 42.35122196,
                    "long": 13.39843823,
                    "ricoverati_con_sintomi": 0,
                    "terapia_intensiva": 0,
                    "totale_ospedalizzati": 0,
                    "isolamento_domiciliare": 0,
                    "totale_positivi": 0,
                    "variazione_totale_positivi": 0,
                    "nuovi_positivi": 0,
                    "dimessi_guariti": 0,
                    "deceduti": 0,
                    "totale_casi": 0,
                    "tamponi": 5,
                    "casi_testati": null,
                    "note_it": null,
                    "note_en": null
                },
                {
                    "stato": "ITA",
                    "codice_regione": 17,
                    "denominazione_regione": "Basilicata",
                    "lat": 40.63947052,
                    "long": 15.80514834,
                    "ricoverati_con_sintomi": 0,
                    "terapia_intensiva": 0,
                    "totale_ospedalizzati": 0,
                    "isolamento_domiciliare": 0,
                    "totale_positivi": 0,
                    "variazione_totale_positivi": 0,
                    "nuovi_positivi": 0,
                    "dimessi_guariti": 0,
                    "deceduti": 0,
                    "totale_casi": 0,
                    "tamponi": 0,
                    "casi_testati": null,
                    "note_it": null,
                    "note_en": null
                },
                [...]
            ]
        }
    ]
}
```

* **Visualizzare tutti i trend regionali filtrati per regione**

`GET https://pdgt-covid.herokuapp.com/api/trend/regionale/regione/:byregid`
```json
// https://pdgt-covid.herokuapp.com/api/trend/regionale/regione/11
{
    "count": 105,
    "data": [
        {
            "data": "2020-02-24T00:00:00Z",
            "info": [
                {
                    "stato": "ITA",
                    "codice_regione": 11,
                    "denominazione_regione": "Marche",
                    "lat": 43.61675973,
                    "long": 13.5188753,
                    "ricoverati_con_sintomi": 0,
                    "terapia_intensiva": 0,
                    "totale_ospedalizzati": 0,
                    "isolamento_domiciliare": 0,
                    "totale_positivi": 0,
                    "variazione_totale_positivi": 0,
                    "nuovi_positivi": 0,
                    "dimessi_guariti": 0,
                    "deceduti": 0,
                    "totale_casi": 0,
                    "tamponi": 16,
                    "casi_testati": null,
                    "note_it": null,
                    "note_en": null
                }
            ]
        },
        {
            "data": "2020-02-25T00:00:00Z",
            "info": [
                {
                    "stato": "ITA",
                    "codice_regione": 11,
                    "denominazione_regione": "Marche",
                    "lat": 43.61675973,
                    "long": 13.5188753,
                    "ricoverati_con_sintomi": 0,
                    "terapia_intensiva": 0,
                    "totale_ospedalizzati": 0,
                    "isolamento_domiciliare": 0,
                    "totale_positivi": 0,
                    "variazione_totale_positivi": 0,
                    "nuovi_positivi": 0,
                    "dimessi_guariti": 0,
                    "deceduti": 0,
                    "totale_casi": 0,
                    "tamponi": 21,
                    "casi_testati": null,
                    "note_it": null,
                    "note_en": null
                }
            ]
        },
        [...]
    ]
}      
```

* **Visualizzare il picco di nuovi positivi più alto tra tutti i trend regionali**

`GET https://pdgt-covid.herokuapp.com/api/trend/regionale/picco/`
```json
{
    "count": 1,
    "data": [
        {
            "data": "2020-03-21T00:00:00Z",
            "info": [
                {
                    "stato": "ITA",
                    "codice_regione": 3,
                    "denominazione_regione": "Lombardia",
                    "lat": 45.46679409,
                    "long": 9.190347404,
                    "ricoverati_con_sintomi": 8258,
                    "terapia_intensiva": 1093,
                    "totale_ospedalizzati": 9351,
                    "isolamento_domiciliare": 8019,
                    "totale_positivi": 17370,
                    "variazione_totale_positivi": 1950,
                    "nuovi_positivi": 3251,
                    "dimessi_guariti": 5050,
                    "deceduti": 3095,
                    "totale_casi": 25515,
                    "tamponi": 66730,
                    "casi_testati": null,
                    "note_it": null,
                    "note_en": null
                }
            ]
        }
    ],
    "status": 200
}
```

* **Visualizzare il picco di nuovi positivi più alto tra tutti i trend regionali filtrato per regione**

`GET https://pdgt-covid.herokuapp.com/api/trend/regionale/picco/:byregid`
```json
// https://pdgt-covid.herokuapp.com/api/trend/regionale/picco/11
{
    "count": 1,
    "data": [
        {
            "data": "2020-03-22T00:00:00Z",
            "info": [
                {
                    "stato": "ITA",
                    "codice_regione": 11,
                    "denominazione_regione": "Marche",
                    "lat": 43.61675973,
                    "long": 13.5188753,
                    "ricoverati_con_sintomi": 816,
                    "terapia_intensiva": 138,
                    "totale_ospedalizzati": 954,
                    "isolamento_domiciliare": 1277,
                    "totale_positivi": 2231,
                    "variazione_totale_positivi": 234,
                    "nuovi_positivi": 268,
                    "dimessi_guariti": 6,
                    "deceduti": 184,
                    "totale_casi": 2421,
                    "tamponi": 6391,
                    "casi_testati": null,
                    "note_it": null,
                    "note_en": null
                }
            ]
        }
    ],
    "status": 200
}
```

###### Utenti ######

* **Visualizzare tutti gli utenti registrati** `(Richiesta Autenticazione)`

`GET https://pdgt-covid.herokuapp.com/api/utenti/`
```json
// richiesta
// Authorization <token>
{
    "count": 4,
    "data": [
        {
            "username": "edoardo",
            "password": "$2a$10$ZLeA29AWNfETxs1.ymJlLeKfeTGQ54dOUjzWnSuT4jMa1VdmeXF3.",
            "is_admin": true,
            "avatar_url": "https://avatars.dicebear.com/api/initials/e.svg"
        },
        {
            "username": "professore",
            "password": "$2a$10$l1zAnEcs6xadNhhaUkcrP.mJ1mK7VmrXqtUKVyncsnd.drkoPQPjC",
            "is_admin": true,
            "avatar_url": "https://avatars.dicebear.com/api/initials/p.svg"
        },
        [...]
    ],
    "status": 200
}
```

* **Visualizzare utente registrato filtrato per username** `(Richiesta Autenticazione)`

`GET https://pdgt-covid.herokuapp.com/api/utenti/:byusername`
```json
// richiesta
// https://pdgt-covid.herokuapp.com/api/utenti/test
// Authorization <token>
// 
{
    "data": {
        "username": "test",
        "password": "$2a$10$mtO95KBoLCAunUy5TzgM9uGEwKSyWHNas6JjR63ofzkEqgpIif1J6",
        "is_admin": false,
        "avatar_url": "https://avatars.dicebear.com/api/initials/t.svg"
    },
    "status": 200
}
```

* **Registrazione utente nel database** `(Richiesta Autenticazione)`

`POST https://pdgt-covid.herokuapp.com/api/utenti/signup`
```json
// richiesta
// Authorization <token>
{
    "username":"mario",
    "password":"segreto",
    "is_admin":false
}
```
```json
// risposta
{
    "info": "Per visualizzare: /utenti/mario",
    "message": "Utente registrato con successo.",
    "status": 200
}
```

* **Accesso utente nel sistema**

`POST https://pdgt-covid.herokuapp.com/api/utenti/signin`
```json
// richiesta
{
    "username":"test",
    "password":"test"
}
```
```json
// risposta
{
    "message": "Utente test loggato con successo.",
    "status": 200,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6ZmFsc2UsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTU5MzAwOTAyMCwidXNlcm5hbWUiOiJ0ZXN0In0.4wCtltOGmB9G0JIY3LQ1PAH1M22U7bdUn5nmueHy7aE"
}
```

* **Eliminazione utente dal sistema per username** `(Richiesta Autorizzazione)`

`DELETE https://pdgt-covid.herokuapp.com/api/utenti/:byusername`
```json
// richiesta
// https://pdgt-covid.herokuapp.com/api/utenti/mario
// Authorization <token>
```
```json
// risposta
{
    "message": "Utente mario eliminato dal database con successo.",
    "status": 200
}
```

------------------------------------------

### Licenza ###
Questo progetto è rilasciato sotto i termini della [licenza MIT](https://github.com/edoardo-conti/pdgt-covid/blob/master/LICENSE).
