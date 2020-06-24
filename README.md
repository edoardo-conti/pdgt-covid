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

###### Sistema di Autenticazione ed Autorizzazione ######
La sicurezza sulla modifica dei dati archiviati è gestita su due layer: autenticazione ed autorizzazione. Il sistema è stato pensato ponendo dei limiti di lettura e scrittura su utenti non autenticati. La registrazione di un nuovo utente è effettuabile solamente da un utente già registrato nel database per evitare spiacevoli inconvenienti. Per quanto riguarda le autorizzazioni, un utente registrato gode di privilegi di lettura superiori rispetto ad un visitatore e vanta la possibilità di effettuare richieste http *POST*. Un Admin possiede tutti i privilegi di lettura e modifica, per tanto oltre a quanto possibile ad un Utente può effettuare richieste http *PATCH* e *DELETE*. Di seguito è riportata una tabella che riassume i permessi delle API:

Visitatore | Utente | Admin
------------ | ------------- | -------------
GET /api/trend/* | GET /api/trend/* | GET /api/trend/*
\- | GET /api/utenti/* | GET /api/utenti/*
POST /api/utenti/signin | POST /api/* | POST /api/*
\- | - | PATCH /api/*
\- | - | DELETE /api/*

Il login di un utente è verificato tramite **JSON Web Token** (JWT), uno standard open che definisce uno schema JSON per lo scambio di informazioni tra vari servizi. Il token generato verrà firmato con una chiave segreta impostata come variabile d'ambiente in Heroku (`JWT_ACCESS_SECRET`) tramite l'algoritmo HMAC. Durante la fase di login le credenziali vengono cryptate secondo l'algoritmo di hashing *bcrypt* evitando di esporre password in chiaro.

###### Client ######
Per poter sfruttare le API messe a disposizione dal web service in questione si è scelto di sviluppare una Web App con tecnologia *React*, una libreria javascript per creare interfacce utente moderne.
L'obiettivo imposto è stato quello di offrire un'interfaccia grafica per interrogare il web service fino ad ora discusso. Il design della GUI è stato realizzato con *material-ui*, componente React per web development semplice e rapido con stile Material by Google.

L'inizializzazione di una web app React è possibile grazie ad un semplice e coninciso comando:
```sh
npx create-react-app pdgt-covid-webapp
```


------------------------------------------

### Dati e Servizi Esterni Sfruttati ###
